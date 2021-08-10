package data

import (
	"database/sql"
	"time"
)

type Currency struct {
	Id        int64     `json:"id"`
	Code      string    `json:"code"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// Define a CurrencyModel struct type which wraps a sql.DB connection pool.
type CurrencyModel struct {
	DB *sql.DB
}

// Insert for inserting new currency into the database
func (c *CurrencyModel) Insert(currency *Currency) error {
	// Define the SQL query for inserting a new record in the currencies table and
	// returning the system-generated data.
	query := `
		INSERT INTO currencies (code) VALUES ($1)
		RETURNING id, updated_at, created_at`
	// Use the QueryRow() method to execute the SQL query on our connection pool
	// passing in the currency code as a parameter and scanning the system-generated
	// id, updated_at and created_at values into the currency struct.
	return c.DB.QueryRow(
		query, currency.Code,
	).Scan(
		&currency.Id, &currency.UpdatedAt, &currency.CreatedAt,
	)
}

// Placeholder for fetching currency from the database
func (c *CurrencyModel) Get(currency *Currency) error {
	return nil
}
