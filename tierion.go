package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func saveCertToTierion(cert Cert) string {

	//cert := struct {
	//	Id           int    `json:"id,string"`
	//	Data         string `json:"data"`
	//	Semail       string `json:"semail"`
	//	Iemail       string `json:"iemail"`
	//	Nonce        string `json:"nonce"`
	//	Confirmed    bool   `json:"confirmed,string"`
	//	BlockchainId string `json:"blockchainId"`
	//}{
	//	Id:           5577006791947779410,
	//	Data:         "test",
	//	Semail:       "b@c.com",
	//	Iemail:       "yousefzoq@gmail.com",
	//	Nonce:        "8674665223082153551",
	//	Confirmed:    false,
	//	BlockchainId: "",
	//}

	url := "https://api.tierion.com/v1/records"

	reqCert := struct {
		DataStoreId int    `json:"datastoreId"`
		Id          int    `json:"id"`
		Data        string `json:"data"`
		Semail      string `json:"semail"`
		Iemail      string `json:"iemail"`
	}{
		DataStoreId: 7485,
		Id:          cert.Id,
		Data:        cert.Data,
		Semail:      cert.Semail,
		Iemail:      cert.Iemail,
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
		Timeout: time.Second * 5, // Maximum of 2 secs
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

	//var tierionResp TerionResp
	var tierionResp struct {
		BlockchainId string `json:"id"`
		AccountId    int    `json:"accountId"`
		DatastoreId  int    `json:"datastoreId"`
		Status       string `json:"status"`
		Data         Cert
		Json      string `json:"json"`
		Sha256    string `json:"sha256"`
		Timestamp int    `json:"timestamp"`
	}
	unmarshalErr := json.Unmarshal(body, &tierionResp)
	if unmarshalErr != nil {
		fmt.Println("Failed to parse response from Terion: ", unmarshalErr)
	}

	fmt.Println("tierionResp.Data.BlockchainId: ", tierionResp.BlockchainId)
	fmt.Println("tierionResp.Status", tierionResp.Status)

	return tierionResp.BlockchainId
}

//func main() {
//	saveCertToTierion()
//}
