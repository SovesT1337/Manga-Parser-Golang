package main

import (
	"bytes"
	"encoding/json"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"strings"

	"x.localhost/scripts/parsers"
	"x.localhost/scripts/telegraph"
	"x.localhost/scripts/database"
	"x.localhost/scripts/schemas"

	"github.com/joho/godotenv"
)

const (
	pollDelay  = 1 * time.Second
	parseMode = "MarkdownV2"
)

var botToken string
var apiURL string
var repo database.ContentRepositoryInterface
var httpClient = &http.Client{Timeout: 20 * time.Second}

func startUp() {
	_ = godotenv.Load()
	
	botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	apiURL = os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "https://api.telegram.org/bot"
	}
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}
	
	if err := database.InitDB("hentai.db"); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	
	repo = database.NewContentRepository()

	log.Println("Bot Initiated")
}

func startPolling() {
	var offset int
	log.Println("Polling Started")
	for {
		log.Println("Checking New Messages")
		updates, err := getUpdates(offset)
		if err != nil {
			log.Printf("Ошибка получения обновлений: %v", err)
			time.Sleep(pollDelay)
			continue
		}
		
		for _, update := range updates {
			offset = update.UpdateID + 1
			go processMessage(update.Message)
		}
		
		time.Sleep(pollDelay)
	}
}


func main() {
	startUp()
	startPolling()
}

// Получение обновлений через long polling
func getUpdates(offset int) ([]schemas.Update, error) {
	url := buildTelegramURL("getUpdates") + fmt.Sprintf("?timeout=60&offset=%d", offset)

	ctx, cancel := context.WithTimeout(context.Background(), 65*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		OK     bool             `json:"ok"`
		Result []schemas.Update `json:"result"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if !response.OK {
		return nil, fmt.Errorf("ошибка API: %s", string(body))
	}

	return response.Result, nil
}

// Обработка входящего сообщения
func processMessage(msg schemas.Message) {
	log.Printf(
		"Новое сообщение [%s]: %s",
		msg.From.Username,
		msg.Text,
	)

	url := strings.TrimSpace(msg.Text)
	if url == "" {
		return
	}

	url = strings.ReplaceAll(url, "/manga/", "/online/")

	var reply string

	exists, err := repo.ExistsByHentaichanURL(url)
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		reply = "Запись уже существует"
	} else {
		name, imageURLs, err := parsers.ParseHentaichan_v2(url)
		
		if err != nil {
			log.Printf("Ошибка: %v\n", err)
			return
		}

		res_url, err := telegraph.CreateTelegraphPage(name, imageURLs)
		if err != nil {
			log.Println("\nError:", err)
		}


		newContent := &database.HContent{
			Name:          name,
			UrlHentaichan: url,
			UrlTelegraph:  res_url,
		}

		if err := repo.Create(newContent); err != nil {
			log.Fatalf("Ошибка создания: %v", err)
		}

		name = escapeMarkdown(name)

		reply = fmt.Sprintf("Название: [%s](%s)\n\nПодписывайся: [NikoSan](https://t.me/+5mlshfg5qQozYjIy)", name, res_url)
	}

	if err := sendTextMessage(msg.Chat.ID, reply); err != nil {
		log.Printf("Ошибка отправки: %v", err)
	}
}

func escapeMarkdown(text string) string {
	chars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range chars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}


// Отправка текстового сообщения
func sendTextMessage(chatID int64, text string) error {
	message := schemas.SendMessage{
		ChatID: chatID,
		Text:   text,
		ParseMode: parseMode,
		LinkPreviewOptions: schemas.LinkPreviewOptions{
							PreferLargeMedia: true,
							ShowAboveText: false,
						},
	}
	
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := buildTelegramURL("sendMessage")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ошибка HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func buildTelegramURL(path string) string {
	return fmt.Sprintf("%s%s/%s", apiURL, botToken, path)
}