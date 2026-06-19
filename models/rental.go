package models

import "time"

type Rental struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	VenueID    int       `json:"venue_id"`
	RentalDate string    `json:"rental_date"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	TotalCost  int       `json:"total_cost"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RentalHistory struct {
	ID         int       `json:"id"`
	VenueID    int       `json:"venue_id"`
	VenueName  string    `json:"venue_name"`
	RentalDate string    `json:"rental_date"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	TotalCost  int       `json:"total_cost"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
