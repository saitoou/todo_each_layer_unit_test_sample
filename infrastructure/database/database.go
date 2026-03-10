package infrastructure

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

func NewDatabase(ctx context.Context, dbUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql DB from gorm: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}
	return db, nil
}
