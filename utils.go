package iitkgp_erp_login

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

func check_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}


func output_body(res *http.Response) string {
	body, err := io.ReadAll(res.Body)
	check_error(err)
	return string(body)
}

func is_file(filename string) bool {
	file, err := os.Open(filename)
	defer file.Close()
	return !errors.Is(err, os.ErrNotExist)
}

func generate_token(ctx *context.Context, cancel context.CancelFunc, conf *oauth2.Config) (*oauth2.Token, error) {
	authURL := conf.AuthCodeURL("aadsjf30958") // should be a randomly generated string
	fmt.Println("Visit this URL for authentication: ", authURL)
	browser.OpenURL(authURL)

	var token *oauth2.Token
	var err error

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") == "aadsjf30958" {
			token, err = conf.Exchange(*ctx, r.URL.Query().Get("code"))
		}
		fmt.Fprintf(w, "Authentication complete. Check your terminal.")
		cancel()
	})

	server := &http.Server{Addr: ":7007"}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			check_error(err)
		}
	}()
	<-(*ctx).Done()
	return token, err
}
