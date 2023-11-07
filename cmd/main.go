package main

import erp "iitkgp_erp_login"

func main() {
	client := erp.Login(true)
	departments := [3]string{"CS", "GG", "CE"}

	for _, department := range departments {
		// go func() {
		erp.GetTimetable(client, department)
		// }()
	}
}
