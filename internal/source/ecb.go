package source

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	_ "github.com/lib/pq"
)

type ECB struct {
	XMLName xml.Name
	Days    []struct {
		Date       string `xml:"time,attr"`
		Currencies []struct {
			Code  string  `xml:"currency,attr"`
			Value float64 `xml:"rate,attr"`
		} `xml:"Cube"`
	} `xml:"Cube>Cube"`
}

func (e *ECB) Get(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := xml.NewDecoder(resp.Body).Decode(e); err != nil {
		return err
	}
	return nil
}

func (e *ECB) Insert(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	for i := range e.Days {
		// Iterate from oldest to newest to make sure that records are inserted
		// in a proper order in case there is more than one date to process.
		// ECB sorts by date by default.
		day := e.Days[len(e.Days)-1-i]

		var latestDate string
		row := db.QueryRow("SELECT MAX(created_at) FROM currencies")
		err := row.Scan(&latestDate)
		if err != nil {
			return err
		}
		if latestDate >= day.Date {
			return nil
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}
		stmt, err := tx.Prepare(
			`INSERT INTO currencies (code, rate, created_at) VALUES ($1, $2, $3)`,
		)
		if err != nil {
			return err
		}

		for _, currency := range day.Currencies {
			_, err := stmt.Exec(currency.Code, currency.Value, day.Date)
			if err != nil {
				return err
			}
		}
		tx.Commit()
		fmt.Printf("Inserted data for: %v\r", day.Date)
	}
	return nil
}
