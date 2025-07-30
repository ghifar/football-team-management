package usecases

import (
	"context"
	"errors"
	"football-team-management/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository interface {
	Register(ctx context.Context, team domain.Team) error
	Update(ctx context.Context, name string, team domain.Team) error
	Delete(ctx context.Context, name string) error
	List(ctx context.Context) ([]domain.Team, error)
	Restore(ctx context.Context, name string) error
}

type PostgresTeamRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresTeamRepo(pool *pgxpool.Pool) *PostgresTeamRepo {
	return &PostgresTeamRepo{pool: pool}
}

func (r *PostgresTeamRepo) Register(ctx context.Context, team domain.Team) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx, `INSERT INTO teams (name, logo, year_founded, stadium_addr, city, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NULL)`,
		team.Name, team.Logo, team.YearFounded, team.StadiumAddr, team.City, now, now)
	return err
}

func (r *PostgresTeamRepo) Update(ctx context.Context, name string, team domain.Team) error {
	now := time.Now()
	cmd, err := r.pool.Exec(ctx, `UPDATE teams SET name=$1, logo=$2, year_founded=$3, stadium_addr=$4, city=$5, updated_at=$6 WHERE name=$7 AND deleted_at IS NULL`,
		team.Name, team.Logo, team.YearFounded, team.StadiumAddr, team.City, now, name)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("team not found")
	}
	return nil
}

func (r *PostgresTeamRepo) Delete(ctx context.Context, name string) error {
	now := time.Now()
	cmd, err := r.pool.Exec(ctx, `UPDATE teams SET deleted_at=$1, updated_at=$2 WHERE name=$3 AND deleted_at IS NULL`, now, now, name)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("team not found")
	}
	return nil
}

func (r *PostgresTeamRepo) List(ctx context.Context) ([]domain.Team, error) {
	rows, err := r.pool.Query(ctx, `SELECT name, logo, year_founded, stadium_addr, city, created_at, updated_at, deleted_at FROM teams WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teams []domain.Team
	for rows.Next() {
		var t domain.Team
		var deletedAt *time.Time
		if err := rows.Scan(&t.Name, &t.Logo, &t.YearFounded, &t.StadiumAddr, &t.City, &t.CreatedAt, &t.UpdatedAt, &deletedAt); err != nil {
			return nil, err
		}
		t.DeletedAt = deletedAt
		teams = append(teams, t)
	}
	return teams, nil
}

func (r *PostgresTeamRepo) Restore(ctx context.Context, name string) error {
	now := time.Now()
	cmd, err := r.pool.Exec(ctx, `UPDATE teams SET deleted_at=NULL, updated_at=$1 WHERE name=$2 AND deleted_at IS NOT NULL`, now, name)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("team not found or not deleted")
	}
	return nil
}
