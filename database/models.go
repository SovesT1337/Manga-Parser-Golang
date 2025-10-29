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
	Status        string     `gorm:"type:varchar(16);index"` // New, Parsed, Sent, Error
	LastError     string     `gorm:"type:text"`
	ScheduledAt   *time.Time `gorm:"index"`
	SentAt        *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
