package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/rsbh/doc2md/internals/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "doc2md",
	Short: "Generate docs from google docs",
	Run:   runRootCmd,
}

func runRootCmd(cmd *cobra.Command, args []string) {
	cfgFile, _ := cmd.Flags().GetString("config")
	tokenFile, _ := cmd.Flags().GetString("token")
	token := auth.GetToken(tokenFile)
	fmt.Println(token)
	viper.SetConfigFile(cfgFile)
}

const defaultPath = "token.json"

func init() {
	rootCmd.Flags().StringP("config", "c", "", "Location for config file")
	rootCmd.Flags().StringP("token", "t", defaultPath, "Location for token file")
	rootCmd.MarkFlagRequired("config")
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
