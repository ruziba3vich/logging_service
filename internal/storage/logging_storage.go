package storage

import (
	"context"

	"github.com/ruziba3vich/logging_service/genprotos/genprotos/logging_service"
	"github.com/ruziba3vich/logging_service/internal/models"
	"gorm.io/gorm"
)

type LoggingStorage struct {
	db *gorm.DB
}

func NewLoggingStorage(db *gorm.DB) *LoggingStorage {
	return &LoggingStorage{db: db}
}

// query := `INSERT INTO logs (message, event_time, level, service) VALUES (?, ?, ?, ?)`
func (s *LoggingStorage) StoreLog(ctx context.Context, l *logging_service.Log) error {
	log := models.Log{
		Message:   l.Message,
		EventTime: l.EventTime.AsTime(),
		Level:     l.Level,
		Service:   l.Service,
	}
	if err := s.db.WithContext(ctx).Create(&log).Error; err != nil {
		return err
	}
	return nil
}
