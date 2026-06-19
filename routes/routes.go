package routes

import (
	"net/http"

	"sport-venue-rental-api/config"
	"sport-venue-rental-api/handlers"
	"sport-venue-rental-api/middlewares"
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

	venueRepository := repositories.NewVenueRepository(config.DB)
	venueService := services.NewVenueService(venueRepository)
	venueHandler := handlers.NewVenueHandler(venueService)

	depositRepository := repositories.NewDepositRepository(config.DB)
	xenditService := services.NewXenditService()
	depositService := services.NewDepositService(userRepository, depositRepository, xenditService)
	depositHandler := handlers.NewDepositHandler(depositService)

	rentalRepository := repositories.NewRentalRepository(config.DB)
	rentalService := services.NewRentalService(userRepository, venueRepository, rentalRepository)
	rentalHandler := handlers.NewRentalHandler(rentalService)

	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	e.POST("/deposits", depositHandler.CreateDeposit, middlewares.JWTMiddleware)
	e.POST("/xendit/callback", depositHandler.XenditCallback)

	e.POST("/venues", venueHandler.Create, middlewares.JWTMiddleware, middlewares.AdminOnly)
	e.GET("/venues", venueHandler.GetAll)
	e.GET("/venues/:id", venueHandler.GetByID)
	e.PUT("/venues/:id", venueHandler.Update, middlewares.JWTMiddleware, middlewares.AdminOnly)

	e.POST("/rentals", rentalHandler.CreateRental, middlewares.JWTMiddleware)
	e.GET("/rentals", rentalHandler.GetRentalHistory, middlewares.JWTMiddleware)
}
