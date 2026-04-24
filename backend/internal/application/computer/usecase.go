package appcomputer

import (
	"context"
	"errors"

	"github.com/google/uuid"
	domainComputer "endurance/internal/domain/computer"
	"endurance/pkg/apperrors"
)

type useCaseImpl struct {
	repo domainComputer.Repository
}

func NewUseCase(repo domainComputer.Repository) UseCase {
	return &useCaseImpl{repo: repo}
}

func (uc *useCaseImpl) GetAll(ctx context.Context, page, limit int, statusFilter string) (*ListOutput, error) {
	items, total, err := uc.repo.FindAll(ctx, page, limit, statusFilter)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	out := make([]*ComputerOutput, 0, len(items))
	for _, c := range items {
		out = append(out, toOutput(c))
	}
	return &ListOutput{Computers: out, Total: total, Page: page, Limit: limit}, nil
}

func (uc *useCaseImpl) GetByLabID(ctx context.Context, labID uuid.UUID) ([]*ComputerOutput, error) {
	items, err := uc.repo.FindByLabID(ctx, labID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	out := make([]*ComputerOutput, 0, len(items))
	for _, c := range items {
		out = append(out, toOutput(c))
	}
	return out, nil
}

func (uc *useCaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*ComputerOutput, error) {
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainComputer.ErrNotFound) {
			return nil, apperrors.NotFound(apperrors.ErrNotFound)
		}
		return nil, apperrors.Internal(err)
	}
	return toOutput(c), nil
}

func (uc *useCaseImpl) Create(ctx context.Context, input CreateInput) (*ComputerOutput, error) {
	c := domainComputer.New(input.LabID, input.Hostname, input.IPAddress, input.MACAddress,
		input.OS, input.CPU, input.RAM, input.Storage)
	if err := uc.repo.Create(ctx, c); err != nil {
		return nil, apperrors.Internal(err)
	}
	return toOutput(c), nil
}

func (uc *useCaseImpl) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*ComputerOutput, error) {
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domainComputer.ErrNotFound) {
			return nil, apperrors.NotFound(apperrors.ErrNotFound)
		}
		return nil, apperrors.Internal(err)
	}
	c.Hostname = input.Hostname
	c.IPAddress = input.IPAddress
	c.MACAddress = input.MACAddress
	c.OS = input.OS
	c.CPU = input.CPU
	c.RAM = input.RAM
	c.Storage = input.Storage
	c.Status = input.Status
	if err := uc.repo.Update(ctx, c); err != nil {
		return nil, apperrors.Internal(err)
	}
	return toOutput(c), nil
}

func (uc *useCaseImpl) UpdateStatus(ctx context.Context, id uuid.UUID, input UpdateStatusInput) error {
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, domainComputer.ErrNotFound) {
			return apperrors.NotFound(apperrors.ErrNotFound)
		}
		return apperrors.Internal(err)
	}
	return uc.repo.UpdateStatus(ctx, id, input.Status)
}

func (uc *useCaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, domainComputer.ErrNotFound) {
			return apperrors.NotFound(apperrors.ErrNotFound)
		}
		return apperrors.Internal(err)
	}
	return uc.repo.Delete(ctx, id)
}
