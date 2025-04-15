package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AsaHero/whereismycity/internal/app"
	"github.com/AsaHero/whereismycity/pkg/config"
	"github.com/AsaHero/whereismycity/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := config.New()

	app, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	go func() {
		logger.Info(cfg.APP, "starting...")
		if err := app.Start(); err != nil {
			logger.Error("failed to start app", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	logger.Info(cfg.APP, "stopping...")

	app.Stop()
}
