# Amped
The missing account switcher for Amp

![Demo](demo.gif) 

[Amp](http://ampcode.com/) is a coding agent that, so far, has been one of the most powerful I've used. I have both a work account, which is part of a paid Amp team, and a personal, free account. Issue is, Amp currently has no way of quickly switching between them. Amped is a CLI tool that extracts and stores Amp's API keys, allowing different API keys to be utilized depending on what account you need to use. API keys are stored securely in the system keyring.

## Installation
Either go to the [releases page](https://github.com/slashtechno/amped/releases) and download a binary for your platform, or use the following command to install it via Go:

```bash
go install github.com/slashtechno/amped@latest
```

## Usage
Amped has four commands, `add`, `switch`, `status`, and `delete`.
- `amped add name`: Adds a new Amp account to the system keyring and stores the name of the account in the accounts list. 
    - Uses Amp's API key found at `~/.local/share/amp/secrets.json` 
        - Instead of being in AppData, even Windows has it under `~/.local/share/amp/secrets.json`. Not sure about MacOS though. Feel free to open an [issue](https://github.com/slashtechno/amped/issues) or make a PR if you're able to confirm.
- `amped switch name`: Switches to the given saved Amp account by updating `~/.local/share/amp/secrets.json` with the API key from the keyring.
- `amped status`: Shows the currently active Amp account, as well as a list of all saved accounts.
- `amped delete name`: Deletes the given Amp account from the system keyring and accounts list
    - You can delete all accounts, including those that are missing from the accounts list but are still in the keyring, by using `amped delete --delete-all`.
    - This does not sign you out of Amp, it only removes the API key from the keyring and accounts list.
