package dto

type VenueRequest struct {
	Name              string `json:"name" validate:"required,min=3"`
	Category          string `json:"category" validate:"required"`
	Location          string `json:"location" validate:"required"`
	StockAvailability int    `json:"stock_availability" validate:"required,gte=0"`
	RentalCost        int    `json:"rental_cost" validate:"required,gt=0"`
}
