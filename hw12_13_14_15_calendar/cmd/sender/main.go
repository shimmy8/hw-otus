package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/logger"
	internalqueue "github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/queue"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/scheduler/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config := config.NewConfig(configFile)
	logg := logger.New(config.Logger.Level, "sender")

	rabbitmq := internalqueue.New(
		config.AMQPConf.Username,
		config.AMQPConf.Password,
		config.AMQPConf.Host,
		config.AMQPConf.Port,
		config.AMQPConf.QueueName,
	)

	if err := rabbitmq.Connect(); err != nil {
		logg.Error(fmt.Sprintf("failed to connect to rabbitmq %v", err))
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	defer rabbitmq.Close()

	logg.Info("sender started")

	rabbitmq.Consume(ctx, func(msg *internalqueue.Message) {
		logg.Info(fmt.Sprintf("Message received %v", msg))
	})
}
