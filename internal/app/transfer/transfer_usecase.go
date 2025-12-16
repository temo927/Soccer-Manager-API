package transfer

import (
	"context"

	"soccer-manager-api/internal/domain"
	infraCache "soccer-manager-api/internal/infrastructure/cache"
	"soccer-manager-api/internal/ports/cache"
	"soccer-manager-api/internal/ports/repository"
)


type TransferUseCase struct {
	transferRepo repository.TransferRepository
	teamRepo     repository.TeamRepository
	playerRepo   repository.PlayerRepository
	cache        cache.Cache
	cacheHelper  *infraCache.CacheHelper
}


func NewTransferUseCase(
	transferRepo repository.TransferRepository,
	teamRepo repository.TeamRepository,
	playerRepo repository.PlayerRepository,
	cache cache.Cache,
) *TransferUseCase {
	return &TransferUseCase{
		transferRepo: transferRepo,
		teamRepo:     teamRepo,
		playerRepo:   playerRepo,
		cache:        cache,
		cacheHelper:  infraCache.NewCacheHelper(cache),
	}
}


type ListPlayerRequest struct {
	AskingPrice float64 `json:"asking_price" binding:"required,gt=0"`
}


type BuyPlayerRequest struct {
	ListingID string `json:"listing_id" binding:"required"`
}


func (uc *TransferUseCase) ListPlayer(ctx context.Context, userID, playerID string, req ListPlayerRequest) (*domain.TransferListing, error) {

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


	existingListing, _ := uc.transferRepo.GetListingByPlayerID(ctx, playerID)
	if existingListing != nil && existingListing.IsActive() {
		return nil, domain.ErrPlayerAlreadyListed
	}


	if req.AskingPrice <= 0 {
		return nil, domain.ErrInvalidAskingPrice
	}


	listing := domain.NewTransferListing(player.ID, req.AskingPrice)
	if err := uc.transferRepo.CreateListing(ctx, listing); err != nil {
		return nil, err
	}


	uc.cacheHelper.InvalidateTransferListCache(ctx)

	return listing, nil
}


func (uc *TransferUseCase) RemoveFromTransferList(ctx context.Context, userID, playerID string) error {

	team, err := uc.teamRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}


	player, err := uc.playerRepo.GetByID(ctx, playerID)
	if err != nil {
		return err
	}


	if !player.IsOwnedBy(team.ID) {
		return domain.ErrPlayerNotOwned
	}


	listing, err := uc.transferRepo.GetListingByPlayerID(ctx, playerID)
	if err != nil {
		return domain.ErrPlayerNotOnTransferList
	}


	listing.Cancel()
	if err := uc.transferRepo.UpdateListing(ctx, listing); err != nil {
		return err
	}


	uc.cacheHelper.InvalidateTransferListCache(ctx)

	return nil
}


func (uc *TransferUseCase) GetTransferList(ctx context.Context, userID string) ([]*domain.TransferListingWithPlayer, error) {

	team, err := uc.teamRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}


	cacheKey := infraCache.CacheKey("transfer_list", "all")
	var listings []*domain.TransferListingWithPlayer
	if err := uc.cacheHelper.Get(ctx, cacheKey, &listings); err == nil {

		return uc.filterOwnPlayers(listings, team.ID.String()), nil
	}


	listings, err = uc.transferRepo.GetActiveListings(ctx, team.ID.String())
	if err != nil {
		return nil, err
	}


	uc.cacheHelper.Set(ctx, cacheKey, listings, 60)

	return listings, nil
}


func (uc *TransferUseCase) BuyPlayer(ctx context.Context, userID, listingID string) (*domain.Transfer, error) {

	buyerTeam, err := uc.teamRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}


	listing, err := uc.transferRepo.GetListingByID(ctx, listingID)
	if err != nil {
		return nil, err
	}

	if !listing.IsActive() {
		return nil, domain.ErrTransferListingNotFound
	}


	player, err := uc.playerRepo.GetByID(ctx, listing.PlayerID.String())
	if err != nil {
		return nil, err
	}


	if player.TeamID != nil && player.IsOwnedBy(buyerTeam.ID) {
		return nil, domain.ErrCannotBuyOwnPlayer
	}


	if !buyerTeam.CanAfford(listing.AskingPrice) {
		return nil, domain.ErrInsufficientBudget
	}


	playerCount, err := uc.teamRepo.GetPlayerCount(ctx, buyerTeam.ID.String())
	if err != nil {
		return nil, err
	}
	if playerCount >= domain.MaxPlayers {
		return nil, domain.ErrTeamFull
	}


	var sellerTeam *domain.Team
	if player.TeamID != nil {
		sellerTeam, err = uc.teamRepo.GetByID(ctx, player.TeamID.String())
		if err != nil {
			return nil, err
		}
	} else {

		return nil, domain.ErrTeamNotFound
	}



	player.Transfer(buyerTeam.ID)
	if err := uc.playerRepo.Update(ctx, player); err != nil {
		return nil, err
	}


	buyerTeam.DeductBudget(listing.AskingPrice)
	sellerTeam.AddBudget(listing.AskingPrice)
	if err := uc.teamRepo.Update(ctx, buyerTeam); err != nil {
		return nil, err
	}
	if err := uc.teamRepo.Update(ctx, sellerTeam); err != nil {
		return nil, err
	}


	listing.MarkAsSold()
	if err := uc.transferRepo.UpdateListing(ctx, listing); err != nil {
		return nil, err
	}


	transfer := domain.NewTransfer(
		player.ID,
		sellerTeam.ID,
		buyerTeam.ID,
		listing.AskingPrice,
	)
	if err := uc.transferRepo.CreateTransfer(ctx, transfer); err != nil {
		return nil, err
	}


	uc.cacheHelper.InvalidateTeamCache(ctx, buyerTeam.ID.String())
	uc.cacheHelper.InvalidateTeamCache(ctx, sellerTeam.ID.String())
	uc.cacheHelper.InvalidateTransferListCache(ctx)

	return transfer, nil
}


func (uc *TransferUseCase) filterOwnPlayers(listings []*domain.TransferListingWithPlayer, teamID string) []*domain.TransferListingWithPlayer {
	filtered := make([]*domain.TransferListingWithPlayer, 0)
	for _, listing := range listings {
		if listing.Player.TeamID == nil || listing.Player.TeamID.String() != teamID {
			filtered = append(filtered, listing)
		}
	}
	return filtered
}

