package models

import (
	"database/sql"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Currencies from database
type Currencies struct {
	Base  string                     `json:"base"`
	Date  string                     `json:"date"`
	Rates map[string]decimal.Decimal `json:"rates"`
}

// LatestRates query db for most updated results for each iso_code
// supports different queries based on request parameters (base, symbols)
func LatestRates(baseParam, symbolsParam string) (*Currencies, error) {
	base := strings.ToUpper(baseParam)
	symbols := strings.ToUpper(symbolsParam)
	if base == "EUR" {
		base = ""
	}

	var rows *sql.Rows
	var err error

	switch {
	case base != "" && symbols != "":
		if strings.Contains(symbols, "EUR") {
			rows, err = db.Query(`SELECT iso_code, 
			rate / (SELECT rate FROM currencies WHERE date = (SELECT max(date) FROM currencies) 
			AND iso_code = $1) AS rate FROM currencies WHERE iso_code !=$1 
			AND iso_code = ANY(string_to_array($2, ',')) 
			AND date = (SELECT max(date) FROM currencies) 
			UNION ALL SELECT 'EUR', 1 / (SELECT rate FROM currencies 
			WHERE date = (SELECT max(date) FROM currencies) 
			AND iso_code = $1)`, base, symbols)
			if err != nil {
				return nil, err
			}
		} else {
			rows, err = db.Query(`SELECT iso_code, 
			rate / (SELECT rate FROM currencies WHERE date = (SELECT max(date) FROM currencies) 
			AND iso_code = $1) AS rate FROM currencies WHERE iso_code !=$1 
			AND iso_code = ANY(string_to_array($2, ',')) 
			AND date = (SELECT max(date) FROM currencies)`, base, symbols)
			if err != nil {
				return nil, err
			}
		}
	case base != "" && symbols == "":
		rows, err = db.Query(`SELECT iso_code, 
		rate / (SELECT rate FROM currencies WHERE date = (SELECT max(date) FROM currencies) 
		AND iso_code = $1) AS rate FROM currencies WHERE iso_code !=$1 
		AND date = (SELECT max(date) FROM currencies) 
		UNION ALL SELECT 'EUR', 1 / (SELECT rate FROM currencies 
		WHERE date = (SELECT max(date) FROM currencies) AND iso_code = $1)`, base)
		if err != nil {
			return nil, err
		}
	case base == "" && symbols != "":
		rows, err = db.Query(`SELECT iso_code, rate 
		FROM currencies WHERE iso_code = ANY(string_to_array($1, ',')) 
		AND date = (SELECT max(date) FROM currencies)`, symbols)
		if err != nil {
			return nil, err
		}
	default:
		rows, err = db.Query(`SELECT iso_code, rate FROM currencies 
		WHERE date = (SELECT max(date) FROM currencies)`)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	currencies := &Currencies{Rates: make(map[string]decimal.Decimal)}
	for rows.Next() {
		var (
			isoCode string
			rate    decimal.Decimal
		)
		if err := rows.Scan(&isoCode, &rate); err != nil {
			return nil, err
		}
		currencies.Rates[isoCode] = rate.Round(5)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	var date time.Time
	err = db.QueryRow("SELECT max(date) FROM currencies").Scan(&date)
	currencies.Date = date.Format("2006-01-02")

	if base != "" {
		currencies.Base = base
	} else {
		currencies.Base = "EUR"
	}

	return currencies, nil
}

// HistoricalRates query db for historical results for each iso_code
// supports different queries based on request parameters (base, symbols, date)
func HistoricalRates(baseParam, symbolsParam, dateParam string) (*Currencies, error) {
	base := strings.ToUpper(baseParam)
	symbols := strings.ToUpper(symbolsParam)
	date := dateParam
	if base == "EUR" {
		base = ""
	}

	var rows *sql.Rows
	var err error

	switch {
	case base != "" && symbols != "":
		if strings.Contains(symbols, "EUR") {
			rows, err = db.Query(`SELECT iso_code, 
			rate / (SELECT rate FROM currencies 
			WHERE date = (SELECT date FROM currencies WHERE date <= $3 
			ORDER BY date DESC LIMIT 1) 
			AND iso_code = $1) AS rate FROM currencies WHERE iso_code !=$1 
			AND iso_code = ANY(string_to_array($2, ',')) 
			AND date = (SELECT date FROM currencies WHERE date <= $3 
			ORDER BY date DESC LIMIT 1) 
			UNION ALL SELECT 'EUR', 1 / (SELECT rate FROM currencies 
			WHERE date = (SELECT date FROM currencies WHERE date <= $3 
			ORDER BY date DESC LIMIT 1) 
			AND iso_code = $1)`, base, symbols, date)
			if err != nil {
				return nil, err
			}
		} else {
			rows, err = db.Query(`SELECT iso_code, 
			rate / (SELECT rate FROM currencies 
			WHERE date = (SELECT date FROM currencies WHERE date <= $3 
			ORDER BY date DESC LIMIT 1) 
			AND iso_code = $1) AS rate FROM currencies WHERE iso_code !=$1 
			AND iso_code = ANY(string_to_array($2, ',')) 
			AND date = (SELECT date FROM currencies 
			WHERE date <= $3 ORDER BY date DESC LIMIT 1)`, base, symbols, date)
			if err != nil {
				return nil, err
			}
		}
	case base != "" && symbols == "":
		rows, err = db.Query(`SELECT iso_code, 
		rate / (SELECT rate FROM currencies 
		WHERE date = (SELECT date FROM currencies WHERE date <= $2 
		ORDER BY date DESC LIMIT 1) 
		AND iso_code = $1) AS rate FROM currencies WHERE iso_code !=$1 
		AND date = (SELECT date FROM currencies WHERE date <= $2 
		ORDER BY date DESC LIMIT 1) 
		UNION ALL SELECT 'EUR', 1 / (SELECT rate FROM currencies 
		WHERE date = (SELECT date FROM currencies 
		WHERE date <= $2 ORDER BY date DESC LIMIT 1) AND iso_code = $1)`, base, date)
		if err != nil {
			return nil, err
		}
	case base == "" && symbols != "":
		rows, err = db.Query(`SELECT iso_code, rate 
		FROM currencies WHERE iso_code = ANY(string_to_array($1, ',')) 
		AND date = (SELECT date FROM currencies WHERE date <= $2 
		ORDER BY date DESC LIMIT 1)`, symbols, date)
		if err != nil {
			return nil, err
		}
	default:
		rows, err = db.Query(`SELECT iso_code, rate FROM currencies 
		WHERE date = (SELECT date FROM currencies WHERE date <= $1 
		ORDER BY date DESC LIMIT 1)`,
			date)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	currencies := &Currencies{Rates: make(map[string]decimal.Decimal)}
	for rows.Next() {
		var (
			isoCode string
			rate    decimal.Decimal
		)
		if err := rows.Scan(&isoCode, &rate); err != nil {
			return nil, err
		}
		currencies.Rates[isoCode] = rate.Round(5)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	var availableDate time.Time
	err = db.QueryRow(`SELECT date FROM currencies 
	WHERE date <= $1 ORDER BY date DESC LIMIT 1`, date).Scan(&availableDate)
	currencies.Date = availableDate.Format("2006-01-02")

	if base != "" {
		currencies.Base = base
	} else {
		currencies.Base = "EUR"
	}

	return currencies, nil
}
