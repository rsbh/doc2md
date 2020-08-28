package cmd

import (
	"github.com/rsbh/doc2md/internals/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getTokenCmd = &cobra.Command{
	Use:   "get-token",
	Short: "Generate token for google auth",
	Run:   runGetTokenCmd,
}

func init() {
	rootCmd.AddCommand(getTokenCmd)
	getTokenCmd.Flags().StringP("out", "o", defaultPath, "Location for output file")
	viper.AutomaticEnv()
}

func runGetTokenCmd(cmd *cobra.Command, args []string) {
	clientID, clientSercet := auth.GetClientCredentials()
	config := auth.GetConfig(clientID, clientSercet)
	outLocation, _ := cmd.Flags().GetString("out")
	auth.SaveToken(outLocation, config)
}
