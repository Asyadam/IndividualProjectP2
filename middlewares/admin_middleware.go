package middlewares

import (
	"net/http"

	"sport-venue-rental-api/utils"

	"github.com/labstack/echo/v4"
)

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get("role")

		if role == nil {
			return utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		}

		if role != "admin" {
			return utils.ErrorResponse(c, http.StatusForbidden, "admin access only")
		}

		return next(c)
	}
}
