# Contributing

Thank you for your interest in contributing to **secretd** ❤️

---

## Philosophy

`secretd` prioritizes:

1. Security over features
2. Simplicity over cleverness
3. Auditability over abstraction

If a feature complicates the security model, it will likely be rejected.

---

## Development Setup

Requirements:

* Go ≥ 1.25
* Linux/macOS

Clone repo:

```bash
git clone https://github.com/saintmili/secretd
cd secretd
go build ./...
```

Run locally:

```bash
go run ./cmd/secretd
```

---

## Coding Guidelines

### Security-first mindset

Every PR must answer:

* Does this increase attack surface?
* Does this expose secrets?
* Can this be abused?

If unsure → open a discussion first.

---

### Code style

* Small packages
* Small functions
* Explicit error handling
* No global state
* No hidden magic

Avoid:

* Reflection
* Unnecessary generics
* Over-engineering

---

### Logging rules

Never log:

* Passwords
* Vault contents
* Encryption keys
* Clipboard contents
* Master password attempts

Violating this rule will result in PR rejection.

---

## Commit Messages

Use Conventional Commits:

```
feat:
fix:
refactor:
docs:
security:
test:
```

Example:

```
security: add HMAC integrity verification
```

---

## Pull Request Process

1. Open an issue (for non-trivial changes)
2. Fork the repo
3. Create a feature branch
4. Submit PR with explanation
5. Security review
6. Merge

All PRs are reviewed manually.

---

## Areas That Need Help

* Tests
* Packaging
* Documentation
* UX improvements
* Security review

---

## Thank You

Security tools improve through community review.
Your help is appreciated 🙌

