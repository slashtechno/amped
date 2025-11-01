package internal

import (
	"encoding/json"
	"os"

	"github.com/zalando/go-keyring"
)

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

func DeleteFromKeyring(user string) error {
	err := keyring.Delete("amped", user)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAllFromKeyring() error {
	err := keyring.DeleteAll("amped")
	if err != nil {
		return err
	}

	return nil
}

func EnsureAccounts(accountsPath string) error {
	_, err := os.Stat(accountsPath)
	if os.IsNotExist(err) {
		emptyAccounts := AmpedAccounts{
			Accounts: []AmpedAccount{},
			Active:   "",
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

func WriteToAccounts(accountsPath string, accounts AmpedAccounts) error {
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

func AppendToAccounts(accountsPath string, account AmpedAccount) error {
	accounts, err := ReadFromAccounts(accountsPath)
	if err != nil {
		return err
	}
	accounts.Accounts = append(accounts.Accounts, account)
	err = WriteToAccounts(accountsPath, accounts)
	if err != nil {
		return err
	}
	return nil
}

func UpdateActiveAccount(accountsPath, name string) error {
	accounts, err := ReadFromAccounts(accountsPath)
	if err != nil {
		return err
	}
	accounts.Active = name
	err = WriteToAccounts(accountsPath, accounts)
	if err != nil {
		return err
	}
	return nil
}

func DeleteFromAccounts(accountsPath, name string) error {
	accounts, err := ReadFromAccounts(accountsPath)
	if err != nil {
		return err
	}
	var updatedAccounts []AmpedAccount
	for _, account := range accounts.Accounts {
		if account.Name != name {
			updatedAccounts = append(updatedAccounts, account)
		}
	}
	err = WriteToAccounts(accountsPath, AmpedAccounts{Accounts: updatedAccounts})
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllThreads(threadsPath string) error {
	err := os.RemoveAll(threadsPath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(threadsPath, 0700)
	if err != nil {
		return err
	}
	return nil
}

func DeleteHistoryFile(historyPath string) error {
	err := os.Remove(historyPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
