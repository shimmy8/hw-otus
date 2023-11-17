package storage

import (
	"errors"
	"time"
)

var (
	ErrEventAlreadyCreated = errors.New("event already created")
	ErrEventNotFoud        = errors.New("event not found")
	ErrDateBusy            = errors.New("date busy")
)

type Event struct {
	ID           string        `db:"id"`
	Title        string        `db:"title"`
	StartDT      time.Time     `db:"start_dt"`
	EndDT        time.Time     `db:"end_dt"`
	Description  string        `db:"description"`
	UserID       string        `db:"user_id"`
	NotifyBefore time.Duration `db:"notify_before"`
}
