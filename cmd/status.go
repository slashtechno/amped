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

	"github.com/fatih/color"
	"github.com/slashtechno/amped/internal"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get information regarding current Amp accounts",
	Long: `Get information regarding current Amp accounts. 
Includes which account is currently active and a list of all saved accounts.`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := internal.ReadFromAccounts(internal.Viper.GetString("accounts"))
		if err != nil {
			fmt.Println("Unable to read accounts list:", err)
			return
		}
		// fmt.Println("Current active account:", list.Active)
		fmt.Println("Saved accounts:")
		for _, account := range list.Accounts {
			if account.Name == list.Active {
				// fmt.Println(" - [X] ", account.Name, "(active)")
				color.Green(" - [X] %s", account.Name)
				continue
			}
			// fmt.Println(" - [ ] ", account.Name)
			// Light grey color for inactive accounts
			color.HiBlack(" - [ ] %s", account.Name)
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

}
