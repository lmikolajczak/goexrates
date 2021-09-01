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

// latestCmd represents the latest command
var latestCmd = &cobra.Command{
	Use:   "latest",
	Short: "Fetch and ingest latest exchange rates published by ECB",
	Long: `Fetch data that are published by ECB on daily basis and ingest them 
into the database.

Data consist of euro foreign exchange reference rates and are usually updated around 
16:00 CET on every working day, except on closing days. They are based on a regular 
daily concertation procedure between central banks across Europe, which normally 
takes place at 14:15 CET.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logFile, err := os.OpenFile("latest.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logFile.Close()
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Record command execution timestamp
		start := time.Now().Round(time.Millisecond)
		log.Print("Get latest rates")
		// Get passed args and options
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			log.Fatal(err, nil)
		}
		dbDsn, err := cmd.Flags().GetString("database")
		if err != nil {
			log.Fatal(err, nil)
		}
		// Get xml data and decode them (external API call)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			log.Fatal(err, nil)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err, nil)
		}
		defer resp.Body.Close()

		var rates source.ECB
		if err := xml.NewDecoder(resp.Body).Decode(&rates); err != nil {
			log.Fatal(err, nil)
		}
		// Open database and establish connection.
		db, err := sql.Open("postgres", dbDsn)
		if err != nil {
			log.Fatal(err, nil)
		}
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = db.PingContext(ctx)
		if err != nil {
			log.Fatal(err, nil)
		}
		// This will always contain only a single day.
		for _, day := range rates.Days {
			// Before we insert any latest rates, check if they are not already
			// in the database. This will allow us to run this command multiple
			// times within a single day in an idempotent manner.
			var latestDate string
			row := db.QueryRow("SELECT MAX(created_at) FROM currencies")
			err := row.Scan(&latestDate)
			if err != nil {
				log.Fatal(err, nil)
			}
			if latestDate >= day.Date {
				log.Print("Rates are up to date")
				log.Print(
					fmt.Sprintf(
						"Done (%v)",
						time.Since(start).Round(time.Millisecond).String(),
					),
				)
				return
			}

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
			for _, currency := range day.Currencies {
				_, err := stmt.Exec(currency.Code, currency.Value, day.Date)
				if err != nil {
					log.Fatal(err)
				}
			}
			tx.Commit()
			log.Print(fmt.Sprintf("Inserted latest data (%v)", day.Date))
		}
		log.Print(
			fmt.Sprintf(
				"Done (%v)",
				time.Since(start).Round(time.Millisecond).String(),
			),
		)
	},
}

func init() {
	loadCmd.AddCommand(latestCmd)

	// Define flags and configuration settings.
	latestCmd.Flags().StringP(
		"url", "u", "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml",
		"A url that points to the data. Optional",
	)
	latestCmd.Flags().StringP(
		"database", "d", os.Getenv("EXRATES_DB_DSN"),
		"A database into which the rates will be inserted. Optional",
	)
}
