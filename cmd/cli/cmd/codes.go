/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

// codesCmd represents the codes command
var codesCmd = &cobra.Command{
	Use:   "codes",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		// Get unique ISO codes based on historical data published by ECB (ISO-4217)
		// and insert them into currency table.
		codes := map[string]string{}
		for _, day := range rates.Days {
			for _, rate := range day.Rates {
				codes[rate.Iso] = rate.Iso
			}
		}
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
		fmt.Println(codes)
	},
}

func init() {
	loadCmd.AddCommand(codesCmd)

	// Flags and configuration settings.
	codesCmd.Flags().StringP(
		"url", "u", "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist.xml",
		"A url that points to the data. Optional",
	)
	codesCmd.Flags().StringP(
		"database", "d", os.Getenv("EXRATES_DB_DSN"),
		"A database into which the codes will be inserted. Optional",
	)
}
