package security

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type hashService struct {
	cost int
}

// NewHashService retorna a implementação de bcrypt.
func NewHashService() *hashService {
	return &hashService{cost: bcrypt.DefaultCost}
}

// Hash gera o hash bcrypt da senha em texto puro.
func (s *hashService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", fmt.Errorf("gerar hash: %w", err)
	}
	return string(bytes), nil
}

// Compare verifica se a senha em texto puro corresponde ao hash armazenado.
func (s *hashService) Compare(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
