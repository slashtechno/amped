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
	"github.com/charmbracelet/log"
	"github.com/slashtechno/amped/internal"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add name",
	Short: "Add a new account",
	Long:  `Add a new Amp account to amped by saving the current logged in account to the keyring with a given name.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if an account with the given name already exists
		saved, err := keyring.Get("amped", args[0])
		if saved != "" {
			log.Fatal("account with the given name already exists", "name", args[0])
		} else if err != keyring.ErrNotFound {
			log.Fatal("error checking for existing account in keyring", "error", err)
		}

		// Extract the API key from the Amp secrets.json file

		apiKey, err := internal.ExtractApiKey(internal.Viper.GetString("secrets"))
		if err != nil {
			log.Fatal("unable to extract api key from amp secrets.json", "error", err)
		}
		log.Debug("extracted api key from amp secrets.json", "apiKey", apiKey)

		err = keyring.Set("amped", args[0], apiKey)
		if err != nil {
			log.Fatal("unable to save api key to keyring", "error", err)
		}
		// Try to retrieve the API key from the keyring
		retrievedApiKey, err := keyring.Get("amped", args[0])
		if err != nil {
			log.Fatal("unable to retrieve api key from keyring", "error", err)
		}
		log.Debug("retrieved api key from keyring", "apiKey", retrievedApiKey)
		log.Info("successfully added account to keyring", "name", args[0])
		err = internal.AppendToAccounts(internal.Viper.GetString("accounts"), internal.AmpedAccount{Name: args[0]})
		if err != nil {
			log.Fatal("unable to add account to accounts list", "error", err)
		}
		log.Debug("successfully added account to accounts list", "name", args[0])
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(addCmd)
}
