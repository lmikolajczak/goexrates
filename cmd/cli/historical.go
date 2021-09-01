package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Luqqk/goexrates/internal/source"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	_ "github.com/lib/pq"
)

// historicalCmd represents the historical command
var historicalCmd = &cobra.Command{
	Use:   "historical",
	Short: "Fetch and ingest historical exchange rates published by ECB",
	Long: `Fetch historical data that are published by ECB and ingest them 
into the database.
	
Data consist of euro foreign exchange reference rates and go back as far as
1999-01-04. Keep in mind that available currencies changed over the years.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logFile, err := os.OpenFile("historical.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logFile.Close()
	},
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		log.Print("Get historical rates")

		url, err := cmd.Flags().GetString("url")
		if err != nil {
			log.Fatal(err)
		}
		dbDsn, err := cmd.Flags().GetString("database")
		if err != nil {
			log.Fatal(err)
		}
		// Get xml data and decode them.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var rates source.ECB
		if err := xml.NewDecoder(resp.Body).Decode(&rates); err != nil {
			log.Fatal(err)
		}
		// Open database and establish connection.
		db, err := sql.Open("postgres", dbDsn)
		if err != nil {
			log.Fatal(err)
		}
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = db.PingContext(ctx)
		if err != nil {
			log.Fatal(err)
		}

		// This is highly inefficient, TODO bulk inserts
		for i := range rates.Days {
			tx, err := db.Begin()
			if err != nil {
				log.Fatal(err)
			}
			stmt, err := tx.Prepare(
				`INSERT INTO currencies (code, rate, created_at) VALUES ($1, $2, $3)`,
			)
			if err != nil {
				log.Fatal(err)
			}
			// Iterate from oldest to newest to make sure that records are inserted
			// in a proper order.
			day := rates.Days[len(rates.Days)-1-i]
			for _, currency := range day.Currencies {
				_, err := stmt.Exec(currency.Code, currency.Value, day.Date)
				if err != nil {
					fmt.Println(err)
				}
			}
			tx.Commit()
			fmt.Printf("Inserted data for: %v\r", day.Date)
		}
		end := time.Now()
		log.Print(
			fmt.Sprintf(
				"Inserted historical data from %v to %v (%v)\n",
				rates.Days[len(rates.Days)-1].Date,
				rates.Days[0].Date,
				end.Sub(start).Round(time.Second),
			),
		)
		log.Print(
			fmt.Sprintf(
				"Done (%v)",
				time.Since(start).Round(time.Millisecond).String(),
			),
		)
	},
}

func init() {
	loadCmd.AddCommand(historicalCmd)

	// Flags and configuration settings.
	historicalCmd.Flags().StringP(
		"url", "u", "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.xml",
		"A url that points to the data. Optional",
	)
	historicalCmd.Flags().StringP(
		"database", "d", os.Getenv("EXRATES_DB_DSN"),
		"A database into which the rates will be inserted. Optional",
	)
}
