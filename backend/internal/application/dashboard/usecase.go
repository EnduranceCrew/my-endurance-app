package dashboard

import (
	"context"

	"endurance/internal/domain/alert"
	"endurance/internal/domain/computer"
	"endurance/internal/domain/lab"
	"endurance/internal/domain/user"
	"endurance/pkg/apperrors"
)

// ── Output DTO ───────────────────────────────────────────────────────────────

type StatsOutput struct {
	TotalLabs      int64 `json:"total_labs"`
	ActiveLabs     int64 `json:"active_labs"`
	MaintenanceLabs int64 `json:"maintenance_labs"`
	TotalComputers int64 `json:"total_computers"`
	OnlineComputers int64 `json:"online_computers"`
	OfflineComputers int64 `json:"offline_computers"`
	ErrorComputers int64 `json:"error_computers"`
	TotalUsers     int64 `json:"total_users"`
	OpenAlerts     int64 `json:"open_alerts"`
	CriticalAlerts int64 `json:"critical_alerts"`
}

// ── Port de entrada ──────────────────────────────────────────────────────────

type UseCase interface {
	GetStats(ctx context.Context) (*StatsOutput, error)
}

// ── Implementação ────────────────────────────────────────────────────────────

type useCaseImpl struct {
	labRepo      lab.Repository
	computerRepo computer.Repository
	userRepo     user.Repository
	alertRepo    alert.Repository
}

func NewUseCase(
	labRepo lab.Repository,
	computerRepo computer.Repository,
	userRepo user.Repository,
	alertRepo alert.Repository,
) UseCase {
	return &useCaseImpl{
		labRepo:      labRepo,
		computerRepo: computerRepo,
		userRepo:     userRepo,
		alertRepo:    alertRepo,
	}
}

func (uc *useCaseImpl) GetStats(ctx context.Context) (*StatsOutput, error) {
	labStatus, err := uc.labRepo.CountByStatus(ctx)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	computerStatus, err := uc.computerRepo.CountByStatus(ctx)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	_, totalUsers, err := uc.userRepo.FindAll(ctx, 1, 1)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	openAlerts, err := uc.alertRepo.CountUnresolved(ctx)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	alertBySeverity, err := uc.alertRepo.CountBySeverity(ctx)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var totalLabs int64
	for _, v := range labStatus {
		totalLabs += v
	}

	var totalComputers int64
	for _, v := range computerStatus {
		totalComputers += v
	}

	return &StatsOutput{
		TotalLabs:        totalLabs,
		ActiveLabs:       labStatus[lab.StatusActive],
		MaintenanceLabs:  labStatus[lab.StatusMaintenance],
		TotalComputers:   totalComputers,
		OnlineComputers:  computerStatus[computer.StatusOnline],
		OfflineComputers: computerStatus[computer.StatusOffline],
		ErrorComputers:   computerStatus[computer.StatusError],
		TotalUsers:       totalUsers,
		OpenAlerts:       openAlerts,
		CriticalAlerts:   alertBySeverity[alert.SeverityCritical],
	}, nil
}
