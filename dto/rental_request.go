package dto

type RentalRequest struct {
	VenueID    int    `json:"venue_id" validate:"required,gt=0"`
	RentalDate string `json:"rental_date" validate:"required"`
	StartTime  string `json:"start_time" validate:"required"`
	EndTime    string `json:"end_time" validate:"required"`
}

type RentalResponse struct {
	ID               int    `json:"id"`
	VenueID          int    `json:"venue_id"`
	VenueName        string `json:"venue_name"`
	RentalDate       string `json:"rental_date"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	TotalCost        int    `json:"total_cost"`
	RemainingDeposit int    `json:"remaining_deposit"`
	Status           string `json:"status"`
}
