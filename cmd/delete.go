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
)

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Args:  cobra.MaximumNArgs(1),
	Short: "Delete a saved account, or all saved accounts with --delete-all",
	Run: func(cmd *cobra.Command, args []string) {
		deleteAll, err := cmd.Flags().GetBool("delete-all")
		if err != nil {
			log.Fatal("unable to parse flags", "error", err)
		}

		if deleteAll {
			if err := internal.DeleteAllFromKeyring(); err != nil {
				log.Fatal("unable to delete accounts from keyring", "error", err)
				return
			}
			if err := internal.DeleteAllFromAccounts(internal.Viper.GetString("accounts")); err != nil {
				log.Fatal("unable to clear accounts list", "error", err)
				return
			}
			log.Info("successfully deleted all accounts")
			return
		} else if len(args) == 0 {
			// Make sure that if --delete-all is not set, an account name is provided
			log.Fatal("please provide an account name to delete or use --delete-all to delete all accounts")
			return
		}

		// Delete specific account
		svc, err := resolveService()
		if err != nil {
			log.Fatal("unable to determine service", "error", err)
		}
		err = internal.DeleteFromKeyring(svc, args[0])
		if err != nil {
			log.Fatal("unable to delete account from keyring", "name", args[0], "error", err)
			return
		}
		log.Info("successfully deleted account from keyring", "name", args[0])
		err = internal.DeleteFromAccounts(internal.Viper.GetString("accounts"), args[0], svc)
		if err != nil {
			log.Fatal("unable to delete account from accounts list", "name", args[0], "error", err)
			return
		}
		log.Debug("successfully deleted account from accounts list", "name", args[0])
	},
}

func init() {
	deleteCmd.Flags().BoolP("delete-all", "a", false, "Delete all saved accounts (scoped to --service if given)")
	rootCmd.AddCommand(deleteCmd)
}
