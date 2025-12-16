package postgres

import (
	"context"
	"database/sql"
	"errors"

	"soccer-manager-api/internal/domain"
	"soccer-manager-api/internal/ports/repository"

	"github.com/jmoiron/sqlx"
)

type teamRepository struct {
	db *sqlx.DB
}


func NewTeamRepository(db *sqlx.DB) repository.TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(ctx context.Context, team *domain.Team) error {
	query := `
		INSERT INTO teams (id, user_id, name, country, budget, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query, team.ID, team.UserID, team.Name, team.Country, team.Budget, team.CreatedAt, team.UpdatedAt)
	return err
}

func (r *teamRepository) GetByID(ctx context.Context, id string) (*domain.Team, error) {
	var team domain.Team
	query := `SELECT id, user_id, name, country, budget, created_at, updated_at FROM teams WHERE id = $1`
	err := r.db.GetContext(ctx, &team, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTeamNotFound
		}
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) GetByUserID(ctx context.Context, userID string) (*domain.Team, error) {
	var team domain.Team
	query := `SELECT id, user_id, name, country, budget, created_at, updated_at FROM teams WHERE user_id = $1`
	err := r.db.GetContext(ctx, &team, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTeamNotFound
		}
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) Update(ctx context.Context, team *domain.Team) error {
	query := `
		UPDATE teams 
		SET name = $1, country = $2, budget = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.ExecContext(ctx, query, team.Name, team.Country, team.Budget, team.UpdatedAt, team.ID)
	return err
}

func (r *teamRepository) GetTotalValue(ctx context.Context, teamID string) (float64, error) {
	var totalValue sql.NullFloat64
	query := `SELECT COALESCE(SUM(market_value), 0) FROM players WHERE team_id = $1`
	err := r.db.GetContext(ctx, &totalValue, query, teamID)
	if err != nil {
		return 0, err
	}
	if !totalValue.Valid {
		return 0, nil
	}
	return totalValue.Float64, nil
}

func (r *teamRepository) GetPlayerCount(ctx context.Context, teamID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM players WHERE team_id = $1`
	err := r.db.GetContext(ctx, &count, query, teamID)
	return count, err
}

