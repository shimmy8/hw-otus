package storage

import (
	"time"
)

type Event struct {
	ID           string
	Title        string
	StartDT      time.Time
	EndDT        time.Time
	Description  string
	UserID       string
	NotifyBefore time.Duration
}
