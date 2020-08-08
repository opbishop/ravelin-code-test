package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Data struct {
	WebsiteUrl         string
	SessionId          string
	ResizeFrom         Dimension
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool // map[fieldId]true
	FormCompletionTime int             // Seconds
}

type Dimension struct {
	Width  string
	Height string
}

type TrackingEvent struct {
	WebsiteUrl string `json:"websiteUrl"`
	SessionId  string `json:"sessionId"`
	EventType string `json:"eventType"`
}

func getRecord(sessionId string, websiteUrl string) Data {
	if record, exist := db[sessionId]; exist {
		return record
	} else {
		return Data{WebsiteUrl: websiteUrl, SessionId: sessionId, CopyAndPaste: make(map[string]bool)}
	}
}

func createDim(height string, width string) Dimension {
	return Dimension{
		Width:  width,
		Height: height,
	}
}

type ResizeEvent struct {
	TrackingEvent
	OldWidth  string `json:"oldWidth"`
	OldHeight string `json:"oldHeight"`
	NewWidth  string `json:"newWidth"`
	NewHeight string `json:"newHeight"`
}

func handleResizeEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		handleCORS(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	handleCORS(w)

	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("New resize event")
	fmt.Println(string(reqBody))

	var event ResizeEvent
	err := json.Unmarshal(reqBody, &event)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	record := getRecord(event.SessionId, event.WebsiteUrl)
	record.ResizeFrom = createDim(event.OldHeight, event.OldWidth)
	record.ResizeTo = createDim(event.NewHeight, event.NewWidth)

	db[event.SessionId] = record
	fmt.Println(record)
}

type CopyPasteEvent struct {
	TrackingEvent
	Pasted    bool   `json:"pasted"`
	FormId    string `json:"formId"`
}

func handleCopyPasteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		handleCORS(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	handleCORS(w)

	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("New copy-paste event")
	//fmt.Println(string(reqBody))

	var event CopyPasteEvent
	err := json.Unmarshal(reqBody, &event)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	record := getRecord(event.SessionId, event.WebsiteUrl)
	record.CopyAndPaste[event.FormId] = event.Pasted

	db[event.SessionId] = record
	fmt.Println(record)
}

type TimerEvent struct {
	TrackingEvent
	Time int `json:"time"`
}

func handleTimerEvent(w http.ResponseWriter, r *http.Request){
	if r.Method == "OPTIONS" {
		handleCORS(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	handleCORS(w)

	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("New timer event")
	fmt.Println(string(reqBody))

	var event TimerEvent
	err := json.Unmarshal(reqBody, &event)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	record := getRecord(event.SessionId, event.WebsiteUrl)
	record.FormCompletionTime = event.Time

	db[event.SessionId] = record
	fmt.Println(record)
}

func handleCORS(w http.ResponseWriter){
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type")

}

func handleRequests() {
	http.HandleFunc("/resize", handleResizeEvent)
	http.HandleFunc("/copypaste", handleCopyPasteEvent)
	http.HandleFunc("/timer", handleTimerEvent)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

// data structure to hold Data structs as they are completed
// maps SessionId (string): Data struct object
var db map[string]Data

func main() {
	db = make(map[string]Data)
	handleRequests()
}
