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

var updateCmd = &cobra.Command{
	Use:   "update name",
	Short: "Update credentials for a saved account",
	Long:  `Update the credentials for a previously saved account with the currently logged-in credentials.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		svc, err := resolveService()
		if err != nil {
			log.Fatal("unable to determine service", "error", err)
		}

		// Check if an account with this name exists for the service
		saved, err := internal.ReadFromKeyring(svc, name)
		if saved == "" {
			if err == keyring.ErrNotFound {
				log.Fatal("no account found with the given name", "name", name, "service", svc)
			}
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
			log.Fatal("unable to update credentials in keyring", "error", err)
		}
		// Try to retrieve the credentials back from the keyring to verify storage
		retrievedValue, err := internal.ReadFromKeyring(svc, name)
		if err != nil {
			log.Fatal("unable to retrieve credentials from keyring", "error", err)
		}
		log.Debug("retrieved credentials from keyring", "name", name, "service", svc, "stored", retrievedValue != "")
		log.Info("successfully updated account in keyring", "name", name, "service", svc)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}