package tests

import (
	"errors"
	"testing"
	"time"

	"sport-venue-rental-api/dto"
	"sport-venue-rental-api/models"
	"sport-venue-rental-api/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user models.User) (models.User, error) {
	args := m.Called(user)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (models.User, error) {
	args := m.Called(email)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id int) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) AddDeposit(userID int, amount int) (models.User, error) {
	args := m.Called(userID, amount)
	return args.Get(0).(models.User), args.Error(1)
}

type MockVenueRepository struct {
	mock.Mock
}

func (m *MockVenueRepository) Create(venue models.Venue) (models.Venue, error) {
	args := m.Called(venue)
	return args.Get(0).(models.Venue), args.Error(1)
}

func (m *MockVenueRepository) FindAll() ([]models.Venue, error) {
	args := m.Called()
	return args.Get(0).([]models.Venue), args.Error(1)
}

func (m *MockVenueRepository) FindByID(id int) (models.Venue, error) {
	args := m.Called(id)
	return args.Get(0).(models.Venue), args.Error(1)
}

func (m *MockVenueRepository) Update(id int, venue models.Venue) (models.Venue, error) {
	args := m.Called(id, venue)
	return args.Get(0).(models.Venue), args.Error(1)
}

type MockRentalRepository struct {
	mock.Mock
}

func (m *MockRentalRepository) HasScheduleConflict(venueID int, rentalDate string, startTime string, endTime string) (bool, error) {
	args := m.Called(venueID, rentalDate, startTime, endTime)
	return args.Bool(0), args.Error(1)
}

func (m *MockRentalRepository) CreateWithDepositDeduction(rental models.Rental) (models.Rental, int, error) {
	args := m.Called(rental)
	return args.Get(0).(models.Rental), args.Int(1), args.Error(2)
}

func (m *MockRentalRepository) FindByUserID(userID int) ([]models.RentalHistory, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.RentalHistory), args.Error(1)
}

func TestCreateRentalSuccess(t *testing.T) {
	userRepo := new(MockUserRepository)
	venueRepo := new(MockVenueRepository)
	rentalRepo := new(MockRentalRepository)

	rentalService := services.NewRentalService(userRepo, venueRepo, rentalRepo)

	user := models.User{
		ID:            1,
		Email:         "user@gmail.com",
		DepositAmount: 300000,
		Role:          "user",
	}

	venue := models.Venue{
		ID:                1,
		Name:              "Lapangan Futsal Garuda",
		Category:          "Futsal",
		Location:          "Medan",
		StockAvailability: 1,
		RentalCost:        100000,
	}

	request := dto.RentalRequest{
		VenueID:    1,
		RentalDate: "2026-06-20",
		StartTime:  "19:00",
		EndTime:    "21:00",
	}

	expectedRental := models.Rental{
		ID:         1,
		UserID:     1,
		VenueID:    1,
		RentalDate: "2026-06-20",
		StartTime:  "19:00",
		EndTime:    "21:00",
		TotalCost:  200000,
		Status:     "booked",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	userRepo.On("FindByID", 1).Return(user, nil)
	venueRepo.On("FindByID", 1).Return(venue, nil)
	rentalRepo.On("HasScheduleConflict", 1, "2026-06-20", "19:00", "21:00").Return(false, nil)
	rentalRepo.On("CreateWithDepositDeduction", mock.Anything).Return(expectedRental, 100000, nil)

	result, err := rentalService.CreateRental(1, request)

	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Lapangan Futsal Garuda", result.VenueName)
	assert.Equal(t, 200000, result.TotalCost)
	assert.Equal(t, 100000, result.RemainingDeposit)
	assert.Equal(t, "booked", result.Status)

	userRepo.AssertExpectations(t)
	venueRepo.AssertExpectations(t)
	rentalRepo.AssertExpectations(t)
}

func TestCreateRentalInsufficientDeposit(t *testing.T) {
	userRepo := new(MockUserRepository)
	venueRepo := new(MockVenueRepository)
	rentalRepo := new(MockRentalRepository)

	rentalService := services.NewRentalService(userRepo, venueRepo, rentalRepo)

	user := models.User{
		ID:            1,
		Email:         "user@gmail.com",
		DepositAmount: 50000,
		Role:          "user",
	}

	venue := models.Venue{
		ID:                1,
		Name:              "Lapangan Futsal Garuda",
		StockAvailability: 1,
		RentalCost:        100000,
	}

	request := dto.RentalRequest{
		VenueID:    1,
		RentalDate: "2026-06-20",
		StartTime:  "19:00",
		EndTime:    "21:00",
	}

	userRepo.On("FindByID", 1).Return(user, nil)
	venueRepo.On("FindByID", 1).Return(venue, nil)
	rentalRepo.On("HasScheduleConflict", 1, "2026-06-20", "19:00", "21:00").Return(false, nil)

	result, err := rentalService.CreateRental(1, request)

	assert.Error(t, err)
	assert.Equal(t, "insufficient deposit amount", err.Error())
	assert.Equal(t, dto.RentalResponse{}, result)

	userRepo.AssertExpectations(t)
	venueRepo.AssertExpectations(t)
	rentalRepo.AssertExpectations(t)
	rentalRepo.AssertNotCalled(t, "CreateWithDepositDeduction", mock.Anything)
}

func TestCreateRentalScheduleConflict(t *testing.T) {
	userRepo := new(MockUserRepository)
	venueRepo := new(MockVenueRepository)
	rentalRepo := new(MockRentalRepository)

	rentalService := services.NewRentalService(userRepo, venueRepo, rentalRepo)

	user := models.User{
		ID:            1,
		Email:         "user@gmail.com",
		DepositAmount: 300000,
		Role:          "user",
	}

	venue := models.Venue{
		ID:                1,
		Name:              "Lapangan Futsal Garuda",
		StockAvailability: 1,
		RentalCost:        100000,
	}

	request := dto.RentalRequest{
		VenueID:    1,
		RentalDate: "2026-06-20",
		StartTime:  "19:00",
		EndTime:    "21:00",
	}

	userRepo.On("FindByID", 1).Return(user, nil)
	venueRepo.On("FindByID", 1).Return(venue, nil)
	rentalRepo.On("HasScheduleConflict", 1, "2026-06-20", "19:00", "21:00").Return(true, nil)

	result, err := rentalService.CreateRental(1, request)

	assert.Error(t, err)
	assert.Equal(t, "venue already booked at this time", err.Error())
	assert.Equal(t, dto.RentalResponse{}, result)

	userRepo.AssertExpectations(t)
	venueRepo.AssertExpectations(t)
	rentalRepo.AssertExpectations(t)
	rentalRepo.AssertNotCalled(t, "CreateWithDepositDeduction", mock.Anything)
}

func TestCreateRentalInvalidDuration(t *testing.T) {
	userRepo := new(MockUserRepository)
	venueRepo := new(MockVenueRepository)
	rentalRepo := new(MockRentalRepository)

	rentalService := services.NewRentalService(userRepo, venueRepo, rentalRepo)

	user := models.User{
		ID:            1,
		Email:         "user@gmail.com",
		DepositAmount: 300000,
		Role:          "user",
	}

	venue := models.Venue{
		ID:                1,
		Name:              "Lapangan Futsal Garuda",
		StockAvailability: 1,
		RentalCost:        100000,
	}

	request := dto.RentalRequest{
		VenueID:    1,
		RentalDate: "2026-06-20",
		StartTime:  "21:00",
		EndTime:    "19:00",
	}

	userRepo.On("FindByID", 1).Return(user, nil)
	venueRepo.On("FindByID", 1).Return(venue, nil)

	result, err := rentalService.CreateRental(1, request)

	assert.Error(t, err)
	assert.Equal(t, "end time must be after start time", err.Error())
	assert.Equal(t, dto.RentalResponse{}, result)

	userRepo.AssertExpectations(t)
	venueRepo.AssertExpectations(t)
	rentalRepo.AssertNotCalled(t, "HasScheduleConflict", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rentalRepo.AssertNotCalled(t, "CreateWithDepositDeduction", mock.Anything)
}

func TestCreateRentalVenueNotFound(t *testing.T) {
	userRepo := new(MockUserRepository)
	venueRepo := new(MockVenueRepository)
	rentalRepo := new(MockRentalRepository)

	rentalService := services.NewRentalService(userRepo, venueRepo, rentalRepo)

	user := models.User{
		ID:            1,
		Email:         "user@gmail.com",
		DepositAmount: 300000,
		Role:          "user",
	}

	request := dto.RentalRequest{
		VenueID:    99,
		RentalDate: "2026-06-20",
		StartTime:  "19:00",
		EndTime:    "21:00",
	}

	userRepo.On("FindByID", 1).Return(user, nil)
	venueRepo.On("FindByID", 99).Return(models.Venue{}, errors.New("sql: no rows in result set"))

	result, err := rentalService.CreateRental(1, request)

	assert.Error(t, err)
	assert.Equal(t, "internal server error", err.Error())
	assert.Equal(t, dto.RentalResponse{}, result)

	userRepo.AssertExpectations(t)
	venueRepo.AssertExpectations(t)
}
