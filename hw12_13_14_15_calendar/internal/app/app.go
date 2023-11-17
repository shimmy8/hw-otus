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
	GetEvent(id string) (*storage.Event, error)
	GetEventsForInterval(startDt time.Time, endDt time.Time, userID string) ([]*storage.Event, error)
	CreateEvent(e *storage.Event) error
	UpdateEvent(id string, e *storage.Event) error
	DeleteEvent(id string) error
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(_ context.Context, _, _ string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}
