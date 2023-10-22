package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os/exec"
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
	data := map[string][]string{
		"user_id": {roll_number},
	}

	res, err := client.PostForm(SECRET_QUESTION_URL, data)
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

func request_otp(client http.Client, roll_number string, logging bool) {
	data := map[string][]string{
		"typeee":  {"SI"},
		"loginid": {roll_number},
	}
	// data.Set("pass", loginDetails.password) this field seems to be unnecessary according to testing

	res, err := client.PostForm(OTP_URL, data)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if logging {
		log.Println("Requested OTP")
	}
}

func session_alive(client http.Client) bool {
	res, err := client.Get(WELCOMEPAGE_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	return res.ContentLength == 1034
}

func login(client http.Client, loginDetails LoginDetails) {
	data := map[string][]string{
		"user_id":      {loginDetails.user_id},
		"password":     {loginDetails.password},
		"answer":       {loginDetails.answer},
		"sessionToken": {loginDetails.sessionToken},
		"requestedUrl": {loginDetails.requestedUrl},
		"email_otp":    {loginDetails.email_otp},
	}
	res, err := client.PostForm(LOGIN_URL, data)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodys := string(body)
	i := strings.Index(bodys, "ssoToken")
	ssoToken := bodys[strings.LastIndex(bodys[:i], "\"")+1 : strings.Index(bodys, "ssoToken")+strings.Index(bodys[i:], "\"")]

	fmt.Println("ERP login complete!")
	err = exec.Command("xdg-open", HOMEPAGE_URL+"?"+ssoToken).Start()
}

func main() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{Jar: jar}

	loginDetails := LoginDetails{
		sessionToken: get_sessiontoken(client, true),
		requestedUrl: HOMEPAGE_URL,
	}

	fmt.Print("Enter Roll No.: ")
	fmt.Scan(&loginDetails.user_id)

	fmt.Print("Enter ERP password: ")
	fmt.Scan(&loginDetails.password)

	fmt.Printf("Your secret question: %s\n", get_secret_question(client, loginDetails.user_id, true))
	fmt.Print("Enter answer to your secret question: ")
	fmt.Scan(&loginDetails.answer)

	if is_otp_required() {
		request_otp(client, loginDetails.user_id, true)
	}

	login(client, loginDetails)

	// fmt.Println(session_alive(client))
}
