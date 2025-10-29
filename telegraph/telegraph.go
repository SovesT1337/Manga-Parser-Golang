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
	Ok     bool       `json:"ok"`
	Result PageResult `json:"result"`
	Error  string     `json:"error,omitempty"`
}

type PageResult struct {
	URL string `json:"url"`
}

type Node struct {
	Tag   string            `json:"tag,omitempty"`
	Attrs map[string]string `json:"attrs,omitempty"`
}

func CreateTelegraphPage(title string, imageURLs []string) (string, error) {
	accessToken := os.Getenv("ACCESS_TOKEN")
	authorName := os.Getenv("AUTHOR_NAME")
	authorURL := os.Getenv("AUTHOR_URL")

	content := make([]Node, 0, len(imageURLs))
	for _, u := range imageURLs {
		content = append(content, Node{Tag: "img", Attrs: map[string]string{"src": u}})
	}

	body, err := json.Marshal(content)
	if err != nil {
		return "", fmt.Errorf("marshal content: %w", err)
	}

	params := url.Values{
		"access_token":   {accessToken},
		"title":          {title},
		"author_name":    {authorName},
		"author_url":     {authorURL},
		"content":        {string(body)},
		"return_content": {"false"},
	}

	resp, err := http.PostForm("https://api.telegra.ph/createPage", params)
	if err != nil {
		return "", fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	var r PageResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}
	if !r.Ok {
		return "", fmt.Errorf("API error: %s", r.Error)
	}
	return r.Result.URL, nil
}
