package app

import (
	"context"
	"time"

	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	GetEvent(ID string) (*storage.Event, error)
	GetEventsForInterval(startDt time.Time, endDt time.Time, userID string) ([]*storage.Event, error)
	CreateEvent(e *storage.Event) error
	UpdateEvent(ID string, e *storage.Event) error
	DeleteEvent(ID string) error
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}
