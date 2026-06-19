package handlers

import (
	"net/http"
	"os"

	"sport-venue-rental-api/dto"
	"sport-venue-rental-api/services"
	"sport-venue-rental-api/utils"

	"github.com/labstack/echo/v4"
)

type DepositHandler struct {
	depositService services.DepositService
}

func NewDepositHandler(depositService services.DepositService) *DepositHandler {
	return &DepositHandler{depositService: depositService}
}

func (h *DepositHandler) CreateDeposit(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	var request dto.DepositRequest

	if err := c.Bind(&request); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	if err := utils.Validate.Struct(request); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	deposit, err := h.depositService.CreateDeposit(userID, request)
	if err != nil {
		if err.Error() == "failed to create payment link" {
			return utils.ErrorResponse(c, http.StatusBadGateway, err.Error())
		}

		if err.Error() == "user not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		}

		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusCreated, "deposit created successfully", deposit)
}

func (h *DepositHandler) XenditCallback(c echo.Context) error {
	callbackToken := c.Request().Header.Get("x-callback-token")
	expectedToken := os.Getenv("XENDIT_CALLBACK_TOKEN")

	if expectedToken != "" && callbackToken != expectedToken {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "invalid callback token")
	}

	var request dto.XenditCallbackRequest

	if err := c.Bind(&request); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	deposit, err := h.depositService.HandleXenditCallback(request)
	if err != nil {
		if err.Error() == "deposit not found" {
			return utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		}

		if err.Error() == "payment is not paid" || err.Error() == "amount mismatch" {
			return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}

		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "xendit callback processed successfully", deposit)
}
