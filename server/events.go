package main

type TrackingEvent interface {
	EventType() string
}

type TimerEvent struct {
	WebsiteUrl string `json:"websiteUrl"`
	SessionId  string `json:"sessionId"`
	Time       int    `json:"time"`
}

func (t TimerEvent) EventType() string {
	return "timerEvent"
}


type CopyPasteEvent struct {
	WebsiteUrl string `json:"websiteUrl"`
	SessionId  string `json:"sessionId"`
	Pasted     bool   `json:"pasted"`
	FormId     string `json:"formId"`
}

func (t CopyPasteEvent) EventType() string {
	return "pasteEvent"
}

type ResizeEvent struct {
	WebsiteUrl string `json:"websiteUrl"`
	SessionId  string `json:"sessionId"`
	OldWidth   string `json:"oldWidth"`
	OldHeight  string `json:"oldHeight"`
	NewWidth   string `json:"newWidth"`
	NewHeight  string `json:"newHeight"`
}

func (t ResizeEvent) EventType() string {
	return "resizeEvent"
}