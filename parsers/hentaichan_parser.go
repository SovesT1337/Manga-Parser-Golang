package parsers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"regexp"
	"strings"
	"time"
)

func HentaichanParser(url string) (string, []string, error) {
	info, err := HentaichanParseAll(url)
	if err != nil {
		return "", nil, err
	}
	return info.Title, info.ImageURLs, nil
}

// HentaichanInfo содержит данные с обеих страниц (/manga/ и /online/)
type HentaichanInfo struct {
	Title      string
	Series     string
	Author     string
	Translator string
	Tags       []string
	ImageURLs  []string
}

// HentaichanParseAll загружает обе страницы (manga и online) и извлекает нужные данные.
func HentaichanParseAll(inputURL string) (*HentaichanInfo, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	mangaURL, onlineURL := derivePairURLs(inputURL)

	mangaHTML, err := httpGetString(client, mangaURL)
	if err != nil {
		return nil, err
	}
	onlineHTML, err := httpGetString(client, onlineURL)
	if err != nil {
		return nil, err
	}

	// meta from manga page
	title := extractTitle(mangaHTML)
	if title == "" {
		// fallback to JSON meta name from online page
		title = extractJSONName(onlineHTML)
	}
	series := extractField(mangaHTML, "Аниме/манга")
	author := extractField(mangaHTML, "Автор")
	translator := extractField(mangaHTML, "Переводчик")
	tags := extractTags(mangaHTML)

	// images from online page
	imgs := extractFullImgArray(onlineHTML)
	if len(imgs) == 0 {
		imgs = extractImgTags(onlineHTML)
	}
	if len(imgs) == 0 {
		imgs = extractAnyQuotedImages(onlineHTML)
	}
	imgs = normalizeURLs(onlineURL, imgs)

	return &HentaichanInfo{
		Title:      title,
		Series:     series,
		Author:     author,
		Translator: translator,
		Tags:       tags,
		ImageURLs:  imgs,
	}, nil
}

func httpGetString(client *http.Client, url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Referer", "https://x5.h-chan.me/")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка сервера: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения: %v", err)
	}
	return string(body), nil
}

func derivePairURLs(input string) (string, string) {
	u, err := neturl.Parse(input)
	if err != nil {
		// fallback: treat input as-is
		return input, input
	}
	// keep slug
	slug := u.Path[strings.LastIndex(u.Path, "/")+1:]
	// detect which variant is provided
	host := u.Scheme + "://" + u.Host
	if strings.Contains(u.Path, "/online/") {
		manga := host + "/manga/" + slug
		online := host + "/online/" + slug
		return manga, online
	}
	if strings.Contains(u.Path, "/manga/") {
		manga := host + "/manga/" + slug
		online := host + "/online/" + slug
		return manga, online
	}
	// unknown path; default to building both
	manga := host + "/manga/" + slug
	online := host + "/online/" + slug
	return manga, online
}

func extractTitle(html string) string {
	// Prefer <h1>..</h1>
	reH1 := regexp.MustCompile(`(?is)<h1[^>]*>(.*?)</h1>`)
	if m := reH1.FindStringSubmatch(html); len(m) == 2 {
		t := strings.TrimSpace(stripHTML(m[1]))
		if t != "" {
			return t
		}
	}
	// Fallback to JSON meta name
	return extractJSONName(html)
}

func extractJSONName(s string) string {
	reName := regexp.MustCompile(`(?is)["']name["']\s*:\s*["']([^"']+)["']`)
	if m := reName.FindStringSubmatch(s); len(m) == 2 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func extractField(html, label string) string {
	// Match a block: <div class="item">label</div> ... <h2> ... </h2>
	re := regexp.MustCompile(`(?is)<div\s+class=["']item["']>\s*` + regexp.QuoteMeta(label) + `\s*</div>\s*<div\s+class=["']item2["']>\s*<h2>(.*?)</h2>`)
	if m := re.FindStringSubmatch(html); len(m) == 2 {
		v := strings.TrimSpace(stripHTML(m[1]))
		return v
	}
	return ""
}

func extractTags(html string) []string {
	// Find tags in side list; ignore "+" and "-" entries
	re := regexp.MustCompile(`(?is)<li\s+class=["']sidetag["'][^>]*>(.*?)</li>`) // capture li content
	var tags []string
	for _, li := range re.FindAllStringSubmatch(html, -1) {
		if len(li) != 2 {
			continue
		}
		// extract all anchors' inner texts
		aRe := regexp.MustCompile(`(?is)<a[^>]*>([^<]+)</a>`)
		as := aRe.FindAllStringSubmatch(li[1], -1)
		if len(as) == 0 {
			continue
		}
		for _, a := range as {
			if len(a) != 2 {
				continue
			}
			text := strings.TrimSpace(stripHTML(a[1]))
			if text == "+" || text == "-" || text == "" {
				continue
			}
			// candidate tag text
			tags = append(tags, text)
			break // take the first non +/- text per li
		}
	}
	// dedupe preserving order
	seen := map[string]struct{}{}
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}
	return out
}

func extractFullImgArray(s string) []string {
	var urls []string
	reImgs := regexp.MustCompile(`(?is)["']fullimg["']\s*:\s*\[(.*?)\]`)
	if m := reImgs.FindStringSubmatch(s); len(m) == 2 {
		arr := "[" + m[1] + "]"
		arr = strings.ReplaceAll(arr, "'", "\"")
		arr = regexp.MustCompile(`,\s*]`).ReplaceAllString(arr, "]")
		_ = json.Unmarshal([]byte(arr), &urls)
	}
	return urls
}

func extractImgTags(s string) []string {
	var urls []string
	reTagImgs := regexp.MustCompile(`(?is)<img[^>]+(?:data-src|src)=["']([^"']+\.(?:jpe?g|png|webp))(?:\?[^"']*)?["']`)
	for _, m := range reTagImgs.FindAllStringSubmatch(s, -1) {
		if len(m) == 2 {
			urls = append(urls, strings.TrimSpace(m[1]))
		}
	}
	return urls
}

func extractAnyQuotedImages(s string) []string {
	var urls []string
	reAnyImgs := regexp.MustCompile(`(?is)\"(https?://[^\"]+\.(?:jpe?g|png|webp))(?:\?[^\"]*)?\"`)
	for _, m := range reAnyImgs.FindAllStringSubmatch(s, -1) {
		if len(m) == 2 {
			urls = append(urls, strings.TrimSpace(m[1]))
		}
	}
	return urls
}

func normalizeURLs(baseURL string, urls []string) []string {
	base, _ := neturl.Parse(baseURL)
	seen := make(map[string]struct{}, len(urls))
	norm := make([]string, 0, len(urls))
	for _, u := range urls {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		if bu, err := neturl.Parse(u); err == nil {
			if base != nil {
				bu = base.ResolveReference(bu)
			}
			u = bu.String()
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		norm = append(norm, u)
	}
	return norm
}

// stripHTML removes HTML tags; simplistic but adequate for titles
func stripHTML(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch r {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}
