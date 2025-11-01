package internal

type AmpSecrets struct {
	APIKeyHTTPSAmpcodeCom string `json:"apiKey@https://ampcode.com/"`
}

type AmpedAccount struct {
	Name string `json:"name"`
}

type AmpedAccounts struct {
	Accounts []AmpedAccount `json:"accounts"`
	Active   string         `json:"active"`
}
