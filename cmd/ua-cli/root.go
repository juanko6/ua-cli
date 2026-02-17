package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ua",
	Short: "ua-cli - University of Alicante Command-Line Companion",
	Long:  `Fast, privacy-focused CLI for University of Alicante services (Schedule, Grades, Campus).`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
