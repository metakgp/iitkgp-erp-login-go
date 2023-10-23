package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/go-ping/ping"
)

type LoginDetails struct {
	user_id      string
	password     string
	answer       string
	sessionToken string
	requestedUrl string
	email_otp    string
}

type ErpCreds struct {
	ROL_NUMBER                 string
	PASSWORD                   string
	SECURITY_QUESTIONS_ANSWERS map[string]string
}


func get_sessiontoken(client http.Client, logging bool) string {
	res, err := client.Get(HOMEPAGE_URL)
	check_error(err)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	check_error(err)

	doc := soup.HTMLParse(string(body))
	sessionToken := doc.Find("input", "id", "sessionToken").Attrs()["value"]
	if logging {
		log.Println("Generated sessionToken")
	}

	return sessionToken
}

func input_creds(client http.Client) LoginDetails {
	loginDetails := LoginDetails{
		requestedUrl: HOMEPAGE_URL,
	}

	fmt.Print("Enter Roll No.: ")
	fmt.Scan(&loginDetails.user_id)

	fmt.Print("Enter ERP password: ")
	fmt.Scan(&loginDetails.password)

	fmt.Printf("Your secret question: %s\n", get_secret_question(client, loginDetails.user_id, true))
	fmt.Print("Enter answer to your secret question: ")
	fmt.Scan(&loginDetails.answer)

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

func request_otp(client http.Client, roll_number string, logging bool) string {
	data := map[string][]string{
		"typeee":  {"SI"},
		"loginid": {roll_number},
	}
	// data.Set("pass", loginDetails.password) this field seems to be unnecessary according to testing

	res, err := client.PostForm(OTP_URL, data)
	check_error(err)
	defer res.Body.Close()

	if logging {
		log.Println("Requested OTP")
	}

	var otp string
	fmt.Print("Enter OTP: ")
	fmt.Scan(&otp)
	return otp
}

func is_session_alive(client http.Client, logging bool) (bool, string) {
	if logging {
		log.Println("Checking token validity...")
	}

	var ssoToken string
	token_byte, err := os.ReadFile(".token")
	check_error(err)
	ssoToken = string(token_byte)

	res, err := client.Get(HOMEPAGE_URL + "?" + ssoToken)
	check_error(err)
	defer res.Body.Close()

	return res.ContentLength != 4145, ssoToken
}

func Login(logging bool) {
	jar, err := cookiejar.New(nil)
	check_error(err)
	client := http.Client{Jar: jar}

	if is_token_file() {
		if logging {
			log.Println("Found token file!")
		}

		is_session_alive, ssoToken := is_session_alive(client, true)

		if is_session_alive {
			if logging {
				log.Println("Token valid!")
			}
			open_browser(HOMEPAGE_URL+"?"+ssoToken)
			return
		} else {
			if logging {
				log.Println("Token invalid!")
			}
		}

	}

	loginDetails := input_creds(client)
	loginDetails.sessionToken = get_sessiontoken(client, true)

	if is_otp_required() {
		loginDetails.email_otp = request_otp(client, loginDetails.user_id, true)
	}

	data := map[string][]string{
		"user_id":      {loginDetails.user_id},
		"password":     {loginDetails.password},
		"answer":       {loginDetails.answer},
		"sessionToken": {loginDetails.sessionToken},
		"requestedUrl": {loginDetails.requestedUrl},
		"email_otp":    {loginDetails.email_otp},
	}

	res, err := client.PostForm(LOGIN_URL, data)
	check_error(err)
	defer res.Body.Close()

	log.Println("ERP login complete!")

	body, err := io.ReadAll(res.Body)
	check_error(err)

	bodys := string(body)
	i := strings.Index(bodys, "ssoToken")
	ssoToken := bodys[strings.LastIndex(bodys[:i], "\"")+1 : strings.Index(bodys, "ssoToken")+strings.Index(bodys[i:], "\"")]

	err = os.WriteFile(".token", []byte(ssoToken), 0666)
	check_error(err)

	open_browser(HOMEPAGE_URL+"?"+ssoToken)
}

func main() {
	Login(true)
}
