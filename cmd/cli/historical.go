package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Luqqk/goexrates/internal/source"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
		logFile, err := os.OpenFile("log/historical.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logFile.Close()
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Get historical rates")
		start := time.Now()

		url, err := cmd.Flags().GetString("url")
		if err != nil {
			log.Fatal(err)
		}
		dsn, err := cmd.Flags().GetString("database")
		if err != nil {
			log.Fatal(err)
		}

		var rates source.ECB
		if err := rates.Get(url); err != nil {
			log.Fatal(err, nil)
		}
		// This is highly inefficient, TODO bulk inserts
		if err := rates.Insert(dsn); err != nil {
			log.Fatal(err)
		}

		log.Println(
			fmt.Sprintf(
				"Inserted historical data from %v to %v (%v)",
				rates.Days[len(rates.Days)-1].Date,
				rates.Days[0].Date,
				time.Since(start).Round(time.Second),
			),
		)
		log.Println(
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
