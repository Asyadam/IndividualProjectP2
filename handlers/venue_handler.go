package handlers

import (
	"net/http"
	"strconv"

	"sport-venue-rental-api/dto"
	"sport-venue-rental-api/services"
	"sport-venue-rental-api/utils"

	"github.com/labstack/echo/v4"
)

type VenueHandler struct {
	venueService services.VenueService
}

func NewVenueHandler(venueService services.VenueService) *VenueHandler {
	return &VenueHandler{venueService: venueService}
}

func (h *VenueHandler) Create(c echo.Context) error {
	var request dto.VenueRequest

	if err := c.Bind(&request); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate.Struct(request); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	venue, err := h.venueService.Create(request)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusCreated, "venue created successfully", venue)
}

func (h *VenueHandler) GetAll(c echo.Context) error {
	venues, err := h.venueService.GetAll()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "success get venues", venues)
}

func (h *VenueHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid venue id")
	}

	venue, err := h.venueService.GetByID(id)
	if err != nil {
		if err.Error() == "venue not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		}

		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "success get venue detail", venue)
}

func (h *VenueHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid venue id")
	}

	var request dto.VenueRequest

	if err := c.Bind(&request); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate.Struct(request); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	venue, err := h.venueService.Update(id, request)
	if err != nil {
		if err.Error() == "venue not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		}

		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "venue updated successfully", venue)
}
