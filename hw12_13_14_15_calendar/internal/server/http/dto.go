package internalhttp

import "time"

type CreateEventRequest struct {
	Title        string `json:"title"`
	StartDT      string `json:"start_dt"` //nolint:tagliatelle
	EndDT        string `json:"end_dt"`   //nolint:tagliatelle
	Description  string `json:"description"`
	UserID       string `json:"user_id"`       //nolint:tagliatelle
	NotifyBefore string `json:"notify_before"` //nolint:tagliatelle
}

type CreateEventReply struct {
	EventID string `json:"event_id"` //nolint:tagliatelle
}

type UpdateEventRequest struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	StartDT      string `json:"start_dt"` //nolint:tagliatelle
	EndDT        string `json:"end_dt"`   //nolint:tagliatelle
	Description  string `json:"description"`
	UserID       string `json:"user_id"`       //nolint:tagliatelle
	NotifyBefore string `json:"notify_before"` //nolint:tagliatelle
}

type UpdateEventReply struct {
	Ok bool `json:"ok"`
}

type DeleteEventRequest struct {
	ID string `json:"event_id"` //nolint:tagliatelle
}

type DeleteEventReply struct {
	Ok bool `json:"ok"`
}

type ListEventsRequest struct {
	StartDT time.Time `json:"start_dt"` //nolint:tagliatelle
	UserID  string    `json:"user_id"`  //nolint:tagliatelle
}

type ListEventsReply struct {
	Events []*ReplyEvent `json:"events"`
}

type ReplyEvent struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	StartDT      time.Time     `json:"start_dt"` //nolint:tagliatelle
	EndDT        time.Time     `json:"end_dt"`   //nolint:tagliatelle
	Description  string        `json:"description"`
	UserID       string        `json:"user_id"`       //nolint:tagliatelle
	NotifyBefore time.Duration `json:"notify_before"` //nolint:tagliatelle
	Notified     bool          `json:"notified"`      //nolint:tagliatelle
}
