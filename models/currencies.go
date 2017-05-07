package models

import (
	"database/sql"
	"strings"
	"time"
)

// Currencies type
type Currencies struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float32 `json:"rates"`
}

// LatestRates query db for most updated results for each currency
func LatestRates() (*Currencies, error) {
	rows, err := db.Query("SELECT currency, rate, ratedate FROM rates WHERE ratedate = (SELECT max(ratedate) FROM rates)")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	currencies := &Currencies{Rates: make(map[string]float32)}
	for rows.Next() {
		var (
			currency string
			rate     float32
			date     time.Time
		)
		if err := rows.Scan(&currency, &rate, &date); err != nil {
			return nil, err
		}
		currencies.Date = date.Format("2006-01-02")
		currencies.Rates[currency] = rate
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	currencies.Base = "EUR"

	return currencies, nil
}

// FilteredRates query db for most updated results for specified currencies (symbols)
func FilteredRates(symbols string) (*Currencies, error) {
	rows, err := db.Query(`SELECT currency, rate, ratedate 
	FROM rates WHERE currency = ANY(string_to_array($1, ',')) 
	AND ratedate = (SELECT max(ratedate) FROM rates)`, symbols)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	currencies := &Currencies{Rates: make(map[string]float32)}
	for rows.Next() {
		var (
			currency string
			rate     float32
			date     time.Time
		)
		if err := rows.Scan(&currency, &rate, &date); err != nil {
			return nil, err
		}
		currencies.Date = date.Format("2006-01-02")
		currencies.Rates[currency] = rate
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	currencies.Base = "EUR"

	return currencies, nil
}

// RecalculatedRates query db for recalculated rates (base)
func RecalculatedRates(base string) (*Currencies, error) {
	rows, err := db.Query(`SELECT currency, 
	rate / (SELECT rate FROM rates WHERE ratedate = (SELECT max(ratedate) FROM rates) 
	AND currency = $1) AS rate, ratedate FROM rates WHERE currency !=$1 
	AND ratedate = (SELECT max(ratedate) FROM rates) 
	UNION ALL SELECT 'EUR', 1 / (SELECT rate FROM rates 
	WHERE ratedate = (SELECT max(ratedate) FROM rates) AND currency = $1), 
	(SELECT max(ratedate) FROM rates)`, base)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	currencies := &Currencies{Rates: make(map[string]float32)}
	for rows.Next() {
		var (
			currency string
			rate     float32
			date     time.Time
		)
		if err := rows.Scan(&currency, &rate, &date); err != nil {
			return nil, err
		}
		currencies.Date = date.Format("2006-01-02")
		currencies.Rates[currency] = rate
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	currencies.Base = base

	return currencies, nil
}

// RecalculatedAndFilteredRates query db for recalculated and filtered rates (base, symbols)
func RecalculatedAndFilteredRates(base, symbols string) (*Currencies, error) {
	var rows *sql.Rows
	var err error
	if strings.Contains(symbols, "EUR") {
		rows, err = db.Query(`SELECT currency, 
		rate / (SELECT rate FROM rates WHERE ratedate = (SELECT max(ratedate) FROM rates) 
		AND currency = $1) AS rate, ratedate FROM rates WHERE currency !=$1 
		AND currency = ANY(string_to_array($2, ',')) 
		AND ratedate = (SELECT max(ratedate) FROM rates) 
		UNION ALL SELECT 'EUR', 1 / (SELECT rate FROM rates 
		WHERE ratedate = (SELECT max(ratedate) FROM rates) 
		AND currency = $1), (SELECT max(ratedate) FROM rates)`, base, symbols)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = db.Query(`SELECT currency, 
		rate / (SELECT rate FROM rates WHERE ratedate = (SELECT max(ratedate) FROM rates) 
		AND currency = $1) AS rate, ratedate FROM rates WHERE currency !=$1 
		AND currency = ANY(string_to_array($2, ',')) 
		AND ratedate = (SELECT max(ratedate) FROM rates)`, base, symbols)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	currencies := &Currencies{Rates: make(map[string]float32)}
	for rows.Next() {
		var (
			currency string
			rate     float32
			date     time.Time
		)
		if err := rows.Scan(&currency, &rate, &date); err != nil {
			return nil, err
		}
		currencies.Date = date.Format("2006-01-02")
		currencies.Rates[currency] = rate
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	currencies.Base = base

	return currencies, nil
}
