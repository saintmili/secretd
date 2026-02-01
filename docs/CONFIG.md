# Configuration

This document explains where `secretd` stores its files and how configuration works.

---

## Directory Layout

By default, `secretd` follows the **XDG Base Directory Specification**.

| Type   | Path                                |
| ------ | ----------------------------------- |
| Vault  | `~/.local/share/secretd/vault.json` |
| Config | `~/.config/secretd/config.toml`     |
| Logs   | `~/.local/state/secretd/`           |

If XDG variables are set, they are respected:

* `$XDG_DATA_HOME`
* `$XDG_CONFIG_HOME`
* `$XDG_STATE_HOME`

---

## Config File

Config file path:

```
~/.config/secretd/config.toml
```

If the file does not exist, defaults are used.

---

## Example Config

```toml
[clipboard]
clear_after_seconds = 15
wayland_backend = "wl-copy"
x11_backend = "xclip"

[vault]
path = "$HOME/.local/share/secretd/vault.json"

[security]
wipe_memory = true
argon2_time = 3
argon2_memory = 65536
argon2_threads = 4
salt_length = 16
key_length = 32
max_failed_unlocks = 5
lockout_seconds = 300

[logging]
enabled = true
file = "$HOME/.local/share/secretd/secretd.log"
```

---

## Clipboard Settings

### clear_after_sconds

How long passwords stay in clipboard before being cleared.

Recommended values:

* 10–30 seconds for laptops
* 5–10 seconds for shared computers

Set to `0` to disable auto-clear (NOT recommended).

### wayland_backend

Wayland clipboard backend.

### x11_backend

X11 clipboard backend

## Vault Settings

## path

Path of vault.

## Security Settings

### argon2_time

Argon2id itterations.

### argon2_memory

Argon2id memory.

### argon2_threads

Argon2id threads.

### salt_length

Length of salt.

### key_length

Length of encKey and macKey.

### max_failed_unlocks

Max attemps to login before lock.

### lockout_seconds

How long vault be locked when max_failed_unlcoks reached in seconds.

## Logging Settings

### enabled

Enables logging service.

Set to `false` if you dont want logging.

### file

Path of log file.

---

## Permissions

Created files use secure permissions:

| File   | Mode           |
| ------ | -------------- |
| Vault  | 0600           |
| Config | 0600           |
| Logs   | 0700 directory |

If permissions are unsafe, `secretd` will refuse to run.

---

## Resetting Configuration

To reset to defaults:

```
rm ~/.config/secretd/config.toml
```

The file will be recreated on next run.

