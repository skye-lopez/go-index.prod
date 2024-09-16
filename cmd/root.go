package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-index.prod",
	Short: "Internal cli tools for serving and maintaing the go-index backend.",
	Long:  "Internal cli tools for serving and maintaing the go-index backend.",
}

func init() {}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Error loading cobra\n%s", err)
	}
}
