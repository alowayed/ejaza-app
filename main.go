package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"strconv"
	"github.com/go-redis/redis"
	"math/rand"
	"encoding/json"
	"time"
	"io/ioutil"
	"bytes"
)

var client *redis.Client

func setupDB() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// ---------- Types

type Cert struct {
	Id        int
	Data      string
	Semail    string
	Iemail    string
	Nonce     string
	Confirmed bool
	BlockchainId string
}

type TerionResp struct {
	BlockchainId string `json:"id"`
	AccountId int `json:"accountId"`
	DatastoreId int `json:"datastoreId"`
	Status string `json:"status"`
	Data Cert
	Json string `json:"json"`
	Sha256 string `json:"sha256"`
	Timestamp int `json:"timestamp"`
}

// ---------- Global Variables

var tpl *template.Template

// Push the given resource to the client.
// func push(w http.ResponseWriter, resource string) {
// 	pusher, ok := w.(http.Pusher)
// 	fmt.Println("push not supported")
// 	if ok {
// 		if err := pusher.Push(resource, nil); err != nil {
// 			fmt.Println("Failed to push")
// 		}
// 	}
// }

// ---------- Routes

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

	id := generateId()

	data := req.Form["data"][0]
	sEmail := req.Form["semail"][0]
	iEmail := req.Form["iemail"][0]

	cert := Cert{
		Id:        id,
		Data:      data,
		Semail:    sEmail,
		Iemail:    iEmail,
		Nonce:     "001",
	}

	strId := strconv.Itoa(cert.Id)

	certJson, err := json.Marshal(cert)
	if err != nil {
		panic (err)
	}

	err = client.Set(strId, certJson, 0).Err()
	if err != nil {
		fmt.Println("getting cert by id: ", strId)
		url := fmt.Sprintf("/cert/id/%v", id)
		http.Redirect(w, req, url, http.StatusSeeOther)

	}
	// TODO: Err handle if cert cannot be added to DB

	saveCertToBlockchain(cert)

	// Redirect to this certificate's page (should be pending confirmation)
	url := fmt.Sprintf("/cert/id/%v", id)
	http.Redirect(w, req, url, http.StatusSeeOther)
}

func generateId() int {

	id := rand.Int()

	// TODO: check if ID in database

	return id
}

func GetCertSubmitted(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

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

}

// ---------- Methods

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))

	setupDB()
}

func preHandler(w *http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	enableCors(w)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func saveCertToBlockchain(cert Cert) {

	url := "https://api.tierion.com/v1/records"

	reqCert := struct {
		DataStoreId int `json:"datastoreId"`
		Id        int	`json:"id"`
		Data      string`json:"data"`
		Semail    string`json:"semail"`
		Iemail    string`json:"iemail"`
	}{
		DataStoreId: 7485,
		Id: cert.Id,
		Data: cert.Data,
		Semail: cert.Semail,
		Iemail: cert.Iemail,
	}

	jsonString, err := json.Marshal(reqCert)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonString))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("X-Username", "ejazaco@gmail.com")
	req.Header.Set("X-Api-Key", "hALcxG45opSVjr2e/qpjOn3260jiKVQEe3eUc/hxAF4=")
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	fmt.Println("resp status: ", res.Status)

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var terionResp TerionResp
	unmarshalErr := json.Unmarshal(body, &terionResp)
	if unmarshalErr != nil {
		fmt.Println("Failed to parse response from Terion: ", unmarshalErr)
	}

	fmt.Println("body: ", terionResp)
}


func main() {

	mux := httprouter.New()

	mux.ServeFiles("/css/*filepath", http.Dir("css"))

	mux.GET("/", GetIndex)
	mux.GET("/cert", GetCert)
	mux.POST("/cert", GetCert)
	mux.GET("/cert/submit", GetCertSubmit)
	mux.POST("/cert/submit", PostCertSubmit)
	mux.GET("/cert/submitted", GetCertSubmitted)
	mux.GET("/cert/id/:id", GetCertById)
	mux.POST("/cert/id/:id", GetCertById)
	mux.GET("/cert/confirm/:id", GetCertConfirmById)

	port := os.Getenv("PORT")
	port = ":" + port

	if port == ":" {
		port = ":8080"
	}

	fmt.Println("Listening on port", port)
	http.ListenAndServe(port, mux)
}
