package sqlstorage

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db *sqlx.DB
}

func New(ctx context.Context, dbURL string) *Storage {
	storage := &Storage{}
	storage.Connect(ctx, dbURL)
	return storage
}

func (s *Storage) Connect(ctx context.Context, dbURL string) error {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalln(err)
	}
	s.db = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
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
	 (id, title, start_dt, end_dt, description, user_id, notify_before) VALUES
	 (:id, :title, :start_dt, :end_dt, :description, :user_id, :notify_before)"`, e)
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

func (s *Storage) GetEvent(ID string) (*storage.Event, error) {
	event := &storage.Event{}
	err := s.db.Get(event, "SELECT * FROM events WHERE id=$1", ID)
	return event, err
}

func (s *Storage) UpdateEvent(ID string, e *storage.Event) error {
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
	) WHERE id=:id`, e)
	err := tx.Commit()
	return err

}

func (s *Storage) DeleteEvent(ID string) error {
	tx := s.db.MustBegin()
	tx.MustExec(`DELETE FROM events WHERE id=$1`, ID)
	err := tx.Commit()
	return err
}
