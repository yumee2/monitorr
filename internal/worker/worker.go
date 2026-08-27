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

type Alerter interface {
	Send(ctx context.Context, text string) error
}

type alertState struct {
	consecutiveFailures int
	down                bool
}

func StartWorker(ctx context.Context, services []config.Service, storage ResultSaver,
	alerter Alerter, failureThreshold int) {
	if failureThreshold <= 0 {
		failureThreshold = 1
	}

	result := make(chan connectionResult, len(services))
	wg := &sync.WaitGroup{}

	for _, service := range services {
		serviceInterval := time.Duration(service.Interval) * time.Second
		wg.Add(1)
		go checkConnection(ctx, wg, service.Name, service.URL, serviceInterval, result)
	}

	states := make(map[string]*alertState, len(services))

	for {
		select {
		case res := <-result:
			fmt.Printf("%s: up=%v, statusCode=%d, err=%v, checkedAt=%d\n",
				res.name, res.up, res.statusCode, res.err, res.checkedAt)
			err := storage.SaveResult(ctx, res.name, res.up, res.statusCode, res.checkedAt)
			if err != nil {
				fmt.Printf("failed to save result: %v\n", err)
			}
			if alerter != nil {
				handleAlert(ctx, alerter, states, res, failureThreshold)
			}
		case <-ctx.Done():
			wg.Wait()
			return
		}
	}
}

func handleAlert(ctx context.Context, alerter Alerter, states map[string]*alertState,
	res connectionResult, failureThreshold int) {
	state, ok := states[res.name]
	if !ok {
		state = &alertState{}
		states[res.name] = state
	}

	if res.up {
		state.consecutiveFailures = 0
		if state.down {
			state.down = false
			text := fmt.Sprintf("%s is back up (status %d)", res.name, res.statusCode)
			if err := alerter.Send(ctx, text); err != nil {
				fmt.Printf("failed to send alert: %v\n", err)
			}
		}
		return
	}

	state.consecutiveFailures++
	if !state.down && state.consecutiveFailures >= failureThreshold {
		state.down = true
		text := fmt.Sprintf("%s is down (status %d, err %v)", res.name, res.statusCode, res.err)
		if err := alerter.Send(ctx, text); err != nil {
			fmt.Printf("faisled to send alert: %v\n", err)
		}
	}
}
