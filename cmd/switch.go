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

var switchCmd = &cobra.Command{
	Use:   "switch name",
	Short: "Switch to a saved account",
	Long:  `Restore credentials for a previously saved account, making it the active account.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		svc, err := resolveService()
		if err != nil {
			log.Fatal("unable to determine service", "error", err)
		}

		// Retrieve the stored credentials string from the keyring
		stored, err := internal.ReadFromKeyring(svc, name)
		if err == keyring.ErrNotFound {
			log.Fatal("no account found with the given name", "name", name, "service", svc)
		} else if err != nil {
			log.Fatal("failed to read from keyring", "error", err, "name", name)
		}
		log.Debug("read credentials from keyring", "name", name, "service", svc)

		// Auto-update the currently active account before switching
		if accounts, err := internal.ReadFromAccounts(internal.Viper.GetString("accounts")); err == nil {
			var activeAccount string
			switch svc {
			case internal.ServiceAmp:
				activeAccount = accounts.ActiveAmp
			case internal.ServiceClaude:
				activeAccount = accounts.ActiveClaude
			}

			if activeAccount != "" && activeAccount != name {
				log.Info("auto-updating currently active account before switching", "name", activeAccount)
				var storedValue string
				var updateErr error
				switch svc {
				case internal.ServiceAmp:
					storedValue, updateErr = internal.ExtractApiKey(internal.Viper.GetString("amp-secrets"))
				case internal.ServiceClaude:
					var storedCreds internal.ClaudeStoredCredentials
					storedCreds, updateErr = internal.ExtractClaudeCredentials(
						internal.Viper.GetString("claude-config"),
						internal.Viper.GetString("claude-creds"),
					)
					if updateErr == nil {
						var storedJSON []byte
						storedJSON, updateErr = json.Marshal(storedCreds)
						storedValue = string(storedJSON)
					}
				}
				if updateErr != nil {
					log.Warn("failed to extract credentials for auto-update (skipping)", "error", updateErr)
				} else if storedValue != "" {
					if err := internal.WriteToKeyring(svc, activeAccount, storedValue); err != nil {
						log.Warn("failed to save auto-updated credentials to keyring", "error", err)
					} else {
						log.Debug("successfully auto-updated active account", "name", activeAccount)
					}
				}
			}
		}

		switch svc {
		case internal.ServiceAmp:
			if err = internal.WriteToAmpSecrets(stored, internal.Viper.GetString("amp-secrets")); err != nil {
				log.Fatal("unable to write api key to Amp secrets.json", "error", err)
			}
			// Make sure the key we just wrote matches what was in the keyring
			verified, err := internal.ExtractApiKey(internal.Viper.GetString("amp-secrets"))
			if err != nil {
				log.Fatal("unable to verify Amp secrets.json after switch", "error", err)
			}
			if verified != stored {
				log.Fatal("api key in Amp secrets.json does not match the keyring after switching account", "name", name)
			}
			log.Debug("verified Amp api key written successfully", "name", name)

		case internal.ServiceClaude:
			// Terminate any running Claude instances to prevent them from refreshing tokens
			// and overwriting the keychain with the old account's credentials after we switch.
			internal.KillClaudeProcesses()

			var claudeCreds internal.ClaudeStoredCredentials
			if err = json.Unmarshal([]byte(stored), &claudeCreds); err != nil {
				log.Fatal("unable to unmarshal stored Claude Code credentials", "error", err)
			}
			if err = internal.WriteToClaudeCredentials(claudeCreds,
				internal.Viper.GetString("claude-config"),
				internal.Viper.GetString("claude-creds"),
			); err != nil {
				log.Fatal("unable to write Claude Code credentials", "error", err)
			}
			email, subType, accessToken := internal.ExtractClaudeAccountDetails(claudeCreds)
			if email != "" || subType != "" {
				log.Info("account details", "email", email, "subscription", subType)
			}

			if accessToken == "" {
				log.Fatal("switch aborted: restored Claude credentials are missing access token")
			}

			if liveEmail, orgName, liveErr := internal.VerifyClaudeToken(accessToken); liveErr != nil {
				log.Warn("live verification failed, but credentials have been restored. Claude Code will attempt to refresh the token on next run.", "error", liveErr)
			} else {
				if liveEmail != "" || orgName != "" {
					log.Info("live verified", "email", liveEmail, "org", orgName)
				} else {
					log.Info("live verified")
				}
			}

		default:
			log.Fatal("unsupported service", "service", svc)
		}

		log.Info("switched account successfully", "name", name, "service", svc)

		// UpdateActiveAccount also updates LastService so bare commands default to this service next time
		if err = internal.UpdateActiveAccount(internal.Viper.GetString("accounts"), name, svc); err != nil {
			log.Fatal("unable to update active account in accounts list", "error", err)
		}

		// It is NOT needed to delete threads/history when switching accounts.

		// // Delete all threads
		// err = internal.DeleteAllThreads(internal.Viper.GetString("threads"))
		// if err != nil {
		// 	log.Fatal("unable to delete threads", "error", err)
		// }
		// log.Info("successfully deleted all threads after switching account")

		// // Delete history file
		// err = internal.DeleteHistoryFile(internal.Viper.GetString("history"))
		// if err != nil {
		// 	log.Fatal("unable to delete history file", "error", err)
		// }
		// log.Info("successfully deleted history file after switching account")
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
