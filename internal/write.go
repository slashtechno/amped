package internal

import (
	"encoding/json"
	"os"

	"github.com/zalando/go-keyring"
)

func WriteToAmpSecrets(secretsPath, apiKey string) error {
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

func DeleteFromKeyring(user string) error {
	err := keyring.Delete("amped", user)
	if err != nil {
		return err
	}
	return nil
}
