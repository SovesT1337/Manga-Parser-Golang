package telegraph

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// Структуры для обработки ответа API
type PageResponse struct {
	Ok     bool         `json:"ok"`
	Result PageResult   `json:"result"`
	Error  string       `json:"error,omitempty"`
}

type PageResult struct {
	URL string `json:"url"` // Основное поле, которое нам нужно
}

type Node struct {
	Tag   string            `json:"tag,omitempty"`
	Attrs map[string]string `json:"attrs,omitempty"`
}

func CreateTelegraphPage(title string, imageURLs []string) (string, error) {
	accessToken := os.Getenv("ACCESS_TOKEN")
	authorName := os.Getenv("AUTHOR_NAME")
	authorURL := os.Getenv("AUTHOR_URL")

	// Формирование контента
	contentNodes := make([]Node, 0, len(imageURLs))
	for _, u := range imageURLs {
		contentNodes = append(contentNodes, Node{
			Tag: "img",
			Attrs: map[string]string{
				"src": u,
			},
		})
	}

	contentBytes, err := json.Marshal(contentNodes)
	if err != nil {
		return "", fmt.Errorf("error marshaling content: %w", err)
	}

	// Отправка запроса
	params := url.Values{
		"access_token":  {accessToken},
		"title":         {title},
		"author_name":   {authorName},
		"author_url":    {authorURL},
		"content":       {string(contentBytes)},
		"return_content": {"false"},
	}

	resp, err := http.PostForm("https://api.telegra.ph/createPage", params)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Декодирование ответа
	var pageResp PageResponse
	if err := json.NewDecoder(resp.Body).Decode(&pageResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !pageResp.Ok {
		return "", fmt.Errorf("API error: %s", pageResp.Error)
	}

	return pageResp.Result.URL, nil
}