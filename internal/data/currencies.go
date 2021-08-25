package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Currency struct {
	Id        int64     `json:"id"`
	Code      string    `json:"code"`
	Rate      float64   `json:"rate"` // TODO arbitrary precision e.g github.com/shopspring/decimal
	CreatedAt time.Time `json:"created_at"`
}

// Define a CurrencyModel struct type which wraps a sql.DB connection pool.
type CurrencyModel struct {
	DB *sql.DB
}

// Insert for inserting new currency into the database.
func (c *CurrencyModel) Insert(currency *Currency) error {
	// Define the SQL query for inserting a new record in the currencies table and
	// returning the system-generated data.
	query := `
		INSERT INTO currencies (code, rate) VALUES ($1, $2)
		RETURNING id, created_at`
	// Use the QueryRow() method to execute the SQL query on our connection pool
	// passing in the currency code as a parameter and scanning the system-generated
	// id, and created_at values into the currency struct.
	return c.DB.QueryRow(
		query, currency.Code, currency.Rate,
	).Scan(
		&currency.Id, &currency.CreatedAt,
	)
}

// Get for fetching currency with the given id from the database.
func (c CurrencyModel) Get(id int64) (*Currency, error) {
	// The Id is always >= 1
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the currency data.
	query := `SELECT id, code, rate, created_at FROM currencies WHERE id = $1`
	// Declare a Movie struct to hold the data returned by the query.
	var currency Currency
	err := c.DB.QueryRow(
		query, id,
	).Scan(
		&currency.Id, &currency.Code, currency.Rate, &currency.CreatedAt,
	)
	// Handle any errors. If there was no matching currency found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &currency, nil
}

// GetLatest method which returns a slice of latest rates available in the database
// and the date of the rates. It is neccessary because the date when the request is
// made does not always match the date of the data (weekends, holidays, etc).
func (c CurrencyModel) GetLatest() ([]*Currency, string, error) {
	// Construct the SQL query to retrieve all movie records.
	query := `
		SELECT id, code, rate, created_at FROM currencies
		WHERE DATE(created_at) = (SELECT MAX(DATE(created_at)) FROM currencies)
		ORDER BY code`
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := c.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, "", err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetLatest() returns.
	defer rows.Close()
	// Initialize an empty slice to hold the currency data.
	currencies := []*Currency{}
	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Currency struct to hold the data for an individual
		// currency.
		var currency Currency
		// Scan the values from the row into the Currency struct.
		err := rows.Scan(
			&currency.Id,
			&currency.Code,
			&currency.Rate,
			&currency.CreatedAt,
		)
		if err != nil {
			return nil, "", err
		}
		// Add the Currency struct to the slice.
		currencies = append(currencies, &currency)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, "", err
	}
	date := currencies[0].CreatedAt.Format("2006-01-02")
	// If everything went OK, then return the slice of currencies.
	return currencies, date, nil
}
