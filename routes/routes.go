package routes

import (
	"net/http"

	"sport-venue-rental-api/config"
	"sport-venue-rental-api/handlers"
	"sport-venue-rental-api/repositories"
	"sport-venue-rental-api/services"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "welcome to sport venue rental api",
		})
	})

	userRepository := repositories.NewUserRepository(config.DB)
	authService := services.NewAuthService(userRepository)
	authHandler := handlers.NewAuthHandler(authService)

	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)
}
