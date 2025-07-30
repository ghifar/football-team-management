package usecases

import (
	"context"
	"errors"
	"football-team-management/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchResultRepository interface {
	Register(ctx context.Context, result domain.MatchResult) error
	Update(ctx context.Context, id int, result domain.MatchResult) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]domain.MatchResult, error)
	GetByMatchID(ctx context.Context, matchID int) (*domain.MatchResult, error)
	GetByID(ctx context.Context, id int) (*domain.MatchResult, error)
	Restore(ctx context.Context, id int) error
}

type PostgresMatchResultRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresMatchResultRepo(pool *pgxpool.Pool) *PostgresMatchResultRepo {
	return &PostgresMatchResultRepo{pool: pool}
}

func (r *PostgresMatchResultRepo) Register(ctx context.Context, result domain.MatchResult) error {
	// Check if match exists
	var matchExists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM matches WHERE id = $1 AND deleted_at IS NULL)`, result.MatchID).Scan(&matchExists)
	if err != nil {
		return err
	}
	if !matchExists {
		return errors.New("match not found")
	}

	// Check if result already exists for this match
	var resultExists bool
	err = r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM match_results WHERE match_id = $1 AND deleted_at IS NULL)`, result.MatchID).Scan(&resultExists)
	if err != nil {
		return err
	}
	if resultExists {
		return errors.New("result already exists for this match")
	}

	// Validate that scores match the number of goals
	homeGoals := 0
	awayGoals := 0
	for _, goal := range result.Goals {
		if goal.Team == "home" {
			homeGoals++
		} else if goal.Team == "away" {
			awayGoals++
		}
	}

	if homeGoals != result.HomeScore {
		return errors.New("home score does not match number of home goals")
	}
	if awayGoals != result.AwayScore {
		return errors.New("away score does not match number of away goals")
	}

	// Start transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	now := time.Now()

	// Insert match result
	var resultID int
	err = tx.QueryRow(ctx, `INSERT INTO match_results (match_id, home_score, away_score, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, NULL) RETURNING id`,
		result.MatchID, result.HomeScore, result.AwayScore, now, now).Scan(&resultID)
	if err != nil {
		return err
	}

	// Insert goals
	for _, goal := range result.Goals {
		_, err = tx.Exec(ctx, `INSERT INTO goals (match_id, scorer, goal_time, team, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, NULL)`,
			result.MatchID, goal.Scorer, goal.GoalTime, goal.Team, now, now)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresMatchResultRepo) Update(ctx context.Context, id int, result domain.MatchResult) error {
	// Check if result exists
	var resultExists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM match_results WHERE id = $1 AND deleted_at IS NULL)`, id).Scan(&resultExists)
	if err != nil {
		return err
	}
	if !resultExists {
		return errors.New("match result not found")
	}

	// Validate that scores match the number of goals
	homeGoals := 0
	awayGoals := 0
	for _, goal := range result.Goals {
		if goal.Team == "home" {
			homeGoals++
		} else if goal.Team == "away" {
			awayGoals++
		}
	}

	if homeGoals != result.HomeScore {
		return errors.New("home score does not match number of home goals")
	}
	if awayGoals != result.AwayScore {
		return errors.New("away score does not match number of away goals")
	}

	// Start transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	now := time.Now()

	// Update match result
	cmd, err := tx.Exec(ctx, `UPDATE match_results SET home_score=$1, away_score=$2, updated_at=$3 WHERE id=$4 AND deleted_at IS NULL`,
		result.HomeScore, result.AwayScore, now, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("match result not found")
	}

	// Delete existing goals for this match
	_, err = tx.Exec(ctx, `DELETE FROM goals WHERE match_id = $1`, result.MatchID)
	if err != nil {
		return err
	}

	// Insert new goals
	for _, goal := range result.Goals {
		_, err = tx.Exec(ctx, `INSERT INTO goals (match_id, scorer, goal_time, team, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, NULL)`,
			result.MatchID, goal.Scorer, goal.GoalTime, goal.Team, now, now)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresMatchResultRepo) Delete(ctx context.Context, id int) error {
	now := time.Now()
	cmd, err := r.pool.Exec(ctx, `UPDATE match_results SET deleted_at=$1, updated_at=$2 WHERE id=$3 AND deleted_at IS NULL`, now, now, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("match result not found")
	}
	return nil
}

func (r *PostgresMatchResultRepo) List(ctx context.Context) ([]domain.MatchResult, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, match_id, home_score, away_score, created_at, updated_at, deleted_at FROM match_results WHERE deleted_at IS NULL ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.MatchResult
	for rows.Next() {
		var result domain.MatchResult
		var deletedAt *time.Time
		if err := rows.Scan(&result.ID, &result.MatchID, &result.HomeScore, &result.AwayScore, &result.CreatedAt, &result.UpdatedAt, &deletedAt); err != nil {
			return nil, err
		}
		result.DeletedAt = deletedAt

		// Get goals for this result
		goals, err := r.getGoalsByMatchID(ctx, result.MatchID)
		if err != nil {
			return nil, err
		}
		result.Goals = goals

		results = append(results, result)
	}
	return results, nil
}

func (r *PostgresMatchResultRepo) GetByMatchID(ctx context.Context, matchID int) (*domain.MatchResult, error) {
	var result domain.MatchResult
	var deletedAt *time.Time
	err := r.pool.QueryRow(ctx, `SELECT id, match_id, home_score, away_score, created_at, updated_at, deleted_at FROM match_results WHERE match_id = $1 AND deleted_at IS NULL`, matchID).
		Scan(&result.ID, &result.MatchID, &result.HomeScore, &result.AwayScore, &result.CreatedAt, &result.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}
	result.DeletedAt = deletedAt

	// Get goals for this result
	goals, err := r.getGoalsByMatchID(ctx, matchID)
	if err != nil {
		return nil, err
	}
	result.Goals = goals

	return &result, nil
}

func (r *PostgresMatchResultRepo) GetByID(ctx context.Context, id int) (*domain.MatchResult, error) {
	var result domain.MatchResult
	var deletedAt *time.Time
	err := r.pool.QueryRow(ctx, `SELECT id, match_id, home_score, away_score, created_at, updated_at, deleted_at FROM match_results WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&result.ID, &result.MatchID, &result.HomeScore, &result.AwayScore, &result.CreatedAt, &result.UpdatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}
	result.DeletedAt = deletedAt

	// Get goals for this result
	goals, err := r.getGoalsByMatchID(ctx, result.MatchID)
	if err != nil {
		return nil, err
	}
	result.Goals = goals

	return &result, nil
}

func (r *PostgresMatchResultRepo) Restore(ctx context.Context, id int) error {
	now := time.Now()
	cmd, err := r.pool.Exec(ctx, `UPDATE match_results SET deleted_at=NULL, updated_at=$1 WHERE id=$2 AND deleted_at IS NOT NULL`, now, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("match result not found or not deleted")
	}
	return nil
}

// Helper method to get goals by match ID
func (r *PostgresMatchResultRepo) getGoalsByMatchID(ctx context.Context, matchID int) ([]domain.Goal, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, match_id, scorer, goal_time, team, created_at, updated_at, deleted_at FROM goals WHERE match_id = $1 AND deleted_at IS NULL ORDER BY goal_time`, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []domain.Goal
	for rows.Next() {
		var goal domain.Goal
		var deletedAt *time.Time
		if err := rows.Scan(&goal.ID, &goal.MatchID, &goal.Scorer, &goal.GoalTime, &goal.Team, &goal.CreatedAt, &goal.UpdatedAt, &deletedAt); err != nil {
			return nil, err
		}
		goal.DeletedAt = deletedAt
		goals = append(goals, goal)
	}
	return goals, nil
}
