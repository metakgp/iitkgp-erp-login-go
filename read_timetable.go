package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// func getAcademicToken(client *http.Client, ssoToken string, sessionToken string) string {
// 	ur, _ := url.Parse("https://erp.iitkgp.ac.in/Academic/getExamOption.htm")

// 	data := url.Values{}
// 	data.Set("ssoToken", ssoToken)

// 	req, err := http.NewRequest("POST", ur.String(), strings.NewReader(data.Encode()))
// 	check_error(err)
// 	defer req.Body.Close()

// 	res, err := client.Do(req)
// 	check_error(err)
// 	defer res.Body.Close()

// 	fmt.Println(res.Header["Set-Cookie"][1])
// 	return res.Header["Set-Cookie"][1]
// }

// func getAcadToken(client *http.Client, ssoToken string, sessionToken string) string {
// 	u, _ := url.Parse("https://erp.iitkgp.ac.in/Acad/timetable_track.jsp?action=first")

// 	data := url.Values{}
// 	data.Set("ssoToken", ssoToken)
// 	data.Set("module_id", "16")
// 	data.Set("menu_id", "245")

// 	req, err := http.NewRequest("POST", u.String(), strings.NewReader(data.Encode()))
// 	check_error(err)
// 	defer req.Body.Close()

// 	res, err := client.Do(req)
// 	check_error(err)
// 	defer res.Body.Close()

// 	return res.Header["Set-Cookie"][1]
// }

func GetTimetable(dept string) string {
	jar, err := cookiejar.New(nil)
	check_error(err)
	client := http.Client{Jar: jar}

	login(&client, true)

	u, _ := url.Parse("https://erp.iitkgp.ac.in/Acad/timetable_track.jsp?action=second&dept=" + dept)

	data := url.Values{}
	data.Set("for_session", "2023-24")
	data.Set("for_semester", "AUTUMN")
	data.Set("dept", dept)

	res, err := client.Post(u.String(), "text/html;charset=UTF-8", strings.NewReader(data.Encode()))
	check_error(err)
	defer res.Body.Close()

	fmt.Println(output_body(res))

	return output_body(res)
}
