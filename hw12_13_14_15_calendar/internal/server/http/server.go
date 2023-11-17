package internalhttp

import (
	"context"
	"net/http"
	"strconv"
)

type Server struct {
	logger Logger
	app    Application
	server *http.Server
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, id, title string) error
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{logger: logger, app: app}
}

func (s *Server) Start(ctx context.Context, host string, port int) error {
	mux := http.NewServeMux()

	helloHandler := http.HandlerFunc(s.Hello)
	mux.Handle("/hello", loggingMiddleware(helloHandler, s.logger))

	addr := host + ":" + strconv.Itoa(port)

	s.logger.Info("Starting HTTP server at " + addr)
	err := http.ListenAndServe(addr, mux)
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}

func (s *Server) Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello-world"))
}
