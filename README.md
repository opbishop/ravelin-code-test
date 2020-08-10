Ravelin Code Test
=================

## Summary
This repo contains a basic HTTP server written using the Go standard library which handles POST requests from 
an HTML/vanilla Javascript UI.

The server accepts three types of request:
* POST to /resize endpoint
* POST to /copypaste endpoint
* POST to /timer endpoint

As each request is received, a Data scruct stored in the db map is updated and printed to stdout.

![Code flow](flowchart.jpg)


### Example JSON Requests
```javascript
{
  "eventType": "pasteEvent",
  "websiteUrl": "https://ravelin.com",
  "sessionId": "123123-123123-123123123",
  "pasted": true,
  "formId": "inputCardNumber"
}

{
    "eventType": "resizeEvent",
    "websiteUrl": "https://ravelin.com",
    "sessionId": "123123-123123-123123123",
    "oldWidth": "100",
    "oldHeight": "100",
    "newWidth": "200",
    "newHeight": "200"
}

{
  "eventType": "timerEvent",
  "websiteUrl": "https://ravelin.com",
  "sessionId": "123123-123123-123123123",
  "time": 72
}
```

## Frontend (JS)
To run:
Open index.html in browser (tested using Chrome version 84.0.4147.105)

Uses listeners to wait for:
1. Page resize
2. Paste into any of the text fields
3. Typing begin in any text field & button click

When any of these events are triggered, the Javascript POSts the above data to the appropriate endpoint. 

## Backend (Go)
To run:


The Backend should:

1. Create a Server
2. Accept POST requests in JSON format similar to those specified above
3. Map the JSON requests to relevant sections of the data struct (specified below)
4. Print the struct for each stage of its construction
5. Also print the struct when it is complete (i.e. when the form submit button has been clicked)

We would like the server to be written to handle multiple requests arriving on
the same session at the same time. We'd also like to see some tests.


### Go Struct
```go
type Data struct {
	WebsiteUrl         string
	SessionId          string
	ResizeFrom         Dimension
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool // map[fieldId]true
	FormCompletionTime int // Seconds
}

type Dimension struct {
	Width  string
	Height string
}
```




