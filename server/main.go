package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func handleNewData(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("New POST received")
	fmt.Println(string(reqBody))
}

func handleRequests() {
	http.HandleFunc("/", handleNewData)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}
