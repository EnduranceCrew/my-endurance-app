package apperrors

import (
	"errors"
	"net/http"
)

// ── Erros de domínio (sem dependência de framework) ─────────────────────────

var (
	ErrNotFound           = errors.New("recurso não encontrado")
	ErrAlreadyExists      = errors.New("recurso já existe")
	ErrInvalidInput       = errors.New("entrada inválida")
	ErrUnauthorized       = errors.New("não autorizado")
	ErrForbidden          = errors.New("acesso negado")
	ErrInvalidCPF         = errors.New("CPF inválido")
	ErrInvalidEmail       = errors.New("e-mail inválido")
	ErrWeakPassword       = errors.New("senha muito fraca")
	ErrInactiveUser       = errors.New("usuário inativo")
	ErrInvalidCredentials = errors.New("e-mail ou senha incorretos")
)

// ── AppError: erro com código HTTP + mensagem legível ────────────────────────

type AppError struct {
	HTTPCode int
	Message  string
	Cause    error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Cause }

// New cria um AppError genérico.
func New(httpCode int, message string, cause error) *AppError {
	return &AppError{HTTPCode: httpCode, Message: message, Cause: cause}
}

// ── Construtores semânticos ──────────────────────────────────────────────────

func BadRequest(cause error) *AppError {
	return New(http.StatusBadRequest, cause.Error(), cause)
}

func Unauthorized(cause error) *AppError {
	return New(http.StatusUnauthorized, cause.Error(), cause)
}

func Forbidden(cause error) *AppError {
	return New(http.StatusForbidden, cause.Error(), cause)
}

func NotFound(cause error) *AppError {
	return New(http.StatusNotFound, cause.Error(), cause)
}

func Conflict(cause error) *AppError {
	return New(http.StatusConflict, cause.Error(), cause)
}

func Internal(cause error) *AppError {
	return New(http.StatusInternalServerError, "erro interno do servidor", cause)
}

// Is permite usar errors.Is com *AppError.
func Is(err, target error) bool { return errors.Is(err, target) }
