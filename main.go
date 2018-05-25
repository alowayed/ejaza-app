package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
)

// ---------- Types

type ejazaServer struct {
}

// ---------- Global Variables

var tpl *template.Template

// ---------- Methods

func (ser ejazaServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}

	data := struct {
		Method      string
		Submissions url.Values
	}{
		req.Method,
		req.Form,
	}
	tpl.ExecuteTemplate(w, "index", data)
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	var ser ejazaServer
	http.ListenAndServe(":8080", ser)
}
