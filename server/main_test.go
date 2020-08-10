package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddCorsHeaders(t *testing.T) {
	// create test request
	req, err := http.NewRequest("GET", "/timer", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleTimerEvent)

	handler.ServeHTTP(rr, req)

	expectedHeaders := map[string][]string{
		"Access-Control-Allow-Methods": {"POST, OPTIONS"},
		"Access-Control-Allow-Headers": {"Content-Type"},
		"Access-Control-Allow-Origin": {"*"},
	}

	for key, value := range expectedHeaders {
		if strings.Join(rr.Header()[key], "") != strings.Join(value, "") {
			t.Errorf("CORS headers incorrectly set: got %v but expected %v",
				strings.Join(value, " "), strings.Join(rr.Header()[key], " "))
		}
	}

}

var timerEvents = []struct {
	reason string
	// input json string
	input string
	// expected http response code
	expectedOutput int
}{
	{"Should return 200 if valid json", "{\n  \"eventType\": \"timeTaken\",\n  \"websiteUrl\": \"https://ravelin.com\",\n  \"sessionId\": \"123123-123123-123123123\",\n  \"time\": 72\n}\n", 200},
	{"Should return 400 if invalid input", "some invalid json", 400},
}

func TestHandleTimerEvent(t *testing.T) {
	for _, testData := range timerEvents {
		t.Run(testData.reason, func(t *testing.T) {
			// create test request
			req, err := http.NewRequest("POST", "/timer", strings.NewReader(testData.input))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handleTimerEvent)

			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != testData.expectedOutput {
				t.Errorf("Timer event handler returned wrong status code: got %v but expected %v",
					status, testData.expectedOutput)
			}
		})
	}
}

var pasteEvents = []struct {
	reason string
	// input json string
	input string
	// expected http response code
	expectedOutput int
}{
	{"Should return 200 if valid json", "{\n  \"eventType\": \"copyAndPaste\",\n  \"websiteUrl\": \"https://ravelin.com\",\n  \"sessionId\": \"123123-123123-123123123\",\n  \"pasted\": true,\n  \"formId\": \"inputCardNumber\"\n}", 200},
	{"Should return 400 if invalid input", "some invalid json", 400},
}

func TestHandleCopyPasteEvent(t *testing.T) {
	for _, testData := range pasteEvents {
		t.Run(testData.reason, func(t *testing.T) {
			// create test request
			req, err := http.NewRequest("POST", "/copypaste", strings.NewReader(testData.input))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handleCopyPasteEvent)

			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != testData.expectedOutput {
				t.Errorf("Timer event handler returned wrong status code: got %v but expected %v",
					status, testData.expectedOutput)
			}
		})
	}
}

var resizeEvents = []struct {
	reason string
	// input json string
	input string
	// expected http response code
	expectedOutput int
}{
	{"Should return 200 if valid json", "{\n    \"eventType\": \"resize\",\n    \"websiteUrl\": \"https://ravelin.com\",\n    \"sessionId\": \"123123-123123-123123123\",\n\t\"oldWidth\": \"100\",\n\t\"oldHeight\": \"100\",\n\t\"newWidth\": \"200\",\n\t\"newHeight\": \"200\"\n}", 200},
	{"Should return 400 if invalid input", "some invalid json", 400},
}

func TestHandleResizeEvent(t *testing.T) {
	for _, testData := range pasteEvents {
		t.Run(testData.reason, func(t *testing.T) {
			// create test request
			req, err := http.NewRequest("POST", "/resize", strings.NewReader(testData.input))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handleResizeEvent)

			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != testData.expectedOutput {
				t.Errorf("Timer event handler returned wrong status code: got %v but expected %v",
					status, testData.expectedOutput)
			}
		})
	}
}


func TestProcessEvents(t *testing.T){
	mockDb := make(map[string]Data)
	mockQueue := make(chan TrackingEvent, 1)



	event := TimerEvent{
		SessionId: "abc",
		Time: 500,
	}

	select {
	case mockQueue <- event:
		close(mockQueue)
		processEvents(mockDb, mockQueue)
	default:
		t.Error("Test failed; unable to send event on mockQueue")
	}

	if _, exist := mockDb[event.SessionId] ; ! exist {
		t.Errorf("Expected a record to be created with session ID %v",
			event.SessionId)
	}

	if mockDb[event.SessionId].FormCompletionTime != event.Time {
		t.Errorf("Event time incorrectly set; got %v but expected %v",
			mockDb[event.SessionId], event.Time)
	}

}
