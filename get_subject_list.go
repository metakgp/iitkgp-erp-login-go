package iitkgp_erp_login

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

func GetSubjectList(client *http.Client, ssoToken string, dept string) string {
	u, _ := url.Parse("https://erp.iitkgp.ac.in/Acad/timetable_track.jsp?action=second&dept=" + dept)

	data := url.Values{}
	data.Set("for_session", "2023-24")
	data.Set("for_semester", "AUTUMN")
	data.Set("dept", dept)

	if Logging {
		log.Println("Getting timetable for", dept)
	}

	res, err := client.Post(u.String(), "text/html;charset=UTF-8", strings.NewReader(data.Encode()))
	check_error(err)
	defer res.Body.Close()

	return output_body(res)
}
