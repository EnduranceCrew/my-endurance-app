// Package auth contém os DTOs e a interface do caso de uso de autenticação.
package auth

import (
	"context"
	"github.com/google/uuid"
)

// ── Input DTOs (vêm do cliente) ──────────────────────────────────────────────

type LoginInput struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterInput struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	CPF      string `json:"cpf"      binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// ── Output DTOs (vão para o cliente) ─────────────────────────────────────────

type TokenOutput struct {
	AccessToken string      `json:"access_token"`
	TokenType   string      `json:"token_type"`
	UserID      uuid.UUID   `json:"user_id"`
	Name        string      `json:"name"`
	Email       string      `json:"email"`
	Role        string      `json:"role"`
	ExpiresIn   int         `json:"expires_in"`
}

// ── Port de entrada (primary port) ──────────────────────────────────────────

type UseCase interface {
	Login(ctx context.Context, input LoginInput) (*TokenOutput, error)
	Register(ctx context.Context, input RegisterInput) (*TokenOutput, error)
}
