package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"sport-venue-rental-api/dto"
	"sport-venue-rental-api/models"
	"sport-venue-rental-api/repositories"
)

type DepositService interface {
	CreateDeposit(userID int, request dto.DepositRequest) (models.DepositHistory, error)
	HandleXenditCallback(request dto.XenditCallbackRequest) (models.DepositHistory, error)
}

type depositService struct {
	userRepository    repositories.UserRepository
	depositRepository repositories.DepositRepository
	xenditService     XenditService
}

func NewDepositService(
	userRepository repositories.UserRepository,
	depositRepository repositories.DepositRepository,
	xenditService XenditService,
) DepositService {
	return &depositService{
		userRepository:    userRepository,
		depositRepository: depositRepository,
		xenditService:     xenditService,
	}
}

func (s *depositService) CreateDeposit(userID int, request dto.DepositRequest) (models.DepositHistory, error) {
	user, err := s.userRepository.FindByID(userID)
	if err == sql.ErrNoRows {
		return models.DepositHistory{}, errors.New("user not found")
	}

	if err != nil {
		log.Printf("[ERROR] failed to find user: %v", err)
		return models.DepositHistory{}, errors.New("internal server error")
	}

	externalID := fmt.Sprintf("deposit-%d-%d", userID, time.Now().Unix())

	paymentLink, err := s.xenditService.CreateInvoice(
		externalID,
		request.Amount,
		user.Email,
		"Top up deposit Sport Venue Rental",
	)
	if err != nil {
		log.Printf("[ERROR] failed to create xendit invoice: %v", err)
		return models.DepositHistory{}, errors.New("failed to create payment link")
	}

	deposit := models.DepositHistory{
		UserID:        userID,
		Amount:        request.Amount,
		PaymentLink:   paymentLink,
		PaymentStatus: "pending",
		ExternalID:    externalID,
	}

	createdDeposit, err := s.depositRepository.Create(deposit)
	if err != nil {
		log.Printf("[ERROR] failed to create deposit history: %v", err)
		return models.DepositHistory{}, errors.New("internal server error")
	}

	return createdDeposit, nil
}

func (s *depositService) HandleXenditCallback(request dto.XenditCallbackRequest) (models.DepositHistory, error) {
	if request.Status != "PAID" && request.Status != "paid" {
		return models.DepositHistory{}, errors.New("payment is not paid")
	}

	deposit, err := s.depositRepository.FindByExternalID(request.ExternalID)
	if err == sql.ErrNoRows {
		return models.DepositHistory{}, errors.New("deposit not found")
	}

	if err != nil {
		log.Printf("[ERROR] failed to find deposit by external id: %v", err)
		return models.DepositHistory{}, errors.New("internal server error")
	}

	if deposit.PaymentStatus == "paid" {
		return deposit, nil
	}

	if request.Amount != deposit.Amount {
		return models.DepositHistory{}, errors.New("amount mismatch")
	}

	paidDeposit, err := s.depositRepository.MarkAsPaid(request.ExternalID)
	if err != nil {
		log.Printf("[ERROR] failed to mark deposit as paid: %v", err)
		return models.DepositHistory{}, errors.New("internal server error")
	}

	_, err = s.userRepository.AddDeposit(deposit.UserID, deposit.Amount)
	if err != nil {
		log.Printf("[ERROR] failed to add user deposit: %v", err)
		return models.DepositHistory{}, errors.New("internal server error")
	}

	return paidDeposit, nil
}
