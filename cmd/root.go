package cmd

import (
	"github.com/spf13/cobra"

	"fmt"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "ascot",
	Short: "Ascot is an AWS toolkit for SecOps and GRC",
	Long: `A suite of various tools to inspect AWS resources that
          security operations and GRC (compliance) teams often
          need to use.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
