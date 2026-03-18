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
	"encoding/json"

	"github.com/charmbracelet/log"
	"github.com/slashtechno/amped/internal"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var addCmd = &cobra.Command{
	Use:   "add name",
	Short: "Add a new account",
	Long:  `Save the currently logged-in account to the keyring under a given name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		svc, err := resolveService()
		if err != nil {
			log.Fatal("unable to determine service", "error", err)
		}

		// Check if an account with this name already exists for the service
		saved, err := internal.ReadFromKeyring(svc, name)
		if saved != "" {
			log.Fatal("account with the given name already exists", "name", name, "service", svc)
		} else if err != keyring.ErrNotFound {
			log.Fatal("error checking for existing account in keyring", "error", err)
		}

		// Read credentials from the appropriate source and store as a string in the keyring.
		// Amp stores a plain API key; Claude stores a JSON-encoded ClaudeStoredCredentials blob.
		var storedValue string
		switch svc {
		case internal.ServiceAmp:
			apiKey, err := internal.ExtractApiKey(internal.Viper.GetString("amp-secrets"))
			if err != nil {
				log.Fatal("unable to extract api key from Amp secrets.json", "error", err)
			}
			log.Debug("extracted api key from Amp secrets.json")
			storedValue = apiKey

		case internal.ServiceClaude:
			stored, err := internal.ExtractClaudeCredentials(
				internal.Viper.GetString("claude-config"),
				internal.Viper.GetString("claude-creds"),
			)
			if err != nil {
				log.Fatal("unable to extract Claude Code credentials", "error", err)
			}
			log.Debug("extracted Claude Code credentials")
			storedJSON, err := json.Marshal(stored)
			if err != nil {
				log.Fatal("unable to marshal Claude Code credentials", "error", err)
			}
			storedValue = string(storedJSON)

		default:
			log.Fatal("unsupported service", "service", svc)
		}

		if err = internal.WriteToKeyring(svc, name, storedValue); err != nil {
			log.Fatal("unable to save credentials to keyring", "error", err)
		}
		// Try to retrieve the credentials back from the keyring to verify storage
		retrievedValue, err := internal.ReadFromKeyring(svc, name)
		if err != nil {
			log.Fatal("unable to retrieve credentials from keyring", "error", err)
		}
		log.Debug("retrieved credentials from keyring", "name", name, "service", svc, "stored", retrievedValue != "")
		log.Info("successfully added account to keyring", "name", name, "service", svc)

		// AppendToAccounts also updates LastService so bare commands default to this service next time
		if err = internal.AppendToAccounts(internal.Viper.GetString("accounts"), internal.Account{Name: name, Service: svc}); err != nil {
			log.Fatal("unable to add account to accounts list", "error", err)
		}
		log.Debug("successfully added account to accounts list", "name", name, "service", svc)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
