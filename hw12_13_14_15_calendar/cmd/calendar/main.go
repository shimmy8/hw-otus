package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/server/http"
	storagebuilder "github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage/builder"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := config.NewConfig(configFile)

	storage := storagebuilder.New(ctx, config.Storage.DB, config.Storage.URL, config.Storage.Timeout)
	logg := logger.New(config.Logger.Level, "server")
	calendar := app.New(logg, storage)
	httpServer := internalhttp.NewServer(logg, calendar)
	grpcServer := internalgrpc.NewServer(logg, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpServer.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}

		if err := grpcServer.Stop(ctx); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	go func() {
		if err := grpcServer.Start(ctx, config.GRPCServer.Host, config.GRPCServer.Port); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
		}
	}()

	if err := httpServer.Start(ctx, config.HTTPServer.Host, config.HTTPServer.Port); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
