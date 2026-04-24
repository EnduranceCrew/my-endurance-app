package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"endurance/config"
	domainLab "endurance/internal/domain/lab"
)

type gormLab struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name          string         `gorm:"not null"`
	Location      string         `gorm:"not null"`
	Capacity      int            `gorm:"not null"`
	Status        string         `gorm:"type:varchar(20);default:'active'"`
	ResponsibleID *uuid.UUID     `gorm:"type:uuid"`
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (gormLab) TableName() string { return "labs" }

func toDomainLab(m *gormLab) *domainLab.Lab {
	return &domainLab.Lab{
		ID:            m.ID,
		Name:          m.Name,
		Location:      m.Location,
		Capacity:      m.Capacity,
		Status:        domainLab.Status(m.Status),
		ResponsibleID: m.ResponsibleID,
		Description:   m.Description,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func fromDomainLab(l *domainLab.Lab) *gormLab {
	return &gormLab{
		ID:            l.ID,
		Name:          l.Name,
		Location:      l.Location,
		Capacity:      l.Capacity,
		Status:        string(l.Status),
		ResponsibleID: l.ResponsibleID,
		Description:   l.Description,
	}
}

type LabRepository struct{}

func NewLabRepository() *LabRepository { return &LabRepository{} }
func (r *LabRepository) db(ctx context.Context) *gorm.DB  { return config.DB.WithContext(ctx) }

func (r *LabRepository) Create(ctx context.Context, l *domainLab.Lab) error {
	return r.db(ctx).Create(fromDomainLab(l)).Error
}

func (r *LabRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainLab.Lab, error) {
	var m gormLab
	if err := r.db(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainLab.ErrNotFound
		}
		return nil, err
	}
	return toDomainLab(&m), nil
}

func (r *LabRepository) FindAll(ctx context.Context, page, limit int) ([]*domainLab.Lab, int64, error) {
	var models []gormLab
	var total int64

	offset := (page - 1) * limit
	r.db(ctx).Model(&gormLab{}).Count(&total)

	if err := r.db(ctx).Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	labs := make([]*domainLab.Lab, 0, len(models))
	for i := range models {
		labs = append(labs, toDomainLab(&models[i]))
	}
	return labs, total, nil
}

func (r *LabRepository) Update(ctx context.Context, l *domainLab.Lab) error {
	return r.db(ctx).Model(&gormLab{}).Where("id = ?", l.ID).Updates(map[string]interface{}{
		"name":           l.Name,
		"location":       l.Location,
		"capacity":       l.Capacity,
		"status":         string(l.Status),
		"description":    l.Description,
		"responsible_id": l.ResponsibleID,
		"updated_at":     l.UpdatedAt,
	}).Error
}

func (r *LabRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db(ctx).Delete(&gormLab{}, "id = ?", id).Error
}

func (r *LabRepository) CountByStatus(ctx context.Context) (map[domainLab.Status]int64, error) {
	type result struct {
		Status string
		Count  int64
	}
	var results []result
	if err := r.db(ctx).Model(&gormLab{}).Select("status, count(*) as count").Group("status").Scan(&results).Error; err != nil {
		return nil, err
	}
	m := make(map[domainLab.Status]int64)
	for _, r := range results {
		m[domainLab.Status(r.Status)] = r.Count
	}
	return m, nil
}
