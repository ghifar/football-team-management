package domain

import (
	"time"
)

// Match represents a football match schedule between two teams
// Fields: match date, match time, home team, away team
// All fields are required for registration

type Match struct {
	ID        int        `json:"id"`
	MatchDate time.Time  `json:"match_date" binding:"required"`
	MatchTime string     `json:"match_time" binding:"required"` // Format: "HH:MM"
	HomeTeam  string     `json:"home_team" binding:"required"`
	AwayTeam  string     `json:"away_team" binding:"required"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// MatchRequest represents the request structure for creating/updating matches
type MatchRequest struct {
	MatchDate string `json:"match_date" binding:"required"` // Format: "YYYY-MM-DD"
	MatchTime string `json:"match_time" binding:"required"` // Format: "HH:MM"
	HomeTeam  string `json:"home_team" binding:"required"`
	AwayTeam  string `json:"away_team" binding:"required"`
}

// ToMatch converts MatchRequest to Match domain model
func (mr *MatchRequest) ToMatch() (*Match, error) {
	// Parse the date string
	date, err := time.Parse("2006-01-02", mr.MatchDate)
	if err != nil {
		return nil, err
	}

	return &Match{
		MatchDate: date,
		MatchTime: mr.MatchTime,
		HomeTeam:  mr.HomeTeam,
		AwayTeam:  mr.AwayTeam,
	}, nil
}

// MatchResponse represents the response structure for matches
type MatchResponse struct {
	ID        int        `json:"id"`
	MatchDate string     `json:"match_date"` // Format: "YYYY-MM-DD"
	MatchTime string     `json:"match_time"` // Format: "HH:MM"
	HomeTeam  string     `json:"home_team"`
	AwayTeam  string     `json:"away_team"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// ToMatchResponse converts Match domain model to MatchResponse
func (m *Match) ToMatchResponse() *MatchResponse {
	return &MatchResponse{
		ID:        m.ID,
		MatchDate: m.MatchDate.Format("2006-01-02"),
		MatchTime: m.MatchTime,
		HomeTeam:  m.HomeTeam,
		AwayTeam:  m.AwayTeam,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}
