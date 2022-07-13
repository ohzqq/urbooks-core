package cmd

import (
	"os"
	"path/filepath"

	"github.com/ohzqq/urbooks-core/urbooks"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var lib string
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "urbooks",
	Short: "",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.urbooks-core.yaml)")
	rootCmd.PersistentFlags().StringVarP(&lib, "library", "l", "", "library by name")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "library by name")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".urbooks-core" (without extension).
		//viper.AddConfigPath(home)
		viper.AddConfigPath(filepath.Join(home, "Code/urbooks-core/tmp/"))
		viper.AddConfigPath(filepath.Join(home, ".config/urbooks/"))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		urbooks.InitConfig(viper.GetStringMapString("library_options"))
		urbooks.InitLibraries(viper.Sub("libraries"), viper.GetStringMapString("libraries"), false)

		urbooks.CfgCdb(viper.Sub("calibre"))
	}
}
