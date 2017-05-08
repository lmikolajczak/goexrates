package models

import (
	"database/sql"
	"strings"
	"time"
)

// Currencies from database
type Currencies struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float32 `json:"rates"`
}

// LatestRates query db for most updated results for each currency
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
			rows, err = db.Query(`SELECT currency, 
			rate / (SELECT rate FROM rates WHERE ratedate = (SELECT max(ratedate) FROM rates) 
			AND currency = $1) AS rate FROM rates WHERE currency !=$1 
			AND currency = ANY(string_to_array($2, ',')) 
			AND ratedate = (SELECT max(ratedate) FROM rates) 
			UNION ALL SELECT 'EUR', 1 / (SELECT rate FROM rates 
			WHERE ratedate = (SELECT max(ratedate) FROM rates) 
			AND currency = $1)`, base, symbols)
			if err != nil {
				return nil, err
			}
		} else {
			rows, err = db.Query(`SELECT currency, 
			rate / (SELECT rate FROM rates WHERE ratedate = (SELECT max(ratedate) FROM rates) 
			AND currency = $1) AS rate FROM rates WHERE currency !=$1 
			AND currency = ANY(string_to_array($2, ',')) 
			AND ratedate = (SELECT max(ratedate) FROM rates)`, base, symbols)
			if err != nil {
				return nil, err
			}
		}
	case base != "" && symbols == "":
		rows, err = db.Query(`SELECT currency, 
		rate / (SELECT rate FROM rates WHERE ratedate = (SELECT max(ratedate) FROM rates) 
		AND currency = $1) AS rate FROM rates WHERE currency !=$1 
		AND ratedate = (SELECT max(ratedate) FROM rates) 
		UNION ALL SELECT 'EUR', 1 / (SELECT rate FROM rates 
		WHERE ratedate = (SELECT max(ratedate) FROM rates) AND currency = $1)`, base)
		if err != nil {
			return nil, err
		}
	case base == "" && symbols != "":
		rows, err = db.Query(`SELECT currency, rate 
		FROM rates WHERE currency = ANY(string_to_array($1, ',')) 
		AND ratedate = (SELECT max(ratedate) FROM rates)`, symbols)
		if err != nil {
			return nil, err
		}
	default:
		rows, err = db.Query(`SELECT currency, rate FROM rates 
		WHERE ratedate = (SELECT max(ratedate) FROM rates)`)
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
		)
		if err := rows.Scan(&currency, &rate); err != nil {
			return nil, err
		}
		currencies.Rates[currency] = rate
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	var date time.Time
	err = db.QueryRow("SELECT max(ratedate) FROM rates").Scan(&date)
	currencies.Date = date.Format("2006-01-02")

	if base != "" {
		currencies.Base = base
	} else {
		currencies.Base = "EUR"
	}

	return currencies, nil
}
