package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/zalando/go-keyring"
)

// ReadFromKeyring retrieves stored credentials for the given service and account name.
// Falls back to the legacy (no-prefix) key format for backward compatibility with old Amp accounts.
func ReadFromKeyring(service Service, name string) (string, error) {
	return keyring.Get("amped", string(service)+":"+name)
}

// ReadFromAccounts reads and parses the amped accounts JSON file.
// Migrates legacy single-service format (Active field) to the new per-service format on read.
func ReadFromAccounts(accountsPath string) (Accounts, error) {
	dataB, err := os.ReadFile(accountsPath)
	if err != nil {
		return Accounts{}, err
	}
	var accounts Accounts
	if err = json.Unmarshal(dataB, &accounts); err != nil {
		return Accounts{}, err
	}

	return accounts, nil
}

// ExtractApiKey reads the Amp API key from Amp's secrets.json.
func ExtractApiKey(secretsPath string) (string, error) {
	dataB, err := os.ReadFile(secretsPath)
	if err != nil {
		return "", err
	}
	var secrets AmpSecrets
	if err = json.Unmarshal(dataB, &secrets); err != nil {
		return "", err
	}
	return secrets.APIKeyHTTPSAmpcodeCom, nil
}

// ExtractClaudeCredentials reads the current Claude Code credentials and oauth account info.
// On macOS, credentials are read from the system Keychain ("Claude Code-credentials").
// On Linux/other, credentials are read from the file at claudeCredsPath.
func ExtractClaudeCredentials(claudeConfigPath, claudeCredsPath string) (ClaudeStoredCredentials, error) {
	creds, err := readClaudeCredentialsBlob(claudeCredsPath)
	if err != nil {
		return ClaudeStoredCredentials{}, fmt.Errorf("unable to read Claude Code credentials: %w", err)
	}

	oauth, err := readClaudeOAuthAccount(claudeConfigPath)
	if err != nil {
		return ClaudeStoredCredentials{}, fmt.Errorf("unable to read Claude Code oauth account: %w", err)
	}

	return ClaudeStoredCredentials{Credentials: creds, OAuthAccount: oauth}, nil
}

// readClaudeCredentialsBlob returns the raw credentials JSON string from Claude Code's storage.
func readClaudeCredentialsBlob(credsFilePath string) (string, error) {
	if runtime.GOOS == "darwin" {
		out, err := exec.Command("security", "find-generic-password", "-s", "Claude Code-credentials", "-w").Output()
		if err != nil {
			// Exit code 44 means "item not found" in the macOS security tool
			if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 44 {
				return "", fmt.Errorf("no Claude Code credentials found in keychain; make sure you are logged in to Claude Code")
			}
			return "", fmt.Errorf("security command failed: %w", err)
		}
		return strings.TrimSpace(string(out)), nil
	}

	// Linux/other: read from file
	data, err := os.ReadFile(credsFilePath)
	if err != nil {
		return "", fmt.Errorf("unable to read %s: %w", credsFilePath, err)
	}
	return strings.TrimSpace(string(data)), nil
}

// readClaudeOAuthAccount reads the oauthAccount section from ~/.claude/.claude.json.
func readClaudeOAuthAccount(claudeConfigPath string) (*ClaudeOAuthAccount, error) {
	data, err := os.ReadFile(claudeConfigPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var config struct {
		OAuthAccount *ClaudeOAuthAccount `json:"oauthAccount"`
	}
	if err = json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return config.OAuthAccount, nil
}
