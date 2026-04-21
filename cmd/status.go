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
	"fmt"

	"github.com/fatih/color"
	"github.com/slashtechno/amped/internal"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show saved accounts and which are active",
	Long: `Show saved accounts and which are currently active.
Use --service to filter to a specific service.`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := internal.ReadFromAccounts(internal.Viper.GetString("accounts"))
		if err != nil {
			fmt.Println("Unable to read accounts list:", err)
			return
		}

		filter := internal.Service(internal.Viper.GetString("service"))

		fmt.Println("Saved accounts:")
		for _, account := range list.Accounts {
			if filter != "" && account.Service != filter {
				continue
			}

			var activeForService string
			switch account.Service {
			case internal.ServiceAmp:
				activeForService = list.ActiveAmp
			case internal.ServiceClaude:
				activeForService = list.ActiveClaude
			}

			if account.Name == activeForService {
				// fmt.Printf(" - [X] %-20s [%s]\n", account.Name, account.Service)
				color.Green(" - [X] %-20s [%s]", account.Name, account.Service)
			} else {
				// Light grey color for inactive accounts
				// fmt.Printf(" - [ ] %-20s [%s]\n", account.Name, account.Service)
				color.HiBlack(" - [ ] %-20s [%s]", account.Name, account.Service)
			}
		}

		if filter == internal.ServiceClaude || filter == "" {
			active := list.ActiveClaude
			if active != "" {
				stored, readErr := internal.ReadFromKeyring(internal.ServiceClaude, active)
				if readErr == nil {
					var creds internal.ClaudeStoredCredentials
					if jsonErr := json.Unmarshal([]byte(stored), &creds); jsonErr == nil {
						_, _, accessToken := internal.ExtractClaudeAccountDetails(creds)
						if accessToken != "" {
							if email, org, verifyErr := internal.VerifyClaudeToken(accessToken); verifyErr != nil {
								color.Yellow("Live Claude auth: invalid (%v)", verifyErr)
							} else if email != "" || org != "" {
								color.Green("Live Claude auth: valid (email=%s org=%s)", email, org)
							} else {
								color.Green("Live Claude auth: valid")
							}
						}
					}
				}
			}
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
