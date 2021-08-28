package cmd

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Luqqk/goexrates/internal/source"
	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

// historicalCmd represents the historical command
var historicalCmd = &cobra.Command{
	Use:   "historical",
	Short: "Fetches and ingests historical exchange rates published by ECB",
	Long: `Fetches historical data that are published by ECB and ingests them 
into the database.
	
Data consist of euro foreign exchange reference rates and go back as far as
1999-01-04. Keep in mind that available currencies changed over the years.`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		url, err := cmd.Flags().GetString("url")
		if err != nil {
			fmt.Println("unable to parse url option")
			os.Exit(1)
		}
		dbDsn, err := cmd.Flags().GetString("database")
		if err != nil {
			fmt.Println("unable to parse db option")
			os.Exit(1)
		}
		// Get xml data and decode them.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			fmt.Println("unable to create request")
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("unable to fetch xml data")
			os.Exit(1)
		}
		defer resp.Body.Close()

		var rates source.ECB
		if err := xml.NewDecoder(resp.Body).Decode(&rates); err != nil {
			fmt.Println("unable to decode xml data")
			os.Exit(1)
		}
		// Open database and establish connection.
		db, err := sql.Open("postgres", dbDsn)
		if err != nil {
			fmt.Println("unable to open database")
			os.Exit(1)
		}
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = db.PingContext(ctx)
		if err != nil {
			fmt.Println("unable to establish database connection")
			os.Exit(1)
		}

		// This is highly inefficient, TODO bulk inserts
		for i := range rates.Days {
			tx, err := db.Begin()
			if err != nil {
				fmt.Println(err)
			}
			stmt, err := tx.Prepare(
				`INSERT INTO currencies (code, rate, created_at) VALUES ($1, $2, $3)`,
			)
			if err != nil {
				fmt.Println(err)
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
		// Print summary of the operation
		fmt.Printf(
			"Inserted historical data from %v to %v in %v\n",
			rates.Days[len(rates.Days)-1].Date,
			rates.Days[0].Date,
			end.Sub(start).Round(time.Second),
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
		"A database into which the codes will be inserted. Optional",
	)
}
