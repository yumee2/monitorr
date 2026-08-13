package main

import (
	"context"
	"log"
	"monitorr/internal/api"
	"monitorr/internal/config"
	"monitorr/internal/storage"
	"monitorr/internal/worker"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const cfgPath = "./config.yaml"

func main() {
	cfg, err := config.ReadConfig(cfgPath)
	if err != nil {
		panic(err)
	}

	repo, err := storage.NewSqliteRepository("data/storage.db")
	if err != nil {
		panic(err)
	}
	defer repo.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())

	doneWorker := make(chan struct{})
	go func() {
		worker.StartWorker(ctx, cfg.Services, repo)
		close(doneWorker)
	}()

	handler := api.NewHandler(repo)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: handler.Routes(),
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Printf("api server error: %v", err)
		}
	}()

	<-quit
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("api server shutdown error: %v", err)
	}

	<-doneWorker
}
