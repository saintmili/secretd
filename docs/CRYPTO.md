# Cryptography in secretd

This document explains the cryptographic design and algorithms used in `secretd`.


## Vault Encryption Flow

```text
User Master Password
│
▼
┌───────────────┐
│ Argon2id    │ ← Derives encryption key using salt
└───────────────┘
│ Key
▼
┌───────────────┐
│ AES-256-GCM │ ← Encrypts vault content + authentication tag
└───────────────┘
│ Ciphertext + Tag + Nonce + Salt
▼
Vault File
(stored 0600)
```


## Key Derivation

- **Algorithm:** Argon2id
- **Purpose:** Derive encryption key from master password
- **Default Parameters:**
  - Memory: 64 MB
  - Iterations: 3
  - Parallelism: 4 threads
- **Security goal:** Slow down brute-force attacks while remaining usable for legitimate users


## Encryption

- **Algorithm:** AES-256-GCM
- **Mode:** Galois/Counter Mode (GCM)
- **Purpose:**
  - AES-256 provides confidentiality
  - GCM provides authentication (prevents silent tampering)
- **Nonce:** Random per vault file, stored with ciphertext
- **Vault MAC:** Additional integrity check over the entire vault file


## Randomness

All randomness is generated using: `crypto/rand`

This uses the operating system’s secure RNG.



## Known Limitations

- Keys exist in process memory during runtime
- Go GC prevents guaranteed memory wiping
- No protection against runtime compromise


## Cryptography Philosophy

Use:

- Few primitives
- Well-audited standards
- No custom cryptography

**Never roll your own crypto.**
