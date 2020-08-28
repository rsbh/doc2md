package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "doc2md",
	Short: "Generate docs from google docs",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Println("Hello Doc2md") },
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
