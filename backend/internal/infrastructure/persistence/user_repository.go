package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"endurance/config"
	domainUser "endurance/internal/domain/user"
)

// gormUser é o modelo GORM (adaptador) — isola GORM do domínio.
type gormUser struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name      string         `gorm:"not null"`
	Email     string         `gorm:"uniqueIndex;not null"`
	CPF       string         `gorm:"uniqueIndex;not null"`
	Password  string         `gorm:"not null"`
	Role      string         `gorm:"type:varchar(10);default:'user'"`
	Active    bool           `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (gormUser) TableName() string { return "users" }

func toDomainUser(m *gormUser) *domainUser.User {
	return &domainUser.User{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		CPF:       m.CPF,
		Password:  m.Password,
		Role:      domainUser.Role(m.Role),
		Active:    m.Active,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func fromDomainUser(u *domainUser.User) *gormUser {
	return &gormUser{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		CPF:      u.CPF,
		Password: u.Password,
		Role:     string(u.Role),
		Active:   u.Active,
	}
}

// ── Repositório ──────────────────────────────────────────────────────────────

type UserRepository struct{}

func NewUserRepository() *UserRepository { return &UserRepository{} }

func (r *UserRepository) db(ctx context.Context) *gorm.DB { return config.DB.WithContext(ctx) }

func (r *UserRepository) Create(ctx context.Context, u *domainUser.User) error {
	return r.db(ctx).Create(fromDomainUser(u)).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainUser.User, error) {
	var m gormUser
	if err := r.db(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainUser.ErrNotFound
		}
		return nil, err
	}
	return toDomainUser(&m), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domainUser.User, error) {
	var m gormUser
	if err := r.db(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainUser.ErrNotFound
		}
		return nil, err
	}
	return toDomainUser(&m), nil
}

func (r *UserRepository) FindByCPF(ctx context.Context, cpf string) (*domainUser.User, error) {
	var m gormUser
	if err := r.db(ctx).Where("cpf = ?", cpf).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainUser.ErrNotFound
		}
		return nil, err
	}
	return toDomainUser(&m), nil
}

func (r *UserRepository) FindAll(ctx context.Context, page, limit int) ([]*domainUser.User, int64, error) {
	var models []gormUser
	var total int64

	offset := (page - 1) * limit
	r.db(ctx).Model(&gormUser{}).Count(&total)

	if err := r.db(ctx).Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	users := make([]*domainUser.User, 0, len(models))
	for i := range models {
		users = append(users, toDomainUser(&models[i]))
	}
	return users, total, nil
}

func (r *UserRepository) Update(ctx context.Context, u *domainUser.User) error {
	return r.db(ctx).Model(&gormUser{}).Where("id = ?", u.ID).Updates(map[string]interface{}{
		"name":       u.Name,
		"email":      u.Email,
		"password":   u.Password,
		"active":     u.Active,
		"updated_at": u.UpdatedAt,
	}).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db(ctx).Delete(&gormUser{}, "id = ?", id).Error
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db(ctx).Model(&gormUser{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) ExistsByCPF(ctx context.Context, cpf string) (bool, error) {
	var count int64
	err := r.db(ctx).Model(&gormUser{}).Where("cpf = ?", cpf).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) UpdateRole(ctx context.Context, id uuid.UUID, role domainUser.Role) error {
	return r.db(ctx).Model(&gormUser{}).Where("id = ?", id).Updates(map[string]interface{}{
		"role":       string(role),
		"updated_at": time.Now().UTC(),
	}).Error
}
