package parsers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	
	"golang.org/x/net/html"
)

func ParseHentaichan(url string) ([]string, error) {
	// Загружаем HTML страницу
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки страницы: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неверный статус: %s", resp.Status)
	}

	// Парсим HTML напрямую из потока
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML: %w", err)
	}

	// Регулярное выражение для поиска массива fullimg
	re := regexp.MustCompile(`"fullimg":\s*\[([^\]]+)\]`)
	var links []string

	// Рекурсивная функция для поиска тегов <script>
	var findScripts func(*html.Node)
	findScripts = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			var scriptContent string
			// Собираем содержимое скрипта
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					scriptContent += c.Data
				}
			}

			// Ищем совпадения в тексте скрипта
			matches := re.FindStringSubmatch(scriptContent)
			if len(matches) > 1 {
				// Разделяем найденные ссылки
				urls := strings.Split(matches[1], ",")
				for _, url := range urls {
					// Очищаем URL
					cleanUrl := strings.TrimSpace(url)
					cleanUrl = strings.Trim(cleanUrl, `"' `)
					cleanUrl = strings.ReplaceAll(cleanUrl, " ", "")
					cleanUrl = strings.ReplaceAll(cleanUrl, "'", "")

					if cleanUrl != "" {
						links = append(links, cleanUrl)
					}
				}
			}
		}

		// Рекурсивно обходим дочерние узлы
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findScripts(c)
		}
	}

	findScripts(doc)
	return links, nil
}