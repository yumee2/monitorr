package worker

import (
	"context"
	"net/http"
	"time"
)

func checkConnection(ctx context.Context, name string, url string, interval time.Duration,
	results chan<- connectionResult) {
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
		results <- connectionResult{up: false, name: name, statusCode: 0, err: err}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		results <- connectionResult{up: false, name: name, statusCode: 0, err: err}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		results <- connectionResult{up: true, name: name, statusCode: resp.StatusCode, err: nil}
		return
	}
	results <- connectionResult{up: false, name: name, statusCode: resp.StatusCode, err: nil}
}
