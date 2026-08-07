package worker

import (
	"context"
	"net/http"
	"sync"
	"time"
)

func checkConnection(ctx context.Context, wg *sync.WaitGroup, name string, url string, interval time.Duration,
	results chan<- connectionResult) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	doCheck(ctx, name, url, results)

	for {
		select {
		case <-ticker.C:
			doCheck(ctx, name, url, results)
		case <-ctx.Done():
			return
		}
	}
}

func doCheck(ctx context.Context, name string, url string, results chan<- connectionResult) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		sendResult(ctx, results, connectionResult{up: false, name: name, statusCode: 0,
			err: err, checkedAt: time.Now()})
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		sendResult(ctx, results, connectionResult{up: false, name: name, statusCode: 0,
			err: err, checkedAt: time.Now()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		sendResult(ctx, results, connectionResult{up: true, name: name, statusCode: resp.StatusCode,
			err: nil, checkedAt: time.Now()})
		return
	}
	sendResult(ctx, results, connectionResult{up: false, name: name, statusCode: resp.StatusCode,
		err: nil, checkedAt: time.Now()})
}

func sendResult(ctx context.Context, results chan<- connectionResult, res connectionResult) {
	select {
	case results <- res:
	case <-ctx.Done():
	}
}
