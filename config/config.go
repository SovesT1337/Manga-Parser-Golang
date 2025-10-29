package config

import (
	"os"
	"strconv"
	"time"

	appErr "go_scripts/internal/errors"
)

// Unified config and loader

type Config struct {
	TelegramAPI                string
	TelegramToken              string
	DatabaseDSN                string
	BotTimeout                 time.Duration
	BotMaxRetries              int
	SchedulerIntervalSec       int
	SchedulerTelegramChannelID int64
	LoggingLevel               string
}

func Load() (*Config, error) {
	c := &Config{}
	c.TelegramAPI = getEnv("TELEGRAM_API", "https://api.telegram.org/bot")
	c.TelegramToken = getEnv("TELEGRAM_BOT_TOKEN", "")
	if c.TelegramToken == "" {
		return nil, appErr.NewValidationError("Отсутствует TELEGRAM_BОT_TOKEN", "Токен обязателен")
	}
	c.DatabaseDSN = getEnv("POSTGRES_DSN", "")
	if c.DatabaseDSN == "" {
		return nil, appErr.NewValidationError("Отсутствует POSTGRES_DSN", "DSN обязателен")
	}
	if t, err := parseIntEnv("BOT_TIMEOUT", "30", "BOT_TIMEOUT"); err == nil {
		c.BotTimeout = time.Duration(t) * time.Second
	} else {
		return nil, err
	}
	if mr, err := parseIntEnv("MAX_RETRIES", "3", "MAX_RETRIES"); err == nil {
		c.BotMaxRetries = mr
	} else {
		return nil, err
	}

	if n, err := parseIntEnv("SCHEDULER_INTERVAL_SEC", "10", "SCHEDULER_INTERVAL_SEC"); err == nil {
		if n <= 0 {
			return nil, appErr.NewValidationError("Неверный SCHEDULER_INTERVAL_SEC", "Должен быть числом > 0")
		}
		c.SchedulerIntervalSec = n
	} else {
		return nil, err
	}
	if v := getEnv("TELEGRAM_CHANNEL_ID", ""); v != "" {
		id, err := parseInt64Env("TELEGRAM_CHANNEL_ID")
		if err != nil {
			return nil, err
		}
		c.SchedulerTelegramChannelID = id
	}

	if c.SchedulerTelegramChannelID == 0 {
		return nil, appErr.NewValidationError("Отсутствует TELEGRAM_CHANNEL_ID", "Должен быть задан для рассылки")
	}

	c.LoggingLevel = getEnv("LOG_LEVEL", "INFO")
	return c, nil
}

// Validation helper methods removed; validation occurs in Load() based on SERVICE

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseIntEnv(key, def, field string) (int, error) {
	v := getEnv(key, def)
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, appErr.NewValidationError("Неверный "+key, field+" должен быть числом")
	}
	return n, nil
}

func parseInt64Env(key string) (int64, error) {
	v := getEnv(key, "")
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, appErr.NewValidationError("Неверный "+key, "Должен быть int64")
	}
	return n, nil
}
