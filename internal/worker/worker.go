package worker

import (
	"context"
	"fmt"
	"monitorr/internal/config"
	"time"
)

type connectionResult struct {
	name       string
	up         bool
	statusCode int
	err        error
	checkedAt  time.Time
}

func StartWorker(ctx context.Context, services []config.Service) {
	result := make(chan connectionResult, len(services))
	for _, service := range services {
		serviceInterval := time.Duration(service.Interval) * time.Second
		go checkConnection(ctx, service.Name, service.URL, serviceInterval, result)
	}

	for {
		select {
		case res := <-result:
			fmt.Printf("%s: up=%v, statusCode=%d, err=%v\n", res.name, res.up, res.statusCode, res.err)
		case <-ctx.Done():
			return
		}
	}
}
