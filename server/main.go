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

type TrackingEvent interface {
	EventType() string

}

func getRecord(sessionId string, websiteUrl string) *Data {
	if record, exist := db[sessionId]; exist {
		return &record
	} else {
		return &Data{WebsiteUrl: websiteUrl, SessionId: sessionId, CopyAndPaste: make(map[string]bool)}
	}
}

func createDim(height string, width string) Dimension {
	return Dimension{
		Width:  width,
		Height: height,
	}
}

type ResizeEvent struct {
	WebsiteUrl string `json:"websiteUrl"`
	SessionId  string `json:"sessionId"`
	OldWidth  string `json:"oldWidth"`
	OldHeight string `json:"oldHeight"`
	NewWidth  string `json:"newWidth"`
	NewHeight string `json:"newHeight"`
}

func (t ResizeEvent) EventType() string {
	return "resizeEvent"
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

	te := TrackingEvent(event)
	c <- te
}

type CopyPasteEvent struct {
	WebsiteUrl string `json:"websiteUrl"`
	SessionId  string `json:"sessionId"`
	Pasted    bool   `json:"pasted"`
	FormId    string `json:"formId"`
}

func (t CopyPasteEvent) EventType() string {
	return "pasteEvent"
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

	te := TrackingEvent(event)
	c <- te
}

type TimerEvent struct {
	WebsiteUrl string `json:"websiteUrl"`
	SessionId  string `json:"sessionId"`
	Time int `json:"time"`
}

func (t TimerEvent) EventType() string {
	return "timerEvent"
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

	te := TrackingEvent(event)
	c <- te
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
var c chan TrackingEvent

func processEvents(){
	for range c {
		trackingEvent := <- c

		switch trackingEvent.EventType() {
		case "resizeEvent":
			resizeEvent := trackingEvent.(ResizeEvent)
			record := getRecord(resizeEvent.SessionId, resizeEvent.WebsiteUrl)
			record.ResizeFrom = createDim(resizeEvent.OldHeight, resizeEvent.OldWidth)
			record.ResizeTo = createDim(resizeEvent.NewHeight, resizeEvent.NewWidth)

			db[resizeEvent.SessionId] = *record
			fmt.Println(record)
		case "pasteEvent":
			pasteEvent := trackingEvent.(CopyPasteEvent)
			record := getRecord(pasteEvent.SessionId, pasteEvent.WebsiteUrl)
			record.CopyAndPaste[pasteEvent.FormId] = pasteEvent.Pasted

			db[pasteEvent.SessionId] = *record
			fmt.Println(record)
		case "timerEvent":
			timerEvent := trackingEvent.(TimerEvent)
			record := getRecord(timerEvent.SessionId, timerEvent.WebsiteUrl)
			record.FormCompletionTime = timerEvent.Time

			db[timerEvent.SessionId] = *record
			fmt.Println(record)
		}

	}

}


func main() {
	db = make(map[string]Data)
	c = make(chan TrackingEvent, 10)
	go processEvents()
	handleRequests()
}
