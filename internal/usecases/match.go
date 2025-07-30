package usecases

import (
	"context"
	"errors"
	"football-team-management/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchRepository interface {
	Register(ctx context.Context, match domain.Match) error
	Update(ctx context.Context, id int, match domain.Match) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]domain.Match, error)
	ListByTeam(ctx context.Context, teamName string) ([]domain.Match, error)
	GetByID(ctx context.Context, id int) (*domain.Match, error)
	Restore(ctx context.Context, id int) error
}

type PostgresMatchRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresMatchRepo(pool *pgxpool.Pool) *PostgresMatchRepo {
	return &PostgresMatchRepo{pool: pool}
}

func (r *PostgresMatchRepo) Register(ctx context.Context, match domain.Match) error {
	// Check if home team exists
	var homeTeamExists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1 AND deleted_at IS NULL)`, match.HomeTeam).Scan(&homeTeamExists)
	if err != nil {
		return err
	}
	if !homeTeamExists {
		return errors.New("home team not found")
	}

	// Check if away team exists
	var awayTeamExists bool
	err = r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1 AND deleted_at IS NULL)`, match.AwayTeam).Scan(&awayTeamExists)
	if err != nil {
		return err
	}
	if !awayTeamExists {
		return errors.New("away team not found")
	}

	// Check if teams are different
	if match.HomeTeam == match.AwayTeam {
		return errors.New("home team and away team cannot be the same")
	}

	// Parse the time string to time.Time
	matchTime, err := time.Parse("15:04", match.MatchTime)
	if err != nil {
		return errors.New("invalid time format. Use HH:MM")
	}

	now := time.Now()
	_, err = r.pool.Exec(ctx, `INSERT INTO matches (match_date, match_time, home_team, away_team, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, NULL)`,
		match.MatchDate, matchTime, match.HomeTeam, match.AwayTeam, now, now)
	return err
}

func (r *PostgresMatchRepo) Update(ctx context.Context, id int, match domain.Match) error {
	// Check if home team exists
	var homeTeamExists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1 AND deleted_at IS NULL)`, match.HomeTeam).Scan(&homeTeamExists)
	if err != nil {
		return err
	}
	if !homeTeamExists {
		return errors.New("home team not found")
	}

	// Check if away team exists
	var awayTeamExists bool
	err = r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1 AND deleted_at IS NULL)`, match.AwayTeam).Scan(&awayTeamExists)
	if err != nil {
		return err
	}
	if !awayTeamExists {
		return errors.New("away team not found")
	}

	// Check if teams are different
	if match.HomeTeam == match.AwayTeam {
		return errors.New("home team and away team cannot be the same")
	}

	// Parse the time string to time.Time
	matchTime, err := time.Parse("15:04", match.MatchTime)
	if err != nil {
		return errors.New("invalid time format. Use HH:MM")
	}

	now := time.Now()
	cmd, err := r.pool.Exec(ctx, `UPDATE matches SET match_date=$1, match_time=$2, home_team=$3, away_team=$4, updated_at=$5 WHERE id=$6 AND deleted_at IS NULL`,
		match.MatchDate, matchTime, match.HomeTeam, match.AwayTeam, now, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("match not found")
	}
	return nil
}

func (r *PostgresMatchRepo) Delete(ctx context.Context, id int) error {
	now := time.Now()
	cmd, err := r.pool.Exec(ctx, `UPDATE matches SET deleted_at=$1, updated_at=$2 WHERE id=$3 AND deleted_at IS NULL`, now, now, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("match not found")
	}
	return nil
}

func (r *PostgresMatchRepo) List(ctx context.Context) ([]domain.Match, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, match_date, match_time, home_team, away_team, created_at, updated_at, deleted_at FROM matches WHERE deleted_at IS NULL ORDER BY match_date, match_time`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var matches []domain.Match
	for rows.Next() {
		var m domain.Match
		var deletedAt *time.Time
		var matchTime time.Time
		if err := rows.Scan(&m.ID, &m.MatchDate, &matchTime, &m.HomeTeam, &m.AwayTeam, &m.CreatedAt, &m.UpdatedAt, &deletedAt); err != nil {
			return nil, err
		}
		m.MatchTime = matchTime.Format("15:04")
		m.DeletedAt = deletedAt
		matches = append(matches, m)
	}
	return matches, nil
}

func (r *PostgresMatchRepo) ListByTeam(ctx context.Context, teamName string) ([]domain.Match, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, match_date, match_time, home_team, away_team, created_at, updated_at, deleted_at FROM matches WHERE (home_team = $1 OR away_team = $1) AND deleted_at IS NULL ORDER BY match_date, match_time`, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var matches []domain.Match
	for rows.Next() {
		var m domain.Match
		var deletedAt *time.Time
		var matchTime time.Time
		if err := rows.Scan(&m.ID, &m.MatchDate, &matchTime, &m.HomeTeam, &m.AwayTeam, &m.CreatedAt, &m.UpdatedAt, &deletedAt); err != nil {
			return nil, err
		}
		m.MatchTime = matchTime.Format("15:04")
		m.DeletedAt = deletedAt
		matches = append(matches, m)
	}
	return matches, nil
}

func (r *PostgresMatchRepo) GetByID(ctx context.Context, id int) (*domain.Match, error) {
	var m domain.Match
	var deletedAt *time.Time
	var matchTime time.Time
	err := r.pool.QueryRow(ctx, `SELECT id, match_date, match_time, home_team, away_team, created_at, updated_at, deleted_at FROM matches WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&m.ID, &m.MatchDate, &matchTime, &m.HomeTeam, &m.AwayTeam, &m.CreatedAt, &m.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}
	m.MatchTime = matchTime.Format("15:04")
	m.DeletedAt = deletedAt
	return &m, nil
}

func (r *PostgresMatchRepo) Restore(ctx context.Context, id int) error {
	now := time.Now()
	cmd, err := r.pool.Exec(ctx, `UPDATE matches SET deleted_at=NULL, updated_at=$1 WHERE id=$2 AND deleted_at IS NOT NULL`, now, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("match not found or not deleted")
	}
	return nil
}
