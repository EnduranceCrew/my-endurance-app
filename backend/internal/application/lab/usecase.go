package applab

import (
	"context"

	"github.com/google/uuid"
	"endurance/internal/domain/computer"
	"endurance/internal/domain/lab"
	"endurance/pkg/apperrors"
)

type useCaseImpl struct {
	repo         lab.Repository
	computerRepo computer.Repository
}

func NewUseCase(repo lab.Repository, computerRepo computer.Repository) UseCase {
	return &useCaseImpl{repo: repo, computerRepo: computerRepo}
}

func (uc *useCaseImpl) GetAll(page, limit int) (*ListOutput, error) {
	ctx := context.Background()

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

func (uc *useCaseImpl) GetByID(id uuid.UUID) (*LabOutput, error) {
	ctx := context.Background()

	l, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NotFound(apperrors.ErrNotFound)
	}

	return toOutput(l), nil
}

func (uc *useCaseImpl) Create(input CreateInput) (*LabOutput, error) {
	ctx := context.Background()

	l := lab.New(input.Name, input.Location, input.Description, input.Capacity, input.ResponsibleID)

	if err := uc.repo.Create(ctx, l); err != nil {
		return nil, apperrors.Internal(err)
	}

	return toOutput(l), nil
}

func (uc *useCaseImpl) Update(id uuid.UUID, input UpdateInput) (*LabOutput, error) {
	ctx := context.Background()

	l, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NotFound(apperrors.ErrNotFound)
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
