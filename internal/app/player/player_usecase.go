package player

import (
	"context"

	"soccer-manager-api/internal/domain"
	infraCache "soccer-manager-api/internal/infrastructure/cache"
	"soccer-manager-api/internal/ports/cache"
	"soccer-manager-api/internal/ports/repository"
)


type PlayerUseCase struct {
	playerRepo repository.PlayerRepository
	teamRepo   repository.TeamRepository
	cache      cache.Cache
	cacheHelper *infraCache.CacheHelper
}


func NewPlayerUseCase(
	playerRepo repository.PlayerRepository,
	teamRepo repository.TeamRepository,
	cache cache.Cache,
) *PlayerUseCase {
	return &PlayerUseCase{
		playerRepo:  playerRepo,
		teamRepo:    teamRepo,
		cache:       cache,
		cacheHelper: infraCache.NewCacheHelper(cache),
	}
}


type UpdatePlayerRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Country   string `json:"country" binding:"required"`
}


func (uc *PlayerUseCase) GetPlayer(ctx context.Context, playerID string) (*domain.Player, error) {
	player, err := uc.playerRepo.GetByID(ctx, playerID)
	if err != nil {
		return nil, err
	}
	return player, nil
}


func (uc *PlayerUseCase) UpdatePlayer(ctx context.Context, userID, playerID string, req UpdatePlayerRequest) (*domain.Player, error) {

	team, err := uc.teamRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}


	player, err := uc.playerRepo.GetByID(ctx, playerID)
	if err != nil {
		return nil, err
	}


	if !player.IsOwnedBy(team.ID) {
		return nil, domain.ErrPlayerNotOwned
	}


	player.FirstName = req.FirstName
	player.LastName = req.LastName
	player.Country = req.Country

	if err := uc.playerRepo.Update(ctx, player); err != nil {
		return nil, err
	}


	uc.cacheHelper.InvalidateTeamCache(ctx, team.ID.String())

	return player, nil
}

