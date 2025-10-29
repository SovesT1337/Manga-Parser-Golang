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
		_ = database.ContentMarkParsed(content.ID, url)
		logger.Info("PROCESSOR", "marked parsed url=%s", content.UrlHentaichan)
		logger.Info("PROCESSOR", "processed url elapsed=%s", time.Since(start))
	}
}
