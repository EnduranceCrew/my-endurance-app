package appuser

import (
	"context"
	"time"

	"github.com/google/uuid"
	domainUser "endurance/internal/domain/user"
)

// ── Input DTOs ───────────────────────────────────────────────────────────────

type UpdateInput struct {
	Name   string `json:"name"  binding:"required,min=2,max=100"`
	Active *bool  `json:"active"`
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password"     binding:"required,min=8"`
}

type ChangeRoleInput struct {
	Role string `json:"role" binding:"required,oneof=admin user"`
}

type PaginationInput struct {
	Page  int `form:"page,default=1"   binding:"min=1"`
	Limit int `form:"limit,default=20" binding:"min=1,max=100"`
}

// ── Output DTOs ──────────────────────────────────────────────────────────────

type UserOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CPF       string    `json:"cpf"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListOutput struct {
	Users []*UserOutput `json:"users"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}

func toOutput(u *domainUser.User) *UserOutput {
	return &UserOutput{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CPF:       u.CPF,
		Role:      string(u.Role),
		Active:    u.Active,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ── Port de entrada ──────────────────────────────────────────────────────────

type UseCase interface {
	GetAll(ctx context.Context, page, limit int) (*ListOutput, error)
	GetByID(ctx context.Context, id uuid.UUID) (*UserOutput, error)
	GetMe(ctx context.Context, id uuid.UUID) (*UserOutput, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*UserOutput, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ChangePassword(ctx context.Context, id uuid.UUID, input ChangePasswordInput) error
	ChangeRole(ctx context.Context, id uuid.UUID, input ChangeRoleInput, requestorID uuid.UUID) (*UserOutput, error)
}
