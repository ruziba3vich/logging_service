package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ruziba3vich/logging_service/pkg/config"
)

func ConnectAndMigrate(cfg config.DbConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open ClickHouse connection: %s", err.Error())
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %s", err.Error())
	}

	// Migrate database
	if err := migrate(db); err != nil {
		return nil, err
	}

	log.Println("Connected to ClickHouse and migration done")
	return db, nil
}

func migrate(db *sql.DB) error {
	migrationBytes, err := os.ReadFile("migrations/init.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %s", err.Error())
	}

	_, err = db.Exec(string(migrationBytes))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %s", err.Error())
	}

	return nil
}
