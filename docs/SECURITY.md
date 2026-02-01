# Security Policy 

## Supported Versions

The **latest release** of `secretd` is the only supported version.

Security fixes are not backported to older versions.

Users are strongly encouraged to upgrade immediately when a new release is published.


## Reporting a Vulnerability

If you discover a security vulnerability, **please do not open a public issue**.

Instead, report it privately:

- **Email:** m.ahrari77@gmail.com  
- **PGP Key(optional):** [Optional PGP key fingerprint for encrypted reports]  
- **Guidelines:**  
  - A clear description of the issue
  - Steps to reproduce
  - Impact assessment (if known)
  - Proof-of-concept (if possible)
  - Version of `secretd` and OS environment.  
  - Do **not** post the vulnerability publicly before a fix is released.  

We aim to acknowledge all reports within **72 hours** and provide a fix or mitigation as soon as possible.


## Security Model Summary

`secretd` is a password manager designed to protect secrets stored on disk.

It provides protection against:

- Stolen or lost devices
- Offline brute-force attacks on the vault file
- Vault tampering and corruption
- Accidental clipboard leaks

It is **not designed** to protect against:

- Malware or a compromised operating system
- Malicious root or administrator users
- Live memory attacks
- Hardware keyloggers or screen capture tools

If your OS is compromised, your secrets are compromised.


## Cryptography Overview

|Purpose|Algorithm|
|--|--|
|Key derivation|Argon2id|
|Encryption|AES-256-GCM|
|Additional authentication|HMAC-SHA-256|
|Randomness|crypto/rand|

Keys are derived from the master password using a unique salt and memory-hard parameters.

Vault integrity is verified before decryption to detect tampering.


## Data Storage

Vault location: `~/.local/share/secretd/vault.json`

File permissions are enforced: `0600 (owner read/write only)`

Logs and config are stored separately.


## Memory Safety Disclaimer

`secretd` attempts to wipe sensitive buffers (passwords, keys) from memory after use.

However, due to Go’s garbage collector and runtime behavior:

**Complete memory zeroing cannot be guaranteed.**

Secrets may temporarily exist in memory until reclaimed by the runtime.

This is a known limitation of memory-safe managed languages.


## Clipboard Risk Notice

Passwords copied to clipboard are automatically cleared after a timeout.

However, during that window:

- Other applications may read clipboard contents
- Clipboard managers may persist history

Clipboard use is optional and configurable.


## Export Warning

Exported vault data (CSV/JSON) excludes passwords, but may still contain:

- Usernames
- URLs
- Titles
- Notes

This metadata may be sensitive.

Treat exported files as confidential.


## Hardening Recommendations

For maximum security:

- Use a strong, unique master password
- Keep your OS updated
- Enable full-disk encryption
- Avoid clipboard managers
- Lock your screen when away
- Backup the vault securely


## Additional Resources

- Threat model: [THREAT_MODEL.md](THREAT_MODEL.md)  
- Cryptography documentation: [CRYPTO.md](CRYPTO.md)  

---

Thank you for helping keep secretd secure ❤️
