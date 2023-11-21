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
	internalhttp "github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
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

	var storage app.Storage

	switch config.Storage.DB {
	case "in-memory":
		storage = memorystorage.New()
	case "db":
		storage = sqlstorage.New(ctx, config.Storage.URL, config.Storage.Timeout)
	default:
		storage = memorystorage.New()
	}

	logg := logger.New(config.Logger.Level, "server")
	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(logg, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx, config.HTTPServer.Host, config.HTTPServer.Port); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
