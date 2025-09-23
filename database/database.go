package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Инициализация базы данных
func InitDB(filepath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Автомиграция моделей
	if err := DB.AutoMigrate(&HContent{}); err != nil {
		return fmt.Errorf("ошибка миграции: %w", err)
	}

	return nil
}