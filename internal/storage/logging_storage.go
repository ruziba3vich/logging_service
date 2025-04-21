package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/ruziba3vich/logging_service/genprotos/genprotos/logging_service"
)

type LoggingStorage struct {
	db *sql.DB
}

func NewLoggingStorage(db *sql.DB) *LoggingStorage {
	return &LoggingStorage{db: db}
}

func (s *LoggingStorage) StoreLog(ctx context.Context, l *logging_service.Log) {
	query := `INSERT INTO logs (message, event_time) VALUES (?, ?)`

	_, err := s.db.ExecContext(ctx, query, l.Message, l.EventTime.AsTime())
	if err != nil {
		log.Printf("Failed to store log: %v. Trying to log this error.", err)
		errorMsg := "Failed to store log: " + err.Error()
		_, innerErr := s.db.ExecContext(ctx, query, errorMsg, time.Now())
		if innerErr != nil {
			log.Printf("Failed to store internal error log too: %v", innerErr)
		}
	}
}
