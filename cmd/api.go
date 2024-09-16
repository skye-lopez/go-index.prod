package cmd

import (
	"github.com/skye-lopez/go-index.prod/api"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "open-api",
	Short: "serves api",
	Long:  "serves api",
	Run:   open,
}

func init() {
	rootCmd.AddCommand(apiCmd)
}

func open(cmd *cobra.Command, args []string) {
	api.Open()
}
