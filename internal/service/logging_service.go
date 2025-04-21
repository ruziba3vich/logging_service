package service

import (
	"context"

	"github.com/ruziba3vich/logging_service/genprotos/genprotos/logging_service"
	"github.com/ruziba3vich/logging_service/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	LoggingService struct {
		loggingStorage *storage.LoggingStorage
		logging_service.UnimplementedLoggingServiceServer
	}
)

func NewLoggingService(loggingStorage *storage.LoggingStorage) *LoggingService {
	return &LoggingService{
		loggingStorage: loggingStorage,
	}
}

func (s *LoggingService) SendLog(ctx context.Context, l *logging_service.Log) (*emptypb.Empty, error) {
	err := s.loggingStorage.StoreLog(ctx, l)
	if err != nil {
		errLog := &logging_service.Log{
			Message:   "Failed to store log: " + err.Error(),
			Level:     "ERROR",
			EventTime: timestamppb.Now(),
			Service:   "LoggingService",
		}

		_ = s.loggingStorage.StoreLog(ctx, errLog)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
