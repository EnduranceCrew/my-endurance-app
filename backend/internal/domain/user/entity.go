// Package user contém a entidade de domínio User — sem dependências externas.
package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// User é a entidade central do domínio. Não importa nada de infraestrutura.
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CPF       string    `json:"cpf"`
	Password  string    `json:"-"`
	Role      Role      `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// New cria um novo usuário com valores-padrão.
func New(name, email, cpf, hashedPassword string, role Role) *User {
	return &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		CPF:       cpf,
		Password:  hashedPassword,
		Role:      role,
		Active:    true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC()
	}
}

// ── Comportamentos de domínio ────────────────────────────────────────────────

func (u *User) IsAdmin() bool { return u.Role == RoleAdmin }

func (u *User) Deactivate() {
	u.Active = false
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) Activate() {
	u.Active = true
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) UpdateName(name string) {
	u.Name = name
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) ChangePassword(hashedPassword string) {
	u.Password = hashedPassword
	u.UpdatedAt = time.Now().UTC()
}

// ErrNotFound é o erro de domínio para usuário não encontrado.
var ErrNotFound = errors.New("usuário não encontrado")
