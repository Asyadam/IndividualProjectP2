package handlers

import (
	"net/http"

	"sport-venue-rental-api/dto"
	"sport-venue-rental-api/services"
	"sport-venue-rental-api/utils"

	"github.com/labstack/echo/v4"
)

type RentalHandler struct {
	rentalService services.RentalService
}

func NewRentalHandler(rentalService services.RentalService) *RentalHandler {
	return &RentalHandler{rentalService: rentalService}
}

func (h *RentalHandler) CreateRental(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	var request dto.RentalRequest

	if err := c.Bind(&request); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate.Struct(request); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	rental, err := h.rentalService.CreateRental(userID, request)
	if err != nil {
		switch err.Error() {
		case "user not found", "venue not found":
			return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		case "venue is not available",
			"invalid rental date format",
			"invalid start time format",
			"invalid end time format",
			"end time must be after start time",
			"minimum rental duration is 1 hour",
			"venue already booked at this time",
			"insufficient deposit amount":
			return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		default:
			return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	return utils.SuccessResponse(c, http.StatusCreated, "rental created successfully", rental)
}

func (h *RentalHandler) GetRentalHistory(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	rentals, err := h.rentalService.GetRentalHistory(userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "success get rental history", rentals)
}
