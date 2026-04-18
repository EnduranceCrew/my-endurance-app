package applab

import (
	"time"

	"github.com/google/uuid"
	domainLab "endurance/internal/domain/lab"
)

// ── Input DTOs ───────────────────────────────────────────────────────────────

type CreateInput struct {
	Name          string     `json:"name"           binding:"required,min=2,max=100"`
	Location      string     `json:"location"       binding:"required"`
	Capacity      int        `json:"capacity"       binding:"required,min=1"`
	Description   string     `json:"description"`
	ResponsibleID *uuid.UUID `json:"responsible_id"`
}

type UpdateInput struct {
	Name          string           `json:"name"           binding:"required,min=2,max=100"`
	Location      string           `json:"location"       binding:"required"`
	Capacity      int              `json:"capacity"       binding:"required,min=1"`
	Status        domainLab.Status `json:"status"`
	Description   string           `json:"description"`
	ResponsibleID *uuid.UUID       `json:"responsible_id"`
}

// ── Output DTOs ──────────────────────────────────────────────────────────────

type LabOutput struct {
	ID            uuid.UUID        `json:"id"`
	Name          string           `json:"name"`
	Location      string           `json:"location"`
	Capacity      int              `json:"capacity"`
	Status        domainLab.Status `json:"status"`
	Description   string           `json:"description"`
	ResponsibleID *uuid.UUID       `json:"responsible_id,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

type ListOutput struct {
	Labs  []*LabOutput `json:"labs"`
	Total int64        `json:"total"`
	Page  int          `json:"page"`
	Limit int          `json:"limit"`
}

func toOutput(l *domainLab.Lab) *LabOutput {
	return &LabOutput{
		ID:            l.ID,
		Name:          l.Name,
		Location:      l.Location,
		Capacity:      l.Capacity,
		Status:        l.Status,
		Description:   l.Description,
		ResponsibleID: l.ResponsibleID,
		CreatedAt:     l.CreatedAt,
		UpdatedAt:     l.UpdatedAt,
	}
}

// ── Port de entrada ──────────────────────────────────────────────────────────

type UseCase interface {
	GetAll(page, limit int) (*ListOutput, error)
	GetByID(id uuid.UUID) (*LabOutput, error)
	Create(input CreateInput) (*LabOutput, error)
	Update(id uuid.UUID, input UpdateInput) (*LabOutput, error)
	Delete(id uuid.UUID) error
}
