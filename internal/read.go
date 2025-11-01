package internal

import (
	"encoding/json"
	"os"

	"github.com/zalando/go-keyring"
)

func ReadFromKeyring(user string) (string, error) {
	retrievedApiKey, err := keyring.Get("amped", user)
	if err != nil {
		return "", err
	}
	return retrievedApiKey, nil
}

func ExtractApiKey(secretsPath string) (string, error) {
	dataB, err := os.ReadFile(secretsPath)
	if err != nil {
		return "", err
	}
	var secrets AmpSecrets
	err = json.Unmarshal(dataB, &secrets)
	if err != nil {
		return "", err
	}
	return secrets.APIKeyHTTPSAmpcodeCom, nil
}
