package bot

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go_scripts/database"
	"go_scripts/internal/fsm"
	"go_scripts/internal/scheduler"
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
	if u.Callback != nil {
		h.handleCallback(ctx, *u.Callback)
		return
	}
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

func (h *Handler) handleCallback(ctx context.Context, cb telegram.CallbackQuery) {
	userID := int(cb.From.ID)
	chatID := cb.From.ID
	// Only admins can act on callbacks
	ok, _ := database.AdminExists(cb.From.ID)
	if !ok {
		return
	}
	data := cb.Data
	if strings.HasPrefix(data, "confirm:") {
		idStr := strings.TrimPrefix(data, "confirm:")
		if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			last, _ := database.ContentLastScheduledAt()
			base := time.Now()
			if last != nil {
				base = *last
			}
			sched := scheduler.NextMoscowSlotAfter(base)
			_ = database.ContentMarkConfirmedAndSchedule(uint(id), sched)
			_ = telegram.SendMessage(h.botURL, chatID, "Пост подтвержден и поставлен в очередь")
			// reset user state optionally
			_ = userID // keep for linter if unused
		}
	} else if strings.HasPrefix(data, "reject:") {
		idStr := strings.TrimPrefix(data, "reject:")
		if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
			_ = database.ContentMarkCancelled(uint(id))
			_ = telegram.SendMessage(h.botURL, chatID, "Пост отклонен")
		}
	}
}
