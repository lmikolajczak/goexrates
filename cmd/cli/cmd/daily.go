package cmd

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Luqqk/goexrates/internal/source"
	"github.com/spf13/cobra"
)

// dailyCmd represents the daily command
var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Fetches and ingests daily exchange rates published by ECB",
	Long: `Fetches data that are published by ECB on daily basis and ingests them 
into the database.

Data consist of euro foreign exchange reference rates and are usually updated around 
16:00 CET on every working day, except on closing days. They are based on a regular 
daily concertation procedure between central banks across Europe, which normally 
takes place at 14:15 CET.`,
	Run: func(cmd *cobra.Command, args []string) {
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			fmt.Println("unable to parse url option")
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
	},
}

func init() {
	fetchCmd.AddCommand(dailyCmd)

	// Flags and configuration settings.
	dailyCmd.Flags().StringP(
		"url", "u", "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml",
		"A url that points to the data. Optional",
	)
}
