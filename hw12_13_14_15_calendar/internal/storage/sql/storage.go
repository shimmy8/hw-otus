package sqlstorage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db      *pgxpool.Pool
	timeout time.Duration

	ctx    context.Context
	cancel context.CancelFunc
}

func New(ctx context.Context, dbURL string, timeout int) *Storage {
	storage := &Storage{}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(timeout))
	storage.timeout = time.Second * time.Duration(timeout)
	storage.ctx = ctx
	storage.cancel = cancel

	storage.Connect(timeoutCtx, dbURL)
	return storage
}

func (s *Storage) Connect(ctx context.Context, dbURL string) error {
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Storage) Close() error {
	s.cancel()
	s.db.Close()
	return nil
}

func (s *Storage) CreateEvent(e *storage.Event) error {
	busy, checkErr := s.checkDateBusy(e.StartDT, e.EndDT, e.UserID, "")
	if checkErr != nil {
		return checkErr
	}
	if busy {
		return storage.ErrDateBusy
	}

	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	_, err := s.db.Exec(
		ctx,
		`INSERT INTO events
			(id, title, start_dt, end_dt, description, user_id, notify_before, notified)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)`,
		e.ID,
		e.Title,
		e.StartDT,
		e.EndDT,
		e.Description,
		e.UserID,
		e.NotifyBefore,
		e.Notified,
	)
	return err
}

func (s *Storage) GetEventsForInterval(startDt time.Time, endDt time.Time, userID string) ([]*storage.Event, error) {
	events := make([]*storage.Event, 0)
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	rows, queryErr := s.db.Query(
		ctx,
		`SELECT * FROM events
		WHERE
			user_id = $3 AND (
				(start_dt >= $1 AND start_dt <= $2)
				OR
				(end_dt >= $1 AND end_dt <= $2)
				OR
				(end_dt <= $1 AND end_dt >= $2)
			)
		`,
		startDt, endDt, userID,
	)
	if queryErr != nil {
		return events, queryErr
	}

	evs, err := pgx.CollectRows(rows, pgx.RowToStructByName[storage.Event])
	for ind := range evs {
		events = append(events, &evs[ind])
	}

	return events, err
}

func (s *Storage) checkDateBusy(startDT time.Time, endDT time.Time, userID string, skipID string) (bool, error) {
	events, err := s.GetEventsForInterval(startDT, endDT, userID)
	if err != nil {
		return false, err
	}

	for _, evt := range events {
		if evt.ID != skipID {
			return true, nil
		}
	}

	return false, err
}

func (s *Storage) GetEvent(id string) (*storage.Event, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	rows, err := s.db.Query(ctx, "SELECT * FROM events WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	event, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[storage.Event])
	return event, err
}

func (s *Storage) UpdateEvent(_ string, e *storage.Event) error {
	busy, checkErr := s.checkDateBusy(e.StartDT, e.EndDT, e.UserID, e.ID)
	if checkErr != nil {
		return checkErr
	}
	if busy {
		return storage.ErrDateBusy
	}

	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	_, err := s.db.Exec(ctx,
		`UPDATE events
		SET
			title = $1,
			start_dt = $2,
			end_dt = $3,
			description = $4,
			user_id = $5,
			notify_before = $6,
			notified = $7
		WHERE id = $8`,
		e.Title,
		e.StartDT,
		e.EndDT,
		e.Description,
		e.UserID,
		e.NotifyBefore,
		e.Notified,
		e.ID,
	)
	return err
}

func (s *Storage) DeleteEvent(id string) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	_, err := s.db.Exec(ctx, "DELETE FROM events WHERE id=$1", id)
	return err
}

func (s *Storage) GetNotifyEvents(startDt time.Time) ([]*storage.Event, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	rows, queryErr := s.db.Query(
		ctx,
		`SELECT * FROM events
		WHERE
			notified = false AND
			notify_before IS NOT NULL AND
			start_dt - notify_before <= $1`,
		startDt,
	)
	if queryErr != nil {
		return nil, queryErr
	}

	events := make([]*storage.Event, 0)
	evs, err := pgx.CollectRows(rows, pgx.RowToStructByName[storage.Event])
	for ind := range evs {
		events = append(events, &evs[ind])
	}

	return events, err
}

func (s *Storage) DeleteOldEvents(olderThanDt time.Time) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	_, err := s.db.Exec(ctx, "DELETE FROM events WHERE end_dt < $1", olderThanDt)
	return err
}
