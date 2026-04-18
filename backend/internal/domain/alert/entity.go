package alert

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Severity string
type Type string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"

	TypeOffline     Type = "offline"
	TypeError       Type = "error"
	TypeMaintenance Type = "maintenance"
	TypeOverload    Type = "overload"
	TypeInfo        Type = "info"
)

type Alert struct {
	ID         uuid.UUID  `json:"id"`
	LabID      uuid.UUID  `json:"lab_id"`
	ComputerID *uuid.UUID `json:"computer_id,omitempty"`
	Type       Type       `json:"type"`
	Severity   Severity   `json:"severity"`
	Message    string     `json:"message"`
	Resolved   bool       `json:"resolved"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func New(labID uuid.UUID, computerID *uuid.UUID, alertType Type, severity Severity, message string) *Alert {
	return &Alert{
		ID:         uuid.New(),
		LabID:      labID,
		ComputerID: computerID,
		Type:       alertType,
		Severity:   severity,
		Message:    message,
		Resolved:   false,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

// Resolve marca o alerta como resolvido.
func (a *Alert) Resolve() {
	now := time.Now().UTC()
	a.Resolved = true
	a.ResolvedAt = &now
	a.UpdatedAt = now
}

func (a *Alert) IsCritical() bool { return a.Severity == SeverityCritical }

// ErrNotFound é o erro de domínio para alerta não encontrado.
var ErrNotFound = errors.New("alerta não encontrado")
