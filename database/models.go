package database

import (
	"time"
	
	
)

// Модель для хранения контента
type HContent struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	UrlHentaichan string `gorm:"unique;not null"` // Уникальный индекс
	UrlTelegraph  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Кастомное имя таблицы
func (HContent) TableName() string {
	return "hentai_content"
}