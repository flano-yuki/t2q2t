package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of t2q2t",
	Long:  "Print the version number of t2q2t",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("t2q2t version: %s\n", VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
