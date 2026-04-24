package computer

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, computer *Computer) error
	FindByID(ctx context.Context, id uuid.UUID) (*Computer, error)
	FindByLabID(ctx context.Context, labID uuid.UUID) ([]*Computer, error)
	FindAll(ctx context.Context, page, limit int, statusFilter string) ([]*Computer, int64, error)
	Update(ctx context.Context, computer *Computer) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error
	CountByStatus(ctx context.Context) (map[Status]int64, error)
	CountByLabIDs(ctx context.Context, labIDs []uuid.UUID) (map[uuid.UUID]map[Status]int64, error)
}
