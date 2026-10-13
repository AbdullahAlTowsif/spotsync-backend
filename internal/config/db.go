package config

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB opens a GORM connection to PostgreSQL and configures the pool.
func ConnectDB(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.DatabaseURL,
		// Safe with poolers like Neon's/Supabase's PgBouncer, which don't support prepared statements in transaction mode.
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		TranslateError: true,
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}

	// Connection pool: GORM wraps database/sql, which manages these.
	sqlDB.SetMaxOpenConns(25)                  // hard cap on concurrent connections
	sqlDB.SetMaxIdleConns(10)                  // connections kept warm for reuse
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // recycle long-lived connections
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // drop connections idle too long

	// Fail fast if the database is unreachable.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}
