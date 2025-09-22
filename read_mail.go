package iitkgp_erp_login

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	RedirectURL = "http://localhost:7007"
	query       = "from:erpkgp@adm.iitkgp.ac.in is:unread subject: otp"
)

type OTPResponse struct {
	Message               string            `json:"msg"`
}

func request_otp(client http.Client, loginParams loginDetails, logging bool) {
	data := url.Values{}
	data.Set("typeee", "SI")
	data.Set("user_id", loginParams.user_id)
	data.Set("password", loginParams.password)
	data.Set("answer", loginParams.answer)

	res, err := client.PostForm(OTP_URL, data)
	check_error(err)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	check_error(err)

	var otpResponse OTPResponse
	err = json.Unmarshal(body, &otpResponse)
	check_error(err)

	if otpResponse.Message == "An OTP(valid for a short time) has been sent to your email id registered with ERP, IIT Kharagpur. Please use that OTP for further processing. "{
		if logging {
			log.Println("Requested OTP")
		}
	}else{
		otpError := fmt.Errorf("failed to request OTP: %w", errors.New(otpResponse.Message))
		check_error(otpError)
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

func fetch_otp(client *http.Client, loginParams loginDetails, logging bool) string {
	if is_file("client_secret.json") || is_file(".token") {
		return fetch_otp_from_mail(client, loginParams, logging)
	} else {
		return fetch_otp_from_input(client, loginParams)
	}
}

func fetch_otp_from_mail(client *http.Client, loginDetails loginDetails, logging bool) string {
	ctx, cancel := context.WithCancel(context.Background())

	conf := oauth2.Config{
		Scopes:      []string{gmail.GmailReadonlyScope},
		Endpoint:    google.Endpoint,
		RedirectURL: RedirectURL,
	}

	secretByte, err := os.ReadFile("client_secret.json")
	check_error(err)

	var secret map[string]map[string]json.RawMessage
	err = json.Unmarshal(secretByte, &secret)
	check_error(err)

	_ = json.Unmarshal(secret["installed"]["client_id"], &conf.ClientID)
	_ = json.Unmarshal(secret["installed"]["client_secret"], &conf.ClientSecret)

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
	request_otp(*client, loginDetails, logging)
	var mailId string

	for {
		if logging {
			log.Println("Waiting for OTP...")
		}
		if mailId = get_msg_id(service); mailId != latestId {
			if logging {
				log.Println("OTP fetched")
			}
			break
		}
		time.Sleep(1 * time.Second)
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

func fetch_otp_from_input(client *http.Client, loginDetails loginDetails) string {
	request_otp(*client, loginDetails, true)
	var otp string
	fmt.Print("Enter OTP: ")
	fmt.Scan(&otp)
	return otp
}
