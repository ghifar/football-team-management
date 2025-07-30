package domain

import "time"

// MatchResult represents the result of a completed football match
// Fields: match ID, home score, away score, goals details
// All fields are required for reporting

type Goal struct {
	ID        int        `json:"id"`
	MatchID   int        `json:"match_id"`
	Scorer    string     `json:"scorer" binding:"required"`    // Player name who scored
	GoalTime  string     `json:"goal_time" binding:"required"` // Format: "MM:SS" or "HH:MM:SS"
	Team      string     `json:"team" binding:"required"`      // Team that scored
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type MatchResult struct {
	ID        int        `json:"id"`
	MatchID   int        `json:"match_id" binding:"required"`
	HomeScore int        `json:"home_score" binding:"required"`
	AwayScore int        `json:"away_score" binding:"required"`
	Goals     []Goal     `json:"goals,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// MatchResultRequest represents the request structure for reporting match results
type MatchResultRequest struct {
	MatchID   int    `json:"match_id" binding:"required"`
	HomeScore int    `json:"home_score" binding:"required"`
	AwayScore int    `json:"away_score" binding:"required"`
	Goals     []Goal `json:"goals,omitempty"`
}

// ToMatchResult converts MatchResultRequest to MatchResult domain model
func (mr *MatchResultRequest) ToMatchResult() *MatchResult {
	return &MatchResult{
		MatchID:   mr.MatchID,
		HomeScore: mr.HomeScore,
		AwayScore: mr.AwayScore,
		Goals:     mr.Goals,
	}
}

// MatchResultResponse represents the response structure for match results
type MatchResultResponse struct {
	ID        int        `json:"id"`
	MatchID   int        `json:"match_id"`
	HomeScore int        `json:"home_score"`
	AwayScore int        `json:"away_score"`
	Goals     []Goal     `json:"goals,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// ToMatchResultResponse converts MatchResult domain model to MatchResultResponse
func (mr *MatchResult) ToMatchResultResponse() *MatchResultResponse {
	return &MatchResultResponse{
		ID:        mr.ID,
		MatchID:   mr.MatchID,
		HomeScore: mr.HomeScore,
		AwayScore: mr.AwayScore,
		Goals:     mr.Goals,
		CreatedAt: mr.CreatedAt,
		UpdatedAt: mr.UpdatedAt,
		DeletedAt: mr.DeletedAt,
	}
}
