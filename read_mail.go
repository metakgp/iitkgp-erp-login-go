package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/knadh/koanf"
	kjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	RedirectURL = "http://localhost:7007"
	query       = "from:erpkgp@adm.iitkgp.ac.in is:unread subject: otp"
)

var env = koanf.New(".")

func request_otp(client http.Client, roll_number string, logging bool) {
	data := map[string][]string{
		"typeee":  {"SI"},
		"loginid": {roll_number},
	}

	res, err := client.PostForm(OTP_URL, data)
	check_error(err)
	defer res.Body.Close()

	if logging {
		log.Println("Requested OTP")
	}
}

func get_msg_id(service *gmail.Service) string {
	results, err := service.Users.Messages.List("me").Q(query).MaxResults(1).Do()
	check_error(err)

	if len(results.Messages) != 0 {
		return results.Messages[0].Id
	}
	return ""
}

func fetch_otp(client *http.Client, roll_number string, logging bool) string {
	if is_file("client_secret.json") || is_file(".token") {
		return fetch_otp_from_mail(client, roll_number, logging)
	} else {
		return fetch_otp_from_input(client, roll_number)
	}
}

func fetch_otp_from_mail(client *http.Client, roll_number string, logging bool) string {

	err := env.Load(file.Provider("client_secret.json"), kjson.Parser())
	check_error(err)

	ctx, cancel := context.WithCancel(context.Background())

	conf := oauth2.Config{
		ClientID:     env.String("installed.client_id"),
		ClientSecret: env.String("installed.client_secret"),
		Scopes:       []string{gmail.GmailReadonlyScope},
		Endpoint:     google.Endpoint,
		RedirectURL:  RedirectURL,
	}

	var token *oauth2.Token

	if is_file(".token") {
		if logging {
			log.Println("Found token file")
		}

		token_byte, err := os.ReadFile(".token")
		check_error(err)

		err = json.Unmarshal(token_byte, &token)
		check_error(err)

	} else {
		token, err = generate_token(&ctx, cancel, &conf)
		check_error(err)

		token_json, err := json.Marshal(*token)
		check_error(err)

		err = os.WriteFile(".token", token_json, 0666)
		check_error(err)
	}

	service, err := gmail.NewService(ctx, option.WithTokenSource(conf.TokenSource(ctx, token)))
	check_error(err)

	latestId := get_msg_id(service)
	request_otp(*client, roll_number, logging)
	var mailId string

	if logging {
		log.Println("Waiting for OTP...")
	}

	for {
		log.Println("...")
		if mailId = get_msg_id(service); mailId != latestId {
			if logging {
				log.Println("OTP fetched")
			}
			break
		}
		time.Sleep(3 * time.Second)
	}

	message, err := service.Users.Messages.Get("me", mailId).Do()
	check_error(err)

	body, err := base64.URLEncoding.DecodeString(message.Payload.Body.Data)
	check_error(err)

	reg := regexp.MustCompile("[0-9]+")
	otp := reg.FindAllString(string(body), -1)[0]

	cancel()
	return otp
}

func fetch_otp_from_input(client *http.Client, roll_number string) string {
	request_otp(*client, roll_number, true)
	var otp string
	fmt.Scan(&otp)
	return otp
}
