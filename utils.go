package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

func check_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func output_body(res http.Response) string {
	body, err := io.ReadAll(res.Body)
	check_error(err)
	return string(body)
}


func is_token_file() bool {
	token_file, err := os.OpenFile(".token", os.O_RDONLY|os.O_CREATE, 0666)
	check_error(err)
	defer token_file.Close()

	token_file_info, err := token_file.Stat()
	check_error(err)

	return token_file_info.Size() != 0
}

func open_browser(url string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Start()
	case "windows", "darwin":
		exec.Command("open", url).Start()
	default:
		fmt.Errorf("unsupported platform")
	}
}
