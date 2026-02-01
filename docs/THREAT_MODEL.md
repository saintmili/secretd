# Threat Model

This document describes what secretd protects against and what it cannot protect against.


## Assets

Sensitive data protected by `secretd`:

- Account passwords
- API tokens
- Private notes
- Account metadata (usernames, URLs)
- Primary asset: encrypted vault file


## Trust Boundaries

Trusted:

- User operating system
- User hardware
- Go runtime
- Standard crypto libraries

Untrusted:

- Anyone with access to vault file
- Backup providers
- Cloud storage
- Other users of the same system


## Attacker Types

**1. Device thief**

Attacker steals laptop/drive and attempts offline attack.
Protection:

- Argon2id key derivation
- AES-256-GCM encryption
- Unique salt per vault

**2. Curious coworker / shared computer**

Attacker accesses user files without password.

- Protection:
- 0600 vault permissions
- Encryption at rest

**3. Backup compromise**

Vault stored in cloud backup is leaked.
Protection:

- Strong encryption
- No plaintext secrets stored

**4. Vault tampering attacker**

Attacker modifies vault file to corrupt or inject data.
Protection:

- Authentication tag (AES-GCM)
- Additional HMAC integrity check
- Verification before decryption

**5. Clipboard snooping applications**

Malicious app reads clipboard contents.
Mitigation:

- Clipboard auto-clear timer
- User configurable timeout
- Residual risk remains during clipboard lifetime.


## Out-of-Scope Threats

**1. Compromised OS**

If malware runs on the system, it can:

- Capture keystrokes
- Dump process memory
- Read clipboard contents
- Replace binaries

No password manager can defend against this scenario.

**2. Root / Administrator Access**

A privileged attacker can:

- Inspect memory
- Read files
- Modify binaries

This threat is explicitly out of scope.

**3. Hardware Attacks**

Not protected against:

- Cold boot attacks
- DMA attacks
- Hardware keyloggers


## Security Philosophy

`secretd` follows a simple philosophy:

- Protect data at rest and during normal use
- Assume the OS is trusted
- Fail securely when uncertain


## Residual Risks

Remaining unavoidable risks:

- Clipboard exposure window
- In-memory secrets during runtime
- User choosing weak master passwords

Users must understand these limitations.

---

Security is a shared responsibility.
