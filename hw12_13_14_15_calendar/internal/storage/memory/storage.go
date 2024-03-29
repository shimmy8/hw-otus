package memorystorage

import (
	"sync"
	"time"

	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events sync.Map
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) CreateEvent(e *storage.Event) error {
	dateBusy := s.checkDateBusy(e.StartDT, e.EndDT, e.UserID)
	if dateBusy {
		return storage.ErrDateBusy
	}

	_, loaded := s.events.LoadOrStore(e.ID, e)
	if loaded {
		return storage.ErrEventAlreadyCreated
	}
	return nil
}

func (s *Storage) GetEventsForInterval(startDt time.Time, endDt time.Time, userID string) ([]*storage.Event, error) {
	events := make([]*storage.Event, 0)

	s.events.Range(func(key, value any) bool {
		event := value.(*storage.Event)

		if event.UserID != userID {
			return true
		}

		// starts within (startDt - endDt)
		if (event.StartDT.After(startDt) || event.StartDT.Equal(startDt)) && event.StartDT.Before(endDt) {
			events = append(events, event)
			return true
		}

		// ends within (startDt - endDt)
		if event.EndDT.After(startDt) && (event.EndDT.Before(endDt) || event.EndDT.Equal(endDt)) {
			events = append(events, event)
			return true
		}

		// lasts within and outside of (startDt - endDt)
		if event.StartDT.Before(startDt) && event.EndDT.After(endDt) {
			events = append(events, event)
			return true
		}

		return true
	})

	return events, nil
}

func (s *Storage) checkDateBusy(startDT time.Time, endDT time.Time, userID string) bool {
	dateEvents, _ := s.GetEventsForInterval(startDT, endDT, userID)
	return len(dateEvents) > 0
}

func (s *Storage) GetEvent(id string) (*storage.Event, error) {
	value, ok := s.events.Load(id)
	if !ok {
		return nil, storage.ErrEventNotFoud
	}
	event := value.(*storage.Event)

	return event, nil
}

func (s *Storage) UpdateEvent(id string, e *storage.Event) error {
	_, loaded := s.events.Load(id)
	if !loaded {
		return storage.ErrEventNotFoud
	}

	s.events.Store(id, e)

	return nil
}

func (s *Storage) DeleteEvent(id string) error {
	s.events.Delete(id)
	return nil
}

func (s *Storage) GetNotifyEvents(startDt time.Time) ([]*storage.Event, error) {
	events := make([]*storage.Event, 0)

	s.events.Range(func(key, value any) bool {
		event := value.(*storage.Event)

		if event.Notified {
			return true
		}

		if event.StartDT.Add(-1 * event.NotifyBefore).Before(startDt) {
			events = append(events, event)
		}
		return true
	})

	return events, nil
}

func (s *Storage) DeleteOldEvents(olderThanDt time.Time) error {
	s.events.Range(func(key, value any) bool {
		event := value.(*storage.Event)

		if event.EndDT.Before(olderThanDt) {
			s.events.Delete(event.ID)
		}
		return true
	})
	return nil
}
