package appuser

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	domainUser "endurance/internal/domain/user"
	"endurance/pkg/apperrors"
)

type HashService interface {
	Hash(password string) (string, error)
	Compare(hashed, plain string) error
}

type useCaseImpl struct {
	repo    domainUser.Repository
	hashSvc HashService
}

func NewUseCase(repo domainUser.Repository, hash HashService) UseCase {
	return &useCaseImpl{repo: repo, hashSvc: hash}
}

func (uc *useCaseImpl) GetAll(ctx context.Context, page, limit int) (*ListOutput, error) {
	users, total, err := uc.repo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	out := make([]*UserOutput, 0, len(users))
	for _, u := range users {
		out = append(out, toOutput(u))
	}

	return &ListOutput{Users: out, Total: total, Page: page, Limit: limit}, nil
}

func (uc *useCaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*UserOutput, error) {
	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainUser.ErrNotFound) {
			return nil, apperrors.NotFound(apperrors.ErrNotFound)
		}
		return nil, apperrors.Internal(err)
	}

	return toOutput(u), nil
}

func (uc *useCaseImpl) GetMe(ctx context.Context, id uuid.UUID) (*UserOutput, error) {
	return uc.GetByID(ctx, id)
}

func (uc *useCaseImpl) ChangeRole(ctx context.Context, id uuid.UUID, input ChangeRoleInput, requestorID uuid.UUID) (*UserOutput, error) {
	if id == requestorID {
		return nil, apperrors.BadRequest(apperrors.ErrSelfRoleChange)
	}
	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainUser.ErrNotFound) {
			return nil, apperrors.NotFound(apperrors.ErrNotFound)
		}
		return nil, apperrors.Internal(err)
	}
	u.Role = domainUser.Role(input.Role)
	u.UpdatedAt = time.Now().UTC()
	if err := uc.repo.UpdateRole(ctx, id, u.Role); err != nil {
		return nil, apperrors.Internal(err)
	}
	return toOutput(u), nil
}

func (uc *useCaseImpl) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*UserOutput, error) {
	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainUser.ErrNotFound) {
			return nil, apperrors.NotFound(apperrors.ErrNotFound)
		}
		return nil, apperrors.Internal(err)
	}

	u.UpdateName(input.Name)
	if input.Active != nil {
		if *input.Active {
			u.Activate()
		} else {
			u.Deactivate()
		}
	}

	if err := uc.repo.Update(ctx, u); err != nil {
		return nil, apperrors.Internal(err)
	}

	return toOutput(u), nil
}

func (uc *useCaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, domainUser.ErrNotFound) {
			return apperrors.NotFound(apperrors.ErrNotFound)
		}
		return apperrors.Internal(err)
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return apperrors.Internal(err)
	}

	return nil
}

func (uc *useCaseImpl) ChangePassword(ctx context.Context, id uuid.UUID, input ChangePasswordInput) error {
	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainUser.ErrNotFound) {
			return apperrors.NotFound(apperrors.ErrNotFound)
		}
		return apperrors.Internal(err)
	}

	if err := uc.hashSvc.Compare(u.Password, input.CurrentPassword); err != nil {
		return apperrors.Unauthorized(fmt.Errorf("senha atual incorreta"))
	}

	hashed, err := uc.hashSvc.Hash(input.NewPassword)
	if err != nil {
		return apperrors.Internal(err)
	}

	u.ChangePassword(hashed)

	if err := uc.repo.Update(ctx, u); err != nil {
		return apperrors.Internal(err)
	}

	return nil
}
