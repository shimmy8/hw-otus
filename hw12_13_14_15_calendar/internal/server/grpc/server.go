package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/server/grpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Server struct {
	logger Logger
	app    *app.App

	server *grpc.Server
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func NewServer(logger Logger, app *app.App) *Server {
	return &Server{logger: logger, app: app}
}

func (s *Server) Start(_ context.Context, host string, port int) error {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(s.loggingInterceptor))
	s.server = grpcServer

	gen.RegisterEventServiceServer(grpcServer, &eventService{app: s.app})

	addr := fmt.Sprintf("%s:%d", host, port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Error("Failed to start grpc server at " + addr)
		return err
	}

	s.logger.Info("Starting GRPC server at " + addr)

	if err := grpcServer.Serve(lis); err != nil {
		s.logger.Error("Failed to start grpc server at " + addr)
		return err
	}

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	if s.server == nil {
		return nil
	}
	s.server.GracefulStop()
	return nil
}

func (s *Server) loggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	startTime := time.Now()
	resp, err := handler(ctx, req)
	if err != nil {
		s.logger.Error(
			fmt.Sprintf("method %q failed: %s", info.FullMethod, err),
		)
	}
	peerAddr := ""
	if peer, ok := peer.FromContext(ctx); ok {
		peerAddr = peer.Addr.String()
	}
	var userAgent []string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userAgent = md.Get("user-agent")
	}

	var msgBuilder strings.Builder
	for _, part := range []string{
		peerAddr,
		info.FullMethod,
		strings.Join(userAgent, ""),
		strconv.Itoa(int(time.Since(startTime).Microseconds())),
	} {
		msgBuilder.WriteString(part)
		msgBuilder.WriteString(" ")
	}

	s.logger.Info(msgBuilder.String())

	return resp, err
}
