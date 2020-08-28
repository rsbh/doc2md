package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/rsbh/doc2md/internals/auth"
	"github.com/rsbh/doc2md/internals/config"
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

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

}

func runRootCmd(cmd *cobra.Command, args []string) {
	cfgFile, _ := cmd.Flags().GetString("config")
	tokenFile, _ := cmd.Flags().GetString("token")
	_ = auth.GetToken(tokenFile)

	readConfig(cfgFile)

	err := os.MkdirAll(configuration.OutDir, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
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
