package appuser

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"endurance/internal/domain/user"
	"endurance/pkg/apperrors"
)

type HashService interface {
	Hash(password string) (string, error)
	Compare(hashed, plain string) error
}

type useCaseImpl struct {
	repo    user.Repository
	hashSvc HashService
}

func NewUseCase(repo user.Repository, hash HashService) UseCase {
	return &useCaseImpl{repo: repo, hashSvc: hash}
}

func (uc *useCaseImpl) GetAll(page, limit int) (*ListOutput, error) {
	ctx := context.Background()

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

func (uc *useCaseImpl) GetByID(id uuid.UUID) (*UserOutput, error) {
	ctx := context.Background()

	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NotFound(apperrors.ErrNotFound)
	}

	return toOutput(u), nil
}

func (uc *useCaseImpl) Update(id uuid.UUID, input UpdateInput) (*UserOutput, error) {
	ctx := context.Background()

	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NotFound(apperrors.ErrNotFound)
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

func (uc *useCaseImpl) Delete(id uuid.UUID) error {
	ctx := context.Background()

	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return apperrors.NotFound(apperrors.ErrNotFound)
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return apperrors.Internal(err)
	}

	return nil
}

func (uc *useCaseImpl) ChangePassword(id uuid.UUID, input ChangePasswordInput) error {
	ctx := context.Background()

	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.NotFound(apperrors.ErrNotFound)
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
