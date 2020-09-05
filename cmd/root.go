package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/rsbh/doc2md/internals/auth"
	"github.com/rsbh/doc2md/internals/config"
	"github.com/rsbh/doc2md/internals/gdrive"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "doc2md",
	Short: "Generate docs from google docs",
	Run:   runRootCmd,
}

const defaultPath = "token.json"
const defaultOutDir = "out"

var configuration config.Configurations

func readConfig(cfgFile string) {
	viper.SetConfigFile(cfgFile)
	viper.SetDefault("outDir", defaultOutDir)
	viper.SetDefault("breakDoc", false)
	viper.SetDefault("supportCodeBlock", false)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

}

func runRootCmd(cmd *cobra.Command, args []string) {
	clientID, clientSercet := auth.GetClientCredentials()
	cfgFile, _ := cmd.Flags().GetString("config")
	tokenFile, _ := cmd.Flags().GetString("token")
	tok := auth.GetToken(tokenFile)
	readConfig(cfgFile)
	c := auth.GetConfig(clientID, clientSercet)
	client := auth.GetClient(c, tok)
	s := &gdrive.Service{}
	s.Init(client)
	s.GetFiles(configuration.FolderID)
}

func init() {
	rootCmd.Flags().StringP("config", "c", "", "Location for config file")
	rootCmd.Flags().StringP("token", "t", defaultPath, "Location for token file")
	rootCmd.MarkFlagRequired("config")
	viper.AutomaticEnv()
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
