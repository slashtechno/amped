/*
Copyright © 2025 Angad Behl <77907286+slashtechno@users.noreply.github.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var defaultAmpedConfigPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "amped",
	Short: "A utility to switch between Amp (http://ampcode.com/) accounts",
	Long: `Switch between Amp (http://ampcode.com/) accounts by switching out ~/.local/share/amp/secrets.json
First, log in with an Amp account using the Amp CLI and then, run 'amped save <name>' to save the account.
Then, you can switch between saved accounts using 'amped switch <name>'
To delete a saved account (won't log you out if it's the active account), use 'amped delete <name>'`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Info("amped called", "ampSecretsPath", viper.GetString("secrets"), "configFile", viper.ConfigFileUsed(), "logLevel", viper.GetString("log"), "args", args)
	},
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
	log.SetLevel(log.DebugLevel)
	// get xdg paths
	defaultAmpSecretsPath, err := xdg.DataFile("amp/secrets.json")
	if err != nil {
		log.Fatal("unable to get xdg data file path for amp secrets.json", "error", err)
	}

	defaultAmpedConfigPath, err = xdg.ConfigFile("amped.json")
	if err != nil {
		log.Fatal("unable to get xdg config file path for amped config.json", "error", err)
	}

	configFileHelp := fmt.Sprintf("config file (default is %s)", defaultAmpedConfigPath)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", configFileHelp)

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("secrets", "s", "", "path to amp secrets.json file")
	viper.BindPFlag("secrets", rootCmd.PersistentFlags().Lookup("secrets"))
	viper.SetDefault("secrets", defaultAmpSecretsPath)

	rootCmd.PersistentFlags().StringP("log", "l", "info", "log level (debug, info, warn, error, fatal, panic)")
	viper.BindPFlag("log", rootCmd.PersistentFlags().Lookup("log"))
	viper.SetDefault("log", "debug")

	cobra.OnInitialize(func() {
		err := setupLogLevel(viper.GetString("log"))
		if err != nil {
			log.Fatal("unable to set up log level", "error", err)
		}
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	fmt.Println("Initializing config")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile(defaultAmpedConfigPath)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func setupLogLevel(logLevel string) error {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("unable to parse log level: %w", err)
	}
	log.SetLevel(level)

	return nil
}
