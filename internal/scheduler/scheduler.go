package scheduler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go_scripts/database"
	"go_scripts/internal/logger"
	"go_scripts/internal/telegram"
)

type Runner struct {
	BotURL      string
	ChannelID   int64
	IntervalSec int
}

func (r *Runner) Run(ctx context.Context) {
	if r.IntervalSec <= 0 {
		r.IntervalSec = 10
	}
	interval := time.Duration(r.IntervalSec) * time.Second
	for {
		select {
		case <-ctx.Done():
			logger.BotInfo("scheduler stopped")
			return
		case <-time.After(interval):
		}
		due, err := database.ContentFindDue(5)
		if err != nil {
			logger.Error("BOT", "due check: %v", err)
			continue
		}
		for _, item := range due {
			if item.UrlTelegraph == "" {
				_ = database.ContentMarkError(item.ID, "empty telegraph url")
				continue
			}
			// Build message text with meta fields
			text := buildMessageText(item)
			// Send message with large preview shown below text
			_ = telegram.SendMessageWithPreview(r.BotURL, r.ChannelID, text, item.UrlTelegraph, true, false)
			_ = database.ContentMarkSent(item.ID)
		}
	}
}

func buildMessageText(item database.Content) string {
	// Compose message using HTML formatting (safer around entities)
	b := strings.Builder{}
	if item.Name != "" && item.UrlTelegraph != "" {
		b.WriteString("<a href=\"")
		b.WriteString(escapeHTML(item.UrlTelegraph))
		b.WriteString("\">")
		b.WriteString(escapeHTML(item.Name))
		b.WriteString("</a>\n\n")
	}
	if item.Series != "" && item.Series != "Оригинальные работы" {
		b.WriteString("<b>Серия:</b> ")
		b.WriteString(escapeHTML(item.Series))
		b.WriteString("\n")
	}
	if item.Author != "" {
		b.WriteString("<b>Автор:</b> ")
		b.WriteString(escapeHTML(item.Author))
		b.WriteString("\n")
	}
	if item.Translator != "" {
		b.WriteString("<b>Переводчик:</b> ")
		b.WriteString(escapeHTML(item.Translator))
		b.WriteString("\n")
	}
	if item.TagsJSON != "" {
		if tags := parseTags(item.TagsJSON); len(tags) > 0 {
			b.WriteString("<b>Теги:</b> ")
			// prefix each tag with an escaped '#', replacing spaces with underscores
			for i, t := range tags {
				if i > 0 {
					b.WriteString(", ")
				}
				nt := normalizeTagString(t)
				b.WriteString("#")
				b.WriteString(escapeHTML(nt))
			}
			b.WriteString("\n")
		}
	}
	return b.String()
}

func parseTags(tagsJSON string) []string {
	var tags []string
	_ = json.Unmarshal([]byte(tagsJSON), &tags)
	return tags
}

// basic HTML escaping
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

// normalizeTagString lowercases and replaces internal spaces with underscores
func normalizeTagString(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	s = strings.ToLower(s)
	s = strings.Join(strings.Fields(s), "_")
	return s
}
