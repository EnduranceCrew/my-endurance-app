package auth

import (
	"context"
	"fmt"
	"strings"

	"endurance/config"
	"endurance/internal/domain/user"
	"endurance/pkg/apperrors"
	"endurance/pkg/validator"
)

// ── Ports que o caso de uso precisa (definidos aqui, implementados na infra) ─

type HashService interface {
	Hash(password string) (string, error)
	Compare(hashed, plain string) error
}

type JWTService interface {
	Generate(userID, role string) (string, error)
}

// ── Implementação do caso de uso ─────────────────────────────────────────────

type useCaseImpl struct {
	userRepo   user.Repository
	hashSvc    HashService
	jwtSvc     JWTService
}

// NewUseCase injeta as dependências via interfaces (inversão de dependência).
func NewUseCase(repo user.Repository, hash HashService, jwt JWTService) UseCase {
	return &useCaseImpl{userRepo: repo, hashSvc: hash, jwtSvc: jwt}
}

func (uc *useCaseImpl) Login(input LoginInput) (*TokenOutput, error) {
	ctx := context.Background()

	u, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, apperrors.Unauthorized(apperrors.ErrInvalidCredentials)
	}

	if !u.Active {
		return nil, apperrors.Unauthorized(apperrors.ErrInactiveUser)
	}

	if err := uc.hashSvc.Compare(u.Password, input.Password); err != nil {
		return nil, apperrors.Unauthorized(apperrors.ErrInvalidCredentials)
	}

	token, err := uc.jwtSvc.Generate(u.ID.String(), string(u.Role))
	if err != nil {
		return nil, apperrors.Internal(fmt.Errorf("gerar token: %w", err))
	}

	return &TokenOutput{
		AccessToken: token,
		TokenType:   "Bearer",
		UserID:      u.ID,
		Name:        u.Name,
		Email:       u.Email,
		Role:        string(u.Role),
		ExpiresIn:   config.App.JWTExpirationHours * 3600,
	}, nil
}

func (uc *useCaseImpl) Register(input RegisterInput) (*TokenOutput, error) {
	ctx := context.Background()

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	// Validações de domínio
	if !validator.IsValidEmail(input.Email) {
		return nil, apperrors.BadRequest(apperrors.ErrInvalidEmail)
	}
	if !validator.IsValidCPF(input.CPF) {
		return nil, apperrors.BadRequest(apperrors.ErrInvalidCPF)
	}
	if !validator.IsStrongPassword(input.Password) {
		return nil, apperrors.BadRequest(apperrors.ErrWeakPassword)
	}

	emailExists, _ := uc.userRepo.ExistsByEmail(ctx, input.Email)
	if emailExists {
		return nil, apperrors.Conflict(apperrors.ErrAlreadyExists)
	}

	cpfExists, _ := uc.userRepo.ExistsByCPF(ctx, validator.SanitizeCPF(input.CPF))
	if cpfExists {
		return nil, apperrors.Conflict(apperrors.ErrAlreadyExists)
	}

	hashed, err := uc.hashSvc.Hash(input.Password)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	newUser := user.New(
		input.Name,
		input.Email,
		validator.SanitizeCPF(input.CPF),
		hashed,
		user.RoleUser,
	)

	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, apperrors.Internal(err)
	}

	token, err := uc.jwtSvc.Generate(newUser.ID.String(), string(newUser.Role))
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return &TokenOutput{
		AccessToken: token,
		TokenType:   "Bearer",
		UserID:      newUser.ID,
		Name:        newUser.Name,
		Email:       newUser.Email,
		Role:        string(newUser.Role),
		ExpiresIn:   config.App.JWTExpirationHours * 3600,
	}, nil
}
