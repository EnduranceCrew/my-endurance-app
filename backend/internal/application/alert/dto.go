package appalert

import (
	"context"
	"time"

	"github.com/google/uuid"
	domainAlert "endurance/internal/domain/alert"
)

type CreateInput struct {
	LabID      uuid.UUID          `json:"lab_id"    binding:"required"`
	ComputerID *uuid.UUID         `json:"computer_id"`
	Type       domainAlert.Type   `json:"type"      binding:"required"`
	Severity   domainAlert.Severity `json:"severity" binding:"required"`
	Message    string             `json:"message"   binding:"required"`
}

type AlertOutput struct {
	ID         uuid.UUID           `json:"id"`
	LabID      uuid.UUID           `json:"lab_id"`
	ComputerID *uuid.UUID          `json:"computer_id,omitempty"`
	Type       domainAlert.Type    `json:"type"`
	Severity   domainAlert.Severity `json:"severity"`
	Message    string              `json:"message"`
	Resolved   bool                `json:"resolved"`
	ResolvedAt *time.Time          `json:"resolved_at,omitempty"`
	CreatedAt  time.Time           `json:"created_at"`
}

type ListOutput struct {
	Alerts []*AlertOutput `json:"alerts"`
	Total  int64          `json:"total"`
	Page   int            `json:"page"`
	Limit  int            `json:"limit"`
}

type BulkResolveInput struct {
	IDs []uuid.UUID `json:"ids" binding:"required,min=1"`
}

func toOutput(a *domainAlert.Alert) *AlertOutput {
	return &AlertOutput{
		ID:         a.ID,
		LabID:      a.LabID,
		ComputerID: a.ComputerID,
		Type:       a.Type,
		Severity:   a.Severity,
		Message:    a.Message,
		Resolved:   a.Resolved,
		ResolvedAt: a.ResolvedAt,
		CreatedAt:  a.CreatedAt,
	}
}

type UseCase interface {
	GetAll(ctx context.Context, onlyOpen bool, page, limit int) (*ListOutput, error)
	GetByLabID(ctx context.Context, labID uuid.UUID, onlyOpen bool) ([]*AlertOutput, error)
	Create(ctx context.Context, input CreateInput) (*AlertOutput, error)
	Resolve(ctx context.Context, id uuid.UUID) error
	BulkResolve(ctx context.Context, input BulkResolveInput) (int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
