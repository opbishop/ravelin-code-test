package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// data structure to hold Data structs as they are completed
// maps SessionId (string): Data struct object
var queue chan TrackingEvent

func main() {
	db := make(map[string]Data)
	queue = make(chan TrackingEvent, 10)
	go processEvents(db, queue)
	handleRequests()
}

// set up API routes for POST requests
func handleRequests() {
	http.HandleFunc("/resize", handleResizeEvent)
	http.HandleFunc("/copypaste", handleCopyPasteEvent)
	http.HandleFunc("/timer", handleTimerEvent)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func addCorsHeaders(w http.ResponseWriter) {
	header := w.Header()
	// TODO don't push to production with the allow-origin wildcard (obviously)
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "POST, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type")

}

func addBadRequestHeader(w http.ResponseWriter, err error){
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
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
		addBadRequestHeader(w, err)
	}

	te := TrackingEvent(event)
	select {
	case queue <- te:
	default:
	}
}

func handleCopyPasteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		addCorsHeaders(w)
		w.WriteHeader(http.StatusOK)
		return
	}
	addCorsHeaders(w)

	reqBody, _ := ioutil.ReadAll(r.Body)

	var event CopyPasteEvent
	err := json.Unmarshal(reqBody, &event)
	if err != nil {
		addBadRequestHeader(w, err)
	}

	te := TrackingEvent(event)
	select {
	case queue <- te:
	default:
	}
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
		addBadRequestHeader(w, err)
	}

	te := TrackingEvent(event)

	select {
	case queue <- te:
	default:
	}

}

func getRecord(db map[string]Data, sessionId string, websiteUrl string) *Data {
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

// goroutine to process events from the channel
// allows handling of multiple POST requests at once
func processEvents(db map[string]Data, queue chan TrackingEvent) {
	select {
	case trackingEvent := <- queue:
		switch trackingEvent.EventType() {
		case "resizeEvent":
			resizeEvent := trackingEvent.(ResizeEvent)
			record := getRecord(db, resizeEvent.SessionId, resizeEvent.WebsiteUrl)
			record.ResizeFrom = createDim(resizeEvent.OldHeight, resizeEvent.OldWidth)
			record.ResizeTo = createDim(resizeEvent.NewHeight, resizeEvent.NewWidth)

			db[resizeEvent.SessionId] = *record
			fmt.Println(record)
		case "pasteEvent":
			pasteEvent := trackingEvent.(CopyPasteEvent)
			record := getRecord(db, pasteEvent.SessionId, pasteEvent.WebsiteUrl)
			record.CopyAndPaste[pasteEvent.FormId] = pasteEvent.Pasted

			db[pasteEvent.SessionId] = *record
			fmt.Println(record)
		case "timerEvent":
			timerEvent := trackingEvent.(TimerEvent)
			record := getRecord(db, timerEvent.SessionId, timerEvent.WebsiteUrl)
			record.FormCompletionTime = timerEvent.Time

			db[timerEvent.SessionId] = *record
			fmt.Println(record)
		}
	default:
	}

}
