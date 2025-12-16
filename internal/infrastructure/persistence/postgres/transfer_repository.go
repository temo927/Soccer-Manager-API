package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"soccer-manager-api/internal/domain"
	"soccer-manager-api/internal/ports/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type transferRepository struct {
	db *sqlx.DB
}


func NewTransferRepository(db *sqlx.DB) repository.TransferRepository {
	return &transferRepository{db: db}
}

func (r *transferRepository) CreateListing(ctx context.Context, listing *domain.TransferListing) error {
	query := `
		INSERT INTO transfer_listings (id, player_id, asking_price, status, listed_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, listing.ID, listing.PlayerID, listing.AskingPrice, listing.Status, listing.ListedAt)
	return err
}

func (r *transferRepository) GetListingByID(ctx context.Context, id string) (*domain.TransferListing, error) {
	var listing domain.TransferListing
	query := `SELECT id, player_id, asking_price, status, listed_at FROM transfer_listings WHERE id = $1`
	err := r.db.GetContext(ctx, &listing, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTransferListingNotFound
		}
		return nil, err
	}
	return &listing, nil
}

func (r *transferRepository) GetListingByPlayerID(ctx context.Context, playerID string) (*domain.TransferListing, error) {
	var listing domain.TransferListing
	query := `SELECT id, player_id, asking_price, status, listed_at FROM transfer_listings WHERE player_id = $1 AND status = 'active'`
	err := r.db.GetContext(ctx, &listing, query, playerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTransferListingNotFound
		}
		return nil, err
	}
	return &listing, nil
}

func (r *transferRepository) GetActiveListings(ctx context.Context, excludeTeamID string) ([]*domain.TransferListingWithPlayer, error) {
	type listingRow struct {
		domain.TransferListing
		PPlayerID     string  `db:"p_id"`
		PlayerTeamID  *string `db:"player_team_id"`
		PlayerFirstName string `db:"player_first_name"`
		PlayerLastName  string `db:"player_last_name"`
		PlayerCountry   string `db:"player_country"`
		PlayerAge        int    `db:"player_age"`
		PlayerPosition   string `db:"player_position"`
		PlayerMarketValue float64 `db:"player_market_value"`
		PlayerCreatedAt   time.Time `db:"player_created_at"`
		PlayerUpdatedAt   time.Time `db:"player_updated_at"`
	}

	var rows []listingRow
	query := `
		SELECT 
			tl.id, tl.player_id, tl.asking_price, tl.status, tl.listed_at,
			p.id as p_id, p.team_id as player_team_id, p.first_name as player_first_name,
			p.last_name as player_last_name, p.country as player_country, p.age as player_age,
			p.position as player_position, p.market_value as player_market_value,
			p.created_at as player_created_at, p.updated_at as player_updated_at
		FROM transfer_listings tl
		INNER JOIN players p ON tl.player_id = p.id
		WHERE tl.status = 'active' AND (p.team_id IS NULL OR p.team_id::text != $1)
		ORDER BY tl.listed_at DESC
	`
	err := r.db.SelectContext(ctx, &rows, query, excludeTeamID)
	if err != nil {
		return nil, err
	}

	listings := make([]*domain.TransferListingWithPlayer, 0, len(rows))
	for _, row := range rows {
		var teamID *uuid.UUID
		if row.PlayerTeamID != nil {
			id, _ := uuid.Parse(*row.PlayerTeamID)
			teamID = &id
		}
		
		listing := &domain.TransferListingWithPlayer{
			TransferListing: row.TransferListing,
			Player: domain.Player{
				ID:          uuid.MustParse(row.PPlayerID),
				TeamID:      teamID,
				FirstName:   row.PlayerFirstName,
				LastName:    row.PlayerLastName,
				Country:     row.PlayerCountry,
				Age:         row.PlayerAge,
				Position:    domain.Position(row.PlayerPosition),
				MarketValue: row.PlayerMarketValue,
				CreatedAt:   row.PlayerCreatedAt,
				UpdatedAt:   row.PlayerUpdatedAt,
			},
		}
		listings = append(listings, listing)
	}

	return listings, nil
}

func (r *transferRepository) UpdateListing(ctx context.Context, listing *domain.TransferListing) error {
	query := `
		UPDATE transfer_listings 
		SET asking_price = $1, status = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, listing.AskingPrice, listing.Status, listing.ID)
	return err
}

func (r *transferRepository) DeleteListing(ctx context.Context, id string) error {
	query := `DELETE FROM transfer_listings WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *transferRepository) CreateTransfer(ctx context.Context, transfer *domain.Transfer) error {
	query := `
		INSERT INTO transfers (id, player_id, seller_team_id, buyer_team_id, transfer_price, transferred_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		transfer.ID, transfer.PlayerID, transfer.SellerTeamID,
		transfer.BuyerTeamID, transfer.TransferPrice, transfer.TransferredAt)
	return err
}

func (r *transferRepository) GetTransferByID(ctx context.Context, id string) (*domain.Transfer, error) {
	var transfer domain.Transfer
	query := `
		SELECT id, player_id, seller_team_id, buyer_team_id, transfer_price, transferred_at 
		FROM transfers WHERE id = $1
	`
	err := r.db.GetContext(ctx, &transfer, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTransferNotFound
		}
		return nil, err
	}
	return &transfer, nil
}

func (r *transferRepository) GetTransfersByTeamID(ctx context.Context, teamID string) ([]*domain.Transfer, error) {
	var transfers []*domain.Transfer
	query := `
		SELECT id, player_id, seller_team_id, buyer_team_id, transfer_price, transferred_at 
		FROM transfers 
		WHERE seller_team_id = $1 OR buyer_team_id = $1
		ORDER BY transferred_at DESC
	`
	err := r.db.SelectContext(ctx, &transfers, query, teamID)
	return transfers, err
}

