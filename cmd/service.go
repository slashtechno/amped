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

	"github.com/slashtechno/amped/internal"
)

// resolveService returns the service from --service, falling back to the last-used service.
// Used by commands that require a specific service (add, switch, delete).
func resolveService() (internal.Service, error) {
	if raw := internal.Viper.GetString("service"); raw != "" {
		svc := internal.Service(raw)
		if svc != internal.ServiceAmp && svc != internal.ServiceClaude {
			return "", fmt.Errorf("unknown service %q; use amp or claude", raw)
		}
		return svc, nil
	}
	accounts, err := internal.ReadFromAccounts(internal.Viper.GetString("accounts"))
	if err != nil {
		return "", err
	}
	if accounts.LastService == "" {
		return "", fmt.Errorf("no service specified; use --service amp or --service claude")
	}
	return accounts.LastService, nil
}
