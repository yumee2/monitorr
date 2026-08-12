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
	checkedAt  int64
}

type ResultSaver interface {
	SaveResult(ctx context.Context, serviceName string, isUp bool,
		statusCode int, checkedAt int64) error
}

func StartWorker(ctx context.Context, services []config.Service, storage ResultSaver) {
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
			fmt.Printf("%s: up=%v, statusCode=%d, err=%v, checkedAt=%d\n",
				res.name, res.up, res.statusCode, res.err, res.checkedAt)
			err := storage.SaveResult(ctx, res.name, res.up, res.statusCode, res.checkedAt)
			if err != nil {
				fmt.Printf("failed to save result: %v\n", err)
			}
		case <-ctx.Done():
			wg.Wait()
			return
		}
	}
}
