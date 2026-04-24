package applab

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"endurance/internal/domain/computer"
	domainLab "endurance/internal/domain/lab"
	"endurance/pkg/apperrors"
)

type useCaseImpl struct {
	repo         domainLab.Repository
	computerRepo computer.Repository
}

func NewUseCase(repo domainLab.Repository, computerRepo computer.Repository) UseCase {
	return &useCaseImpl{repo: repo, computerRepo: computerRepo}
}

func (uc *useCaseImpl) GetAll(ctx context.Context, page, limit int) (*ListOutput, error) {
	labs, total, err := uc.repo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	labIDs := make([]uuid.UUID, 0, len(labs))
	for _, l := range labs {
		labIDs = append(labIDs, l.ID)
	}

	counts, err := uc.computerRepo.CountByLabIDs(ctx, labIDs)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	out := make([]*LabOutput, 0, len(labs))
	for _, l := range labs {
		o := toOutput(l)
		if cm, ok := counts[l.ID]; ok {
			for _, v := range cm {
				o.ComputerCount += v
			}
			o.OnlineCount = cm[computer.StatusOnline]
		}
		out = append(out, o)
	}

	return &ListOutput{Labs: out, Total: total, Page: page, Limit: limit}, nil
}

func (uc *useCaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*LabOutput, error) {
	l, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainLab.ErrNotFound) {
			return nil, apperrors.NotFound(apperrors.ErrNotFound)
		}
		return nil, apperrors.Internal(err)
	}

	return toOutput(l), nil
}

func (uc *useCaseImpl) Create(ctx context.Context, input CreateInput) (*LabOutput, error) {
	l := domainLab.New(input.Name, input.Location, input.Description, input.Capacity, input.ResponsibleID)

	if err := uc.repo.Create(ctx, l); err != nil {
		return nil, apperrors.Internal(err)
	}

	return toOutput(l), nil
}

func (uc *useCaseImpl) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*LabOutput, error) {
	l, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainLab.ErrNotFound) {
			return nil, apperrors.NotFound(apperrors.ErrNotFound)
		}
		return nil, apperrors.Internal(err)
	}

	l.Name = input.Name
	l.Location = input.Location
	l.Capacity = input.Capacity
	l.Description = input.Description
	l.ResponsibleID = input.ResponsibleID
	l.Status = input.Status

	if err := uc.repo.Update(ctx, l); err != nil {
		return nil, apperrors.Internal(err)
	}

	return toOutput(l), nil
}

func (uc *useCaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, domainLab.ErrNotFound) {
			return apperrors.NotFound(apperrors.ErrNotFound)
		}
		return apperrors.Internal(err)
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return apperrors.Internal(err)
	}

	return nil
}
