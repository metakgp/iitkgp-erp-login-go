package iitkgp_erp_login

import (
	"fmt"
	"net/http"
	"net/url"
)

func GetMyTimetable(client *http.Client, ssoToken string) {
	u, err := url.Parse("https://erp.iitkgp.ac.in/Acad/student/view_stud_time_table.jsp")
	check_error(err)

	res, err := client.Get(u.String())
	check_error(err)
	defer res.Body.Close()

	fmt.Println(output_body(res))
}
