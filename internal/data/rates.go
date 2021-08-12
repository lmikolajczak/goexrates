package data

import (
	"database/sql"
	"time"
)

type Rate struct {
	Id         int64     `json:"id"`
	CurrencyId int64     `json:"currency_id"`
	Rate       float64   `json:"rate"` // TODO arbitrary precision e.g github.com/shopspring/decimal
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// Define a RateModel struct type which wraps a sql.DB connection pool.
type RateModel struct {
	DB *sql.DB
}

// Insert for inserting new rate into the database
func (r *RateModel) Insert(rate *Rate) error {
	// Define the SQL query for inserting a new record in the rates table and
	// returning the system-generated data.
	query := `
		INSERT INTO rates (currency_id, rate) VALUES ($1, $2)
		RETURNING id, updated_at, created_at`
	// Use the QueryRow() method to execute the SQL query on our connection pool
	// passing in the rate code and rate as a parameters and scanning the system-
	// generated id, updated_at and created_at values into the rate struct.
	return r.DB.QueryRow(
		query, rate.CurrencyId, rate.Rate,
	).Scan(
		&rate.Id, &rate.UpdatedAt, &rate.CreatedAt,
	)
}

// Placeholder for fetching rate from the database
func (r *RateModel) Get(rate *Rate) error {
	return nil
}
