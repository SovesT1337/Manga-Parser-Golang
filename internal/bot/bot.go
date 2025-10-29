package bot

import (
	"context"
	"time"

	"go_scripts/internal/fsm"
	"go_scripts/internal/logger"
	"go_scripts/internal/telegram"
)

type Bot struct {
	botURL  string
	handler *Handler
}

func New(botURL string, manager *fsm.Manager) *Bot {
	return &Bot{botURL: botURL, handler: NewHandler(botURL, manager)}
}

func (b *Bot) Run(ctx context.Context) {
	offset := 0
	for {
		select {
		case <-ctx.Done():
			logger.BotInfo("bot run stopped")
			return
		default:
		}

		updates, err := telegram.GetUpdates(b.botURL, offset)
		if err != nil {
			logger.TelegramError("get updates failed: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		for _, u := range updates {
			offset = u.UpdateID + 1
			b.handler.Handle(context.Background(), u)
		}
	}
}
