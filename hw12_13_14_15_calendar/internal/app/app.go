package app

import (
	"context"
	"time"

	"github.com/google/uuid"
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

func (a *App) CreateEvent(
	_ context.Context,
	userID string,
	title string,
	startDt time.Time,
	endDt time.Time,
	description string,
	notifyBefore time.Duration,
) (*storage.Event, error) {
	event := &storage.Event{
		ID:           uuid.New().String(),
		Title:        title,
		StartDT:      startDt,
		EndDT:        endDt,
		Description:  description,
		UserID:       userID,
		NotifyBefore: notifyBefore,
	}

	if err := a.storage.CreateEvent(event); err != nil {
		return nil, err
	}

	return event, nil
}

func (a *App) UpdateEvent(
	_ context.Context,
	eventID string,
	userID string,
	title string,
	startDt time.Time,
	endDt time.Time,
	description string,
	notifyBefore time.Duration,
) (*storage.Event, error) {
	event := &storage.Event{
		ID:           eventID,
		Title:        title,
		StartDT:      startDt,
		EndDT:        endDt,
		Description:  description,
		UserID:       userID,
		NotifyBefore: notifyBefore,
	}

	if err := a.storage.UpdateEvent(eventID, event); err != nil {
		return nil, err
	}

	return event, nil
}

func (a *App) DeleteEvent(_ context.Context, eventID string) error {
	return a.storage.DeleteEvent(eventID)
}

func (a *App) ListEventsForDay(_ context.Context, userID string, startDt time.Time) ([]*storage.Event, error) {
	dayStartDt := time.Date(startDt.Year(), startDt.Month(), startDt.Day(), 0, 0, 0, 0, time.UTC)
	dayEndDt := dayStartDt.Add(time.Hour * 24).Add(time.Microsecond * -1)

	return a.storage.GetEventsForInterval(dayStartDt, dayEndDt, userID)
}

func (a *App) ListEventsForWeek(_ context.Context, userID string, startDt time.Time) ([]*storage.Event, error) {
	weekStartDt := time.Date(startDt.Year(), startDt.Month(), startDt.Day(), 0, 0, 0, 0, time.UTC)
	weekEndDt := weekStartDt.Add(time.Hour * 24 * 7).Add(time.Microsecond * -1)

	return a.storage.GetEventsForInterval(weekStartDt, weekEndDt, userID)
}

func (a *App) ListEventsForMonth(_ context.Context, userID string, startDt time.Time) ([]*storage.Event, error) {
	monthStartDt := time.Date(startDt.Year(), startDt.Month(), 1, 0, 0, 0, 0, time.UTC)
	monthEndDt := monthStartDt.AddDate(0, 1, 0)

	return a.storage.GetEventsForInterval(monthStartDt, monthEndDt, userID)
}
