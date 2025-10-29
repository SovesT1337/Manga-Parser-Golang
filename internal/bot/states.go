package bot

import (
	"context"

	"go_scripts/database"
	"go_scripts/internal/telegram"
)

func (h *Handler) handleAwaitLink(ctx context.Context, chatID int64, userID int, text string) {
	if !looksLikeHTTPURL(text) {
		_ = telegram.SendMessage(h.botURL, chatID, "Пришлите корректную ссылку (http/https).")
		return
	}
	if exists, _ := database.ContentExistsByURL(text); exists {
		_ = telegram.SendMessage(h.botURL, chatID, "Такая ссылка уже есть в базе.")
		return
	}
	_ = telegram.SendMessage(h.botURL, chatID, "Обрабатываю ссылку...")
	_, _ = database.ContentCreateNew(text)
}
