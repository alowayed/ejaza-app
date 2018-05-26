package main

import (
	"html/template"
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

// ---------- Methods

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

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
}
