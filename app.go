package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var port string

const blobStorageOutputBindingName = "storage"

func init() {
	port = os.Getenv("APP_PORT")
	if port == "" {
		log.Fatalf("missing environment variable %s", "APP_PORT")
	}
}

func main() {
	http.HandleFunc("/eventhubs-input", func(rw http.ResponseWriter, req *http.Request) {
		var _time TheTime
		err := json.NewDecoder(req.Body).Decode(&_time)
		if err != nil {
			fmt.Println("error reading message from event hub binding", err)
			rw.WriteHeader(500)
			return
		}
		fmt.Printf("time from Event Hubs '%s'\n", _time.Time)
		rw.WriteHeader(200)
		err = json.NewEncoder(rw).Encode(Response{To: []string{blobStorageOutputBindingName}, Data: _time})
		if err != nil {
			fmt.Printf("unable to respond'%s'\n", err)
		}
	})
	http.ListenAndServe(":"+port, nil)
}

type Response struct {
	To   []string    `json:"to"`
	Data interface{} `json:"data"`
}

type TheTime struct {
	Time string `json:"time"`
}
