package domain

import "time"

// Team represents a football team under company XYZ
// Fields: name, logo, year founded, stadium address, city
// All fields are required for registration

type Team struct {
	Name        string     `json:"name" binding:"required"`
	Logo        string     `json:"logo" binding:"required"`
	YearFounded int        `json:"year_founded" binding:"required"`
	StadiumAddr string     `json:"stadium_addr" binding:"required"`
	City        string     `json:"city" binding:"required"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
