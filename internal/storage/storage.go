package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type storage struct {
	db *sql.DB
}

func NewSqliteRepository(dsnURI string) (*storage, error) {
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		return nil, err
	}
	return &storage{db: db}, nil
}

func (s *storage)

func (s *storage) Close() error {
	return s.db.Close()
}
