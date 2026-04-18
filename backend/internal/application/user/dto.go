package appuser

import (
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
	}
}

// ── Port de entrada ──────────────────────────────────────────────────────────

type UseCase interface {
	GetAll(page, limit int) (*ListOutput, error)
	GetByID(id uuid.UUID) (*UserOutput, error)
	Update(id uuid.UUID, input UpdateInput) (*UserOutput, error)
	Delete(id uuid.UUID) error
	ChangePassword(id uuid.UUID, input ChangePasswordInput) error
}
