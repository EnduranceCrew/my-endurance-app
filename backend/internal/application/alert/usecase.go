package appalert

import (
	"context"

	"github.com/google/uuid"
	"endurance/internal/domain/alert"
	"endurance/pkg/apperrors"
)

type useCaseImpl struct {
	repo alert.Repository
}

func NewUseCase(repo alert.Repository) UseCase {
	return &useCaseImpl{repo: repo}
}

func (uc *useCaseImpl) GetAll(onlyOpen bool, page, limit int) (*ListOutput, error) {
	ctx := context.Background()
	items, total, err := uc.repo.FindAll(ctx, onlyOpen, page, limit)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	out := make([]*AlertOutput, 0, len(items))
	for _, a := range items {
		out = append(out, toOutput(a))
	}
	return &ListOutput{Alerts: out, Total: total, Page: page, Limit: limit}, nil
}

func (uc *useCaseImpl) GetByLabID(labID uuid.UUID, onlyOpen bool) ([]*AlertOutput, error) {
	ctx := context.Background()
	items, err := uc.repo.FindByLabID(ctx, labID, onlyOpen)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	out := make([]*AlertOutput, 0, len(items))
	for _, a := range items {
		out = append(out, toOutput(a))
	}
	return out, nil
}

func (uc *useCaseImpl) Create(input CreateInput) (*AlertOutput, error) {
	ctx := context.Background()
	a := alert.New(input.LabID, input.ComputerID, input.Type, input.Severity, input.Message)
	if err := uc.repo.Create(ctx, a); err != nil {
		return nil, apperrors.Internal(err)
	}
	return toOutput(a), nil
}

func (uc *useCaseImpl) Resolve(id uuid.UUID) error {
	ctx := context.Background()
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return apperrors.NotFound(apperrors.ErrNotFound)
	}
	return uc.repo.Resolve(ctx, id)
}

func (uc *useCaseImpl) Delete(id uuid.UUID) error {
	ctx := context.Background()
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return apperrors.NotFound(apperrors.ErrNotFound)
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *useCaseImpl) BulkResolve(input BulkResolveInput) (int64, error) {
	ctx := context.Background()
	count, err := uc.repo.BulkResolve(ctx, input.IDs)
	if err != nil {
		return 0, apperrors.Internal(err)
	}
	return count, nil
}
