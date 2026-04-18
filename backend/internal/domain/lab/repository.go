package lab

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, lab *Lab) error
	FindByID(ctx context.Context, id uuid.UUID) (*Lab, error)
	FindAll(ctx context.Context, page, limit int) ([]*Lab, int64, error)
	Update(ctx context.Context, lab *Lab) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByStatus(ctx context.Context) (map[Status]int64, error)
}
