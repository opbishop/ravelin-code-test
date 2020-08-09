package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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

func handleResizeEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		addCorsHeaders(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	addCorsHeaders(w)

	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("New resize event")
	fmt.Println(string(reqBody))

	var event ResizeEvent
	err := json.Unmarshal(reqBody, &event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	te := TrackingEvent(event)
	c <- te
}

func handleCopyPasteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		addCorsHeaders(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	addCorsHeaders(w)

	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("New copy-paste event")
	//fmt.Println(string(reqBody))

	var event CopyPasteEvent
	err := json.Unmarshal(reqBody, &event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	te := TrackingEvent(event)
	c <- te
}

func handleTimerEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		addCorsHeaders(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	addCorsHeaders(w)

	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("New timer event")
	fmt.Println(string(reqBody))

	var event TimerEvent
	err := json.Unmarshal(reqBody, &event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	te := TrackingEvent(event)
	c <- te
}

func addCorsHeaders(w http.ResponseWriter) {
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

// goroutine to process events from the channel
// allows handling of multiple POST requests at once
func processEvents() {
	for range c {
		trackingEvent := <-c

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

// data structure to hold Data structs as they are completed
// maps SessionId (string): Data struct object
var db map[string]Data
var c chan TrackingEvent

func main() {
	db = make(map[string]Data)
	c = make(chan TrackingEvent, 10)
	go processEvents()
	handleRequests()
}
