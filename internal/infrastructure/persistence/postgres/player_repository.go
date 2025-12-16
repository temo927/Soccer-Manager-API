package postgres

import (
	"context"
	"database/sql"
	"errors"

	"soccer-manager-api/internal/domain"
	"soccer-manager-api/internal/ports/repository"

	"github.com/jmoiron/sqlx"
)

type playerRepository struct {
	db *sqlx.DB
}


func NewPlayerRepository(db *sqlx.DB) repository.PlayerRepository {
	return &playerRepository{db: db}
}

func (r *playerRepository) Create(ctx context.Context, player *domain.Player) error {
	query := `
		INSERT INTO players (id, team_id, first_name, last_name, country, age, position, market_value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		player.ID, player.TeamID, player.FirstName, player.LastName,
		player.Country, player.Age, player.Position, player.MarketValue,
		player.CreatedAt, player.UpdatedAt)
	return err
}

func (r *playerRepository) CreateBatch(ctx context.Context, players []*domain.Player) error {
	if len(players) == 0 {
		return nil
	}


	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO players (id, team_id, first_name, last_name, country, age, position, market_value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, player := range players {
		_, err := stmt.ExecContext(ctx,
			player.ID,
			player.TeamID,
			player.FirstName,
			player.LastName,
			player.Country,
			player.Age,
			player.Position,
			player.MarketValue,
			player.CreatedAt,
			player.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *playerRepository) GetByID(ctx context.Context, id string) (*domain.Player, error) {
	var player domain.Player
	query := `
		SELECT id, team_id, first_name, last_name, country, age, position, market_value, created_at, updated_at 
		FROM players WHERE id = $1
	`
	err := r.db.GetContext(ctx, &player, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPlayerNotFound
		}
		return nil, err
	}
	return &player, nil
}

func (r *playerRepository) GetByTeamID(ctx context.Context, teamID string) ([]*domain.Player, error) {
	var players []*domain.Player
	query := `
		SELECT id, team_id, first_name, last_name, country, age, position, market_value, created_at, updated_at 
		FROM players WHERE team_id = $1
		ORDER BY position, last_name, first_name
	`
	err := r.db.SelectContext(ctx, &players, query, teamID)
	return players, err
}

func (r *playerRepository) Update(ctx context.Context, player *domain.Player) error {
	query := `
		UPDATE players 
		SET team_id = $1, first_name = $2, last_name = $3, country = $4, 
		    age = $5, position = $6, market_value = $7, updated_at = $8
		WHERE id = $9
	`
	_, err := r.db.ExecContext(ctx, query,
		player.TeamID, player.FirstName, player.LastName, player.Country,
		player.Age, player.Position, player.MarketValue, player.UpdatedAt, player.ID)
	return err
}

func (r *playerRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM players WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *playerRepository) GetByTeamIDAndPosition(ctx context.Context, teamID string, position domain.Position) ([]*domain.Player, error) {
	var players []*domain.Player
	query := `
		SELECT id, team_id, first_name, last_name, country, age, position, market_value, created_at, updated_at 
		FROM players WHERE team_id = $1 AND position = $2
		ORDER BY last_name, first_name
	`
	err := r.db.SelectContext(ctx, &players, query, teamID, position)
	return players, err
}
