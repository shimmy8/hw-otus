package internalhttp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestCUDServer(t *testing.T) {
	logg := logger.New("DEBUG", "testhttp")
	appStorage := memorystorage.New()
	app := app.New(logg, appStorage)
	server := NewServer(logg, app)

	t.Run("test event create", func(t *testing.T) {
		body := []byte(`{
			"title": "Test event",
			"user_id": "123",
			"start_dt": "2023-11-21 15:00:00",
			"end_dt": "2023-11-21 17:30:00",
			"description": "Some event desc",
			"notify_before": "1h"
		}`)

		req := httptest.NewRequest(http.MethodPost, "/events/create", io.NopCloser(bytes.NewReader(body)))
		rec := httptest.NewRecorder()

		server.CreateEvent(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		require.Equal(t, http.StatusCreated, res.StatusCode)

		createdResp := CreateEventReply{}
		err := json.Unmarshal(rec.Body.Bytes(), &createdResp)
		require.NoError(t, err)

		eventID := createdResp.EventID

		storedEvent, err := appStorage.GetEvent(eventID)
		require.NoError(t, err)

		require.Equal(t, "123", storedEvent.UserID)
		require.Equal(t, "Test event", storedEvent.Title)
		require.Equal(t, "2023-11-21 15:00:00", storedEvent.StartDT.Format(time.DateTime))
		require.Equal(t, "2023-11-21 17:30:00", storedEvent.EndDT.Format(time.DateTime))
		require.Equal(t, "1h0m0s", storedEvent.NotifyBefore.String())
	})

	t.Run("test event update", func(t *testing.T) {
		storeErr := appStorage.CreateEvent(&storage.Event{
			ID:           "123",
			Title:        "Test title",
			StartDT:      time.Now(),
			EndDT:        time.Now(),
			Description:  "Some desc",
			UserID:       "tes_user",
			NotifyBefore: time.Duration(time.Duration.Hours(1)),
		})
		require.NoError(t, storeErr)

		body := []byte(`{
			"id": "123",
			"title": "Test event",
			"user_id": "updated user id",
			"start_dt": "2023-11-21 15:00:00",
			"end_dt": "2023-11-21 17:30:00",
			"description": "Some event desc",
			"notify_before": "1h"
		}`)

		req := httptest.NewRequest(http.MethodPost, "/events/update", io.NopCloser(bytes.NewReader(body)))
		rec := httptest.NewRecorder()

		server.UpdateEvent(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		require.Equal(t, http.StatusOK, res.StatusCode)

		updatedResp := UpdateEventReply{}
		err := json.Unmarshal(rec.Body.Bytes(), &updatedResp)
		require.NoError(t, err)

		ok := updatedResp.Ok
		require.Equal(t, true, ok)

		storedEvent, err := appStorage.GetEvent("123")
		require.NoError(t, err)

		require.Equal(t, "updated user id", storedEvent.UserID)
		require.Equal(t, "Test event", storedEvent.Title)
		require.Equal(t, "2023-11-21 15:00:00", storedEvent.StartDT.Format(time.DateTime))
		require.Equal(t, "2023-11-21 17:30:00", storedEvent.EndDT.Format(time.DateTime))
		require.Equal(t, "1h0m0s", storedEvent.NotifyBefore.String())
	})

	t.Run("test event delete", func(t *testing.T) {
		storeErr := appStorage.CreateEvent(&storage.Event{
			ID:           "delete-me",
			Title:        "Test title",
			StartDT:      time.Now(),
			EndDT:        time.Now(),
			Description:  "Some desc",
			UserID:       "tes_user",
			NotifyBefore: time.Duration(time.Duration.Hours(1)),
		})
		require.NoError(t, storeErr)

		body := []byte(`{"event_id": "delete-me"}`)

		req := httptest.NewRequest(http.MethodPost, "/events/delete", io.NopCloser(bytes.NewReader(body)))
		rec := httptest.NewRecorder()

		server.DeleteEvent(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		require.Equal(t, http.StatusOK, res.StatusCode)

		deletedResp := DeleteEventReply{}
		err := json.Unmarshal(rec.Body.Bytes(), &deletedResp)
		require.NoError(t, err)

		ok := deletedResp.Ok
		require.Equal(t, true, ok)

		_, getErr := appStorage.GetEvent("delete-me")
		require.ErrorIs(t, getErr, storage.ErrEventNotFoud)
	})
}

func TestServerEventsList(t *testing.T) {
	logg := logger.New("DEBUG", "testhttp")
	appStorage := memorystorage.New()
	app := app.New(logg, appStorage)
	server := NewServer(logg, app)
	t.Run("test events list", func(t *testing.T) {
		eventDates := []time.Time{
			time.Date(2023, 11, 1, 10, 0, 0, 0, time.UTC),
			time.Date(2023, 11, 3, 10, 0, 0, 0, time.UTC),
			time.Date(2023, 11, 10, 10, 0, 0, 0, time.UTC),
		}

		for ind, startDt := range eventDates {
			appStorage.CreateEvent(&storage.Event{
				ID:           "event#" + strconv.Itoa(ind),
				Title:        "Test title",
				StartDT:      startDt,
				EndDT:        startDt.Add(time.Hour * 1),
				Description:  "Some desc",
				UserID:       "test_user",
				NotifyBefore: time.Duration(time.Duration.Hours(1)),
			})
		}

		tests := []struct {
			startDt        string
			period         string
			expectedEvents []string
		}{
			{startDt: "2023-11-01", period: "day", expectedEvents: []string{"event#0"}},
			{startDt: "2023-11-01", period: "week", expectedEvents: []string{"event#0", "event#1"}},
			{startDt: "2023-11-01", period: "month", expectedEvents: []string{"event#0", "event#1", "event#2"}},
		}

		for _, tt := range tests {
			tt := tt
			t.Run("request period "+tt.period, func(t *testing.T) {
				url := "/events/list/" + tt.period + "?user_id=test_user&start_dt=" + tt.startDt

				req := httptest.NewRequest(http.MethodGet, url, nil)
				rec := httptest.NewRecorder()

				switch tt.period {
				case "day":
					server.ListEventsForDay(rec, req)
				case "week":
					server.ListEventsForWeek(rec, req)
				case "month":
					server.ListEventsForMonth(rec, req)
				}

				res := rec.Result()
				defer res.Body.Close()

				require.Equal(t, http.StatusOK, res.StatusCode)

				listResp := ListEventsReply{}
				err := json.Unmarshal(rec.Body.Bytes(), &listResp)
				require.NoError(t, err)

				respEventsID := []string{}
				for _, evt := range listResp.Events {
					respEventsID = append(respEventsID, evt.ID)
				}

				require.Equal(t, len(tt.expectedEvents), len(listResp.Events))
				require.Equal(t, tt.expectedEvents, respEventsID)
			})
		}
	})
}
