package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/config"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/mq"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/logger"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/mailer"
	"github.com/hibiken/asynq"
)

func main() {
	config := config.InitConfig()

	logger.SetLogrusLogger(config.App)

	opt := asynq.RedisClientOpt{
		Addr: config.Redis.ServerAddress,
	}

	mailer := mailer.NewBrevoMailer(config.Smtp)
	consumer := mq.NewTaskConsumer(opt, mailer, config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumer.Start(ctx); err != nil {
			logger.Log.Fatalf("failed to start worker: %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	cancel()

	<-stop

	logger.Log.Info("shutting down worker...")
	consumer.Shutdown()

	logger.Log.Info("worker shutdown gracefully")
}
