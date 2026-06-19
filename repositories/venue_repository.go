package repositories

import (
	"database/sql"

	"sport-venue-rental-api/models"
)

type VenueRepository interface {
	Create(venue models.Venue) (models.Venue, error)
	FindAll() ([]models.Venue, error)
	FindByID(id int) (models.Venue, error)
	Update(id int, venue models.Venue) (models.Venue, error)
}

type venueRepository struct {
	db *sql.DB
}

func NewVenueRepository(db *sql.DB) VenueRepository {
	return &venueRepository{db: db}
}

func (r *venueRepository) Create(venue models.Venue) (models.Venue, error) {
	query := `
		INSERT INTO venues (name, category, location, stock_availability, rental_cost)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, category, location, stock_availability, rental_cost, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		venue.Name,
		venue.Category,
		venue.Location,
		venue.StockAvailability,
		venue.RentalCost,
	).Scan(
		&venue.ID,
		&venue.Name,
		&venue.Category,
		&venue.Location,
		&venue.StockAvailability,
		&venue.RentalCost,
		&venue.CreatedAt,
		&venue.UpdatedAt,
	)

	if err != nil {
		return models.Venue{}, err
	}

	return venue, nil
}

func (r *venueRepository) FindAll() ([]models.Venue, error) {
	query := `
		SELECT id, name, category, location, stock_availability, rental_cost, created_at, updated_at
		FROM venues
		ORDER BY id ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var venues []models.Venue

	for rows.Next() {
		var venue models.Venue

		err := rows.Scan(
			&venue.ID,
			&venue.Name,
			&venue.Category,
			&venue.Location,
			&venue.StockAvailability,
			&venue.RentalCost,
			&venue.CreatedAt,
			&venue.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		venues = append(venues, venue)
	}

	return venues, nil
}

func (r *venueRepository) FindByID(id int) (models.Venue, error) {
	var venue models.Venue

	query := `
		SELECT id, name, category, location, stock_availability, rental_cost, created_at, updated_at
		FROM venues
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&venue.ID,
		&venue.Name,
		&venue.Category,
		&venue.Location,
		&venue.StockAvailability,
		&venue.RentalCost,
		&venue.CreatedAt,
		&venue.UpdatedAt,
	)

	if err != nil {
		return models.Venue{}, err
	}

	return venue, nil
}

func (r *venueRepository) Update(id int, venue models.Venue) (models.Venue, error) {
	query := `
		UPDATE venues
		SET name = $1,
		    category = $2,
		    location = $3,
		    stock_availability = $4,
		    rental_cost = $5,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING id, name, category, location, stock_availability, rental_cost, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		venue.Name,
		venue.Category,
		venue.Location,
		venue.StockAvailability,
		venue.RentalCost,
		id,
	).Scan(
		&venue.ID,
		&venue.Name,
		&venue.Category,
		&venue.Location,
		&venue.StockAvailability,
		&venue.RentalCost,
		&venue.CreatedAt,
		&venue.UpdatedAt,
	)

	if err != nil {
		return models.Venue{}, err
	}

	return venue, nil
}
