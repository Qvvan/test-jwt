package utils

import (
	"fmt"
	"strings"
)

// Коды ошибок PostgreSQL
const (
	PGErrDuplicateCode   = "SQLSTATE 23505" // Нарушение уникального ограничения
	PGErrNotNullCode     = "SQLSTATE 23502" // Нарушение NOT NULL ограничения
	PGErrForeignKey      = "SQLSTATE 23503" // Нарушение внешнего ключа
	PGErrCheckViolation  = "SQLSTATE 23514" // Нарушение CHECK ограничения
	PGErrUnexpectedError = "SQLSTATE 99999"
)

// Базовые ошибки
var (
	ErrDuplicate           = fmt.Errorf("duplicate record")
	ErrNotNullViolation    = fmt.Errorf("not null constraint violation")
	ErrForeignKeyViolation = fmt.Errorf("foreign key constraint violation")
	ErrCheckViolation      = fmt.Errorf("check constraint violation")
	ErrUnexpectedError     = fmt.Errorf("unexpected database error")
	ErrNotFound            = fmt.Errorf("record not found")
)

// CustomError представляет кастомную ошибку
type CustomError struct {
	Code    string
	Message error
	Log     string
}

// Реализация интерфейса error
func (e *CustomError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Log)
}

// Вспомогательные функции создания CustomError
func NewCreateError(table, operation string, err error) *CustomError {
	return &CustomError{
		Code:    PGErrDuplicateCode, // Используем PG код ошибки
		Message: ErrDuplicate,
		Log:     fmt.Sprintf("Table: %s, Operation: %s, Error: %v", table, operation, err),
	}
}

func NewDeleteError(table, id string, err error) *CustomError {
	return &CustomError{
		Code:    PGErrForeignKey,    // Используем PG код ошибки
		Message: ErrUnexpectedError, // Подставляем ошибку ErrUnexpectedError
		Log:     fmt.Sprintf("Table: %s, ID: %s, Error: %v", table, id, err),
	}
}

func NewUpdateError(table, id string, err error) *CustomError {
	return &CustomError{
		Code:    PGErrCheckViolation, // Используем PG код ошибки
		Message: ErrUnexpectedError,  // Подставляем ошибку ErrUnexpectedError
		Log:     fmt.Sprintf("Table: %s, ID: %s, Error: %v", table, id, err),
	}
}

func NewQueryError(table, action string, err error) *CustomError {
	return &CustomError{
		Code:    PGErrNotNullCode,   // Используем PG код ошибки
		Message: ErrUnexpectedError, // Подставляем ошибку ErrUnexpectedError
		Log:     fmt.Sprintf("Table: %s, Action: %s, Error: %v", table, action, err),
	}
}

// ParsePostgresError преобразует PostgreSQL ошибку в CustomError
func ParsePostgresError(err error) *CustomError {
	if err == nil {
		return nil
	}

	switch {
	case containsErrorCode(err, PGErrDuplicateCode):
		return &CustomError{
			Code:    PGErrDuplicateCode, // Используем PG код ошибки
			Message: ErrDuplicate,       // Подставляем ошибку ErrDuplicate
			Log:     "PostgreSQL duplicate key error",
		}
	case containsErrorCode(err, PGErrNotNullCode):
		return &CustomError{
			Code:    PGErrNotNullCode,    // Используем PG код ошибки
			Message: ErrNotNullViolation, // Подставляем ошибку ErrNotNullViolation
			Log:     "PostgreSQL not null constraint violation",
		}
	case containsErrorCode(err, PGErrForeignKey):
		return &CustomError{
			Code:    PGErrForeignKey,        // Используем PG код ошибки
			Message: ErrForeignKeyViolation, // Подставляем ошибку ErrForeignKeyViolation
			Log:     "PostgreSQL foreign key constraint violation",
		}
	case containsErrorCode(err, PGErrCheckViolation):
		return &CustomError{
			Code:    PGErrCheckViolation, // Используем PG код ошибки
			Message: ErrCheckViolation,   // Подставляем ошибку ErrCheckViolation
			Log:     "PostgreSQL check constraint violation",
		}
	default:
		return &CustomError{
			Code:    PGErrUnexpectedError, // Добавляем новый PG код ошибки, если необходимо
			Message: ErrUnexpectedError,   // Подставляем ошибку ErrUnexpectedError
			Log:     fmt.Sprintf("Unexpected database error: %v", err),
		}
	}
}

// containsErrorCode проверяет, содержит ли ошибка определенный код PostgreSQL
func containsErrorCode(err error, code string) bool {
	return err != nil && strings.Contains(err.Error(), code)
}
