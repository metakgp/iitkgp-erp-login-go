package iitkgp_erp_login

import (
	"bufio"
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

	"github.com/go-ping/ping"
	"golang.org/x/term"
)

const logging = true

type loginDetails struct {
	user_id      string
	password     string
	answer       string
	requestedUrl string
	email_otp    string
}

type erpCreds struct {
	RollNumber               string            `json:"roll_number"`
	Password                 string            `json:"password"`
	SecurityQuestionsAnswers map[string]string `json:"answers"`
}

func input_creds(client *http.Client, logging bool) loginDetails {
	loginParams := loginDetails{
		requestedUrl: HOMEPAGE_URL,
	}

	if is_file("erpcreds.json") {
		log.Println("Found ERP Credentials file")

		creds_file, err := os.Open("erpcreds.json")
		check_error(err)
		defer creds_file.Close()

		scanner := bufio.NewScanner(creds_file)
		scanner.Split(bufio.ScanLines)

		var creds_byte []byte

		for scanner.Scan() {
			creds_byte = append(creds_byte, scanner.Bytes()...)
		}

		var erp_creds erpCreds

		err = json.Unmarshal(creds_byte, &erp_creds)
		check_error(err)

		loginParams.user_id = erp_creds.RollNumber
		loginParams.password = erp_creds.Password
		loginParams.answer = erp_creds.SecurityQuestionsAnswers[get_secret_question(client, erp_creds.RollNumber, logging)]
	} else {
		fmt.Print("Enter Roll No.: ")
		fmt.Scan(&loginParams.user_id)

		fmt.Print("Enter ERP Password: ")
		byte_password, err := term.ReadPassword(int(syscall.Stdin))
		check_error(err)
		loginParams.password = string(byte_password)
		fmt.Println()

		fmt.Printf("Your secret question: %s\n", get_secret_question(client, loginParams.user_id, logging))
		fmt.Print("Enter answer to your secret question: ")
		byte_answer, err := term.ReadPassword(int(syscall.Stdin))
		check_error(err)

		loginParams.answer = string(byte_answer)
		fmt.Println()
	}
	return loginParams
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

func normalize_session_token(rawToken string) (string, string, bool) {
	token := strings.TrimSpace(rawToken)
	if token == "" {
		return "", "", false
	}

	if strings.HasPrefix(token, "ssoToken=") {
		cookieToken := strings.TrimPrefix(token, "ssoToken=")
		if cookieToken == "" {
			return "", "", false
		}
		return token, cookieToken, true
	}

	if strings.Contains(token, "=") {
		return "", "", false
	}

	return "ssoToken=" + token, token, true
}

func clear_session_file(logging bool) {
	err := os.Remove(".session")
	if err != nil {
		if !os.IsNotExist(err) {
			check_error(err)
		}
		return
	}

	if logging {
		log.Println("Cleared stale session file")
	}
}

func set_sso_cookie(client *http.Client, ssoToken string) {
	_, cookieToken, ok := normalize_session_token(ssoToken)
	if !ok {
		check_error(fmt.Errorf("invalid ssoToken format"))
	}

	u, err := url.Parse("https://erp.iitkgp.ac.in/")
	check_error(err)

	client.Jar.SetCookies(u, []*http.Cookie{{Name: "ssoToken", Value: cookieToken}})
}

func is_invalid_session_response(res *http.Response, responseBody string) bool {
	if res.StatusCode >= http.StatusBadRequest {
		return true
	}

	body := strings.ToLower(responseBody)
	finalURL := strings.ToLower(res.Request.URL.String())

	invalidURLMarkers := []string{
		"/ssoadministration/auth.htm",
		"/ssoadministration/logout.htm",
		"/ssoadministration/logoutmsg.htm",
		"/iit_erp3/welcome.jsp",
	}
	for _, marker := range invalidURLMarkers {
		if strings.Contains(finalURL, marker) {
			return true
		}
	}

	invalidBodyMarkers := []string{
		"you are idle for a long time",
		"session has been invalidated due to security reasons",
		"checking idle time, please wait",
		"<title>logout - erp</title>",
		"window.open('logoutmsg.htm','_top')",
		"click here to re-login",
		"click here to relogin",
		"invalid session",
		"name=\"user_id\"",
		"name='user_id'",
		"/ssoadministration/auth.htm",
	}
	for _, marker := range invalidBodyMarkers {
		if strings.Contains(body, marker) {
			return true
		}
	}

	return false
}

func extract_sso_token(responseBody string) (string, bool) {
	start := strings.Index(responseBody, "ssoToken=")
	if start == -1 {
		return "", false
	}

	tokenAndRest := responseBody[start:]
	end := strings.IndexAny(tokenAndRest, "\"'&<>\r\n\t ")
	if end == -1 {
		return tokenAndRest, true
	}

	if end <= len("ssoToken=") {
		return "", false
	}

	return tokenAndRest[:end], true
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
	ssoToken, _, ok := normalize_session_token(string(session_byte))
	if !ok {
		if logging {
			log.Println("Session file is malformed")
		}
		clear_session_file(logging)
		return false, ""
	}

	set_sso_cookie(client, ssoToken)

	res, err := client.Get(HOMEPAGE_URL + "?" + ssoToken)
	check_error(err)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	check_error(err)

	sessionInvalid := is_invalid_session_response(res, string(body))

	if logging {
		if !sessionInvalid {
			log.Println("Session valid")
		} else {
			log.Println("Session invalid")
		}
	}

	if sessionInvalid {
		clear_session_file(logging)
		return false, ""
	}

	return true, ssoToken
}

func ERPSession() (*http.Client, string) {
	jar, err := cookiejar.New(nil)
	check_error(err)
	client := http.Client{Jar: jar}

	var ssoToken string
	var isSession bool
	isSession, ssoToken = is_session_alive(&client, logging)

	if !isSession {
		loginParams := input_creds(&client, logging)

		if true {
			if logging {
				log.Println("OTP is required")
			}
			loginParams.email_otp = fetch_otp(&client, loginParams, logging)
		}

		data := url.Values{}
		data.Set("user_id", loginParams.user_id)
		data.Set("password", loginParams.password)
		data.Set("answer", loginParams.answer)
		data.Set("requestedUrl", loginParams.requestedUrl)
		data.Set("email_otp", loginParams.email_otp)

		res, err := client.PostForm(LOGIN_URL, data)
		check_error(err)
		defer res.Body.Close()

		log.Println("ERP login complete!")
		body, err := io.ReadAll(res.Body)
		check_error(err)

		parsedToken, ok := extract_sso_token(string(body))
		if !ok {
			check_error(fmt.Errorf("failed to extract ssoToken from ERP login response"))
		}
		ssoToken = parsedToken

	}

	ssoToken, _, ok := normalize_session_token(ssoToken)
	if !ok {
		check_error(fmt.Errorf("invalid ssoToken generated during login"))
	}

	if !isSession {
		err = os.WriteFile(".session", []byte(ssoToken), 0666)
		check_error(err)
	}

	set_sso_cookie(&client, ssoToken)

	return &client, ssoToken
}
