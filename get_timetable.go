package iitkgp_erp_login

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func GetTimetable(client *http.Client, dept string) string {
	// jar, err := cookiejar.New(nil)
	// check_error(err)
	// client := http.Client{Jar: jar}

	u, _ := url.Parse("https://erp.iitkgp.ac.in/Acad/timetable_track.jsp?action=second&dept=" + dept)

	data := url.Values{}
	data.Set("for_session", "2023-24")
	data.Set("for_semester", "AUTUMN")
	data.Set("dept", dept)

	log.Println("Getting timetable for ", dept)

	res, err := client.Post(u.String(), "text/html;charset=UTF-8", strings.NewReader(data.Encode()))
	check_error(err)
	defer res.Body.Close()

	fmt.Println(output_body(res))

	return output_body(res)
}
