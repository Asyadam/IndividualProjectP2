package repositories

import (
	"database/sql"

	"sport-venue-rental-api/models"
)

type RentalRepository interface {
	HasScheduleConflict(venueID int, rentalDate string, startTime string, endTime string) (bool, error)
	CreateWithDepositDeduction(rental models.Rental) (models.Rental, int, error)
	FindByUserID(userID int) ([]models.RentalHistory, error)
}

type rentalRepository struct {
	db *sql.DB
}

func NewRentalRepository(db *sql.DB) RentalRepository {
	return &rentalRepository{db: db}
}

func (r *rentalRepository) HasScheduleConflict(venueID int, rentalDate string, startTime string, endTime string) (bool, error) {
	var exists bool

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM rentals
			WHERE venue_id = $1
			  AND rental_date = $2::date
			  AND status = 'booked'
			  AND start_time < $4::time
			  AND end_time > $3::time
		)
	`

	err := r.db.QueryRow(query, venueID, rentalDate, startTime, endTime).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *rentalRepository) CreateWithDepositDeduction(rental models.Rental) (models.Rental, int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.Rental{}, 0, err
	}
	defer tx.Rollback()

	var remainingDeposit int

	deductQuery := `
		UPDATE users
		SET deposit_amount = deposit_amount - $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		  AND deposit_amount >= $1
		RETURNING deposit_amount
	`

	err = tx.QueryRow(deductQuery, rental.TotalCost, rental.UserID).Scan(&remainingDeposit)
	if err != nil {
		return models.Rental{}, 0, err
	}

	createRentalQuery := `
		INSERT INTO rentals (user_id, venue_id, rental_date, start_time, end_time, total_cost, status)
		VALUES ($1, $2, $3::date, $4::time, $5::time, $6, $7)
		RETURNING 
			id,
			user_id,
			venue_id,
			to_char(rental_date, 'YYYY-MM-DD'),
			to_char(start_time, 'HH24:MI'),
			to_char(end_time, 'HH24:MI'),
			total_cost,
			status,
			created_at,
			updated_at
	`

	err = tx.QueryRow(
		createRentalQuery,
		rental.UserID,
		rental.VenueID,
		rental.RentalDate,
		rental.StartTime,
		rental.EndTime,
		rental.TotalCost,
		rental.Status,
	).Scan(
		&rental.ID,
		&rental.UserID,
		&rental.VenueID,
		&rental.RentalDate,
		&rental.StartTime,
		&rental.EndTime,
		&rental.TotalCost,
		&rental.Status,
		&rental.CreatedAt,
		&rental.UpdatedAt,
	)

	if err != nil {
		return models.Rental{}, 0, err
	}

	err = tx.Commit()
	if err != nil {
		return models.Rental{}, 0, err
	}

	return rental, remainingDeposit, nil
}

func (r *rentalRepository) FindByUserID(userID int) ([]models.RentalHistory, error) {
	query := `
		SELECT
			r.id,
			r.venue_id,
			v.name,
			to_char(r.rental_date, 'YYYY-MM-DD'),
			to_char(r.start_time, 'HH24:MI'),
			to_char(r.end_time, 'HH24:MI'),
			r.total_cost,
			r.status,
			r.created_at
		FROM rentals r
		JOIN venues v ON v.id = r.venue_id
		WHERE r.user_id = $1
		ORDER BY r.id DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rentals []models.RentalHistory

	for rows.Next() {
		var rental models.RentalHistory

		err := rows.Scan(
			&rental.ID,
			&rental.VenueID,
			&rental.VenueName,
			&rental.RentalDate,
			&rental.StartTime,
			&rental.EndTime,
			&rental.TotalCost,
			&rental.Status,
			&rental.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		rentals = append(rentals, rental)
	}

	return rentals, nil
}
