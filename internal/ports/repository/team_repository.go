package repository

import (
	"context"

	"soccer-manager-api/internal/domain"
)


type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error
	GetByID(ctx context.Context, id string) (*domain.Team, error)
	GetByUserID(ctx context.Context, userID string) (*domain.Team, error)
	Update(ctx context.Context, team *domain.Team) error
	GetTotalValue(ctx context.Context, teamID string) (float64, error)
	GetPlayerCount(ctx context.Context, teamID string) (int, error)
}

