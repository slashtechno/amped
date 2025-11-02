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
	"github.com/slashtechno/amped/internal"
	"github.com/spf13/cobra"
)

var cfgFile string

var defaultAmpedConfigPath string

// defaultAmpSecretsPath doesn't need to be glboal
var defaultAmpedAccountsPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "amped",
	Short: "A utility to switch between Amp (http://ampcode.com/) accounts",
	Long: fmt.Sprintf(`Switch between Amp (http://ampcode.com/) accounts by switching out %s
First, log in with an Amp account using the Amp CLI and then, run 'amped add <name>' to save the account.
Then, you can switch between saved accounts using 'amped switch <name>'
To delete a saved account (won't log you out if it's the active account), use 'amped delete <name>'`, defaultAmpedAccountsPath),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// log.Debug("amped called", "ampSecretsPath", internal.Viper.GetString("secrets"), "configFile", internal.Viper.ConfigFileUsed(), "logLevel", internal.Viper.GetString("log"), "args", args)
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
	// log.SetLevel(log.InfoLevel)
	// get xdg paths

	defaultAmpedConfigPath, err := xdg.ConfigFile("amped.json")
	if err != nil {
		log.Fatal("unable to get xdg config file path for amped config.json", "error", err)
	}

	// defaultAmpSecretsPath, err := xdg.DataFile("amp/secrets.json")
	// if err != nil {
	// 	log.Fatal("unable to get xdg data file path for amp secrets.json", "error", err)
	// }

	// Amp's data directory, on both Windows and Linux, is under ~/.local/share/amp
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("unable to get user home directory", "error", err)
	}
	defaultAmpSecretsPath := fmt.Sprintf("%s/.local/share/amp/secrets.json", home)

	defaultAmpedAccountsPath, err = xdg.StateFile("amped.json")
	if err != nil {
		log.Fatal("unable to get xdg state file path for amped accounts.json", "error", err)
	}

	configFileHelp := fmt.Sprintf("config file (default is %s)", defaultAmpedConfigPath)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", configFileHelp)
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("secrets", "s", "", fmt.Sprintf("path to Amp secrets file (default is %s)", defaultAmpSecretsPath))
	rootCmd.MarkFlagFilename("secrets", "json")
	internal.Viper.BindPFlag("secrets", rootCmd.PersistentFlags().Lookup("secrets"))
	internal.Viper.SetDefault("secrets", defaultAmpSecretsPath)

	// rootCmd.PersistentFlags().String("threads", "", fmt.Sprintf("path to Amp threads directory (default is %s)", defaultAmpThreadsPath))
	// rootCmd.MarkFlagDirname("threads")
	// internal.Viper.BindPFlag("threads", rootCmd.PersistentFlags().Lookup("threads"))
	// internal.Viper.SetDefault("threads", defaultAmpThreadsPath)

	// rootCmd.PersistentFlags().String("history", "", fmt.Sprintf("path to Amp history file (default is %s)", defaultAmpHistoryPath))
	// rootCmd.MarkFlagFilename("history", "json")
	// internal.Viper.BindPFlag("history", rootCmd.PersistentFlags().Lookup("history"))
	// internal.Viper.SetDefault("history", defaultAmpHistoryPath)

	rootCmd.PersistentFlags().String("accounts", "", fmt.Sprintf("path to Amp list of accounts file (default is %s)", defaultAmpedAccountsPath))
	rootCmd.MarkFlagFilename("accounts", "json")
	internal.Viper.BindPFlag("accounts", rootCmd.PersistentFlags().Lookup("accounts"))
	internal.Viper.SetDefault("accounts", defaultAmpedAccountsPath)

	rootCmd.PersistentFlags().StringP("log", "l", "info", "log level (debug, info, warn, error, fatal, panic)")
	internal.Viper.BindPFlag("log", rootCmd.PersistentFlags().Lookup("log"))
	internal.Viper.SetDefault("log", "info")

	cobra.OnInitialize(func() {
		err := setupLogLevel(internal.Viper.GetString("log"))
		if err != nil {
			log.Fatal("unable to set up log level", "error", err)
		}
	})

	cobra.OnInitialize(func() {
		err := internal.EnsureAccounts(internal.Viper.GetString("accounts"))
		if err != nil {
			log.Fatal("unable to ensure accounts file exists", "error", err)
		}
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		internal.Viper.SetConfigFile(cfgFile)
	} else {
		internal.Viper.SetConfigFile(defaultAmpedConfigPath)
	}

	internal.Viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := internal.Viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", internal.Viper.ConfigFileUsed())
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
