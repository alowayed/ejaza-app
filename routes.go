package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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

	fmt.Fprint(w, Certs)

	// io.WriteString(w)
}

func GetCertSubmitted(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

}

func GetCertById(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

}

func GetCertConfirmById(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

}
