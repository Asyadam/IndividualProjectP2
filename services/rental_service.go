package services

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"sport-venue-rental-api/dto"
	"sport-venue-rental-api/models"
	"sport-venue-rental-api/repositories"
)

type RentalService interface {
	CreateRental(userID int, request dto.RentalRequest) (dto.RentalResponse, error)
	GetRentalHistory(userID int) ([]models.RentalHistory, error)
}

type rentalService struct {
	userRepository   repositories.UserRepository
	venueRepository  repositories.VenueRepository
	rentalRepository repositories.RentalRepository
}

func NewRentalService(
	userRepository repositories.UserRepository,
	venueRepository repositories.VenueRepository,
	rentalRepository repositories.RentalRepository,
) RentalService {
	return &rentalService{
		userRepository:   userRepository,
		venueRepository:  venueRepository,
		rentalRepository: rentalRepository,
	}
}

func (s *rentalService) CreateRental(userID int, request dto.RentalRequest) (dto.RentalResponse, error) {
	user, err := s.userRepository.FindByID(userID)
	if err == sql.ErrNoRows {
		return dto.RentalResponse{}, errors.New("user not found")
	}

	if err != nil {
		log.Printf("[ERROR] failed to find user: %v", err)
		return dto.RentalResponse{}, errors.New("internal server error")
	}

	venue, err := s.venueRepository.FindByID(request.VenueID)
	if err == sql.ErrNoRows {
		return dto.RentalResponse{}, errors.New("venue not found")
	}

	if err != nil {
		log.Printf("[ERROR] failed to find venue: %v", err)
		return dto.RentalResponse{}, errors.New("internal server error")
	}

	if venue.StockAvailability <= 0 {
		return dto.RentalResponse{}, errors.New("venue is not available")
	}

	_, err = time.Parse("2006-01-02", request.RentalDate)
	if err != nil {
		return dto.RentalResponse{}, errors.New("invalid rental date format")
	}

	startTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		return dto.RentalResponse{}, errors.New("invalid start time format")
	}

	endTime, err := time.Parse("15:04", request.EndTime)
	if err != nil {
		return dto.RentalResponse{}, errors.New("invalid end time format")
	}

	if !endTime.After(startTime) {
		return dto.RentalResponse{}, errors.New("end time must be after start time")
	}

	duration := endTime.Sub(startTime).Hours()
	if duration < 1 {
		return dto.RentalResponse{}, errors.New("minimum rental duration is 1 hour")
	}

	totalCost := int(duration * float64(venue.RentalCost))

	hasConflict, err := s.rentalRepository.HasScheduleConflict(
		request.VenueID,
		request.RentalDate,
		request.StartTime,
		request.EndTime,
	)
	if err != nil {
		log.Printf("[ERROR] failed to check schedule conflict: %v", err)
		return dto.RentalResponse{}, errors.New("internal server error")
	}

	if hasConflict {
		return dto.RentalResponse{}, errors.New("venue already booked at this time")
	}

	if user.DepositAmount < totalCost {
		return dto.RentalResponse{}, errors.New("insufficient deposit amount")
	}

	rental := models.Rental{
		UserID:     userID,
		VenueID:    request.VenueID,
		RentalDate: request.RentalDate,
		StartTime:  request.StartTime,
		EndTime:    request.EndTime,
		TotalCost:  totalCost,
		Status:     "booked",
	}

	createdRental, remainingDeposit, err := s.rentalRepository.CreateWithDepositDeduction(rental)
	if err != nil {
		log.Printf("[ERROR] failed to create rental: %v", err)
		return dto.RentalResponse{}, errors.New("internal server error")
	}

	response := dto.RentalResponse{
		ID:               createdRental.ID,
		VenueID:          createdRental.VenueID,
		VenueName:        venue.Name,
		RentalDate:       createdRental.RentalDate,
		StartTime:        createdRental.StartTime,
		EndTime:          createdRental.EndTime,
		TotalCost:        createdRental.TotalCost,
		RemainingDeposit: remainingDeposit,
		Status:           createdRental.Status,
	}

	return response, nil
}

func (s *rentalService) GetRentalHistory(userID int) ([]models.RentalHistory, error) {
	rentals, err := s.rentalRepository.FindByUserID(userID)
	if err != nil {
		log.Printf("[ERROR] failed to get rental history: %v", err)
		return nil, errors.New("internal server error")
	}

	return rentals, nil
}
