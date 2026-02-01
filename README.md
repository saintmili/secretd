# secretd 

**secretd** is a minimalist, local-first, terminal-based password manager focused on **security, simplicity, and auditability**.

No cloud.
No sync.
No telemetry.
Just your secrets — encrypted on disk.

[![Go](https://img.shields.io/badge/go-1.25-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

---

## Features

* Strong encryption (Argon2id + AES-256-GCM)
* Local-first design
* XDG-compliant file layout
* Automatic clipboard clearing
* Tamper detection
* Single-instance locking
* Simple, scriptable CLI
* No network access

---

## Security Model

`secretd` protects against:

* Offline attacks on stolen vault files
* Vault tampering
* Accidental clipboard leaks

It does **not** protect against:

* Malware or compromised OS
* Root/admin attackers
* Keyloggers or screen capture

If your system is compromised, your secrets are compromised.

For details, see [`docs/SECURITY.md`](docs/SECURITY.md).

---

## Installation

### From source

```bash
git clone https://github.com/saintmili/secretd
cd secretd
make install
```

---

## Quick Start

Initialize vault:

```bash
secretd init
```

Add an entry:

```bash
secretd add
```

List entries:

```bash
secretd list
```

Show password:

```bash
secretd show github
```

---

## Configuration

Config file:

```
~/.config/secretd/config.toml
```

If missing, defaults are used.

See [`docs/CONFIG.md`](docs/CONFIG.md) for all options.

---

## File Locations

| Type   | Path                                |
| ------ | ----------------------------------- |
| Vault  | `~/.local/share/secretd/vault.json` |
| Config | `~/.config/secretd/config.toml`     |
| Logs   | `~/.local/state/secretd/`           |

All files use strict permissions.

---

## Documentation

* [Security Policy](docs/SECURITY.md)
* [Threat Model](docs/THREAT_MODEL.md)
* [Cryptography](docs/CRYPTO.md)
* [Configuration](docs/CONFIG.md)
* [Contributing](docs/CONTRIBUTING.md)

---

## Philosophy

* Security > convenience
* Explicit > magic
* Small codebase > features
* Easy to audit

If a feature complicates the threat model, it doesn’t belong.

---

## License

MIT License.
See [LICENSE](LICENSE).

---

> **Warning**
> This project is security-sensitive software.
> Always review the code before trusting it with real secrets.

