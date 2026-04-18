package lab

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive      Status = "active"
	StatusInactive    Status = "inactive"
	StatusMaintenance Status = "maintenance"
)

type Lab struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Location      string     `json:"location"`
	Capacity      int        `json:"capacity"`
	Status        Status     `json:"status"`
	ResponsibleID *uuid.UUID `json:"responsible_id,omitempty"`
	Description   string     `json:"description"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func New(name, location, description string, capacity int, responsibleID *uuid.UUID) *Lab {
	return &Lab{
		ID:            uuid.New(),
		Name:          name,
		Location:      location,
		Capacity:      capacity,
		Status:        StatusActive,
		ResponsibleID: responsibleID,
		Description:   description,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
}

func (l *Lab) SetMaintenance() {
	l.Status = StatusMaintenance
	l.UpdatedAt = time.Now().UTC()
}

func (l *Lab) Activate() {
	l.Status = StatusActive
	l.UpdatedAt = time.Now().UTC()
}

func (l *Lab) Deactivate() {
	l.Status = StatusInactive
	l.UpdatedAt = time.Now().UTC()
}

func (l *Lab) IsOperational() bool {
	return l.Status == StatusActive
}

// ErrNotFound é o erro de domínio para laboratório não encontrado.
var ErrNotFound = errors.New("laboratório não encontrado")
