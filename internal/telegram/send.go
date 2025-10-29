package telegram

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	appErr "go_scripts/internal/errors"
	"go_scripts/internal/logger"
)

type sendMessage struct {
	ChatId      int64               `json:"chat_id"`
	Text        string              `json:"text"`
	ParseMode   string              `json:"parse_mode,omitempty"`
	LinkPreview *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	ReplyMarkup interface{}         `json:"reply_markup,omitempty"`
}

// LinkPreviewOptions mirrors Telegram Bot API link_preview_options
type LinkPreviewOptions struct {
	IsDisabled       bool   `json:"is_disabled,omitempty"`
	URL              string `json:"url,omitempty"`
	PreferSmallMedia bool   `json:"prefer_small_media,omitempty"`
	PreferLargeMedia bool   `json:"prefer_large_media,omitempty"`
	ShowAboveText    bool   `json:"show_above_text,omitempty"`
}

func SendMessage(botURL string, chatID int64, text string) error {
	body := sendMessage{ChatId: chatID, Text: text, ParseMode: "HTML"}
	buf, _ := json.Marshal(body)
	resp, err := http.Post(botURL+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		logger.TelegramError("Отправка сообщения: %v", err)
		return appErr.NewNetworkError("Ошибка отправки", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		logger.TelegramError("Ошибка API (код %d): %s", resp.StatusCode, string(body))
		return appErr.NewTelegramError("Ошибка API Telegram", nil)
	}
	logger.TelegramInfo("Сообщение отправлено")
	return nil
}

func SendMessageWithPreview(botURL string, chatID int64, text string, previewURL string, large bool, showAbove bool) error {
	body := sendMessage{
		ChatId:    chatID,
		Text:      text,
		ParseMode: "HTML",
		LinkPreview: &LinkPreviewOptions{
			URL:              previewURL,
			PreferLargeMedia: large,
			ShowAboveText:    showAbove,
		},
	}
	buf, _ := json.Marshal(body)
	resp, err := http.Post(botURL+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		logger.TelegramError("Отправка сообщения: %v", err)
		return appErr.NewNetworkError("Ошибка отправки", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		logger.TelegramError("Ошибка API (код %d): %s", resp.StatusCode, string(body))
		return appErr.NewTelegramError("Ошибка API Telegram", nil)
	}
	logger.TelegramInfo("Сообщение отправлено")
	return nil
}

func SendMessageWithPreviewAndKeyboard(botURL string, chatID int64, text string, previewURL string, large bool, showAbove bool, markup InlineKeyboardMarkup) error {
	body := sendMessage{
		ChatId:    chatID,
		Text:      text,
		ParseMode: "HTML",
		LinkPreview: &LinkPreviewOptions{
			URL:              previewURL,
			PreferLargeMedia: large,
			ShowAboveText:    showAbove,
		},
		ReplyMarkup: markup,
	}
	buf, _ := json.Marshal(body)
	resp, err := http.Post(botURL+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		logger.TelegramError("Отправка сообщения: %v", err)
		return appErr.NewNetworkError("Ошибка отправки", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		logger.TelegramError("Ошибка API (код %d): %s", resp.StatusCode, string(body))
		return appErr.NewTelegramError("Ошибка API Telegram", nil)
	}
	logger.TelegramInfo("Сообщение отправлено")
	return nil
}
