package iitkgp_erp_login

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"syscall"
	"time"

	// "github.com/anaskhan96/soup"
	"github.com/go-ping/ping"
	"golang.org/x/term"
)

const logging = true

type LoginDetails struct {
	user_id      string
	password     string
	answer       string
	requestedUrl string
	email_otp    string
}

type ErpCreds struct {
	RollNumber               string            `json:"roll_number"`
	Password                 string            `json:"password"`
	SecurityQuestionsAnswers map[string]string `json:"answers"`
}

func input_creds(client *http.Client, logging bool) LoginDetails {
	loginDetails := LoginDetails{
		requestedUrl: HOMEPAGE_URL,
	}

	if is_file("erpcreds.json") {
		log.Println("Found ERP Credentials file")

		creds_byte, err := os.ReadFile("erpcreds.json")
		check_error(err)

		var erp_creds ErpCreds

		err = json.Unmarshal(creds_byte, &erp_creds)
		check_error(err)

		loginDetails.user_id = erp_creds.RollNumber
		loginDetails.password = erp_creds.Password
		loginDetails.answer = erp_creds.SecurityQuestionsAnswers[get_secret_question(client, erp_creds.RollNumber, logging)]
	} else {
		fmt.Print("Enter Roll No.: ")
		fmt.Scan(&loginDetails.user_id)

		fmt.Print("Enter ERP Password: ")
		byte_password, err := term.ReadPassword(int(syscall.Stdin))
		check_error(err)
		loginDetails.password = string(byte_password)
		fmt.Println()

		fmt.Printf("Your secret question: %s\n", get_secret_question(client, loginDetails.user_id, logging))
		fmt.Print("Enter answer to your secret question: ")
		byte_answer, err := term.ReadPassword(int(syscall.Stdin))
		check_error(err)

		loginDetails.answer = string(byte_answer)
		fmt.Println()
	}
	return loginDetails
}

func get_secret_question(client *http.Client, roll_number string, logging bool) string {
	data := map[string][]string{
		"user_id": {roll_number},
	}

	res, err := client.PostForm(SECRET_QUESTION_URL, data)
	check_error(err)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	check_error(err)

	if logging {
		log.Println("Fetched Security Question")
	}

	return string(body)
}

func is_otp_required() bool {
	pinger, err := ping.NewPinger(PING_URL)
	check_error(err)
	pinger.Count = 1
	pinger.Timeout = time.Duration(4 * float64(time.Second))

	err = pinger.Run()
	check_error(err)

	return pinger.Statistics().PacketsRecv != 1
}

func is_session_alive(client *http.Client, logging bool) (bool, string) {
	if !is_file(".session") {
		return false, ""
	}

	if logging {
		log.Println("Found session file")
		log.Println("Checking session validity...")
	}

	session_byte, err := os.ReadFile(".session")
	check_error(err)
	ssoToken := string(session_byte)

	res, err := client.Get(HOMEPAGE_URL + "?" + ssoToken)
	check_error(err)
	defer res.Body.Close()

	if logging {
		if res.ContentLength != 4145 {
			log.Println("Session valid")
		} else {
			log.Println("Session invalid")
		}
	}

	return res.ContentLength != 4145, ssoToken
}

func ERPSession() (*http.Client, string) {
	jar, err := cookiejar.New(nil)
	check_error(err)
	client := http.Client{Jar: jar}

	var ssoToken string
	var isSession bool
	isSession, ssoToken = is_session_alive(&client, logging)

	if !isSession {
		loginDetails := input_creds(&client, logging)

		if true {
			if logging {
				log.Println("OTP is required")
			}
			loginDetails.email_otp = fetch_otp(&client, loginDetails.user_id, logging)
		}

		data := url.Values{}
		data.Set("user_id", loginDetails.user_id)
		data.Set("password", loginDetails.password)
		data.Set("answer", loginDetails.answer)
		data.Set("requestedUrl", loginDetails.requestedUrl)
		data.Set("email_otp", loginDetails.email_otp)

		res, err := client.PostForm(LOGIN_URL, data)
		check_error(err)
		defer res.Body.Close()

		log.Println("ERP login complete!")
		body, err := io.ReadAll(res.Body)
		check_error(err)

		bodys := string(body)
		i := strings.Index(bodys, "ssoToken")
		ssoToken = bodys[strings.LastIndex(bodys[:i], "\"")+1 : strings.Index(bodys, "ssoToken")+strings.Index(bodys[i:], "\"")]

		err = os.WriteFile(".session", []byte(ssoToken), 0666)
		check_error(err)

	}

	u, err := url.Parse("https://erp.iitkgp.ac.in/")
	check_error(err)

	client.Jar.SetCookies(u, []*http.Cookie{{Name: "ssoToken", Value: ssoToken[9:]}})

	return &client, ssoToken
}
