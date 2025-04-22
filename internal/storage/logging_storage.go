package storage

import (
	"context"

	"github.com/ruziba3vich/logging_service/genprotos/genprotos/logging_service"
	"github.com/ruziba3vich/logging_service/internal/models"
	"gorm.io/gorm"
)

type LoggingStorage struct {
	db          *gorm.DB
	insertQuery string
}

func NewLoggingStorage(db *gorm.DB, insertQuery string) *LoggingStorage {
	return &LoggingStorage{db: db, insertQuery: insertQuery}
}

// query := `INSERT INTO logs (message, event_time, level, service) VALUES (?, ?, ?, ?)`
func (s *LoggingStorage) StoreLog(ctx context.Context, l *logging_service.Log) error {
	log := models.Log{
		Message:   l.Message,
		EventTime: l.EventTime.AsTime().Format("2006-01-02 15:04:05"), // format for ClickHouse
		Level:     l.Level,
		Service:   l.Service,
	}
	if err := s.db.WithContext(ctx).Create(&log).Error; err != nil {
		return err
	}
	return nil
}
