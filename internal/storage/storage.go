package storage

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type ServiceStatus struct {
	ServiceName string
	IsUp        bool
	StatusCode  int
	CheckedAt   int64
}

type Storage struct {
	db *sql.DB
}

func NewSqliteRepository(dsnURI string) (*Storage, error) {
	if dir := filepath.Dir(dsnURI); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.Exec("PRAGMA busy_timeout=5000;"); err != nil {
		db.Close()
		return nil, err
	}
	_, err = db.Exec(`
              CREATE TABLE IF NOT EXISTS checks (
                      id          INTEGER PRIMARY KEY AUTOINCREMENT,
                      name        TEXT NOT NULL,
                      is_up       BOOLEAN NOT NULL,
                      status_code INTEGER NOT NULL,
                      checked_at  INTEGER NOT NULL
              );
              CREATE INDEX IF NOT EXISTS idx_services_name_checked_at ON checks (name, checked_at);
      `)
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveResult(ctx context.Context, serviceName string, isUp bool, statusCode int, checkedAt int64) error {
	_, err := s.db.Exec("INSERT INTO checks (name, is_up, status_code, checked_at) VALUES (?, ?, ?, ?)", serviceName, isUp, statusCode, checkedAt)
	return err
}

func (s *Storage) FindAll(ctx context.Context, interval time.Duration) (map[string][]ServiceStatus, error) {
	since := time.Now().Add(-interval).Unix()
	rows, err := s.db.QueryContext(ctx, "SELECT name, is_up, status_code, checked_at FROM checks WHERE checked_at >= ? ORDER BY name, checked_at", since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grouped := make(map[string][]ServiceStatus)

	for rows.Next() {
		var status ServiceStatus
		if err := rows.Scan(&status.ServiceName, &status.IsUp, &status.StatusCode, &status.CheckedAt); err != nil {
			return nil, err
		}
		grouped[status.ServiceName] = append(grouped[status.ServiceName], status)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return grouped, nil
}

func (s *Storage) CountUptime(ctx context.Context, serviceName string, interval time.Duration) (int, error) {
	since := time.Now().Add(-interval).Unix()
	rows, err := s.db.QueryContext(ctx, "SELECT COUNT(*) FROM checks WHERE name = ? AND is_up = TRUE AND checked_at >= ?", serviceName, since)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	total, err := s.FindTotalByName(ctx, serviceName, interval)
	if err != nil {
		return 0, err
	}
	if len(total) == 0 {
		return 0, nil
	}
	return (count * 100) / len(total), nil
}

func (s *Storage) FindTotalByName(ctx context.Context, serviceName string, interval time.Duration) ([]ServiceStatus, error) {
	since := time.Now().Add(-interval).Unix()
	rows, err := s.db.QueryContext(ctx, "SELECT name, is_up, status_code, checked_at FROM checks WHERE name = ? AND checked_at >= ? ORDER BY checked_at", serviceName, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []ServiceStatus
	for rows.Next() {
		var status ServiceStatus
		if err := rows.Scan(&status.ServiceName, &status.IsUp, &status.StatusCode, &status.CheckedAt); err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return statuses, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
