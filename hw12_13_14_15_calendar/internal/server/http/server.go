package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	logger Logger
	app    *app.App

	server *http.Server
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func NewServer(logger Logger, app *app.App) *Server {
	return &Server{logger: logger, app: app}
}

func (s *Server) Start(ctx context.Context, host string, port int) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", s.Hello)
	mux.HandleFunc("/events", s.GetEvent)
	mux.HandleFunc("/events/create", s.CreateEvent)
	mux.HandleFunc("/events/update", s.UpdateEvent)
	mux.HandleFunc("/events/delete", s.DeleteEvent)
	mux.HandleFunc("/events/list/day", s.ListEventsForDay)
	mux.HandleFunc("/events/list/week", s.ListEventsForWeek)
	mux.HandleFunc("/events/list/month", s.ListEventsForMonth)

	addr := host + ":" + strconv.Itoa(port)
	server := &http.Server{
		Addr:              addr,
		Handler:           loggingMiddleware(mux, s.logger),
		ReadHeaderTimeout: 3 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	s.logger.Info("Starting HTTP server at " + addr)
	err := server.ListenAndServe()
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}

func (s *Server) Hello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("hello-world"))
}

func (s *Server) GetEvent(w http.ResponseWriter, r *http.Request) {
	eventID := r.URL.Query().Get("id")

	event, err := s.app.GetEvent(r.Context(), eventID)
	if err != nil {
		s.logger.Info(fmt.Sprintf("get event error %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := &ReplyEvent{
		ID:           event.ID,
		Title:        event.Title,
		StartDT:      event.StartDT,
		EndDT:        event.EndDT,
		Description:  event.Description,
		UserID:       event.UserID,
		NotifyBefore: event.NotifyBefore,
		Notified:     event.Notified,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
	rqEvent := CreateEventRequest{}

	err := parseRequestBody(r, &rqEvent)
	if err != nil {
		s.logger.Info(fmt.Sprintf("parse request body error %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDt, err := time.Parse(time.DateTime, rqEvent.StartDT)
	if err != nil {
		s.logger.Info(fmt.Sprintf("parse request startDT error %v, body: %v", err, rqEvent))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	endDt, err := time.Parse(time.DateTime, rqEvent.EndDT)
	if err != nil {
		s.logger.Info(fmt.Sprintf("parse request endDT error %v, body: %v", err, rqEvent))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	notifyBefore, err := time.ParseDuration(rqEvent.NotifyBefore)
	if err != nil {
		s.logger.Info(fmt.Sprintf("parse request notifyBefore error %v, body: %v", err, rqEvent))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := s.app.CreateEvent(
		r.Context(),
		rqEvent.UserID,
		rqEvent.Title,
		startDt,
		endDt,
		rqEvent.Description,
		notifyBefore,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.logger.Info(fmt.Sprintf("create event error %v, body: %v", err, rqEvent))
		return
	}

	w.WriteHeader(http.StatusCreated)
	rep := CreateEventReply{EventID: event.ID}
	json.NewEncoder(w).Encode(rep)
}

func (s *Server) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	rqEvent := UpdateEventRequest{}

	err := parseRequestBody(r, &rqEvent)
	if err != nil {
		s.logger.Info(fmt.Sprintf("parse request body error %v", err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	startDt, err := time.Parse(time.DateTime, rqEvent.StartDT)
	if err != nil {
		s.logger.Info(fmt.Sprintf("parse request startDT error %v, body: %v", err, rqEvent))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	endDt, err := time.Parse(time.DateTime, rqEvent.EndDT)
	if err != nil {
		s.logger.Info(fmt.Sprintf("parse request endDT error %v, body: %v", err, rqEvent))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	notifyBefore, err := time.ParseDuration(rqEvent.NotifyBefore)
	if err != nil {
		s.logger.Info(fmt.Sprintf("parse request notifyBefore error %v, body: %v", err, rqEvent))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, updErr := s.app.UpdateEvent(
		r.Context(),
		rqEvent.ID,
		rqEvent.UserID,
		rqEvent.Title,
		startDt,
		endDt,
		rqEvent.Description,
		notifyBefore,
	)
	if updErr != nil {
		s.logger.Info(fmt.Sprintf("update event error %v, body: %v", err, rqEvent))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rep := UpdateEventReply{Ok: true}
	json.NewEncoder(w).Encode(rep)
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	rqEvent := DeleteEventRequest{}

	err := parseRequestBody(r, &rqEvent)
	if err != nil {
		s.logger.Info(fmt.Sprintf("parse delete body error %v, body: %v", err, rqEvent))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	delErr := s.app.DeleteEvent(r.Context(), rqEvent.ID)
	if delErr != nil {
		s.logger.Info(fmt.Sprintf("delete event error %v", delErr))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rep := DeleteEventReply{Ok: true}
	json.NewEncoder(w).Encode(rep)
	w.WriteHeader(http.StatusOK)
}

const (
	day   = "day"
	week  = "week"
	month = "month"
)

func (s *Server) ListEventsForDay(w http.ResponseWriter, r *http.Request) {
	s.listEvents(w, r, day)
}

func (s *Server) ListEventsForWeek(w http.ResponseWriter, r *http.Request) {
	s.listEvents(w, r, week)
}

func (s *Server) ListEventsForMonth(w http.ResponseWriter, r *http.Request) {
	s.listEvents(w, r, month)
}

func (s *Server) listEvents(w http.ResponseWriter, r *http.Request, period string) {
	startDt, err := time.Parse(time.DateOnly, r.URL.Query().Get("start_dt"))
	if err != nil {
		http.Error(w, "invalid request params", http.StatusBadRequest)
		return
	}
	userID := r.URL.Query().Get("user_id")

	var events []*storage.Event
	var listErr error

	switch period {
	case day:
		events, listErr = s.app.ListEventsForDay(r.Context(), userID, startDt)
	case week:
		events, listErr = s.app.ListEventsForWeek(r.Context(), userID, startDt)
	case month:
		events, listErr = s.app.ListEventsForMonth(r.Context(), userID, startDt)
	default:
		http.Error(w, "unknown period", http.StatusInternalServerError)
		return
	}

	if listErr != nil {
		http.Error(w, listErr.Error(), http.StatusInternalServerError)
		return
	}

	resEvents := make([]*ReplyEvent, len(events))
	for ind, evt := range events {
		resEvents[ind] = &ReplyEvent{
			ID:           evt.ID,
			Title:        evt.Title,
			StartDT:      evt.StartDT,
			EndDT:        evt.EndDT,
			Description:  evt.Description,
			UserID:       evt.UserID,
			NotifyBefore: evt.NotifyBefore,
			Notified:     evt.Notified,
		}
	}

	rep := ListEventsReply{Events: resEvents}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rep)
}

func parseRequestBody(r *http.Request, v any) error {
	res, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(res, &v)
	if err != nil {
		return err
	}

	return nil
}
