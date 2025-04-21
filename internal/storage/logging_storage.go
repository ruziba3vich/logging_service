package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ruziba3vich/logging_service/genprotos/genprotos/logging_service"
)

type LoggingStorage struct {
	db          *sql.DB
	insertQuery string
}

func NewLoggingStorage(db *sql.DB, insertQuery string) *LoggingStorage {
	return &LoggingStorage{db: db, insertQuery: insertQuery}
}

// query := `INSERT INTO logs (message, event_time, level, service) VALUES (?, ?, ?, ?)`
func (s *LoggingStorage) StoreLog(ctx context.Context, l *logging_service.Log) error {
	_, err := s.db.ExecContext(ctx, s.insertQuery,
		l.Message,
		l.EventTime.AsTime(),
		l.Level,
		l.Service,
	)

	if err != nil {
		return fmt.Errorf("failed to store log: %s", err.Error())
	}
	return nil
}
