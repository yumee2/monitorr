package worker

import (
	"context"
	"fmt"
	"monitorr/internal/config"
	"sync"
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
	wg := &sync.WaitGroup{}

	for _, service := range services {
		serviceInterval := time.Duration(service.Interval) * time.Second
		wg.Add(1)
		go checkConnection(ctx, wg, service.Name, service.URL, serviceInterval, result)
	}

	for {
		select {
		case res := <-result:
			fmt.Printf("%s: up=%v, statusCode=%d, err=%v, checkedAt=%s\n", res.name, res.up, res.statusCode, res.err, res.checkedAt)
		case <-ctx.Done():
			wg.Wait()
			return
		}
	}
}
