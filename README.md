# Telegram Hentai Bot

Telegram бот для парсинга контента с hentaichan и публикации на Telegraph.

## Функциональность

- Парсинг манги с hentaichan
- Создание страниц на Telegraph
- Сохранение в локальную базу данных SQLite
- Отправка ссылок в Telegram

## Установка

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd go_scripts
```

2. Установите зависимости:
```bash
go mod tidy
```

3. Создайте файл `.env` с переменными окружения:
```env
TELEGRAM_BOT_TOKEN=your_bot_token
ACCESS_TOKEN=your_telegraph_access_token
AUTHOR_NAME=your_author_name
AUTHOR_URL=your_author_url
API_URL=https://api.telegram.org/bot
```

4. Запустите бота:
```bash
go run bot.go
```

## Структура проекта

```
├── bot.go                 # Основной файл бота
├── database/              # Работа с базой данных
│   ├── database.go
│   ├── models.go
│   ├── operations.go
│   └── repository.go
├── parsers/               # Парсеры контента
│   ├── hentaichan_parser.go
│   └── hentaichan_parser_v2.go
├── schemas/               # Структуры данных
│   └── schemas.go
└── telegraph/             # Интеграция с Telegraph
    └── telegraph_poster.go
```

## Использование

Отправьте боту ссылку на мангу с hentaichan, и он:
1. Проверит, не обрабатывалась ли уже эта ссылка
2. Распарсит название и изображения
3. Создаст страницу на Telegraph
4. Сохранит в базу данных
5. Отправит ссылку на созданную страницу

## Требования

- Go 1.23+
- Telegram Bot Token
- Telegraph Access Token
