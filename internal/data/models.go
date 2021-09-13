package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Create a Models struct which wraps the CurrencyModel.
type Models struct {
	Currencies CurrencyModel
}

// For ease of use, we also add a New() method which returns a Models
// struct containing the initialized models.
func NewModels(db *sql.DB) Models {
	return Models{
		Currencies: CurrencyModel{DB: db},
	}
}
