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

// Placeholder for inserting new currency into the database
func (c *CurrencyModel) Insert(currency *Currency) error {
	return nil
}

// Placeholder for fetching currency from the database
func (c *CurrencyModel) Get(currency *Currency) error {
	return nil
}
