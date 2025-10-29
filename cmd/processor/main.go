package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/joho/godotenv"

	"go_scripts/config"
	"go_scripts/database"
	"go_scripts/internal/logger"
	"go_scripts/parsers"
	"go_scripts/telegraph"
)

func main() {
	logger.Info("PROCESSOR", "starting processor service")
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		logger.Error("PROCESSOR", "config: %v", err)
		os.Exit(1)
	}
	if err := database.InitDB(cfg.DatabaseDSN); err != nil {
		logger.DatabaseError("init db: %v", err)
		os.Exit(1)
	}

	for {
		content, err := database.ContentClaimNew()
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		start := time.Now()
		logger.Info("PROCESSOR", "processing url=%s", content.UrlHentaichan)
		info, err := parsers.HentaichanParseAll(content.UrlHentaichan)
		if err != nil {
			_ = database.ContentMarkError(content.ID, err.Error())
			logger.Error("PROCESSOR", "error parsing url=%s: %v", content.UrlHentaichan, err)
			continue
		}
		// store meta (series, author, translator, tags)
		tagsJSONBytes, _ := json.Marshal(info.Tags)
		_ = database.ContentUpdateMeta(content.ID, info.Title, info.Series, info.Author, info.Translator, string(tagsJSONBytes))
		logger.Info("PROCESSOR", "parsed url=%s, title=%s, series=%s, author=%s, translator=%s, tags=%v, images=%d", content.UrlHentaichan, info.Title, info.Series, info.Author, info.Translator, info.Tags, len(info.ImageURLs))
		url, err := telegraph.CreateTelegraphPage(info.Title, info.ImageURLs)
		if err != nil {
			_ = database.ContentMarkError(content.ID, err.Error())
			logger.Error("PROCESSOR", "error creating telegraph page: %v", err)
			continue
		}
		logger.Info("PROCESSOR", "created telegraph page url=%s", url)
		last, _ := database.ContentLastScheduledAt()
		base := time.Now()
		if last != nil {
			base = *last
		}
		sched := nextMoscowSlotAfter(base)
		logger.Info("PROCESSOR", "scheduling for %s", sched)
		_ = database.ContentMarkParsed(content.ID, url, sched)
		logger.Info("PROCESSOR", "marked parsed url=%s", content.UrlHentaichan)
		logger.Info("PROCESSOR", "processed url elapsed=%s", time.Since(start))
	}
}

func nextMoscowSlotAfter(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Europe/Moscow")
	tt := t.In(loc)
	y, m, d := tt.Date()
	slots := []time.Time{
		time.Date(y, m, d, 12, 0, 0, 0, loc),
		time.Date(y, m, d, 17, 0, 0, 0, loc),
		time.Date(y, m, d, 21, 0, 0, 0, loc),
	}
	for _, s := range slots {
		if s.After(tt) {
			return s
		}
	}
	nd := tt.Add(24 * time.Hour)
	y, m, d = nd.Date()
	return time.Date(y, m, d, 12, 0, 0, 0, loc)
}
