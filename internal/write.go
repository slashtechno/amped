package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/zalando/go-keyring"
)

// WriteToKeyring stores credentials for the given service and account name.
func WriteToKeyring(service Service, name, value string) error {
	return keyring.Set("amped", string(service)+":"+name, value)
}

// DeleteFromKeyring removes the keyring entry for the given service and account name.
func DeleteFromKeyring(service Service, name string) error {
	return keyring.Delete("amped", string(service)+":"+name)
}

// DeleteAllFromKeyring deletes all keyring entries under the "amped" service.
func DeleteAllFromKeyring() error {
	return keyring.DeleteAll("amped")
}

// EnsureAccounts creates the accounts JSON file with an empty structure if it doesn't exist.
func EnsureAccounts(accountsPath string) error {
	_, err := os.Stat(accountsPath)
	if os.IsNotExist(err) {
		emptyAccounts := Accounts{
			Accounts: []Account{},
		}
		dataB, err := json.MarshalIndent(emptyAccounts, "", "  ")
		if err != nil {
			return err
		}
		err = os.WriteFile(accountsPath, dataB, 0600)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteToAccounts writes the full accounts structure to the accounts JSON file.
func WriteToAccounts(accountsPath string, accounts Accounts) error {
	dataB, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(accountsPath, dataB, 0600)
	if err != nil {
		return err
	}
	return nil
}

// AppendToAccounts adds a new account entry to the accounts file and updates LastService.
func AppendToAccounts(accountsPath string, account Account) error {
	accounts, err := ReadFromAccounts(accountsPath)
	if err != nil {
		return err
	}
	accounts.Accounts = append(accounts.Accounts, account)
	accounts.LastService = account.Service
	err = WriteToAccounts(accountsPath, accounts)
	if err != nil {
		return err
	}
	return nil
}

// UpdateActiveAccount sets the active account for the given service and updates LastService.
func UpdateActiveAccount(accountsPath, name string, service Service) error {
	accounts, err := ReadFromAccounts(accountsPath)
	if err != nil {
		return err
	}
	switch service {
	case ServiceAmp:
		accounts.ActiveAmp = name
	case ServiceClaude:
		accounts.ActiveClaude = name
	default:
		return fmt.Errorf("unsupported service: %s", service)
	}
	accounts.LastService = service
	err = WriteToAccounts(accountsPath, accounts)
	if err != nil {
		return err
	}
	return nil
}

// DeleteFromAccounts removes the named account for the given service from the accounts file.
func DeleteFromAccounts(accountsPath, name string, service Service) error {
	accounts, err := ReadFromAccounts(accountsPath)
	if err != nil {
		return err
	}
	var updatedAccounts []Account
	for _, account := range accounts.Accounts {
		if !(account.Name == name && account.Service == service) {
			updatedAccounts = append(updatedAccounts, account)
		}
	}
	accounts.Accounts = updatedAccounts
	// Clear the active reference if the deleted account was the active one
	switch service {
	case ServiceAmp:
		if accounts.ActiveAmp == name {
			accounts.ActiveAmp = ""
		}
	case ServiceClaude:
		if accounts.ActiveClaude == name {
			accounts.ActiveClaude = ""
		}
	}
	err = WriteToAccounts(accountsPath, accounts)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAllFromAccounts clears all accounts from the accounts file.
func DeleteAllFromAccounts(accountsPath string) error {
	return WriteToAccounts(accountsPath, Accounts{})
}

// WriteToAmpSecrets writes an Amp API key to Amp's secrets.json.
func WriteToAmpSecrets(apiKey, secretsPath string) error {
	secrets := AmpSecrets{
		APIKeyHTTPSAmpcodeCom: apiKey,
	}
	dataB, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(secretsPath, dataB, 0600)
	if err != nil {
		return err
	}
	return nil
}

// WriteToClaudeCredentials restores Claude Code credentials and oauth account info.
// On macOS, credentials are written to the system Keychain.
// On Linux/other, credentials are written to the file at claudeCredsPath.
func WriteToClaudeCredentials(stored ClaudeStoredCredentials, claudeConfigPath, claudeCredsPath string) error {
	if err := writeClaudeCredentialsBlob(stored.Credentials, claudeCredsPath); err != nil {
		return fmt.Errorf("unable to write Claude Code credentials: %w", err)
	}
	if stored.OAuthAccount != nil {
		if err := writeClaudeOAuthAccount(claudeConfigPath, stored.OAuthAccount); err != nil {
			return fmt.Errorf("unable to write Claude Code oauth account: %w", err)
		}
	}
	return nil
}

// writeClaudeCredentialsBlob writes the raw credentials JSON string to Claude Code's storage.
func writeClaudeCredentialsBlob(credentials, credsFilePath string) error {
	if runtime.GOOS == "darwin" {
		// -U updates the entry if it already exists
		username := os.Getenv("USER")
		if username == "" {
			username = "user"
		}
		cmd := exec.Command("security", "add-generic-password", "-U",
			"-s", "Claude Code-credentials",
			"-a", username,
			"-w", credentials,
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("security command failed: %w (output: %s)", err, strings.TrimSpace(string(out)))
		}
		return nil
	}

	// Linux/other: write to file
	if err := os.MkdirAll(filepath.Dir(credsFilePath), 0700); err != nil {
		return err
	}
	return os.WriteFile(credsFilePath, []byte(credentials), 0600)
}

// writeClaudeOAuthAccount updates only the oauthAccount field in ~/.claude/.claude.json,
// preserving all other existing fields in the file.
func writeClaudeOAuthAccount(claudeConfigPath string, account *ClaudeOAuthAccount) error {
	// Read the existing config as a raw map so we don't clobber unrelated fields
	config := make(map[string]interface{})
	data, err := os.ReadFile(claudeConfigPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if len(data) > 0 {
		if err = json.Unmarshal(data, &config); err != nil {
			return err
		}
	}

	config["oauthAccount"] = account

	out, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(claudeConfigPath), 0700); err != nil {
		return err
	}
	return os.WriteFile(claudeConfigPath, out, 0600)
}

// DeleteAllThreads removes all files in the threads directory and recreates it empty.
func DeleteAllThreads(threadsPath string) error {
	if err := os.RemoveAll(threadsPath); err != nil {
		return err
	}
	return os.MkdirAll(threadsPath, 0700)
}

// DeleteHistoryFile removes the history file, if it exists.
func DeleteHistoryFile(historyPath string) error {
	err := os.Remove(historyPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
