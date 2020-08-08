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

type ResizeEvent struct {
	*TrackingEvent
	OldWidth  string `json:"oldWidth"`
	OldHeight string `json:"oldHeight"`
	NewWidth  string `json:"newWidth"`
	NewHeight string `json:"newHeight"`
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

func handleResizeEvent(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("New resize event")
	fmt.Println(string(reqBody))

	var event ResizeEvent
	json.Unmarshal(reqBody, &event)

	record := db[event.SessionId]
	record.ResizeFrom = createDim(event.OldHeight, event.OldWidth)
	record.ResizeTo = createDim(event.NewHeight, event.NewWidth)

	db[event.SessionId] = record

}

type CopyPasteEvent struct {
	*TrackingEvent
	Pasted    bool   `json:"pasted"`
	FormId    string `json:"inputCardNumber"`
}

func handleCopyPasteEvent(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("New copy-paste event")
	fmt.Println(string(reqBody))

	var event CopyPasteEvent
	json.Unmarshal(reqBody, &event)

	record := getRecord(event.SessionId, event.WebsiteUrl)
	record.CopyAndPaste[event.FormId] = event.Pasted

	db[event.SessionId] = record
}


type TimerEvent struct {
	*TrackingEvent
	Time int `json:"time"`
}

func handleTimerEvent(w http.ResponseWriter, r *http.Request){
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("New copy-paste event")
	fmt.Println(string(reqBody))

	var event TimerEvent
	json.Unmarshal(reqBody, &event)

	record := getRecord(event.SessionId, event.WebsiteUrl)
	record.FormCompletionTime = event.Time

	db[event.SessionId] = record
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
