package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

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

// ResolveClaudeConfigPath returns the path to Claude Code's config file.
// Newer Claude Code uses ~/.claude.json; older versions use ~/.claude/.claude.json.
func ResolveClaudeConfigPath(defaultPath string) string {
	legacy := strings.Replace(defaultPath, ".claude.json", ".claude/.claude.json", 1)
	if data, err := os.ReadFile(legacy); err == nil {
		var v map[string]json.RawMessage
		if json.Unmarshal(data, &v) == nil {
			if _, ok := v["oauthAccount"]; ok {
				return legacy
			}
		}
	}
	return defaultPath
}

// ExtractClaudeCredentials reads the current Claude Code credentials and oauthAccount from the config.
// On macOS, credentials are read from the system Keychain ("Claude Code-credentials").
// On Linux/other, credentials are read from the file at claudeCredsPath.
func ExtractClaudeCredentials(claudeConfigPath, claudeCredsPath string) (ClaudeStoredCredentials, error) {
	creds, err := readClaudeCredentialsBlob(claudeCredsPath)
	if err != nil {
		return ClaudeStoredCredentials{}, fmt.Errorf("unable to read Claude Code credentials: %w", err)
	}

	var oauthAccount string
	configData, err := os.ReadFile(ResolveClaudeConfigPath(claudeConfigPath))
	if err == nil && len(configData) > 0 {
		var config map[string]any
		if err := json.Unmarshal(configData, &config); err == nil {
			if oauth, ok := config["oauthAccount"]; ok {
				oauthBytes, _ := json.Marshal(oauth)
				oauthAccount = string(oauthBytes)
			}
		}
	}

	return ClaudeStoredCredentials{Credentials: creds, OAuthAccount: oauthAccount}, nil
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

// ExtractClaudeAccountDetails pulls the email and subscriptionType from stored credentials
// for display purposes, and returns the accessToken for use in live verification.
func ExtractClaudeAccountDetails(stored ClaudeStoredCredentials) (email, subscriptionType, accessToken string) {
	var credsBlob struct {
		ClaudeAiOauth struct {
			AccessToken      string `json:"accessToken"`
			SubscriptionType string `json:"subscriptionType"`
		} `json:"claudeAiOauth"`
	}
	if err := json.Unmarshal([]byte(stored.Credentials), &credsBlob); err == nil {
		accessToken = credsBlob.ClaudeAiOauth.AccessToken
		subscriptionType = credsBlob.ClaudeAiOauth.SubscriptionType
	}
	if stored.OAuthAccount != "" {
		var oauthAccount struct {
			EmailAddress string `json:"emailAddress"`
		}
		if err := json.Unmarshal([]byte(stored.OAuthAccount), &oauthAccount); err == nil {
			email = oauthAccount.EmailAddress
		}
	}
	return
}

// VerifyClaudeToken makes a live request to the Anthropic API to confirm the access token
// is valid, returning the organization name associated with the account.
func VerifyClaudeToken(accessToken string) (orgName string, err error) {
	req, err := http.NewRequest("GET", "https://api.anthropic.com/api/oauth/claude_cli/roles", nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}
	var body struct {
		OrgName string `json:"organization_name"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return
	}
	return body.OrgName, nil
}

