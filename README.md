# Amped
The missing account switcher for Amp and Claude Code

![Demo](demo.gif)

[Amp](http://ampcode.com/) and [Claude Code](https://claude.ai/code) are coding agents with no built-in way to switch between accounts. Amped saves and restores credentials for both, storing secrets securely in the system keyring.

## Installation
Either go to the [releases page](https://github.com/slashtechno/amped/releases) and download a binary for your platform, or use the following command to install it via Go:

```bash
go install github.com/slashtechno/amped@latest
```

## Usage
Amped has four commands: `add`, `switch`, `status`, and `delete`. All commands accept `--service amp` or `--service claude` (shorthand `-S`). If `--service` is omitted, the last-used service is reused automatically.

- `amped add name --service <amp|claude>`: Saves the currently logged-in account to the keyring under the given name.
    - For Amp, reads the API key from `~/.local/share/amp/secrets.json`.
    - For Claude Code, reads credentials from the system keychain (macOS) or `~/.claude/.credentials.json` (Linux). Windows is not supported for Claude Code accounts.
- `amped switch name`: Restores credentials for a previously saved account, making it active. For Claude Code, also reads account info from `~/.claude.json` (or `~/.claude/.claude.json` on older versions) and merges it back on switch.
- `amped status`: Shows all saved accounts and which are currently active. Use `--service` to filter.
- `amped delete name`: Removes the given account from the keyring and accounts list.
    - `amped delete --delete-all` removes all saved accounts entirely.
    - This does not sign you out of Amp or Claude Code, it only removes the saved credentials from amped.
