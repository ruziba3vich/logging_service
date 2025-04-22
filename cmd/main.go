package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/fx"
	"google.golang.org/grpc"

	"github.com/ruziba3vich/logging_service/genprotos/genprotos/logging_service"
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
			storage.NewLoggingStorage,
			service.NewLoggingService,
			consumer.NewLogConsumer,
			newGrpcServer,
		),
		fx.Invoke(registerHooks),
	)

	app.Run()
}

// Create a new gRPC server and register the logging service
func newGrpcServer(loggingService *service.LoggingService) *grpc.Server {
	server := grpc.NewServer()
	logging_service.RegisterLoggingServiceServer(server, loggingService)
	return server
}

// Register application lifecycle hooks
func registerHooks(
	lc fx.Lifecycle,
	logConsumer *consumer.LogConsumer,
	db *sql.DB,
	grpcServer *grpc.Server,
	cfg *config.Config,
) {
	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Println("Starting logging service...")

			if err := logConsumer.Start(ctx); err != nil {
				return fmt.Errorf("failed to start consumer: %s", err.Error())
			}

			listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
			if err != nil {
				return fmt.Errorf("failed to listen on port %s: %s", cfg.GRPCPort, err.Error())
			}

			log.Printf("gRPC server listening on port %s", cfg.GRPCPort)

			go func() {
				if err := grpcServer.Serve(listener); err != nil {
					log.Fatalf("Failed to serve gRPC: %v", err)
				}
			}()

			log.Println("Logging service started")
			return nil
		},
		OnStop: func(context.Context) error {
			log.Println("Stopping logging service...")

			cancel()

			logConsumer.Stop()

			grpcServer.GracefulStop()

			if err := db.Close(); err != nil {
				log.Printf("Error closing database connection: %v", err)
			}

			log.Println("Logging service stopped")
			return nil
		},
	})

	// Setup signal handling for graceful shutdown
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals

		log.Println("Received shutdown signal")
		cancel()
	}()
}
