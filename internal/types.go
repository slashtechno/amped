package internal

// Service identifies which coding tool's accounts are being managed.
type Service string

const (
	ServiceAmp    Service = "amp"
	ServiceClaude Service = "claude"
)

// AmpSecrets represents the structure of Amp's secrets.json.
type AmpSecrets struct {
	APIKeyHTTPSAmpcodeCom string `json:"apiKey@https://ampcode.com/"`
}

// ClaudeStoredCredentials packages the credentials blob and OAuth account info
// needed to fully restore a Claude Code account.
type ClaudeStoredCredentials struct {
	// Credentials is the raw JSON blob from Claude Code's credential storage.
	// On macOS this comes from the Keychain entry "Claude Code-credentials".
	// On Linux this comes from ~/.claude/.credentials.json.
	Credentials string `json:"credentials"`
	// OAuthAccount is the oauthAccount JSON section from ~/.claude.json.
	// On switch, this is merged back into the live config, leaving other settings intact.
	OAuthAccount string `json:"oauthAccount,omitempty"`
}

// Account represents a single saved account entry in amped's accounts file.
type Account struct {
	Name    string  `json:"name"`
	Service Service `json:"service"`
}

// Accounts is the top-level structure of amped's accounts JSON file.
type Accounts struct {
	Accounts     []Account `json:"accounts"`
	ActiveAmp    string    `json:"activeAmp"`
	ActiveClaude string    `json:"activeClaude"`
	// LastService is used when no --service flag is given, so the last-used service is re-applied automatically.
	LastService Service `json:"lastService"`
}
