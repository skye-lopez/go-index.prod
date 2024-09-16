package cmd

import (
	"github.com/skye-lopez/go-index.prod/idx"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch-index",
	Short: "fetches the latest index info and updates the db.",
	Long:  "fetches the latest index info and updates the db.",
	Run:   fetch,
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}

func fetch(cmd *cobra.Command, args []string) {
	idx.FetchIdx()
}
