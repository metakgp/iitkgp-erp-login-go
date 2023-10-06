package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/anaskhan96/soup"
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

func main() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{Jar: jar}
	fmt.Println(get_sessiontoken(client, true))
	fmt.Println(get_secret_question(client, "20CS10020", true))
}
