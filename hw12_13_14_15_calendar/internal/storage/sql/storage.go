package sqlstorage

import (
	"context"
	"time"

	// use pgx as driver.
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db     *sqlx.DB
	ctx    context.Context
	cancel context.CancelFunc
}

func New(ctx context.Context, dbURL string, timeout int) *Storage {
	storage := &Storage{}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(timeout))
	storage.ctx = timeoutCtx
	storage.cancel = cancel

	storage.Connect(dbURL)
	return storage
}

func (s *Storage) Connect(dbURL string) error {
	db := sqlx.MustConnect("pgx", dbURL)

	s.db = db
	return nil
}

func (s *Storage) Close() error {
	s.cancel()
	err := s.db.Close()
	return err
}

func (s *Storage) CreateEvent(e *storage.Event) error {
	busy := s.checkDateBusy(e.StartDT, e.EndDT, e.UserID)
	if busy {
		return storage.ErrDateBusy
	}

	tx := s.db.MustBegin()
	tx.NamedExec(`INSERT INTO events
	 (id, title, start_dt, end_dt, description, user_id, notify_before, notified) VALUES
	 (:id, :title, :start_dt, :end_dt, :description, :user_id, :notify_before, :notified)"`, e)
	err := tx.Commit()
	return err
}

func (s *Storage) GetEventsForInterval(startDt time.Time, endDt time.Time, userID string) ([]*storage.Event, error) {
	events := make([]*storage.Event, 0)
	rows, err := s.db.Queryx(`
	SELECT * FROM events
	WHERE
		user_id = '$3' AND (
			(start_dt >= $1 AND start_dt <= $2)
			OR
			(end_dt >= $1 AND end_dt <= $2)
			OR
			(end_dt <= $1 AND end_dt >= $2)
		)
	`, startDt, endDt, userID)
	if err != nil {
		return events, err
	}

	for rows.Next() {
		event := storage.Event{}
		err := rows.StructScan(&event)
		if err != nil {
			return events, err
		}
		events = append(events, &event)
	}
	return events, nil
}

func (s *Storage) checkDateBusy(startDT time.Time, endDT time.Time, userID string) bool {
	events, _ := s.GetEventsForInterval(startDT, endDT, userID)
	return len(events) > 0
}

func (s *Storage) GetEvent(id string) (*storage.Event, error) {
	event := &storage.Event{}
	err := s.db.Get(event, "SELECT * FROM events WHERE id=$1", id)
	return event, err
}

func (s *Storage) UpdateEvent(_ string, e *storage.Event) error {
	busy := s.checkDateBusy(e.StartDT, e.EndDT, e.UserID)
	if busy {
		return storage.ErrDateBusy
	}
	tx := s.db.MustBegin()
	tx.NamedExec(`UPDATE events
	SET (
		title=:title,
		start_dt=:start_dt,
		end_dt=:end_dt,
		description=:description,
		user_id=:user_id,
		notify_before=:notify_before
		notified=:notified
	) WHERE id=:id`, e)
	err := tx.Commit()
	return err
}

func (s *Storage) DeleteEvent(id string) error {
	tx := s.db.MustBegin()
	tx.MustExec(`DELETE FROM events WHERE id=$1`, id)
	err := tx.Commit()
	return err
}

func (s *Storage) GetNotifyEvents(startDt time.Time) ([]*storage.Event, error) {
	events := make([]*storage.Event, 0)
	rows, err := s.db.Queryx(`
	SELECT * FROM events
	WHERE
		notified = false AND
		notify_before != "" AND
		start_dt - notify_before <= $1
	`, startDt)
	if err != nil {
		return events, err
	}

	for rows.Next() {
		event := storage.Event{}
		err := rows.StructScan(&event)
		if err != nil {
			return events, err
		}
		events = append(events, &event)
	}
	return events, nil
}

func (s *Storage) DeleteOldEvents(olderThanDt time.Time) error {
	tx := s.db.MustBegin()
	tx.MustExec(`DELETE FROM events WHERE end_dt < $1`, olderThanDt)
	err := tx.Commit()
	return err
}
