package main

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
	"github.com/pkg/browser"
	"golang.org/x/term"
)

type LoginDetails struct {
	user_id  string
	password string
	answer   string
	// sessionToken string
	requestedUrl string
	email_otp    string
}

type ErpCreds struct {
	RollNumber               string            `json:"roll_number"`
	Password                 string            `json:"password"`
	SecurityQuestionsAnswers map[string]string `json:"answers"`
}

// func get_sessiontoken(client http.Client, logging bool) string {
// 	res, err := client.Get(HOMEPAGE_URL)
// 	check_error(err)
// 	defer res.Body.Close()

// 	body, err := io.ReadAll(res.Body)
// 	check_error(err)

// 	doc := soup.HTMLParse(string(body))
// 	sessionToken := doc.Find("input", "id", "sessionToken").Attrs()["value"]
// 	if logging {
// 		log.Println("Generated sessionToken")
// 	}
// 	return sessionToken
// }

func input_creds(client http.Client, logging bool) LoginDetails {
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

func get_secret_question(client http.Client, roll_number string, logging bool) string {
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

func is_session_alive(client http.Client, logging bool) (bool, string) {
	if logging {
		log.Println("Checking session validity...")
	}

	var ssoToken string
	session_byte, err := os.ReadFile(".session")
	check_error(err)
	ssoToken = string(session_byte)

	res, err := client.Get(HOMEPAGE_URL + "?" + ssoToken)
	check_error(err)
	defer res.Body.Close()

	return res.ContentLength != 4145, ssoToken
}

func Login(logging bool) {
	jar, err := cookiejar.New(nil)
	check_error(err)
	client := http.Client{Jar: jar}

	if is_file(".session") {
		if logging {
			log.Println("Found session file")
		}
		is_session_alive, ssoToken := is_session_alive(client, logging)

		if is_session_alive {
			if logging {
				log.Println("Session valid")
			}
			browser.OpenURL(HOMEPAGE_URL + "?" + ssoToken)
			return
		} else {
			if logging {
				log.Println("Session invalid!")
			}
		}
	}

	loginDetails := input_creds(client, logging)

	if is_otp_required() {
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

	// sessionToken := res.Header["Set-Cookie"][0]

	log.Println("ERP login complete!")

	body, err := io.ReadAll(res.Body)
	check_error(err)

	bodys := string(body)
	i := strings.Index(bodys, "ssoToken")
	ssoToken := bodys[strings.LastIndex(bodys[:i], "\"")+1 : strings.Index(bodys, "ssoToken")+strings.Index(bodys[i:], "\"")]

	err = os.WriteFile(".session", []byte(ssoToken), 0666)
	check_error(err)

	browser.OpenURL(HOMEPAGE_URL + "?" + ssoToken)

	// getTimetable(&client, ssoToken, sessionToken, "CS")

}

func main() {
	Login(true)
}
