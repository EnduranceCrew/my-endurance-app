package appcomputer

import (
	"context"

	"github.com/google/uuid"
	"endurance/internal/domain/computer"
	"endurance/pkg/apperrors"
)

type useCaseImpl struct {
	repo computer.Repository
}

func NewUseCase(repo computer.Repository) UseCase {
	return &useCaseImpl{repo: repo}
}

func (uc *useCaseImpl) GetAll(page, limit int) (*ListOutput, error) {
	ctx := context.Background()
	items, total, err := uc.repo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	out := make([]*ComputerOutput, 0, len(items))
	for _, c := range items {
		out = append(out, toOutput(c))
	}
	return &ListOutput{Computers: out, Total: total, Page: page, Limit: limit}, nil
}

func (uc *useCaseImpl) GetByLabID(labID uuid.UUID) ([]*ComputerOutput, error) {
	ctx := context.Background()
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

func (uc *useCaseImpl) GetByID(id uuid.UUID) (*ComputerOutput, error) {
	ctx := context.Background()
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NotFound(apperrors.ErrNotFound)
	}
	return toOutput(c), nil
}

func (uc *useCaseImpl) Create(input CreateInput) (*ComputerOutput, error) {
	ctx := context.Background()
	c := computer.New(input.LabID, input.Hostname, input.IPAddress, input.MACAddress,
		input.OS, input.CPU, input.RAM, input.Storage)
	if err := uc.repo.Create(ctx, c); err != nil {
		return nil, apperrors.Internal(err)
	}
	return toOutput(c), nil
}

func (uc *useCaseImpl) Update(id uuid.UUID, input UpdateInput) (*ComputerOutput, error) {
	ctx := context.Background()
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NotFound(apperrors.ErrNotFound)
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

func (uc *useCaseImpl) UpdateStatus(id uuid.UUID, input UpdateStatusInput) error {
	ctx := context.Background()
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return apperrors.NotFound(apperrors.ErrNotFound)
	}
	return uc.repo.UpdateStatus(ctx, id, input.Status)
}

func (uc *useCaseImpl) Delete(id uuid.UUID) error {
	ctx := context.Background()
	if _, err := uc.repo.FindByID(ctx, id); err != nil {
		return apperrors.NotFound(apperrors.ErrNotFound)
	}
	return uc.repo.Delete(ctx, id)
}
