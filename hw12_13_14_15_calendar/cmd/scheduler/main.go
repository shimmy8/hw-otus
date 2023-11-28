package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/logger"
	internalqueue "github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/queue"
	storagebuilder "github.com/shimmy8/hw-otus/hw12_13_14_15_calendar/internal/storage/builder"
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
	storage := storagebuilder.New(ctx, config.Storage.DB, config.Storage.URL, config.Storage.Timeout)
	logg := logger.New(config.Logger.Level, "scheduler")

	rabbitmq := internalqueue.New(
		config.AMQPConf.Username,
		config.AMQPConf.Password,
		config.AMQPConf.Host,
		config.AMQPConf.Port,
		config.AMQPConf.QueueName,
	)

	if err := rabbitmq.Connect(); err != nil {
		logg.Error(fmt.Sprintf("failed to connect to rabbitmq %v", err))
		rabbitmq.Close()
		os.Exit(1) //nolint:gocritic
	}
	defer rabbitmq.Close()

	notifyPeriod := time.Second * (time.Duration(config.Scheduler.NotifyCheckPeriod))
	go sendNotifications(ctx, notifyPeriod, rabbitmq, storage, logg)

	removePeriod := time.Second * (time.Duration(config.Scheduler.RemoveCheckPeriod))
	go removeOldNotifications(ctx, removePeriod, storage, logg)

	logg.Info("scheduler started")

	<-ctx.Done()
}

func sendNotifications(
	ctx context.Context,
	period time.Duration,
	provider *internalqueue.Provider,
	storage storagebuilder.Storage,
	logger *logger.Logger,
) {
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			events, err := storage.GetNotifyEvents(time.Now())
			if err != nil {
				logger.Error(fmt.Sprintf("failed to get events %v", err))
				continue
			}

			logger.Info(fmt.Sprintf("%d events found", len(events)))

			for _, evt := range events {
				msg := &internalqueue.Message{
					EventID:      evt.ID,
					EventTitle:   evt.Title,
					EventStartDt: evt.StartDT,
					UserID:       evt.UserID,
				}

				if err := provider.Publish(ctx, msg); err != nil {
					logger.Error(fmt.Sprintf("failed to publish event %v", err))
					continue
				}

				evt.Notified = true
				if err := storage.UpdateEvent(evt.ID, evt); err != nil {
					logger.Error(fmt.Sprintf("failed to update event %v", err))
					continue
				}

				logger.Info(fmt.Sprintf("notification published %v", evt))
			}
		}
	}
}

func removeOldNotifications(
	ctx context.Context,
	period time.Duration,
	storage storagebuilder.Storage,
	logger *logger.Logger,
) {
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			removeDt := time.Now().AddDate(-1, 0, 0)
			if err := storage.DeleteOldEvents(removeDt); err != nil {
				logger.Error(fmt.Sprintf("failed to remove events %v", err))
			}
		}
	}
}
