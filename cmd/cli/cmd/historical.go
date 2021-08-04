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
	fetchCmd.AddCommand(historicalCmd)

	// Flags and configuration settings.
	historicalCmd.Flags().StringP(
		"url", "u", "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.xml",
		"A url that points to the data. Optional",
	)
}
