package main

import (
	"context"
	"monitorr/internal/config"
	"monitorr/internal/worker"
)

const cfgPath = "./config.yaml"

func main() {
	cfg, err := config.ReadConfig(cfgPath)
	if err != nil {
		panic(err)
	}

	worker.StartWorker(context.TODO(), cfg.Services)
}
