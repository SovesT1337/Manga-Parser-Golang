package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type ErrorType string

const (
	ErrorTypeValidation ErrorType = "validation"
	ErrorTypeDatabase   ErrorType = "database"
	ErrorTypeTelegram   ErrorType = "telegram"
	ErrorTypeNetwork    ErrorType = "network"
	ErrorTypeInternal   ErrorType = "internal"
	ErrorTypeUser       ErrorType = "user"
)

type AppError struct {
	Type    ErrorType
	Message string
	Details string
	UserMsg string
	Code    string
	Context map[string]interface{}
	Inner   error
	Stack   string
}

func (e *AppError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Inner)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}
func (e *AppError) Unwrap() error { return e.Inner }
func (e *AppError) GetUserMessage() string {
	if e.UserMsg != "" {
		return e.UserMsg
	}
	return e.Message
}

func NewAppError(t ErrorType, msg string, inner error) *AppError {
	return &AppError{Type: t, Message: msg, Inner: inner, Stack: getStackTrace()}
}
func NewValidationError(msg, details string) *AppError {
	return &AppError{Type: ErrorTypeValidation, Message: msg, Details: details, UserMsg: "Проверьте введенные данные", Stack: getStackTrace()}
}
func NewDatabaseError(msg string, inner error) *AppError {
	return &AppError{Type: ErrorTypeDatabase, Message: msg, Inner: inner, UserMsg: "Ошибка работы с базой данных", Stack: getStackTrace()}
}
func NewTelegramError(msg string, inner error) *AppError {
	return &AppError{Type: ErrorTypeTelegram, Message: msg, Inner: inner, UserMsg: "Ошибка связи с Telegram", Stack: getStackTrace()}
}
func NewNetworkError(msg string, inner error) *AppError {
	return &AppError{Type: ErrorTypeNetwork, Message: msg, Inner: inner, UserMsg: "Ошибка сети", Stack: getStackTrace()}
}
func NewInternalError(msg string, inner error) *AppError {
	return &AppError{Type: ErrorTypeInternal, Message: msg, Inner: inner, UserMsg: "Внутренняя ошибка системы", Stack: getStackTrace()}
}
func NewUserError(msg string) *AppError {
	return &AppError{Type: ErrorTypeUser, Message: msg, UserMsg: msg, Stack: getStackTrace()}
}

func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = map[string]interface{}{}
	}
	e.Context[key] = value
	return e
}
func (e *AppError) WithUserMessage(m string) *AppError { e.UserMsg = m; return e }
func (e *AppError) WithCode(c string) *AppError        { e.Code = c; return e }

func getStackTrace() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	lines := strings.Split(string(buf[:n]), "\n")
	if len(lines) > 6 {
		lines = lines[6:]
	}
	return strings.Join(lines, "\n")
}

func HandleError(err error) string {
	if err == nil {
		return ""
	}
	if a, ok := err.(*AppError); ok {
		return a.GetUserMessage()
	}
	return "Произошла ошибка, попробуйте позже"
}
