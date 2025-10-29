package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) error {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect db: %w", err)
	}
	if err := DB.AutoMigrate(&Content{}); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
