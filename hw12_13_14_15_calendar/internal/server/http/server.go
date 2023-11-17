package internalhttp

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"
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

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
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
