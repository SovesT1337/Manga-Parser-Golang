package bot

import (
	"context"
	"strings"

	"go_scripts/internal/fsm"
	"go_scripts/internal/telegram"
)

type Handler struct {
	botURL  string
	manager *fsm.Manager
}

func NewHandler(botURL string, manager *fsm.Manager) *Handler {
	return &Handler{botURL: botURL, manager: manager}
}

func (h *Handler) Handle(ctx context.Context, u telegram.Update) {
	if u.Message == nil {
		return
	}
	chatID := u.Message.Chat.ID
	text := u.Message.Text
	userID := int(chatID)
	if u.Message.From != nil {
		userID = int(u.Message.From.ID)
	}

	if strings.HasPrefix(text, "/") {
		h.handleCommand(ctx, chatID, userID, text)
		return
	}

	state, _ := h.manager.Get(userID)
	switch state.Type {
	case fsm.StateAwaitLink:
		h.handleAwaitLink(ctx, chatID, userID, text)
	default:
		return
	}
}
