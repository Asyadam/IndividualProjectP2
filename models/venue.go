package models

import "time"

type Venue struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Category          string    `json:"category"`
	Location          string    `json:"location"`
	StockAvailability int       `json:"stock_availability"`
	RentalCost        int       `json:"rental_cost"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
