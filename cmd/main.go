package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/HandyDaddy/facts/internal/application"
	"github.com/HandyDaddy/facts/internal/config"
)

func main() {
	configPath := findConfigPath()
	cfg, err := config.Parse(configPath)
	if err != nil {
		log.Fatal(err)
	}

	cfg.HttpClient.Token = os.Getenv("BEARERTOKEN")
	if cfg.HttpClient.Token == "" {
		log.Fatal("BEARERTOKEN environment variable not set")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := application.NewFactService(cfg)
	service.Start(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal")

	service.Shutdown()
}

func findConfigPath() string {
	const (
		devConfig  = "config/dev.config.toml"
		prodConfig = "config/config.toml"
	)

	if os.Getenv("CFG") == "DEV" {
		return devConfig
	}

	return prodConfig
}
