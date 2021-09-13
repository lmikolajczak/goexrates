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
		logFile, err := os.OpenFile("logs/latest.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logFile.Close()
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Get latest rates")
		start := time.Now().Round(time.Millisecond)

		url, err := cmd.Flags().GetString("url")
		if err != nil {
			log.Fatal(err, nil)
		}
		dsn, err := cmd.Flags().GetString("database")
		if err != nil {
			log.Fatal(err, nil)
		}

		var rates source.ECB
		if err := rates.Get(url); err != nil {
			log.Fatal(err, nil)
		}
		if err := rates.Insert(dsn); err != nil {
			log.Fatal(err)
		}
		log.Println(
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
