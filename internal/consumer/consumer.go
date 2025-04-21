package consumer

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ruziba3vich/logging_service/genprotos/genprotos/logging_service"
	"github.com/ruziba3vich/logging_service/internal/service"
	"github.com/ruziba3vich/logging_service/pkg/config"
	"google.golang.org/protobuf/proto"
)

type LogConsumer struct {
	consumer       *kafka.Consumer
	loggingService *service.LoggingService
	topic          string
	isRunning      bool
}

func NewLogConsumer(
	cfg *config.Config,
	loggingService *service.LoggingService,
) (*LogConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    cfg.ConsumerCfg.Brokers,
		"group.id":             cfg.ConsumerCfg.GroupID,
		"auto.offset.reset":    cfg.ConsumerCfg.AutoOffsetReset,
		"enable.auto.commit":   cfg.ConsumerCfg.EnableAutoCommit,
		"max.poll.interval.ms": cfg.ConsumerCfg.MaxPoolInterval,
		"session.timeout.ms":   cfg.ConsumerCfg.SessionTimeOut,
	})
	if err != nil {
		return nil, err
	}

	return &LogConsumer{
		consumer:       c,
		loggingService: loggingService,
		topic:          cfg.ConsumerCfg.Topic,
		isRunning:      false,
	}, nil
}

func (lc *LogConsumer) Start(ctx context.Context) error {
	if lc.isRunning {
		return nil
	}

	err := lc.consumer.Subscribe(lc.topic, nil)
	if err != nil {
		return err
	}

	lc.isRunning = true

	go func() {
		for {
			select {
			case <-ctx.Done():
				lc.Stop()
				return
			default:
				msg, err := lc.consumer.ReadMessage(-1)
				if err != nil {
					log.Printf("Error reading message: %v", err)
					continue
				}

				pbLog := &logging_service.Log{}
				if err := proto.Unmarshal(msg.Value, pbLog); err != nil {
					log.Printf("Error unmarshaling protobuf message: %v", err)
					lc.consumer.CommitMessage(msg)
					continue
				}

				_, err = lc.loggingService.SendLog(ctx, pbLog)
				if err != nil {
					log.Printf("Failed to process log: %v", err)
					continue
				}

				_, err = lc.consumer.CommitMessage(msg)
				if err != nil {
					log.Printf("Failed to commit offset: %v", err)
				}
			}
		}
	}()

	return nil
}

func (lc *LogConsumer) Stop() {
	if !lc.isRunning {
		return
	}
	lc.isRunning = false
	lc.consumer.Close()
}
