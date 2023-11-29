package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func makeHTTPRequest(url string, method string, rqBody []byte) (status int, body []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(rqBody))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("%v", err)
	}

	resData, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	return res.StatusCode, resData
}

func parseEventID(resp []byte) (string, error) {
	respStruct := struct {
		EventID string `json:"event_id"` //nolint:tagliatelle
	}{}

	if err := json.Unmarshal(resp, &respStruct); err != nil {
		return "", err
	}

	return respStruct.EventID, nil
}

type Event struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	StartDT      time.Time     `json:"start_dt"` //nolint:tagliatelle
	EndDT        time.Time     `json:"end_dt"`   //nolint:tagliatelle
	Description  string        `json:"description"`
	UserID       string        `json:"user_id"`       //nolint:tagliatelle
	NotifyBefore time.Duration `json:"notify_before"` //nolint:tagliatelle
	Notified     bool          `json:"notified"`      //nolint:tagliatelle
}

func parseEvent(resp []byte) (*Event, error) {
	event := &Event{}

	if err := json.Unmarshal(resp, event); err != nil {
		return nil, err
	}
	return event, nil
}

func TestAddEventHttp(t *testing.T) {
	t.Run("Test add date busy delete event", func(t *testing.T) {
		body := []byte(`{
			"title": "Test event",
			"user_id": "123",
			"start_dt": "2023-11-21 15:00:00",
			"end_dt": "2023-11-21 17:30:00",
			"description": "Some event desc",
			"notify_before": "10h"
		}`)

		status, resp := makeHTTPRequest("http://localhost:8080/events/create", http.MethodPost, body)
		require.Equal(t, http.StatusCreated, status)

		eventID, err := parseEventID(resp)
		require.NoError(t, err)

		status2, resp2 := makeHTTPRequest("http://localhost:8080/events/create", http.MethodPost, body)
		require.Equal(t, http.StatusInternalServerError, status2)
		require.Contains(t, string(resp2), "date busy")

		delBody := []byte(fmt.Sprintf(`{"event_id": "%s"}`, eventID))
		statusDel, _ := makeHTTPRequest("http://localhost:8080/events/delete", http.MethodPost, delBody)
		require.Equal(t, http.StatusOK, statusDel)
	})

	t.Run("Test bad requests", func(t *testing.T) {
		brokenBodies := [][]byte{
			[]byte(`{
				"title": "Test broken event",
				"user_id": "123",
				"start_dt": "2024-11-21---15:00:00",
				"end_dt": "2024-11-21 17:30:00",
				"description": "Some event desc",
				"notify_before": "1h"
			}`),
			[]byte(`{
				"title": "Test broken  event",
				"user_id": "123",
				"start_dt": "2024-11-21 15:00:00",
				"end_dt": "2024-11-21 dddd",
				"description": "Some event desc",
				"notify_before": "1h"
			}`),
			[]byte(`{
				"title": "Test broken  event",
				"user_id": "123",
				"start_dt": "2024-11-21 15:00:00",
				"end_dt": "2024-11-21 dddd",
				"description": "Some event desc",
				"notify_before": "1h555ww"
			}`),
			[]byte(`{
				"title": "Test broken  event",
				"user_id": "123...`),
		}

		for _, bb := range brokenBodies {
			body := bb

			status, _ := makeHTTPRequest("http://localhost:8080/events/create", http.MethodPost, body)
			require.Equal(t, http.StatusBadRequest, status)
		}
	})
}

func TestListEvents(t *testing.T) {
	t.Run("test list events", func(t *testing.T) {
		eventDates := []time.Time{
			time.Date(2023, 11, 1, 10, 0, 0, 0, time.UTC),
			time.Date(2023, 11, 3, 10, 0, 0, 0, time.UTC),
			time.Date(2023, 11, 10, 10, 0, 0, 0, time.UTC),
		}
		eventIDs := make([]string, len(eventDates))
		// cleanup
		defer func() {
			for _, eID := range eventIDs {
				delBody := []byte(fmt.Sprintf(`{"event_id": "%s"}`, eID))
				makeHTTPRequest("http://localhost:8080/events/delete", http.MethodPost, delBody)
			}
		}()

		// create test events
		for ind, startDt := range eventDates {
			body := []byte(fmt.Sprintf(`{
				"title": "event#%d",
				"user_id": "test_user",
				"start_dt": "%s",
				"end_dt": "%s",
				"description": "Some event desc",
				"notify_before": "1h"
			}`,
				ind,
				startDt.Format(time.DateTime),
				startDt.Add(time.Hour*1).Format(time.DateTime),
			))
			status, resp := makeHTTPRequest("http://localhost:8080/events/create", http.MethodPost, body)
			require.Equal(t, http.StatusCreated, status)

			eventID, err := parseEventID(resp)
			require.NoError(t, err)

			eventIDs[ind] = eventID
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
				url := "http://localhost:8080/events/list/" + tt.period + "?user_id=test_user&start_dt=" + tt.startDt

				status, res := makeHTTPRequest(url, http.MethodGet, nil)
				require.Equal(t, http.StatusOK, status)

				stringresp := string(res)
				for _, expextedEvTitle := range tt.expectedEvents {
					require.Contains(t, stringresp, expextedEvTitle)
				}
			})
		}
	})
}

func TestEventNotified(t *testing.T) {
	t.Run("check sender receives event", func(t *testing.T) {
		now := time.Now().UTC()

		body := []byte(fmt.Sprintf(`{
			"title": "event_to_notify",
			"user_id": "test_user_123",
			"start_dt": "%s",
			"end_dt": "%s",
			"description": "Some event desc",
			"notify_before": "10m"
		}`,
			now.Add(time.Minute*10).Format(time.DateTime),
			now.Add(time.Minute*30).Format(time.DateTime),
		))

		status, resp := makeHTTPRequest("http://localhost:8080/events/create", http.MethodPost, body)
		require.Equal(t, http.StatusCreated, status)

		eventID, err := parseEventID(resp)
		require.NoError(t, err)

		// celanup
		defer func() {
			delBody := []byte(fmt.Sprintf(`{"event_id": "%s"}`, eventID))
			makeHTTPRequest("http://localhost:8080/events/delete", http.MethodPost, delBody)
		}()

		// wait for event to be picked up by scheduler
		time.Sleep(time.Second * 10)

		// get created event
		statusGet, resp := makeHTTPRequest(
			fmt.Sprintf("http://localhost:8080/events?id=%s", eventID),
			http.MethodGet, nil,
		)
		require.Equal(t, http.StatusOK, statusGet)
		// parse response
		event, err := parseEvent(resp)
		require.NoError(t, err)
		// check if sender changed Notified to true
		require.True(t, event.Notified)
	})
}
