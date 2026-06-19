package services

import (
	"database/sql"
	"errors"
	"log"

	"sport-venue-rental-api/dto"
	"sport-venue-rental-api/models"
	"sport-venue-rental-api/repositories"
)

type VenueService interface {
	Create(request dto.VenueRequest) (models.Venue, error)
	GetAll() ([]models.Venue, error)
	GetByID(id int) (models.Venue, error)
	Update(id int, request dto.VenueRequest) (models.Venue, error)
}

type venueService struct {
	venueRepository repositories.VenueRepository
}

func NewVenueService(venueRepository repositories.VenueRepository) VenueService {
	return &venueService{venueRepository: venueRepository}
}

func (s *venueService) Create(request dto.VenueRequest) (models.Venue, error) {
	venue := models.Venue{
		Name:              request.Name,
		Category:          request.Category,
		Location:          request.Location,
		StockAvailability: request.StockAvailability,
		RentalCost:        request.RentalCost,
	}

	createdVenue, err := s.venueRepository.Create(venue)
	if err != nil {
		log.Printf("[ERROR] failed to create venue: %v", err)
		return models.Venue{}, errors.New("internal server error")
	}

	return createdVenue, nil
}

func (s *venueService) GetAll() ([]models.Venue, error) {
	venues, err := s.venueRepository.FindAll()
	if err != nil {
		log.Printf("[ERROR] failed to get venues: %v", err)
		return nil, errors.New("internal server error")
	}

	return venues, nil
}

func (s *venueService) GetByID(id int) (models.Venue, error) {
	venue, err := s.venueRepository.FindByID(id)
	if err == sql.ErrNoRows {
		return models.Venue{}, errors.New("venue not found")
	}

	if err != nil {
		log.Printf("[ERROR] failed to get venue by id: %v", err)
		return models.Venue{}, errors.New("internal server error")
	}

	return venue, nil
}

func (s *venueService) Update(id int, request dto.VenueRequest) (models.Venue, error) {
	_, err := s.venueRepository.FindByID(id)
	if err == sql.ErrNoRows {
		return models.Venue{}, errors.New("venue not found")
	}

	if err != nil {
		log.Printf("[ERROR] failed to check venue before update: %v", err)
		return models.Venue{}, errors.New("internal server error")
	}

	venue := models.Venue{
		Name:              request.Name,
		Category:          request.Category,
		Location:          request.Location,
		StockAvailability: request.StockAvailability,
		RentalCost:        request.RentalCost,
	}

	updatedVenue, err := s.venueRepository.Update(id, venue)
	if err != nil {
		log.Printf("[ERROR] failed to update venue: %v", err)
		return models.Venue{}, errors.New("internal server error")
	}

	return updatedVenue, nil
}
