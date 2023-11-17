package memorystorage

import (
	"strconv"
	"testing"
	"time"

	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func createNewEvent(id string) *storage.Event {
	return &storage.Event{
		ID:           id,
		Title:        "",
		StartDT:      time.Now().Add(-24 * time.Hour),
		EndDT:        time.Now().Add(-23 * time.Hour),
		Description:  "",
		UserID:       "123",
		NotifyBefore: time.Second,
	}
}

func TestStorage(t *testing.T) {
	t.Run("test add event", func(t *testing.T) {
		t.Parallel()
		s := New()

		e := createNewEvent("123")
		err := s.CreateEvent(e)

		require.NoError(t, err)

		storedEvent, err := s.GetEvent("123")
		require.NoError(t, err)

		require.Equal(t, e, storedEvent)
	})

	t.Run("test update event", func(t *testing.T) {
		t.Parallel()
		s := New()

		e := createNewEvent("123")
		_ = s.CreateEvent(e)

		e1 := createNewEvent("123")
		e1.Title = "Changed title"

		err := s.UpdateEvent("123", e1)
		require.NoError(t, err)

		storedEvent, err := s.GetEvent("123")
		require.NoError(t, err)

		require.Equal(t, storedEvent.Title, e1.Title)
	})

	t.Run("test delete event", func(t *testing.T) {
		t.Parallel()
		s := New()

		e := createNewEvent("123")
		_ = s.CreateEvent(e)

		err := s.DeleteEvent("123")
		require.NoError(t, err)

		_, getErr := s.GetEvent("123")
		require.ErrorIs(t, getErr, storage.ErrEventNotFoud)
	})

	t.Run("test already created err", func(t *testing.T) {
		t.Parallel()
		s := New()

		e := createNewEvent("123")
		_ = s.CreateEvent(e)

		e1 := createNewEvent("123")
		// override dates to avoid ErrDateBusy
		e1.StartDT = time.Now().Add(-2 * time.Hour)
		e1.EndDT = time.Now().Add(-1 * time.Hour)

		err := s.CreateEvent(e1)
		require.ErrorIs(t, err, storage.ErrEventAlreadyCreated)
	})

	t.Run("test date busy err", func(t *testing.T) {
		t.Parallel()
		s := New()

		e := createNewEvent("123")
		_ = s.CreateEvent(e)

		e1 := createNewEvent("12345")
		err := s.CreateEvent(e1)
		require.ErrorIs(t, err, storage.ErrDateBusy)
	})

	t.Run("test get events in interval", func(t *testing.T) {
		t.Parallel()
		s := New()

		now := time.Now()

		e1 := createNewEvent("e1")
		e1.StartDT = now.Add(-time.Hour)
		e1.EndDT = now.Add(-30 * time.Minute)
		s.CreateEvent(e1)

		e2 := createNewEvent("e2")
		e2.StartDT = now.Add(-20 * time.Minute)
		e2.EndDT = now.Add(10 * time.Minute)
		s.CreateEvent(e2)

		e3 := createNewEvent("e3")
		e3.StartDT = now.Add(20 * time.Minute)
		e3.EndDT = now.Add(time.Hour)
		s.CreateEvent(e3)

		tests := []struct {
			startDt        time.Time
			endDt          time.Time
			expectedEvents []*storage.Event
		}{
			{startDt: now.Add(-2 * time.Hour), endDt: now.Add(-65 * time.Minute), expectedEvents: []*storage.Event{}},
			{startDt: now.Add(-40 * time.Minute), endDt: now.Add(-10 * time.Minute), expectedEvents: []*storage.Event{e1, e2}},
			{startDt: now.Add(-10 * time.Minute), endDt: now, expectedEvents: []*storage.Event{e2}},
			{startDt: now.Add(5 * time.Minute), endDt: now.Add(30 * time.Minute), expectedEvents: []*storage.Event{e2, e3}},
			{startDt: now.Add(40 * time.Minute), endDt: now.Add(80 * time.Minute), expectedEvents: []*storage.Event{e3}},
			{startDt: now.Add(2 * time.Hour), endDt: now.Add(3 * time.Hour), expectedEvents: []*storage.Event{}},
		}

		for i, ts := range tests {
			ts := ts
			t.Run("interval #"+strconv.Itoa(i), func(t *testing.T) {
				t.Parallel()
				events, err := s.GetEventsForInterval(ts.startDt, ts.endDt, "123")
				require.NoError(t, err)
				require.Equal(t, len(ts.expectedEvents), len(events))
				require.Subset(t, ts.expectedEvents, events)
			})
		}
	})
}
