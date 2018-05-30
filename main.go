package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ---------- Types

type Cert struct {
	Id        int
	Data      string
	Semail    string
	Iemail    string
	Nonce     string
	Confirmed bool
}

// ---------- Global Variables

var tpl *template.Template
var Certs []Cert

// ---------- Routes

func GetIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	tpl.ExecuteTemplate(w, "index", nil)
}

func GetCert(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	tpl.ExecuteTemplate(w, "index", nil)
}

func GetCertSubmit(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	tpl.ExecuteTemplate(w, "certSubmit", nil)
}

func PostCertSubmit(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	err := req.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}

	id := len(Certs)
	data := req.Form["data"][0]
	sEmail := req.Form["semail"][0]
	iEmail := req.Form["iemail"][0]

	cert := Cert{
		Id:        id,
		Data:      data,
		Semail:    sEmail,
		Iemail:    iEmail,
		Nonce:     "001",
		Confirmed: false,
	}

	Certs = append(Certs, cert)

	// TODO: Err handle if cert cannor be added to DB

	// Redirect to this certificate's page (should be pending confirmation)
	url := fmt.Sprintf("/cert/id/%v", id)
	http.Redirect(w, req, url, http.StatusSeeOther)
}

func GetCertSubmitted(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

}

func GetCertById(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	// TODO: check if id is in database

	io.WriteString(w, fmt.Sprintf("%v", Certs))
	// fmt.Fprint(w, Certs)
}

func GetCertConfirmById(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

}

// ---------- Methods

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

// func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	return events.APIGatewayProxyResponse{
// 		StatusCode: 200,
// 		Body:       "Hello AWS Lambda and Netlify",
// 	}, nil
// }

func main() {

	mux := httprouter.New()
	mux.GET("/", GetIndex)
	mux.GET("/cert", GetCert)
	mux.GET("/cert/submit", GetCertSubmit)
	mux.POST("/cert/submit", PostCertSubmit)
	mux.GET("/cert/submitted", GetCertSubmitted)
	mux.GET("/cert/id/:id", GetCertById)
	mux.GET("/cert/confirm/:id", GetCertConfirmById)

	http.ListenAndServe(":8080", mux)

	// lambda.Start(handler)
}
