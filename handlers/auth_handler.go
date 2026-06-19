package handlers

import (
	"net/http"

	"sport-venue-rental-api/dto"
	"sport-venue-rental-api/services"
	"sport-venue-rental-api/utils"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var request dto.RegisterRequest

	if err := c.Bind(&request); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate.Struct(request); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := h.authService.Register(request)
	if err != nil {
		if err.Error() == "email already registered" {
			return utils.ErrorResponse(c, http.StatusConflict, err.Error())
		}

		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusCreated, "register success", user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var request dto.LoginRequest

	if err := c.Bind(&request); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate.Struct(request); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	response, err := h.authService.Login(request)
	if err != nil {
		if err.Error() == "invalid email or password" {
			return utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		}

		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "login success", response)
}
