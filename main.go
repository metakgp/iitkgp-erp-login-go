package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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
	req, err := http.NewRequest("GET", HOMEPAGE_URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc := soup.HTMLParse(string(body))
	sessionToken := doc.Find("input", "id", "sessionToken").Attrs()["value"]

	if logging {
		log.Println("Generated sessionToken")
	}

	return sessionToken
}

func get_secret_question(client http.Client, roll_number string, logging bool) string {
	data := url.Values{}
	data.Set("user_id", roll_number)

	req, err := http.NewRequest("POST", SECRET_QUESTION_URL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if logging {
		log.Println("Fetched Security Question")
	}

	return string(body)
}

func get_login_details(roll_number string, password string, secret_answer string, sessionToken string) LoginDetails {
	loginDetails := LoginDetails{
		user_id:      roll_number,
		password:     password,
		answer:       secret_answer,
		sessionToken: sessionToken,
		requestedUrl: HOMEPAGE_URL,
	}

	return loginDetails
}

func is_otp_required() bool {
	pinger, err := ping.NewPinger(PING_URL)
	if err != nil {
		log.Fatal(err)
	}
	pinger.Count = 1
	pinger.Timeout = time.Duration(4 * float64(time.Second))

	err = pinger.Run()
	if err != nil {
		log.Fatal(err)
	}

	return pinger.Statistics().PacketsRecv != 1
}

func request_otp(client http.Client, loginDetails LoginDetails, logging bool) {
	data := url.Values{}
	data.Set("typeee", "SI")
	data.Set("loginid", loginDetails.user_id)
	// data.Set("pass", loginDetails.password) this field seems to be unnecessary according to testing

	req, err := http.NewRequest("POST", OTP_URL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if logging {
		log.Println("Requested OTP")
	}
}

func session_alive(client http.Client) bool {
	req, err := http.NewRequest("GET", WELCOMEPAGE_URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	return res.ContentLength == 1034
}

func main() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{Jar: jar}
	fmt.Println(get_sessiontoken(client, true))
	// fmt.Println(get_secret_question(client, "20CS10020", true))
	// loginDetails := get_login_details("20CS10020", "password", "answer", "token")
	// fmt.Println(is_otp_required())
	// request_otp(client, loginDetails, true)
	// fmt.Println(session_alive(client))
}
