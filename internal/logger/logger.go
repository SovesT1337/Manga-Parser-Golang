package logger

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	level  LogLevel
	logger *log.Logger
}

var (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[37m"
)

func New(level LogLevel) *Logger {
	return &Logger{level: level, logger: log.New(os.Stdout, "", 0)}
}

func (l *Logger) getColor(level LogLevel) string {
	switch level {
	case DEBUG:
		return colorGray
	case INFO:
		return colorBlue
	case WARN:
		return colorYellow
	case ERROR:
		return colorRed
	default:
		return colorReset
	}
}

func (l *Logger) getLevelName(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func getCallerInfo() (string, int, bool) {
	for i := 3; i <= 8; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			parts := strings.Split(file, "/")
			fileName := parts[len(parts)-1]
			if fileName != "logger.go" {
				return fileName, line, true
			}
		}
	}
	return "", 0, false
}

func maskSensitiveData(message string) string {
	tokenRegex := regexp.MustCompile(`bot\d+:[A-Za-z0-9_-]{35}`)
	message = tokenRegex.ReplaceAllString(message, "bot***:***")

	apiKeyRegex := regexp.MustCompile(`[A-Za-z0-9]{32,}`)
	message = apiKeyRegex.ReplaceAllStringFunc(message, func(match string) string {
		if len(match) > 8 {
			return match[:4] + "***" + match[len(match)-4:]
		}
		return "***"
	})
	return message
}

func (l *Logger) formatMessage(level LogLevel, context, message string) string {
	if level != ERROR {
		message = maskSensitiveData(message)
	}
	timestamp := time.Now().Format("15:04:05")
	color := l.getColor(level)
	levelName := l.getLevelName(level)
	file, line, ok := getCallerInfo()

	var contextStr string
	if context != "" {
		contextStr = fmt.Sprintf("[%s] ", context)
	}
	if ok {
		return fmt.Sprintf("%s%s %s[%s]%s %s:%d %s%s%s",
			colorGray, timestamp, color, levelName, colorReset, file, line, contextStr, message, colorReset,
		)
	}
	return fmt.Sprintf("%s%s %s[%s]%s %s%s%s",
		colorGray, timestamp, color, levelName, colorReset, contextStr, message, colorReset,
	)
}

func (l *Logger) log(level LogLevel, context, format string, args ...interface{}) {
	if level < l.level {
		return
	}
	message := fmt.Sprintf(format, args...)
	formatted := l.formatMessage(level, context, message)
	l.logger.Println(formatted)
}

func (l *Logger) Debug(context, format string, args ...interface{}) {
	l.log(DEBUG, context, format, args...)
}
func (l *Logger) Info(context, format string, args ...interface{}) {
	l.log(INFO, context, format, args...)
}
func (l *Logger) Warn(context, format string, args ...interface{}) {
	l.log(WARN, context, format, args...)
}
func (l *Logger) Error(context, format string, args ...interface{}) {
	l.log(ERROR, context, format, args...)
}
func (l *Logger) SetLevel(level LogLevel) { l.level = level }

var defaultLogger = New(INFO)

func Debug(context, format string, args ...interface{}) {
	defaultLogger.Debug(context, format, args...)
}
func Info(context, format string, args ...interface{}) { defaultLogger.Info(context, format, args...) }
func Warn(context, format string, args ...interface{}) { defaultLogger.Warn(context, format, args...) }
func Error(context, format string, args ...interface{}) {
	defaultLogger.Error(context, format, args...)
}
func SetLevel(level LogLevel) { defaultLogger.SetLevel(level) }

func BotInfo(format string, args ...interface{})       { Info("BOT", format, args...) }
func BotError(format string, args ...interface{})      { Error("BOT", format, args...) }
func TelegramInfo(format string, args ...interface{})  { Info("TELEGRAM", format, args...) }
func TelegramError(format string, args ...interface{}) { Error("TELEGRAM", format, args...) }
func TelegramWarn(format string, args ...interface{})  { Warn("TELEGRAM", format, args...) }
func DatabaseInfo(format string, args ...interface{})  { Info("DATABASE", format, args...) }
func DatabaseError(format string, args ...interface{}) { Error("DATABASE", format, args...) }
func UserInfo(userID int, format string, args ...interface{}) {
	Info(fmt.Sprintf("USER_%d", userID), format, args...)
}
func UserError(userID int, format string, args ...interface{}) {
	Error(fmt.Sprintf("USER_%d", userID), format, args...)
}
func AdminInfo(adminID int, format string, args ...interface{}) {
	Info(fmt.Sprintf("ADMIN_%d", adminID), format, args...)
}
func AdminError(adminID int, format string, args ...interface{}) {
	Error(fmt.Sprintf("ADMIN_%d", adminID), format, args...)
}
func HTTPInfo(method, endpoint string, statusCode int, format string, args ...interface{}) {
	Info(fmt.Sprintf("HTTP_%s_%s_%d", method, endpoint, statusCode), format, args...)
}
func HTTPError(method, endpoint string, format string, args ...interface{}) {
	Error(fmt.Sprintf("HTTP_%s_%s", method, endpoint), format, args...)
}
