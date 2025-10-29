package database

import (
	"time"
)

type Content struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	Series        string
	Author        string
	Translator    string
	TagsJSON      string `gorm:"type:text"`
	UrlHentaichan string `gorm:"uniqueIndex;not null"`
	UrlTelegraph  string
	Status        string     `gorm:"type:varchar(16);index"` // New, Processing, Parsed, Confirmed, Cancelled, Sent, Error
	LastError     string     `gorm:"type:text"`
	ScheduledAt   *time.Time `gorm:"index"`
	SentAt        *time.Time
	ReviewSentAt  *time.Time `gorm:"index"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Administrator struct {
	ID             uint   `gorm:"primaryKey"`
	TelegramUserID int64  `gorm:"uniqueIndex;not null"`
	Username       string `gorm:"type:varchar(255)"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
