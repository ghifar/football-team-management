package usecases

import (
	"context"
	"errors"
	"football-team-management/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PlayerRepository interface {
	Register(ctx context.Context, player domain.Player) error
	Update(ctx context.Context, name string, player domain.Player) error
	Delete(ctx context.Context, name string) error
	List(ctx context.Context) ([]domain.Player, error)
	ListByTeam(ctx context.Context, teamName string) ([]domain.Player, error)
	GetByName(ctx context.Context, name string) (*domain.Player, error)
	Restore(ctx context.Context, name string) error
}

type PostgresPlayerRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresPlayerRepo(pool *pgxpool.Pool) *PostgresPlayerRepo {
	return &PostgresPlayerRepo{pool: pool}
}

func (r *PostgresPlayerRepo) Register(ctx context.Context, player domain.Player) error {
	// Check if team exists
	var teamExists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1 AND deleted_at IS NULL)`, player.TeamName).Scan(&teamExists)
	if err != nil {
		return err
	}
	if !teamExists {
		return errors.New("team not found")
	}

	// Check if jersey number is already taken in the team
	var jerseyExists bool
	err = r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM players WHERE team_name = $1 AND jersey_number = $2 AND deleted_at IS NULL)`,
		player.TeamName, player.JerseyNumber).Scan(&jerseyExists)
	if err != nil {
		return err
	}
	if jerseyExists {
		return errors.New("jersey number already taken in this team")
	}

	// Check if player already exists
	var playerExists bool
	err = r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM players WHERE name = $1 AND deleted_at IS NULL)`, player.Name).Scan(&playerExists)
	if err != nil {
		return err
	}
	if playerExists {
		return errors.New("player already exists")
	}

	now := time.Now()
	_, err = r.pool.Exec(ctx, `INSERT INTO players (name, height, weight, position, jersey_number, team_name, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NULL)`,
		player.Name, player.Height, player.Weight, player.Position, player.JerseyNumber, player.TeamName, now, now)
	return err
}

func (r *PostgresPlayerRepo) Update(ctx context.Context, name string, player domain.Player) error {
	// Check if team exists
	var teamExists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1 AND deleted_at IS NULL)`, player.TeamName).Scan(&teamExists)
	if err != nil {
		return err
	}
	if !teamExists {
		return errors.New("team not found")
	}

	// Check if jersey number is already taken by another player in the team
	var jerseyExists bool
	err = r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM players WHERE team_name = $1 AND jersey_number = $2 AND name != $3 AND deleted_at IS NULL)`,
		player.TeamName, player.JerseyNumber, name).Scan(&jerseyExists)
	if err != nil {
		return err
	}
	if jerseyExists {
		return errors.New("jersey number already taken in this team")
	}

	now := time.Now()
	cmd, err := r.pool.Exec(ctx, `UPDATE players SET name=$1, height=$2, weight=$3, position=$4, jersey_number=$5, team_name=$6, updated_at=$7 WHERE name=$8 AND deleted_at IS NULL`,
		player.Name, player.Height, player.Weight, player.Position, player.JerseyNumber, player.TeamName, now, name)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("player not found")
	}
	return nil
}

func (r *PostgresPlayerRepo) Delete(ctx context.Context, name string) error {
	now := time.Now()
	cmd, err := r.pool.Exec(ctx, `UPDATE players SET deleted_at=$1, updated_at=$2 WHERE name=$3 AND deleted_at IS NULL`, now, now, name)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("player not found")
	}
	return nil
}

func (r *PostgresPlayerRepo) List(ctx context.Context) ([]domain.Player, error) {
	rows, err := r.pool.Query(ctx, `SELECT name, height, weight, position, jersey_number, team_name, created_at, updated_at, deleted_at FROM players WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var players []domain.Player
	for rows.Next() {
		var p domain.Player
		var deletedAt *time.Time
		if err := rows.Scan(&p.Name, &p.Height, &p.Weight, &p.Position, &p.JerseyNumber, &p.TeamName, &p.CreatedAt, &p.UpdatedAt, &deletedAt); err != nil {
			return nil, err
		}
		p.DeletedAt = deletedAt
		players = append(players, p)
	}
	return players, nil
}

func (r *PostgresPlayerRepo) ListByTeam(ctx context.Context, teamName string) ([]domain.Player, error) {
	rows, err := r.pool.Query(ctx, `SELECT name, height, weight, position, jersey_number, team_name, created_at, updated_at, deleted_at FROM players WHERE team_name = $1 AND deleted_at IS NULL`, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var players []domain.Player
	for rows.Next() {
		var p domain.Player
		var deletedAt *time.Time
		if err := rows.Scan(&p.Name, &p.Height, &p.Weight, &p.Position, &p.JerseyNumber, &p.TeamName, &p.CreatedAt, &p.UpdatedAt, &deletedAt); err != nil {
			return nil, err
		}
		p.DeletedAt = deletedAt
		players = append(players, p)
	}
	return players, nil
}

func (r *PostgresPlayerRepo) GetByName(ctx context.Context, name string) (*domain.Player, error) {
	var p domain.Player
	var deletedAt *time.Time
	err := r.pool.QueryRow(ctx, `SELECT name, height, weight, position, jersey_number, team_name, created_at, updated_at, deleted_at FROM players WHERE name = $1 AND deleted_at IS NULL`, name).
		Scan(&p.Name, &p.Height, &p.Weight, &p.Position, &p.JerseyNumber, &p.TeamName, &p.CreatedAt, &p.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}
	p.DeletedAt = deletedAt
	return &p, nil
}

func (r *PostgresPlayerRepo) Restore(ctx context.Context, name string) error {
	now := time.Now()
	cmd, err := r.pool.Exec(ctx, `UPDATE players SET deleted_at=NULL, updated_at=$1 WHERE name=$2 AND deleted_at IS NOT NULL`, now, name)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("player not found or not deleted")
	}
	return nil
}
