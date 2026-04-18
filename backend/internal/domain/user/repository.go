package user

import (
	"context"

	"github.com/google/uuid"
)

// Repository é o PORT de saída (secondary port) — define o contrato
// que qualquer adaptador de persistência deve implementar.
// O domínio NÃO sabe como o dado é armazenado.
type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByCPF(ctx context.Context, cpf string) (*User, error)
	FindAll(ctx context.Context, page, limit int) ([]*User, int64, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByCPF(ctx context.Context, cpf string) (bool, error)
}
