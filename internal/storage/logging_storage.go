package storage

import (
	"context"

	"github.com/ruziba3vich/logging_service/genprotos/genprotos/logging_service"
	"github.com/ruziba3vich/logging_service/internal/models"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type LoggingStorage struct {
	db *gorm.DB
}

func NewLoggingStorage(db *gorm.DB) *LoggingStorage {
	return &LoggingStorage{db: db}
}

// StoreLog inserts a new log entry into the database
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

// GetLastNLogs retrieves the most recent N logs with optional filtering
func (s *LoggingStorage) GetLastNLogs(ctx context.Context, req *logging_service.GetLastNLogsRequest) (*logging_service.GetLogsResponse, error) {
	var logs []models.Log

	query := s.db.WithContext(ctx).Model(&models.Log{})

	if req.LevelFilter != "" {
		query = query.Where("level = ?", req.LevelFilter)
	}

	if req.ServiceFilter != "" {
		query = query.Where("service = ?", req.ServiceFilter)
	}

	if err := query.Order("event_time DESC").Limit(int(req.N)).Find(&logs).Error; err != nil {
		return nil, err
	}

	return convertModelLogsToPbLogs(logs), nil
}

// GetLogsInTimeRange retrieves logs within a specific time range with optional filtering
func (s *LoggingStorage) GetLogsInTimeRange(ctx context.Context, req *logging_service.GetLogsInTimeRangeRequest) (*logging_service.GetLogsResponse, error) {
	var logs []models.Log

	query := s.db.WithContext(ctx).Model(&models.Log{})

	query = query.Where("event_time BETWEEN ? AND ?",
		req.StartTime.AsTime(),
		req.EndTime.AsTime())

	if req.LevelFilter != "" {
		query = query.Where("level = ?", req.LevelFilter)
	}

	if req.ServiceFilter != "" {
		query = query.Where("service = ?", req.ServiceFilter)
	}

	if err := query.Order("event_time DESC").Find(&logs).Error; err != nil {
		return nil, err
	}

	return convertModelLogsToPbLogs(logs), nil
}

// Helper function to convert from model logs to protobuf logs
func convertModelLogsToPbLogs(modelLogs []models.Log) *logging_service.GetLogsResponse {
	pbLogs := make([]*logging_service.Log, len(modelLogs))

	for i, log := range modelLogs {
		pbLogs[i] = &logging_service.Log{
			Message:   log.Message,
			EventTime: timestamppb.New(log.EventTime),
			Level:     log.Level,
			Service:   log.Service,
		}
	}

	return &logging_service.GetLogsResponse{
		Logs: pbLogs,
	}
}
