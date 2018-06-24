package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"encoding/json"
	"math/rand"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
)

// -------------------- Global Variables
var tpl *template.Template
var client *redis.Client

// -------------------- Types

type Cert struct {
	Id           int    `json:"id,string"`
	Data         string `json:"data"`
	Semail       string `json:"semail"`
	Iemail       string `json:"iemail"`
	Nonce        string `json:"nonce"`
	Confirmed    bool   `json:"confirmed,string"`
	BlockchainId string `json:"blockchainId"`
}

type TerionResp struct {
	BlockchainId string `json:"id"`
	AccountId    int    `json:"accountId"`
	DatastoreId  int    `json:"datastoreId"`
	Status       string `json:"status"`
	Data         Cert
	Json         string `json:"json"`
	Sha256       string `json:"sha256"`
	Timestamp    int    `json:"timestamp"`
}

// -------------------- Routes

func GetIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	preHandler(&w, req, ps)
	tpl.ExecuteTemplate(w, "index", nil)
}

func GetCert(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	preHandler(&w, req, ps)
	tpl.ExecuteTemplate(w, "cert", nil)
}

func GetCertSubmit(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	tpl.ExecuteTemplate(w, "certSubmit", nil)
}

func PostCertSubmit(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	err := req.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}

	// Create the certificate and pick a random id and nonce for it
	id := generateId()
	nonce := strconv.Itoa(rand.Int())

	data := req.Form["data"][0]
	sEmail := req.Form["semail"][0]
	iEmail := req.Form["iemail"][0]

	cert := Cert{
		Id:     id,
		Data:   data,
		Semail: sEmail,
		Iemail: iEmail,
		Nonce:  nonce,
	}

	// Send an email to
	sendConfirmationEmail(cert)

	// Save the cert to the database
	strId := strconv.Itoa(cert.Id)

	certJson, err := json.Marshal(cert)
	if err != nil {
		panic(err)
	}

	err = client.Set(strId, certJson, 0).Err()
	if err != nil {
		url := fmt.Sprintf("/")
		http.Redirect(w, req, url, http.StatusSeeOther)
		return
	}
	// TODO: Err handle if cert cannot be added to DB

	// Redirect to this certificate's page (should be pending confirmation)
	url := fmt.Sprintf("/cert/id/%v", id)
	http.Redirect(w, req, url, http.StatusSeeOther)
}

func GetCertById(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	preHandler(&w, req, ps)

	id := ps.ByName("id")
	certId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("certId from param: ", certId)
	}

	itemJson := client.Get(id).Val()
	if itemJson == "" {
		fmt.Println("Cert does not exist")
		tpl.ExecuteTemplate(w, "cert", id)
		return
	}

	item := &Cert{}
	unmarshalErr := json.Unmarshal([]byte(itemJson), item)
	if unmarshalErr != nil {
		fmt.Println(unmarshalErr)
	}

	//fmt.Println("item", item.String())
	//if item.Err() != redis.Nil {
	//	fmt.Println("cannot find cert with id: ", certId)
	//	// TODO
	//}

	tpl.ExecuteTemplate(w, "certById", item)
}

func GetCertConfirmById(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	preHandler(&w, req, ps)

	id := ps.ByName("id")
	nonce := ps.ByName("nonce")
	certId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("certId from param: ", certId)
	}

	itemJson := client.Get(id).Val()
	if itemJson == "" {
		fmt.Println("Cert does not exist")
		tpl.ExecuteTemplate(w, "cert", id)
		return
	}

	cert := &Cert{}
	unmarshalErr := json.Unmarshal([]byte(itemJson), cert)
	if unmarshalErr != nil {
		fmt.Println(unmarshalErr)
	}

	// Check the nonce
	if (nonce != cert.Nonce) {
		fmt.Println("Nonce does not match what's in the database")
		tpl.ExecuteTemplate(w, "cert", id)
		return
	}

	cert.Confirmed = true
	cert.BlockchainId = saveCertToBlockchain(*cert)

	// Save the cert to the database
	strId := strconv.Itoa(cert.Id)

	certJson, err := json.Marshal(cert)
	if err != nil {
		panic(err)
	}

	if err = client.Del(strId).Err(); err == nil {
		fmt.Println("Failed to delete cert with id: ", strId)
	}
	err = client.Set(strId, certJson, 0).Err()
	if err != nil {
		url := fmt.Sprintf("/cert/id/%v", id)
		http.Redirect(w, req, url, http.StatusSeeOther)
	}
	// TODO: Err handle if cert cannot be added to DB

	// Redirect to this certificate's page (should be pending confirmation)
	url := fmt.Sprintf("/cert/id/%v", id)
	http.Redirect(w, req, url, http.StatusSeeOther)
}

// -------------------- Methods

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	setupDB()
}

func setupDB() {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = "localhost:6379"
	}

	client = redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func preHandler(w *http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	enableCors(w)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func generateId() int {

	id := rand.Int()

	// TODO: check if ID in database

	return id
}

// -------------------- Methods

func main() {

	mux := httprouter.New()

	mux.ServeFiles("/css/*filepath", http.Dir("css"))

	mux.GET("/", GetIndex)
	mux.GET("/cert", GetCert)
	mux.POST("/cert", GetCert)
	mux.GET("/cert/submit", GetCertSubmit)
	mux.POST("/cert/submit", PostCertSubmit)
	mux.GET("/cert/id/:id", GetCertById)
	mux.POST("/cert/id/:id", GetCertById)
	mux.GET("/cert/confirm/:id/nonce/:nonce", GetCertConfirmById)

	port := os.Getenv("PORT")
	port = ":" + port

	if port == ":" {
		port = ":8080"
	}

	fmt.Println("Listening on port", port)
	http.ListenAndServe(port, mux)
}
