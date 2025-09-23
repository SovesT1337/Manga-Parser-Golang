package parsers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func ParseHentaichan_v2(url string) (string, []string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", nil, fmt.Errorf("Ошибка загрузки: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("Ошибка сервера: %s\n", resp.Status)
	}

	scanner := bufio.NewScanner(resp.Body)
	lineNum := 1
	var name, fullimgData string
	var urls []string

	for scanner.Scan() {
		line := scanner.Text()
		
		if lineNum == 306 {
			name, err = extractName(line)
		}
		
		if lineNum == 312 {
			if strings.Contains(line, "\"fullimg\":") {
				fullimgData = extractArrayContent(line)
			} 
			
			if strings.Contains(line, "]") && fullimgData != "" {
				urls = parseImageUrls(fullimgData)
				break // Прекращаем после получения полного массива
			}
		}
		
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return "", nil, fmt.Errorf("Ошибка чтения: %v\n", err)
	}
	
	err = checkURLs(urls)
	if err := scanner.Err(); err != nil {
		return "", nil, fmt.Errorf("Не открывается ссылка: %v\n", err)
	}

	return name, urls, nil
}

// Извлекает название из строки
func extractName(line string) (string, error) {
	idx := strings.Index(line, "\"name\":")
	if idx == -1 {
		return "", fmt.Errorf("Название не найдено")
	}

	valuePart := line[idx+len("\"name\":"):]
	valuePart = strings.TrimSpace(valuePart)
	
	// Удаляем окружающие кавычки и запятые
	if len(valuePart) > 0 && (valuePart[0] == '"' || valuePart[0] == '\'') {
		end := strings.LastIndexAny(valuePart, "\"")
		if end > 1 {
			fmt.Println(valuePart)
			return valuePart[1:end], nil
		}
		return valuePart[1:], nil
	}
	
	return valuePart, nil
}

// Извлекает начало массива из строки
func extractArrayContent(line string) string {
	startIdx := strings.Index(line, "[")
	if startIdx == -1 {
		return "["
	}
	return line[startIdx:]
}

// Парсит JSON-массив URL
func parseImageUrls(data string) []string {
	// Заменяем одинарные кавычки на двойные для валидного JSON
	data = strings.ReplaceAll(data, "'", "\"")
	
	var urls []string
	if err := json.Unmarshal([]byte(data), &urls); err != nil {
		fmt.Printf("Ошибка парсинга: %v\n", err)
		return nil
	}
	return urls
}


func checkURLs(imageURLs []string) error {
	for _, link := range imageURLs {
		resp, err := http.Get(link)
		if err != nil {
			return fmt.Errorf("Ошибка загрузки %s: %v\n", link, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Ошибка сервера %s: %s\n", link, resp.Status)
		}
	}
	return nil
}

