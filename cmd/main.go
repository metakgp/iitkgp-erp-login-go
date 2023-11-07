package main

import (
	erp "iitkgp_erp_login"
	"log"
	"net/http"
)

func main() {
	client, ssoToken := erp.ERPSession(true)
	// browser.OpenURL(erp.HOMEPAGE_URL+"?"+ssoToken)
	getSubjectList(client, ssoToken)
	// erp.GetMyTimetable(client, ssoToken)

}

func getSubjectList(client *http.Client, ssoToken string) {
	depts := [3]string{"CS", "GG", "CE"}

	chann := make(chan string, 1)

	for _, dept := range depts {
		dept := dept
		go func() {
			s := erp.GetSubjectList(client, ssoToken, dept)
			chann <- s
		}()
	}

	for _, dept := range depts {
		log.Println("Fetched for", dept, "department", len(<-chann))
	}
}
