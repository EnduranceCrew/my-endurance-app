package appalert

import (
	"context"
	"errors"

	"github.com/google/uuid"
	domainAlert "endurance/internal/domain/alert"
	"endurance/pkg/apperrors"
)

type useCaseImpl struct {
	repo domainAlert.Repository
}

func NewUseCase(repo domainAlert.Repository) UseCase {
	return &useCaseImpl{repo: repo}
}

func (uc *useCaseImpl) GetAll(ctx context.Context, onlyOpen bool, page, limit int) (*ListOutput, error) {
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

func (uc *useCaseImpl) GetByLabID(ctx context.Context, labID uuid.UUID, onlyOpen bool) ([]*AlertOutput, error) {
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

func (uc *useCaseImpl) Create(ctx context.Context, input CreateInput) (*AlertOutput, error) {
	a := domainAlert.New(input.LabID, input.ComputerID, input.Type, input.Severity, input.Message)
	if err := uc.repo.Create(ctx, a); err != nil {
		return nil, apperrors.Internal(err)
	}
	return toOutput(a), nil
}

func (uc *useCaseImpl) Resolve(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, domainAlert.ErrNotFound) {
			return apperrors.NotFound(apperrors.ErrNotFound)
		}
		return apperrors.Internal(err)
	}
	return uc.repo.Resolve(ctx, id)
}

func (uc *useCaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, domainAlert.ErrNotFound) {
			return apperrors.NotFound(apperrors.ErrNotFound)
		}
		return apperrors.Internal(err)
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *useCaseImpl) BulkResolve(ctx context.Context, input BulkResolveInput) (int64, error) {
	count, err := uc.repo.BulkResolve(ctx, input.IDs)
	if err != nil {
		return 0, apperrors.Internal(err)
	}
	return count, nil
}
