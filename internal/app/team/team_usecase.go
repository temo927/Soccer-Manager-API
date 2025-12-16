package team

import (
	"context"

	"soccer-manager-api/internal/domain"
	infraCache "soccer-manager-api/internal/infrastructure/cache"
	"soccer-manager-api/internal/ports/cache"
	"soccer-manager-api/internal/ports/repository"
)


type TeamUseCase struct {
	teamRepo   repository.TeamRepository
	playerRepo repository.PlayerRepository
	cache      cache.Cache
	cacheHelper *infraCache.CacheHelper
}


func NewTeamUseCase(
	teamRepo repository.TeamRepository,
	playerRepo repository.PlayerRepository,
	cache cache.Cache,
) *TeamUseCase {
	return &TeamUseCase{
		teamRepo:    teamRepo,
		playerRepo:  playerRepo,
		cache:       cache,
		cacheHelper: infraCache.NewCacheHelper(cache),
	}
}


type UpdateTeamRequest struct {
	Name    string `json:"name" binding:"required"`
	Country string `json:"country" binding:"required"`
}


func (uc *TeamUseCase) GetTeam(ctx context.Context, userID string) (*domain.TeamWithValue, error) {

	cacheKey := infraCache.CacheKey("team", userID)
	var teamWithValue domain.TeamWithValue
	if err := uc.cacheHelper.Get(ctx, cacheKey, &teamWithValue); err == nil {
		return &teamWithValue, nil
	}


	team, err := uc.teamRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}


	totalValue, err := uc.teamRepo.GetTotalValue(ctx, team.ID.String())
	if err != nil {
		return nil, err
	}

	teamWithValue = domain.TeamWithValue{
		Team:      *team,
		TotalValue: totalValue,
	}


	uc.cacheHelper.Set(ctx, cacheKey, teamWithValue, 300)

	return &teamWithValue, nil
}


func (uc *TeamUseCase) UpdateTeam(ctx context.Context, userID string, req UpdateTeamRequest) (*domain.Team, error) {
	team, err := uc.teamRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	team.Name = req.Name
	team.Country = req.Country

	if err := uc.teamRepo.Update(ctx, team); err != nil {
		return nil, err
	}


	uc.cacheHelper.InvalidateTeamCache(ctx, team.ID.String())

	return team, nil
}


func (uc *TeamUseCase) GetTeamPlayers(ctx context.Context, userID string) ([]*domain.Player, error) {

	team, err := uc.teamRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}


	cacheKey := infraCache.CacheKey("team:players", team.ID.String())
	var players []*domain.Player
	if err := uc.cacheHelper.Get(ctx, cacheKey, &players); err == nil {
		return players, nil
	}


	players, err = uc.playerRepo.GetByTeamID(ctx, team.ID.String())
	if err != nil {
		return nil, err
	}


	uc.cacheHelper.Set(ctx, cacheKey, players, 300)

	return players, nil
}

