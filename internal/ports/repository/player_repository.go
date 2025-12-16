package repository

import (
	"context"

	"soccer-manager-api/internal/domain"
)


type PlayerRepository interface {
	Create(ctx context.Context, player *domain.Player) error
	CreateBatch(ctx context.Context, players []*domain.Player) error
	GetByID(ctx context.Context, id string) (*domain.Player, error)
	GetByTeamID(ctx context.Context, teamID string) ([]*domain.Player, error)
	Update(ctx context.Context, player *domain.Player) error
	Delete(ctx context.Context, id string) error
	GetByTeamIDAndPosition(ctx context.Context, teamID string, position domain.Position) ([]*domain.Player, error)
}

