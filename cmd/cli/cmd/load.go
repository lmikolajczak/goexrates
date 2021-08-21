package cmd

import (
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load entails commands to load data",
	Long: `Currently it allows to load:

* historical exchange rates from ECB
* daily exchange rates from ECB`,
	// There's no direct action associated with it
	// Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(loadCmd)
}
