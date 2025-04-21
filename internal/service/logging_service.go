package service

import (
	"context"

	"github.com/ruziba3vich/logging_service/genprotos/genprotos/logging_service"
	"github.com/ruziba3vich/logging_service/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
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
	return nil, s.loggingStorage.StoreLog(ctx, l)
}
