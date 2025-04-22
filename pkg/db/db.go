package db

import (
	"fmt"

	"github.com/ruziba3vich/logging_service/pkg/config"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

func ConnectAndMigrate(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s",
		cfg.DbCfg.User,
		cfg.DbCfg.Password,
		cfg.DbCfg.Host,
		cfg.DbCfg.Port,
		cfg.DbCfg.Database,
	)

	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open ClickHouse connection: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return db, nil
}

func migrate(db *gorm.DB) error {
	rawSQL := `
	CREATE TABLE IF NOT EXISTS logs (
		id UUID DEFAULT generateUUIDv4(),
		message String,
		event_time DateTime,
		level String,
		service String,
		received_at DateTime DEFAULT now()
	) ENGINE = MergeTree
	PARTITION BY toYYYYMM(event_time)
	ORDER BY (event_time)
	TTL event_time + INTERVAL 30 DAY
	SETTINGS index_granularity = 8192;
	`

	return db.Exec(rawSQL).Error
}
