package domain

import "time"

// Player represents a football player under a team
// Fields: name, height, weight, position, jersey number
// All fields are required for registration

type PlayerPosition string

const (
	PositionForward    PlayerPosition = "penyerang"
	PositionMidfielder PlayerPosition = "gelandang"
	PositionDefender   PlayerPosition = "bertahan"
	PositionGoalkeeper PlayerPosition = "penjaga gawang"
)

type Player struct {
	Name         string         `json:"name" binding:"required"`
	Height       int            `json:"height" binding:"required"` // in cm
	Weight       int            `json:"weight" binding:"required"` // in kg
	Position     PlayerPosition `json:"position" binding:"required"`
	JerseyNumber int            `json:"jersey_number" binding:"required"`
	TeamName     string         `json:"team_name" binding:"required"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    *time.Time     `json:"deleted_at,omitempty"`
}
