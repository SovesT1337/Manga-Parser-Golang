package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	appErr "go_scripts/internal/errors"
	"go_scripts/internal/logger"
)

func GetUpdates(botURL string, offset int) ([]Update, error) {
	if offset < 0 {
		return nil, appErr.NewValidationError("Неверный offset", "offset должен быть неотрицательным числом")
	}
	resp, err := http.Get(botURL + "/getUpdates?timeout=25&offset=" + strconv.Itoa(offset))
	if err != nil {
		logger.TelegramError("Ошибка HTTP запроса: %v", err)
		return nil, appErr.NewNetworkError("Ошибка получения обновлений", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		logger.TelegramError("Ошибка API (код %d)", resp.StatusCode)
		return nil, appErr.NewTelegramError("Ошибка API Telegram", nil).WithCode(strconv.Itoa(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.TelegramError("Ошибка чтения тела ответа: %v", err)
		return nil, appErr.NewNetworkError("Ошибка чтения ответа", err)
	}
	var tr telegramResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		logger.TelegramError("Ошибка парсинга ответа: %v", err)
		return nil, appErr.NewTelegramError("Ошибка парсинга JSON", err)
	}
	logger.TelegramInfo("Получено %d обновлений", len(tr.Result))
	return tr.Result, nil
}
