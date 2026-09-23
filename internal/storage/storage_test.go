package storage

import (
	"context"
	"testing"
	"time"
)

func newTestStorage(t *testing.T) *Storage {
	t.Helper()

	s, err := NewSqliteRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory storage: %v", err)
	}
	t.Cleanup(func() {
		_ = s.Close()
	})
	return s
}

func TestCountUptime(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name        string
		checks      []ServiceStatus
		serviceName string
		interval    time.Duration
		want        int
	}{
		{
			name: "все проверки успешны -> 100%",
			checks: []ServiceStatus{
				{ServiceName: "api", IsUp: true, StatusCode: 200, CheckedAt: now.Unix()},
				{ServiceName: "api", IsUp: true, StatusCode: 200, CheckedAt: now.Add(-1 * time.Minute).Unix()},
			},
			serviceName: "api",
			interval:    time.Hour,
			want:        100,
		},
		{
			name: "половина проверок неуспешна -> 50%",
			checks: []ServiceStatus{
				{ServiceName: "api", IsUp: true, StatusCode: 200, CheckedAt: now.Unix()},
				{ServiceName: "api", IsUp: false, StatusCode: 500, CheckedAt: now.Add(-1 * time.Minute).Unix()},
			},
			serviceName: "api",
			interval:    time.Hour,
			want:        50,
		},
		{
			name: "все проверки неуспешны -> 0%",
			checks: []ServiceStatus{
				{ServiceName: "api", IsUp: false, StatusCode: 500, CheckedAt: now.Unix()},
			},
			serviceName: "api",
			interval:    time.Hour,
			want:        0,
		},
		{
			name:        "нет данных за период -> 0%, без деления на ноль",
			checks:      nil,
			serviceName: "api",
			interval:    time.Hour,
			want:        0,
		},
		{
			name: "старые проверки вне интервала не учитываются",
			checks: []ServiceStatus{
				{ServiceName: "api", IsUp: false, StatusCode: 500, CheckedAt: now.Add(-2 * time.Hour).Unix()}, // вне окна
				{ServiceName: "api", IsUp: true, StatusCode: 200, CheckedAt: now.Unix()},                      // в окне
			},
			serviceName: "api",
			interval:    time.Hour,
			want:        100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestStorage(t)

			for _, c := range tt.checks {
				if err := s.SaveResult(ctx, c.ServiceName, c.IsUp, c.StatusCode, c.CheckedAt); err != nil {
					t.Fatalf("SaveResult failed: %v", err)
				}
			}

			got, err := s.CountUptime(ctx, tt.serviceName, tt.interval)
			if err != nil {
				t.Fatalf("CountUptime returned error: %v", err)
			}
			if got != tt.want {
				t.Errorf("CountUptime() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSaveResultAndFindTotalByName(t *testing.T) {
	ctx := context.Background()
	s := newTestStorage(t)
	now := time.Now().Unix()

	if err := s.SaveResult(ctx, "api", true, 200, now); err != nil {
		t.Fatalf("SaveResult failed: %v", err)
	}

	got, err := s.FindTotalByName(ctx, "api", time.Hour)
	if err != nil {
		t.Fatalf("FindTotalByName failed: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].ServiceName != "api" || !got[0].IsUp || got[0].StatusCode != 200 {
		t.Errorf("unexpected result: %+v", got[0])
	}
}
