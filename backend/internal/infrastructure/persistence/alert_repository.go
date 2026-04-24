package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"endurance/config"
	domainAlert "endurance/internal/domain/alert"
)

type gormAlert struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey"`
	LabID      uuid.UUID      `gorm:"type:uuid;not null"`
	ComputerID *uuid.UUID     `gorm:"type:uuid"`
	Type       string         `gorm:"type:varchar(30)"`
	Severity   string         `gorm:"type:varchar(20)"`
	Message    string         `gorm:"not null"`
	Resolved   bool           `gorm:"default:false"`
	ResolvedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (gormAlert) TableName() string { return "alerts" }

func toDomainAlert(m *gormAlert) *domainAlert.Alert {
	return &domainAlert.Alert{
		ID:         m.ID,
		LabID:      m.LabID,
		ComputerID: m.ComputerID,
		Type:       domainAlert.Type(m.Type),
		Severity:   domainAlert.Severity(m.Severity),
		Message:    m.Message,
		Resolved:   m.Resolved,
		ResolvedAt: m.ResolvedAt,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func fromDomainAlert(a *domainAlert.Alert) *gormAlert {
	return &gormAlert{
		ID:         a.ID,
		LabID:      a.LabID,
		ComputerID: a.ComputerID,
		Type:       string(a.Type),
		Severity:   string(a.Severity),
		Message:    a.Message,
		Resolved:   a.Resolved,
		ResolvedAt: a.ResolvedAt,
	}
}

type AlertRepository struct{}

func NewAlertRepository() *AlertRepository { return &AlertRepository{} }
func (r *AlertRepository) db() *gorm.DB    { return config.DB }

func (r *AlertRepository) Create(_ context.Context, a *domainAlert.Alert) error {
	return r.db().Create(fromDomainAlert(a)).Error
}

func (r *AlertRepository) BulkResolve(_ context.Context, ids []uuid.UUID) (int64, error) {
	now := time.Now().UTC()
	result := r.db().Model(&gormAlert{}).
		Where("id IN ? AND resolved = false", ids).
		Updates(map[string]interface{}{
			"resolved":    true,
			"resolved_at": now,
			"updated_at":  now,
		})
	return result.RowsAffected, result.Error
}

func (r *AlertRepository) FindByID(_ context.Context, id uuid.UUID) (*domainAlert.Alert, error) {
	var m gormAlert
	if err := r.db().First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainAlert.ErrNotFound
		}
		return nil, err
	}
	return toDomainAlert(&m), nil
}

func (r *AlertRepository) FindByLabID(_ context.Context, labID uuid.UUID, onlyOpen bool) ([]*domainAlert.Alert, error) {
	q := r.db().Where("lab_id = ?", labID)
	if onlyOpen {
		q = q.Where("resolved = false")
	}
	var models []gormAlert
	if err := q.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*domainAlert.Alert, 0, len(models))
	for i := range models {
		out = append(out, toDomainAlert(&models[i]))
	}
	return out, nil
}

func (r *AlertRepository) FindAll(_ context.Context, onlyOpen bool, page, limit int) ([]*domainAlert.Alert, int64, error) {
	q := r.db().Model(&gormAlert{})
	if onlyOpen {
		q = q.Where("resolved = false")
	}
	var total int64
	q.Count(&total)

	var models []gormAlert
	if err := q.Offset((page-1)*limit).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]*domainAlert.Alert, 0, len(models))
	for i := range models {
		out = append(out, toDomainAlert(&models[i]))
	}
	return out, total, nil
}

func (r *AlertRepository) Resolve(_ context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	return r.db().Model(&gormAlert{}).Where("id = ?", id).Updates(map[string]interface{}{
		"resolved":    true,
		"resolved_at": now,
		"updated_at":  now,
	}).Error
}

func (r *AlertRepository) Delete(_ context.Context, id uuid.UUID) error {
	return r.db().Delete(&gormAlert{}, "id = ?", id).Error
}

func (r *AlertRepository) CountUnresolved(_ context.Context) (int64, error) {
	var count int64
	return count, r.db().Model(&gormAlert{}).Where("resolved = false").Count(&count).Error
}

func (r *AlertRepository) CountBySeverity(_ context.Context) (map[domainAlert.Severity]int64, error) {
	type result struct {
		Severity string
		Count    int64
	}
	var results []result
	if err := r.db().Model(&gormAlert{}).Where("resolved = false").
		Select("severity, count(*) as count").Group("severity").Scan(&results).Error; err != nil {
		return nil, err
	}
	m := make(map[domainAlert.Severity]int64)
	for _, r := range results {
		m[domainAlert.Severity(r.Severity)] = r.Count
	}
	return m, nil
}
