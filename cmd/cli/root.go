package main

import (
	"os"

	"github.com/spf13/cobra"
)

var logFile *os.File

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "A CLI that helps to manage the application",
	Long: `A CLI that helps to manage common/recurring tasks
in the application. Such as:

* loading and managing data from different sources`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
