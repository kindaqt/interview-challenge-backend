package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Rental struct {
	ID              int       `db:"id"`
	UserID          int       `db:"user_id"`
	Name            string    `db:"name"`
	Type            string    `db:"type"`
	Description     string    `db:"description"`
	Sleeps          int       `db:"sleeps"`
	PricePerDay     int64     `db:"price_per_day"`
	HomeCity        string    `db:"home_city"`
	HomeState       string    `db:"home_state"`
	HomeZip         string    `db:"home_zip"`
	HomeCountry     string    `db:"home_country"`
	VehicleMake     string    `db:"vehicle_make"`
	VehicleModel    string    `db:"vehicle_model"`
	VehicleYear     int       `db:"vehicle_year"`
	VehicleLength   float64   `db:"vehicle_length"`
	Created         time.Time `db:"created"`
	Updated         time.Time `db:"updated"`
	Lat             float64   `db:"lat"`
	Lng             float64   `db:"lng"`
	PrimaryImageURL string    `db:"primary_image_url"`
}

func (d *Database) QueryRentalByID(tx *sql.Tx, rentalID int) (*Rental, error) {
	row := tx.QueryRow(`
		SELECT *
		FROM rentals
		WHERE id = $1
	`, rentalID)

	var rental Rental
	if err := row.Scan(
		&rental.ID,
		&rental.UserID,
		&rental.Name,
		&rental.Type,
		&rental.Description,
		&rental.Sleeps,
		&rental.PricePerDay,
		&rental.HomeCity,
		&rental.HomeState,
		&rental.HomeZip,
		&rental.HomeCountry,
		&rental.VehicleMake,
		&rental.VehicleModel,
		&rental.VehicleYear,
		&rental.VehicleLength,
		&rental.Created,
		&rental.Updated,
		&rental.Lat,
		&rental.Lng,
		&rental.PrimaryImageURL,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Rental not found
		}
		return nil, err
	}

	return &rental, nil
}

type QueryRentalsParams struct {
	PriceMin *int
	PriceMax *int
	Limit    *int
	Offset   *int
	IDs      string
	Near     string
	Sort     string
}

func (d *Database) QueryRentals(tx *sql.Tx, params *QueryRentalsParams) ([]*Rental, error) {
	// Construct the SQL query based on the provided parameters
	query := "SELECT * FROM rentals WHERE 1=1"

	if params.PriceMin != nil {
		query += fmt.Sprintf(" AND price_per_day >= %d", *params.PriceMin)
	}
	if params.PriceMax != nil {
		query += fmt.Sprintf(" AND price_per_day <= %d", *params.PriceMax)
	}
	if params.IDs != "" {
		query += fmt.Sprintf(" AND id IN (%s)", params.IDs)
	}
	if params.Near != "" {
		radius := 50
		query += fmt.Sprintf(`
			AND earth_box(ll_to_earth(%s), %d) @> ll_to_earth(lat, lng)
			AND earth_distance(ll_to_earth(%s), ll_to_earth(lat, lng)) <= %d * 1000
		`, params.Near, radius, params.Near, radius)
	}
	if params.Sort != "" {
		var value string

		if params.Sort == "price" {
			value = "price_per_day"
		} else {
			value = params.Sort
		}

		query += " ORDER BY " + value
	}
	if params.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *params.Limit)
	}
	if params.Offset != nil {
		query += fmt.Sprintf(" OFFSET %d", *params.Offset)
	}

	// Execute the query using the constructed SQL statement and arguments
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rentals []*Rental
	for rows.Next() {
		var rental Rental
		err := rows.Scan(
			&rental.ID,
			&rental.UserID,
			&rental.Name,
			&rental.Type,
			&rental.Description,
			&rental.Sleeps,
			&rental.PricePerDay,
			&rental.HomeCity,
			&rental.HomeState,
			&rental.HomeZip,
			&rental.HomeCountry,
			&rental.VehicleMake,
			&rental.VehicleModel,
			&rental.VehicleYear,
			&rental.VehicleLength,
			&rental.Created,
			&rental.Updated,
			&rental.Lat,
			&rental.Lng,
			&rental.PrimaryImageURL,
		)
		if err != nil {
			return nil, err
		}

		rentals = append(rentals, &rental)
	}

	return rentals, nil
}
