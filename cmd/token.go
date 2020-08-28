package cmd

import (
	"github.com/rsbh/doc2md/internals/auth"
	"github.com/spf13/cobra"
)

var getTokenCmd = &cobra.Command{
	Use:   "get-token",
	Short: "Generate token for google auth",
	Run:   runCmd,
}

func init() {
	rootCmd.AddCommand(getTokenCmd)
	getTokenCmd.Flags().StringP("out", "o", "", "Location for output file")
}

const defaultPath = "token.json"

func runCmd(cmd *cobra.Command, args []string) {
	clientID, clientSercet := auth.GetClientCredentials()
	config := auth.GetConfig(clientID, clientSercet)
	outLocation, _ := cmd.Flags().GetString("out")
	if outLocation == "" {
		outLocation = defaultPath
	}
	auth.SaveToken(outLocation, config)
}
