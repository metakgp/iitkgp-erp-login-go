package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"

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

func main() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{Jar: jar}
	fmt.Println(get_sessiontoken(client, true))
}
