package memorystorage

import (
	"errors"
	"sync"

	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]*storage.Event
}

var (
	ErrAlreadyCreated = errors.New("event already created")
	ErrNotFoud        = errors.New("event not found")
)

func New() *Storage {
	return &Storage{
		events: make(map[string]*storage.Event),
	}
}

func (s *Storage) GetEvent(ID string) (*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, inStorage := s.events[ID]

	if !inStorage {
		return nil, ErrNotFoud
	}
	return e, nil
}

func (s *Storage) CreateEvent(e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.events[e.ID] != nil {
		return ErrAlreadyCreated
	}
	s.events[e.ID] = e
	return nil
}

func (s *Storage) UpdateEvent(e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return nil
}

func (s *Storage) DeleteEvent(ID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, inStorage := s.events[ID]
	if !inStorage {
		return ErrNotFoud
	}
	delete(s.events, ID)

	return nil
}
