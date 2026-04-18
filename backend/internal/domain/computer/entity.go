package computer

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusOnline  Status = "online"
	StatusOffline Status = "offline"
	StatusError   Status = "error"
	StatusIdle    Status = "idle"
)

type Computer struct {
	ID         uuid.UUID  `json:"id"`
	LabID      uuid.UUID  `json:"lab_id"`
	Hostname   string     `json:"hostname"`
	IPAddress  string     `json:"ip_address"`
	MACAddress string     `json:"mac_address"`
	Status     Status     `json:"status"`
	OS         string     `json:"os"`
	CPU        string     `json:"cpu"`
	RAM        string     `json:"ram"`
	Storage    string     `json:"storage"`
	LastSeen   *time.Time `json:"last_seen"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func New(labID uuid.UUID, hostname, ip, mac, os, cpu, ram, storage string) *Computer {
	return &Computer{
		ID:         uuid.New(),
		LabID:      labID,
		Hostname:   hostname,
		IPAddress:  ip,
		MACAddress: mac,
		Status:     StatusOffline,
		OS:         os,
		CPU:        cpu,
		RAM:        ram,
		Storage:    storage,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

func (c *Computer) MarkOnline() {
	now := time.Now().UTC()
	c.Status = StatusOnline
	c.LastSeen = &now
	c.UpdatedAt = now
}

func (c *Computer) MarkOffline() {
	c.Status = StatusOffline
	c.UpdatedAt = time.Now().UTC()
}

func (c *Computer) MarkError() {
	c.Status = StatusError
	c.UpdatedAt = time.Now().UTC()
}

func (c *Computer) IsOnline() bool { return c.Status == StatusOnline }

// ErrNotFound é o erro de domínio para computador não encontrado.
var ErrNotFound = errors.New("computador não encontrado")
