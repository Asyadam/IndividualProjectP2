package middlewares

import (
	"net/http"
	"os"
	"strings"

	"sport-venue-rental-api/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			return utils.ErrorResponse(c, http.StatusUnauthorized, "missing authorization header")
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims := &utils.JwtCustomClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		return next(c)
	}
}
