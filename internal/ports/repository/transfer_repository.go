package repository

import (
	"context"

	"soccer-manager-api/internal/domain"
)


type TransferRepository interface {
	CreateListing(ctx context.Context, listing *domain.TransferListing) error
	GetListingByID(ctx context.Context, id string) (*domain.TransferListing, error)
	GetListingByPlayerID(ctx context.Context, playerID string) (*domain.TransferListing, error)
	GetActiveListings(ctx context.Context, excludeTeamID string) ([]*domain.TransferListingWithPlayer, error)
	UpdateListing(ctx context.Context, listing *domain.TransferListing) error
	DeleteListing(ctx context.Context, id string) error

	CreateTransfer(ctx context.Context, transfer *domain.Transfer) error
	GetTransferByID(ctx context.Context, id string) (*domain.Transfer, error)
	GetTransfersByTeamID(ctx context.Context, teamID string) ([]*domain.Transfer, error)
}

