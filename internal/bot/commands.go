package bot

import (
	"context"

	"go_scripts/internal/fsm"
	"go_scripts/internal/logger"
	"go_scripts/internal/telegram"
)

func (h *Handler) handleCommand(ctx context.Context, chatID int64, userID int, text string) {
	switch text {
	case "/start":
		h.handleStart(ctx, chatID, userID)
	case "/cancel":
		h.handleCancel(ctx, chatID, userID)
	default:
		// ignore unknown commands for now
	}
}

func (h *Handler) handleStart(ctx context.Context, chatID int64, userID int) {
	cur, _ := h.manager.Get(userID)
	logger.UserInfo(userID, "/start prev_state=%v", cur)
	h.manager.Set(userID, fsm.AwaitLink())
	_ = telegram.SendMessage(h.botURL, chatID, "Привет! Пришли ссылку для парсера.")
}

func (h *Handler) handleCancel(ctx context.Context, chatID int64, userID int) {
	cur, _ := h.manager.Get(userID)
	h.manager.Set(userID, fsm.Start())
	logger.UserInfo(userID, "/cancel prev_state=%v", cur)
	_ = telegram.SendMessage(h.botURL, chatID, "Отменено")
}
