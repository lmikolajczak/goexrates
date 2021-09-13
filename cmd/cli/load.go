package main

import (
	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load entails commands to load data",
	Long: `Load entails commands to load following data into specified database:

* historical exchange rates from ECB
* latest (published daily) exchange rates from ECB`,
}

func init() {
	rootCmd.AddCommand(loadCmd)
}
