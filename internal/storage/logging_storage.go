package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/ruziba3vich/logging_service/genprotos/genprotos/logging_service"
)

type LoggingStorage struct {
	db          *sql.DB
	insertQuery string
}

func NewLoggingStorage(db *sql.DB, insertQuery string) *LoggingStorage {
	return &LoggingStorage{db: db, insertQuery: insertQuery}
}

func (s *LoggingStorage) StoreLog(ctx context.Context, l *logging_service.Log) {
	_, err := s.db.ExecContext(ctx, s.insertQuery,
		l.Message,
		l.EventTime.AsTime(),
		l.Level,
		l.Service,
	)

	if err != nil {
		log.Printf("Failed to store log: %v. Attempting to log the error internally.", err)

		_, innerErr := s.db.ExecContext(ctx, s.insertQuery,
			"Failed to store log: "+err.Error(),
			time.Now(),
			"ERROR",
			"LoggingService",
		)

		if innerErr != nil {
			log.Printf("Failed to store internal error log too: %v", innerErr)
		}
	}
}
