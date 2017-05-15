package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	// save newest data to db
	// add sql.Open()
	/*txn, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.Prepare(pq.CopyIn("rates", "currency", "rate", "date"))
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
	}*/
	return nil
}

func main() {
	fmt.Println("============ Update DB with newest data ============")

	currencies, date, err := getECBData(url)
	if err != nil {
		log.Fatal(err)
	}

	err = updateDb(currencies, date)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("============ DB successfully updated ============")
}
