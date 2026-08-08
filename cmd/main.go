package main

import (
	"context"
	"monitorr/internal/config"
	"monitorr/internal/worker"
	"os"
	"os/signal"
	"syscall"
)

const cfgPath = "./config.yaml"

func main() {
	cfg, err := config.ReadConfig(cfgPath)
	if err != nil {
		panic(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())

	doneWorker := make(chan struct{})
	go func() {
		worker.StartWorker(ctx, cfg.Services)
		close(doneWorker)
	}()

	<-quit
	cancel()
	<-doneWorker
}
