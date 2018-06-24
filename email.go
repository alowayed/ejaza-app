package main

import (
	"fmt"
	"log"
	"net/smtp"
	"strconv"

	"github.com/matcornic/hermes"
)

const (
	MAIL_HOST      = "smtp.gmail.com"
	EJAZA_EMAIL    = "ejazaco@gmail.com"
	EJAZA_PASSWORD = "Vacation!"
)

var her hermes.Hermes

func init() {
	// Configure hermes by setting a theme and your product info
	her = hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "Ejaza App",
			Link: "https://ejaza.herokuapp.com/",
			// Optional product logo
			//Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
		},
	}
}

func sendConfirmationEmail(cert Cert) {

	// TODO: Add check that all fields are correct in cert

	//cert := struct {
	//	Id     int
	//	Iemail string
	//	Semail string
	//	Data   string
	//	Nonce  string
	//}{
	//	Id:     1,
	//	Iemail: "yousefzoq@gmail.com",
	//	Semail: "student@uni.co",
	//	Data:   "some data",
	//	Nonce:  "123",
	//}

	confirmationEmail := hermes.Email{
		Body: hermes.Body{
			Name: "Registrars",
			Intros: []string{
				"One of your former students would like to have their transcript on the Blockchain." +
				cert.Semail + "     claims that this data is accurate:",
				cert.Data,
			},
			Actions: []hermes.Action{
				{
					Instructions: "If what " + cert.Semail + " claims is true, hit the confirm button " +
						"upload it to the Blockchain and save it forever!",
					Button: hermes.Button{
						Color: "#22BC66", // Optional action button color
						Text:  "Confirm Transcript",
						Link:  "http://localhost:8080/cert/confirm/" + strconv.Itoa(cert.Id) + "/nonce/" + cert.Nonce,
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	emailBody, err := her.GenerateHTML(confirmationEmail)
	//fmt.Println(emailBody)
	if err != nil {
		panic(err) // Tip: Handle error with something else than a panic ;)
	}

	recipient := cert.Iemail

	auth := smtp.PlainAuth("", EJAZA_EMAIL, EJAZA_PASSWORD, MAIL_HOST)
	fmt.Println("auth: ", auth)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "A Student wants to Confirm a Transcript" + "!\n"
	msg := []byte(subject + mime + "\n" + emailBody)

	err = smtp.SendMail(MAIL_HOST+":587", auth, EJAZA_EMAIL, []string{recipient}, msg)
	fmt.Println("Sent mail")
	if err != nil {
		log.Fatal(err)
	}
}

//func main() {
//	sendConfirmationEmail()
//}
