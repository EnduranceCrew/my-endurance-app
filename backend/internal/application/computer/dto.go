package appcomputer

import (
	"time"

	"github.com/google/uuid"
	domainComputer "endurance/internal/domain/computer"
)

// ── Input DTOs ───────────────────────────────────────────────────────────────

type CreateInput struct {
	LabID      uuid.UUID `json:"lab_id"   binding:"required"`
	Hostname   string    `json:"hostname" binding:"required"`
	IPAddress  string    `json:"ip_address"`
	MACAddress string    `json:"mac_address"`
	OS         string    `json:"os"`
	CPU        string    `json:"cpu"`
	RAM        string    `json:"ram"`
	Storage    string    `json:"storage"`
}

type UpdateInput struct {
	Hostname   string                   `json:"hostname" binding:"required"`
	IPAddress  string                   `json:"ip_address"`
	MACAddress string                   `json:"mac_address"`
	OS         string                   `json:"os"`
	CPU        string                   `json:"cpu"`
	RAM        string                   `json:"ram"`
	Storage    string                   `json:"storage"`
	Status     domainComputer.Status    `json:"status"`
}

type UpdateStatusInput struct {
	Status domainComputer.Status `json:"status" binding:"required"`
}

// ── Output DTOs ──────────────────────────────────────────────────────────────

type ComputerOutput struct {
	ID         uuid.UUID             `json:"id"`
	LabID      uuid.UUID             `json:"lab_id"`
	Hostname   string                `json:"hostname"`
	IPAddress  string                `json:"ip_address"`
	MACAddress string                `json:"mac_address"`
	Status     domainComputer.Status `json:"status"`
	OS         string                `json:"os"`
	CPU        string                `json:"cpu"`
	RAM        string                `json:"ram"`
	Storage    string                `json:"storage"`
	LastSeen   *time.Time            `json:"last_seen"`
	CreatedAt  time.Time             `json:"created_at"`
}

type ListOutput struct {
	Computers []*ComputerOutput `json:"computers"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
}

func toOutput(c *domainComputer.Computer) *ComputerOutput {
	return &ComputerOutput{
		ID:         c.ID,
		LabID:      c.LabID,
		Hostname:   c.Hostname,
		IPAddress:  c.IPAddress,
		MACAddress: c.MACAddress,
		Status:     c.Status,
		OS:         c.OS,
		CPU:        c.CPU,
		RAM:        c.RAM,
		Storage:    c.Storage,
		LastSeen:   c.LastSeen,
		CreatedAt:  c.CreatedAt,
	}
}

// ── Port de entrada ──────────────────────────────────────────────────────────

type UseCase interface {
	GetAll(page, limit int) (*ListOutput, error)
	GetByLabID(labID uuid.UUID) ([]*ComputerOutput, error)
	GetByID(id uuid.UUID) (*ComputerOutput, error)
	Create(input CreateInput) (*ComputerOutput, error)
	Update(id uuid.UUID, input UpdateInput) (*ComputerOutput, error)
	UpdateStatus(id uuid.UUID, input UpdateStatusInput) error
	Delete(id uuid.UUID) error
}
