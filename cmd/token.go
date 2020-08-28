package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getTokenCmd = &cobra.Command{
	Use:   "get-token",
	Short: "Generate token for google auth",
	Run: func(cmd *cobra.Command, args []string) {
		outLocation, _ := cmd.Flags().GetString("out")
		fmt.Println("Generate Token", outLocation)
	},
}

func init() {
	rootCmd.AddCommand(getTokenCmd)
	getTokenCmd.Flags().StringP("out", "o", "", "Location for output file")

}
