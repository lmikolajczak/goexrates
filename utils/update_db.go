package main

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/lib/pq"
)

const url = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

type ecbCurrencies struct {
	Currencies []struct {
		Currency string  `xml:"currency,attr"`
		Rate     float32 `xml:"rate,attr"`
	} `xml:"Cube>Cube>Cube"`
}

type ecbUpdateDate struct {
	Date struct {
		Time string `xml:"time,attr"`
	} `xml:"Cube>Cube"`
}

func getECBData(url string) (*ecbCurrencies, *ecbUpdateDate, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()

	var currencies ecbCurrencies
	var updateDate ecbUpdateDate

	data, err := ioutil.ReadAll(response.Body)

	if err := xml.Unmarshal(data, &currencies); err != nil {
		return nil, nil, err
	}
	if err := xml.Unmarshal(data, &updateDate); err != nil {
		return nil, nil, err
	}
	return &currencies, &updateDate, nil
}

func updateDb(currencies *ecbCurrencies, updateDate *ecbUpdateDate) error {
	db, err := sql.Open("postgres", "postgres://user:password@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	txn, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.Prepare(pq.CopyIn("currencies", "iso_code", "rate", "date"))
	if err != nil {
		return err
	}

	for _, currency := range currencies.Currencies {
		_, err = stmt.Exec(currency.Currency, currency.Rate, updateDate.Date.Time)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}

func checkDbDate(ecbDataDate *ecbUpdateDate) error {
	db, err := sql.Open("postgres", "postgres://user:password@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var dbMaxDate string
	err = db.QueryRow("SELECT max(date) FROM currencies").Scan(&dbMaxDate)
	switch {
	case err == sql.ErrNoRows:
		return err
	case err != nil:
		return err
	default:
		if ecbDataDate.Date.Time > dbMaxDate {
			return nil
		}
		err = errors.New("Db is already up-to-date")
		return err
	}
}

func main() {
	fmt.Println("============ Update DB with newest data ============")
	start := time.Now()

	currencies, date, err := getECBData(url)
	if err != nil {
		log.Fatal(err)
	}

	err = checkDbDate(date)
	if err != nil {
		log.Fatal(err)
	}

	err = updateDb(currencies, date)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()
	execTime := end.Sub(start)
	fmt.Printf("============ DB successfully updated in %s ==========\n", execTime)
}
