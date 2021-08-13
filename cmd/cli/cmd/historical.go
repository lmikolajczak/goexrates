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
	"github.com/spf13/cobra"
)

// historicalCmd represents the historical command
var historicalCmd = &cobra.Command{
	Use:   "historical",
	Short: "Fetches and ingests daily exchange rates published by ECB",
	Long: `Fetches historical data that are published by ECB and ingests them 
into the database.
	
Data consist of euro foreign exchange reference rates and go back as far as
1999-01-04. Keep in mind that available currencies changed over the years.`,
	Run: func(cmd *cobra.Command, args []string) {
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
		// Initialize http client with proper timeout.
		netClient := &http.Client{
			Timeout: time.Second * 60,
		}
		// Get xml data and decode them.
		resp, err := netClient.Get(url)
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
		// Insert data into the database.
		fmt.Println(rates.Days)
		// Open database and establish connection.
		db, err := sql.Open("postgres", dbDsn)
		if err != nil {
			fmt.Println("unable to open database")
			os.Exit(1)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = db.PingContext(ctx)
		if err != nil {
			fmt.Println("unable to establish database connection")
			os.Exit(1)
		}
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
