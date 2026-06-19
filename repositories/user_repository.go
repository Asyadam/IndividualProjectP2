package repositories

import (
	"database/sql"

	"sport-venue-rental-api/models"
)

type UserRepository interface {
	Create(user models.User) (models.User, error)
	FindByEmail(email string) (models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user models.User) (models.User, error) {
	query := `
		INSERT INTO users (username, email, password, deposit_amount, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, email, deposit_amount, role, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.Password,
		user.DepositAmount,
		user.Role,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.DepositAmount,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *userRepository) FindByEmail(email string) (models.User, error) {
	var user models.User

	query := `
		SELECT id, username, email, password, deposit_amount, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.DepositAmount,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
