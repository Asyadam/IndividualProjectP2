package services

import (
	"database/sql"
	"errors"
	"log"

	"sport-venue-rental-api/dto"
	"sport-venue-rental-api/models"
	"sport-venue-rental-api/repositories"
	"sport-venue-rental-api/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(request dto.RegisterRequest) (models.User, error)
	Login(request dto.LoginRequest) (dto.LoginResponse, error)
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{userRepository: userRepository}
}

func (s *authService) Register(request dto.RegisterRequest) (models.User, error) {
	_, err := s.userRepository.FindByEmail(request.Email)
	if err == nil {
		return models.User{}, errors.New("email already registered")
	}

	if err != sql.ErrNoRows {
		log.Printf("[ERROR] failed to check email: %v", err)
		return models.User{}, errors.New("internal server error")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[ERROR] failed to hash password: %v", err)
		return models.User{}, errors.New("internal server error")
	}

	user := models.User{
		Username:      request.Username,
		Email:         request.Email,
		Password:      string(hashedPassword),
		DepositAmount: 0,
		Role:          "user",
	}

	createdUser, err := s.userRepository.Create(user)
	if err != nil {
		log.Printf("[ERROR] failed to create user: %v", err)
		return models.User{}, errors.New("internal server error")
	}

	return createdUser, nil
}

func (s *authService) Login(request dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := s.userRepository.FindByEmail(request.Email)
	if err == sql.ErrNoRows {
		return dto.LoginResponse{}, errors.New("invalid email or password")
	}

	if err != nil {
		log.Printf("[ERROR] failed to find user by email: %v", err)
		return dto.LoginResponse{}, errors.New("internal server error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return dto.LoginResponse{}, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		log.Printf("[ERROR] failed to generate token: %v", err)
		return dto.LoginResponse{}, errors.New("internal server error")
	}

	return dto.LoginResponse{
		Token: token,
	}, nil
}
