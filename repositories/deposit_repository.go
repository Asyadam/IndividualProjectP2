package repositories

import (
	"database/sql"

	"sport-venue-rental-api/models"
)

type DepositRepository interface {
	Create(deposit models.DepositHistory) (models.DepositHistory, error)
	FindByExternalID(externalID string) (models.DepositHistory, error)
	MarkAsPaid(externalID string) (models.DepositHistory, error)
}

type depositRepository struct {
	db *sql.DB
}

func NewDepositRepository(db *sql.DB) DepositRepository {
	return &depositRepository{db: db}
}

func (r *depositRepository) Create(deposit models.DepositHistory) (models.DepositHistory, error) {
	query := `
		INSERT INTO deposit_histories (user_id, amount, payment_link, payment_status, external_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, amount, payment_link, payment_status, external_id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		deposit.UserID,
		deposit.Amount,
		deposit.PaymentLink,
		deposit.PaymentStatus,
		deposit.ExternalID,
	).Scan(
		&deposit.ID,
		&deposit.UserID,
		&deposit.Amount,
		&deposit.PaymentLink,
		&deposit.PaymentStatus,
		&deposit.ExternalID,
		&deposit.CreatedAt,
		&deposit.UpdatedAt,
	)

	if err != nil {
		return models.DepositHistory{}, err
	}

	return deposit, nil
}

func (r *depositRepository) FindByExternalID(externalID string) (models.DepositHistory, error) {
	var deposit models.DepositHistory

	query := `
		SELECT id, user_id, amount, payment_link, payment_status, external_id, created_at, updated_at
		FROM deposit_histories
		WHERE external_id = $1
	`

	err := r.db.QueryRow(query, externalID).Scan(
		&deposit.ID,
		&deposit.UserID,
		&deposit.Amount,
		&deposit.PaymentLink,
		&deposit.PaymentStatus,
		&deposit.ExternalID,
		&deposit.CreatedAt,
		&deposit.UpdatedAt,
	)

	if err != nil {
		return models.DepositHistory{}, err
	}

	return deposit, nil
}

func (r *depositRepository) MarkAsPaid(externalID string) (models.DepositHistory, error) {
	var deposit models.DepositHistory

	query := `
		UPDATE deposit_histories
		SET payment_status = 'paid',
		    updated_at = CURRENT_TIMESTAMP
		WHERE external_id = $1
		RETURNING id, user_id, amount, payment_link, payment_status, external_id, created_at, updated_at
	`

	err := r.db.QueryRow(query, externalID).Scan(
		&deposit.ID,
		&deposit.UserID,
		&deposit.Amount,
		&deposit.PaymentLink,
		&deposit.PaymentStatus,
		&deposit.ExternalID,
		&deposit.CreatedAt,
		&deposit.UpdatedAt,
	)

	if err != nil {
		return models.DepositHistory{}, err
	}

	return deposit, nil
}
