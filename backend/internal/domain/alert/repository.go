package alert

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, alert *Alert) error
	FindByID(ctx context.Context, id uuid.UUID) (*Alert, error)
	FindByLabID(ctx context.Context, labID uuid.UUID, onlyOpen bool) ([]*Alert, error)
	FindAll(ctx context.Context, onlyOpen bool, page, limit int) ([]*Alert, int64, error)
	Resolve(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountUnresolved(ctx context.Context) (int64, error)
	CountBySeverity(ctx context.Context) (map[Severity]int64, error)
}
