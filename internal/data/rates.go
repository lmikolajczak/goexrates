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

// Placeholder for inserting new rate into the database
func (r *RateModel) Insert(rate *Rate) error {
	return nil
}

// Placeholder for fetching rate from the database
func (r *RateModel) Get(rate *Rate) error {
	return nil
}
