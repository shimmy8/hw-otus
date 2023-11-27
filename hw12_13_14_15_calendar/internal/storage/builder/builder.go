package storagebuilder

import (
	"context"
	"time"

	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
)

type Storage interface {
	GetEvent(id string) (*storage.Event, error)
	GetEventsForInterval(startDt time.Time, endDt time.Time, userID string) ([]*storage.Event, error)
	CreateEvent(e *storage.Event) error
	UpdateEvent(id string, e *storage.Event) error
	DeleteEvent(id string) error
	GetNotifyEvents(startDt time.Time) ([]*storage.Event, error)
	DeleteOldEvents(olderThanDt time.Time) error
}

func New(ctx context.Context, db string, url string, timeout int) Storage {
	var storage Storage

	switch db {
	case "in-memory":
		storage = memorystorage.New()
	case "db":
		storage = sqlstorage.New(ctx, url, timeout)
	default:
		storage = memorystorage.New()
	}
	return storage
}
