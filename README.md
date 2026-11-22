# Trimlet Auth

The two-factor codes (TOTP) of your accounts on the computer, next to the phone app, in a small window
driven only by the keyboard.

- When a service turns on two-factor authentication, paste its secret key or `otpauth://` address here
  as well as scanning the QR code with the phone: both then make the same codes.
- Enter copies a code, and it leaves the clipboard after 30 seconds.
- The accounts are in one file encrypted with a master password (Argon2id + AES-256-GCM),
  `trimletauth/vault.trimletauth` in the app settings folder; the system keychain remembers the
  password, so it is typed only the first time and on a new computer.
- **Sync between your computers** (optional) through trimlet.dev, with an invite code: the server only
  keeps the encrypted vault and never sees the master password. Sync through a private GitHub
  repository is being set up.

## Download

The app for **macOS** (Mac with Apple chip: M1 and later) and **Windows** is on the
[Releases](https://github.com/trimlet/auth/releases) page. It is a beta (0.x), and not signed yet: the first
time, macOS asks to confirm (right click on the app → Open), Windows shows SmartScreen (More info →
Run anyway).

## Build

Written in Go with [Fyne](https://fyne.io). To build it, with Go installed (and on a Mac the Xcode
command line tools):

```
go run fyne.io/tools/cmd/fyne@v1.7.3 package --release
```

or `go run .` to try it. The console window comes from [trimlet/gui](https://github.com/trimlet/gui),
which `go` fetches by itself.

## How it is made

The code is written with Claude Opus (Anthropic), under the direction and review of a developer with
almost 30 years of experience. Every feature starts from a real, repeated need, and most ideas are
turned down: the goal is to stay small.

## License

GPL-3.0, see [LICENSE](LICENSE).
