package main

import (
	"log"
	"math/rand"
	"net/http"
	"practice_1/handlers"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	http.HandleFunc("/", handlers.MainPageHandle)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
