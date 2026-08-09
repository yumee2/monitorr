package storage

import (
	"context"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type ServiceStatus struct {
	ServiceName string
	IsUp        bool
	StatusCode  int
	CheckedAt   time.Time
}

type storage struct {
	db *sql.DB
}

func NewSqliteRepository(dsnURI string) (*storage, error) {
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
              CREATE TABLE IF NOT EXISTS services (
                      id          INTEGER PRIMARY KEY AUTOINCREMENT,
                      name        TEXT NOT NULL,
                      is_up       BOOLEAN NOT NULL,
                      status_code INTEGER NOT NULL,
                      checked_at  DATETIME NOT NULL
              );
              CREATE INDEX IF NOT EXISTS idx_services_name_checked_at ON services (name, checked_at);
      `)
	if err != nil {
		db.Close()
		return nil, err
	}
	return &storage{db: db}, nil
}

func (s *storage) SaveResult(ctx context.Context, serviceName string, isUp bool, statusCode int, checkedAt time.Time) error {
	_, err := s.db.Exec("INSERT INTO services (name, is_up, status_code, checked_at) VALUES (?, ?, ?, ?)", serviceName, isUp, statusCode, checkedAt)
	return err
}

func (s *storage) FindAll(ctx context.Context, interval time.Duration) (map[string][]ServiceStatus, error) {
	since := time.Now().Add(-interval)

	rows, err := s.db.QueryContext(ctx, "SELECT name, is_up, status_code, checked_at FROM services WHERE checked_at >= ? ORDER BY name, checked_at", since)
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

func (s *storage) Close() error {
	return s.db.Close()
}
