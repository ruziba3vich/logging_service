package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/fx"

	"github.com/ruziba3vich/logging_service/internal/consumer"
	"github.com/ruziba3vich/logging_service/internal/service"
	"github.com/ruziba3vich/logging_service/internal/storage"
	"github.com/ruziba3vich/logging_service/pkg/config"
	"github.com/ruziba3vich/logging_service/pkg/db"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.New,
			db.ConnectAndMigrate,
			newLoggingStorage,
			service.NewLoggingService,
			consumer.NewLogConsumer,
		),
		fx.Invoke(registerHooks),
	)

	app.Run()
}

// LoggingStorage provider
func newLoggingStorage(db *sql.DB) *storage.LoggingStorage {
	insertQuery := `
		INSERT INTO logs (message, event_time, level, service) 
		VALUES (?, ?, ?, ?)
	`
	return storage.NewLoggingStorage(db, insertQuery)
}

// Register application lifecycle hooks
func registerHooks(
	lc fx.Lifecycle,
	logConsumer *consumer.LogConsumer,
	db *sql.DB,
) {
	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Println("Starting logging service...")

			if err := logConsumer.Start(ctx); err != nil {
				return fmt.Errorf("failed to start consumer: %s", err.Error())
			}

			log.Println("Logging service started")
			return nil
		},
		OnStop: func(context.Context) error {
			log.Println("Stopping logging service...")

			cancel()

			logConsumer.Stop()

			if err := db.Close(); err != nil {
				log.Printf("Error closing database connection: %v", err)
			}

			log.Println("Logging service stopped")
			return nil
		},
	})

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals

		log.Println("Received shutdown signal")
		cancel()
	}()
}
