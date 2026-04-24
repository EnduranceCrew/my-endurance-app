package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"endurance/config"
	domainComputer "endurance/internal/domain/computer"
)

type gormComputer struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey"`
	LabID      uuid.UUID      `gorm:"type:uuid;not null"`
	Hostname   string         `gorm:"not null"`
	IPAddress  string
	MACAddress string
	Status     string         `gorm:"type:varchar(20);default:'offline'"`
	OS         string
	CPU        string
	RAM        string
	Storage    string
	LastSeen   *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (gormComputer) TableName() string { return "computers" }

func toDomainComputer(m *gormComputer) *domainComputer.Computer {
	return &domainComputer.Computer{
		ID:         m.ID,
		LabID:      m.LabID,
		Hostname:   m.Hostname,
		IPAddress:  m.IPAddress,
		MACAddress: m.MACAddress,
		Status:     domainComputer.Status(m.Status),
		OS:         m.OS,
		CPU:        m.CPU,
		RAM:        m.RAM,
		Storage:    m.Storage,
		LastSeen:   m.LastSeen,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func fromDomainComputer(c *domainComputer.Computer) *gormComputer {
	return &gormComputer{
		ID:         c.ID,
		LabID:      c.LabID,
		Hostname:   c.Hostname,
		IPAddress:  c.IPAddress,
		MACAddress: c.MACAddress,
		Status:     string(c.Status),
		OS:         c.OS,
		CPU:        c.CPU,
		RAM:        c.RAM,
		Storage:    c.Storage,
	}
}

type ComputerRepository struct{}

func NewComputerRepository() *ComputerRepository { return &ComputerRepository{} }
func (r *ComputerRepository) db() *gorm.DB        { return config.DB }

func (r *ComputerRepository) Create(_ context.Context, c *domainComputer.Computer) error {
	return r.db().Create(fromDomainComputer(c)).Error
}

func (r *ComputerRepository) FindByID(_ context.Context, id uuid.UUID) (*domainComputer.Computer, error) {
	var m gormComputer
	if err := r.db().First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainComputer.ErrNotFound
		}
		return nil, err
	}
	return toDomainComputer(&m), nil
}

func (r *ComputerRepository) FindByLabID(_ context.Context, labID uuid.UUID) ([]*domainComputer.Computer, error) {
	var models []gormComputer
	if err := r.db().Where("lab_id = ?", labID).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*domainComputer.Computer, 0, len(models))
	for i := range models {
		out = append(out, toDomainComputer(&models[i]))
	}
	return out, nil
}

func (r *ComputerRepository) FindAll(_ context.Context, page, limit int) ([]*domainComputer.Computer, int64, error) {
	var models []gormComputer
	var total int64
	offset := (page - 1) * limit
	r.db().Model(&gormComputer{}).Count(&total)
	if err := r.db().Offset(offset).Limit(limit).Order("hostname ASC").Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]*domainComputer.Computer, 0, len(models))
	for i := range models {
		out = append(out, toDomainComputer(&models[i]))
	}
	return out, total, nil
}

func (r *ComputerRepository) Update(_ context.Context, c *domainComputer.Computer) error {
	return r.db().Model(&gormComputer{}).Where("id = ?", c.ID).Updates(map[string]interface{}{
		"hostname":    c.Hostname,
		"ip_address":  c.IPAddress,
		"mac_address": c.MACAddress,
		"status":      string(c.Status),
		"os":          c.OS,
		"cpu":         c.CPU,
		"ram":         c.RAM,
		"storage":     c.Storage,
		"updated_at":  c.UpdatedAt,
	}).Error
}

func (r *ComputerRepository) UpdateStatus(_ context.Context, id uuid.UUID, status domainComputer.Status) error {
	now := time.Now().UTC()
	updates := map[string]interface{}{"status": string(status), "updated_at": now}
	if status == domainComputer.StatusOnline {
		updates["last_seen"] = now
	}
	return r.db().Model(&gormComputer{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ComputerRepository) Delete(_ context.Context, id uuid.UUID) error {
	return r.db().Delete(&gormComputer{}, "id = ?", id).Error
}

func (r *ComputerRepository) CountByStatus(_ context.Context) (map[domainComputer.Status]int64, error) {
	type result struct {
		Status string
		Count  int64
	}
	var results []result
	if err := r.db().Model(&gormComputer{}).Select("status, count(*) as count").Group("status").Scan(&results).Error; err != nil {
		return nil, err
	}
	m := make(map[domainComputer.Status]int64)
	for _, r := range results {
		m[domainComputer.Status(r.Status)] = r.Count
	}
	return m, nil
}
