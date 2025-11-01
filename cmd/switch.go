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

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch name",
	Short: "Switch to a given saved Amp account",
	Long: `Switch between saved Amp (http://ampcode.com/) accounts by switching out ~/.local/share/amp/secrets.json
Provide the name of the saved account to switch to.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		// Check if the account exists
		key, err := internal.ReadFromKeyring(name)
		if err == keyring.ErrNotFound {
			log.Fatal("no account found with the given name", "name", name)
		} else if err != nil {
			log.Fatal("failed to read from keyring", "error", err, "accountName", name)
		}
		log.Debug("read api key from keyring", "accountName", name, "apiKey", key)

		err = internal.WriteToAmpSecrets(key, internal.Viper.GetString("secrets"))
		if err != nil {
			log.Fatal("unable to write api key to amp secrets.json", "error", err)
		}
		log.Info("switched amp account successfully", "accountName", name)

		// Update the current account in the accounts list
		err = internal.UpdateActiveAccount(internal.Viper.GetString("accounts"), name)
		if err != nil {
			log.Fatal("unable to update current account in accounts list", "error", err)
		}
		log.Debug("successfully updated current account in accounts list", "accountName", name)

		// Make sure that the keyring api key and amp secrets.json api key match
		verifiedApiKey, err := internal.ExtractApiKey(internal.Viper.GetString("secrets"))
		if err != nil {
			log.Fatal("unable to extract api key from amp secrets.json for verification", "error", err)
		}
		if verifiedApiKey != key {
			log.Fatal("api key in amp secrets.json does not match the one in keyring after switching account", "accountName", name)
		}
		log.Debug("verified that api key in amp secrets.json matches the one in keyring", "accountName", name)

		// It is NOT needed to delete threads/history when switching accounts.

		// // Delete all threads
		// err = internal.DeleteAllThreads(internal.Viper.GetString("threads"))
		// if err != nil {
		// 	log.Fatal("unable to delete amp threads", "error", err)
		// }
		// log.Info("successfully deleted all amp threads after switching account")

		// // Delete history file
		// err = internal.DeleteHistoryFile(internal.Viper.GetString("history"))
		// if err != nil {
		// 	log.Fatal("unable to delete amp history file", "error", err)
		// }
		// log.Info("successfully deleted amp history file after switching account")
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
