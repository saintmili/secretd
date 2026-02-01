package app

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/saintmili/secretd/internal/config"
	"github.com/saintmili/secretd/internal/crypto"
	"github.com/saintmili/secretd/internal/doctor"
	"github.com/saintmili/secretd/internal/storage"
	"github.com/saintmili/secretd/internal/vault"
)

func PrintUsage() error {
	fmt.Println(`secretd 🔐 — password manager

		Usage:
		  secretd init
		  secretd unlock
		  secretd add <title> [--generate <length>]
		  secretd list
		  secretd show <title> [--reveal]
		  secretd update <title>
		  secretd delete <title>
		  secretd change-master-password
		  secretd generate [length]
		  secretd export <json/csv>
		  secretd config <show|edit>
		  secretd version
		`)
	return nil
}

func Init(app *App) error {
	if err := app.Storage.Load(); err == nil {
		return ErrVaultExists
	}

	pw1, err := app.Vault.ReadPassword("Set master password: ")
	if err != nil {
		return ErrFailedReadPassword
	}
	pw2, err := app.Vault.ReadPassword("Confirm master password: ")
	if err != nil {
		return ErrFailedReadPassword
	}
	defer zeroBytes(pw1)
	defer zeroBytes(pw2)

	if !bytes.Equal(pw1, pw2) {
		return ErrMismatchPassword
	}

	if len(pw1) == 0 {
		return ErrEmptyPassword
	}

	salt, err := app.Vault.GenerateSalt()
	if err != nil {
		return err
	}
	defer zeroBytes(salt)

	encKey, macKey, err := app.Vault.DeriveKeys(pw1, salt)
	if err != nil {
		return err
	}
	defer zeroBytes(encKey)
	defer zeroBytes(macKey)

	v := vault.New()

	sess := &Session{
		Vault:  v,
		EncKey: encKey,
		MacKey: macKey,
		File: &storage.VaultFile{
			Version: 1,
			Salt:    crypto.EncodeSalt(salt),
		},
	}
	defer sess.Close()

	if err := app.SaveSession(sess); err != nil {
		return err
	}

	fmt.Println("Vault initialized successfully 🔐")
	app.Logger.Info("Vault initialized")
	return nil
}

func Unlock(app *App) error {
	sess, err := app.OpenSession()
	if err != nil {
		return err
	}
	defer sess.Close()
	fmt.Printf("Vault unlocked 🔓 (%d entries)\n", len(sess.Vault.Entries))
	return nil
}

func Add(app *App) error {
	if len(os.Args) < 3 {
		return fmt.Errorf("Usage: secretd add <title> [--generate <length>]")
	}

	title := os.Args[2]
	sess, err := app.OpenSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	var pw []byte
	defer zeroBytes(pw)

	autoClipboard := false

	generate, length := parseGenerateArgs(os.Args[3:])
	if generate {
		pw, err = generatePassword(length)
		if err != nil {
			return err
		}
		fmt.Printf("Generated password (%d chars)\n", len(pw))
		autoClipboard = true
	} else {
		pw, err = crypto.ReadPassword("Password: ")
		if err != nil {
			return ErrFailedReadPassword
		}
	}

	entry := &vault.Entry{
		ID:       uuid.NewString(),
		Title:    title,
		Username: readLine("Username: "),
		Password: append([]byte(nil), pw...),
		URL:      readLine("URL: "),
		Notes:    readLine("Notes: "),
	}

	sess.Vault.Entries = append(sess.Vault.Entries, entry)

	if err := app.SaveSession(sess); err != nil {
		return err
	}

	if autoClipboard {
		if err := app.Clipboard.Copy(string(entry.Password)); err != nil {
			fmt.Println("[!] Clipboard error:", err)
		} else {
			fmt.Printf("[+] Generated password copied to clipboard (clears in %ds)\n", app.Clipboard.Timeout)
		}
	}

	fmt.Println("Entry added successfully 🔐")
	return nil
}

func List(app *App) error {
	sess, err := app.OpenSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	for _, e := range sess.Vault.Entries {
		fmt.Println("-", e.Title)
	}

	return nil
}

func Show(app *App) error {
	if len(os.Args) < 3 {
		return fmt.Errorf("Usage: secretd show <title> [--reveal]")
	}

	title := os.Args[2]
	reveal := len(os.Args) >= 4 && os.Args[3] == "--reveal"

	sess, err := app.OpenSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	for _, e := range sess.Vault.Entries {
		if strings.EqualFold(e.Title, title) {
			pwCopy := append([]byte(nil), e.Password...) // temp buffer to zero
			defer zeroBytes(pwCopy)

			fmt.Println("Title:", e.Title)
			fmt.Println("Username:", e.Username)
			if reveal {
				fmt.Println("Password:", string(pwCopy))
			} else {
				fmt.Println("Password:", strings.Repeat("*", 16))
			}
			fmt.Println("URL:", e.URL)
			fmt.Println("Notes:", e.Notes)
			if err := app.Clipboard.Copy(string(e.Password)); err != nil {
				fmt.Println("[!] Clipboard error:", err)
			} else {
				fmt.Printf("[+] Generated password copied to clipboard (clears in %ds)\n", app.Clipboard.Timeout)
			}

			return nil
		}
	}

	return fmt.Errorf("Entry not found")
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?/~"

func Generate() error {
	length := 16
	if len(os.Args) >= 3 {
		fmt.Sscanf(os.Args[2], "%d", &length)
	}

	pw, err := generatePassword(length)
	if err != nil {
		return err
	}
	defer zeroBytes(pw)

	fmt.Println(string(pw))
	return nil
}

func Update(app *App) error {
	if len(os.Args) < 3 {
		return fmt.Errorf("Usage: secretd update <title> [--generate <length>]")
	}

	title := os.Args[2]
	sess, err := app.OpenSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	index := -1
	for i := range sess.Vault.Entries {
		if strings.EqualFold(sess.Vault.Entries[i].Title, title) {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("Entry not found: %v", title)
	}

	entry := sess.Vault.Entries[index]
	fmt.Println("Updating entry:", entry.Title)
	fmt.Println("(Press Enter to keep current value)")

	entry.Username = readOptional("Username", entry.Username)
	entry.URL = readOptional("URL", entry.URL)
	entry.Notes = readOptional("Notes", entry.Notes)

	// Password handling
	fmt.Print("Change password? (y/N): ")
	ans, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	ans = strings.TrimSpace(ans)
	if strings.ToLower(ans) == "y" {
		var pw []byte
		defer zeroBytes(pw)

		generate, length := parseGenerateArgs(os.Args[3:])
		if generate {
			pw, err = generatePassword(length)
			if err != nil {
				return err
			}
			entry.Password = append([]byte(nil), pw...)

			if err := app.Clipboard.Copy(string(entry.Password)); err == nil {
				fmt.Printf("[+] Generated password copied to clipboard (clears in %ds)\n", app.Clipboard.Timeout)
			}
			fmt.Printf("Generated password (%d chars)\n", len(pw))
		} else {
			pw, err = crypto.ReadPassword("New password: ")
			if err != nil {
				return ErrFailedReadPassword
			}
			entry.Password = append([]byte(nil), pw...)
		}
	}

	sess.Vault.Entries[index] = entry

	if err := app.SaveSession(sess); err != nil {
		return err
	}

	fmt.Println("Entry updated successfully ✏️")
	return nil
}

func Delete(app *App) error {
	if len(os.Args) < 3 {
		return fmt.Errorf("Usage: secretd delete <title>")
	}

	title := os.Args[2]
	sess, err := app.OpenSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	index := -1
	for i := range sess.Vault.Entries {
		if strings.EqualFold(sess.Vault.Entries[i].Title, title) {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("Entry not found: %v", title)
	}

	e := sess.Vault.Entries[index]
	fmt.Println("About to delete entry:")
	fmt.Println("Title:", e.Title)
	fmt.Println("Username:", e.Username)
	fmt.Println("URL:", e.URL)

	fmt.Print("Type 'yes' to confirm deletion: ")
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "yes" {
		return fmt.Errorf("Deletion cancelled.")
	}

	// Remove entry
	sess.Vault.Entries = append(sess.Vault.Entries[:index], sess.Vault.Entries[index+1:]...)

	if err := app.SaveSession(sess); err != nil {
		return err
	}

	fmt.Println("Entry deleted permanently 🗑️")
	return nil
}

func ChangeMasterPassword(app *App) error {
	fmt.Println("🔐 Change Master Password")

	// Load vault file (without decrypting yet)
	if err := app.Storage.Load(); err != nil {
		return fmt.Errorf("Failed to load vault: %v", err)
	}

	// Ask for current master password
	oldPW, err := crypto.ReadPassword("Current master password: ")
	if err != nil {
		return ErrFailedReadPassword
	}
	defer zeroBytes(oldPW)

	salt, err := crypto.DecodeSalt(app.Storage.VaultFile.Salt)
	if err != nil {
		return err
	}
	defer zeroBytes(salt)

	// Derive old key
	oldEncKey, oldMacKey, err := app.Vault.DeriveKeys(oldPW, salt)
	if err != nil {
		return err
	}
	defer zeroBytes(oldEncKey)
	defer zeroBytes(oldMacKey)

	// MAC verification
	if !app.Vault.VerifyMAC(oldMacKey, app.Storage.VaultFile.MAC, app.Storage.VaultFile.Nonce, app.Storage.VaultFile.Ciphertext) {
		return ErrFailedVaultIntegrity
	}

	// Decrypt vault
	plain, err := app.Vault.Decrypt(oldEncKey, app.Storage.VaultFile.Nonce, app.Storage.VaultFile.Ciphertext)
	if err != nil {
		return ErrWrongMasterPassword
	}
	defer zeroBytes(plain)

	// Parse vault
	v := &vault.Vault{}
	if err := json.Unmarshal(plain, &v); err != nil {
		return ErrVaultCorrupted
	}

	// Ask for new password twice
	newPW1, err := crypto.ReadPassword("New master password: ")
	if err != nil {
		return ErrFailedReadPassword
	}
	newPW2, err := crypto.ReadPassword("Confirm new master password: ")
	if err != nil {
		return ErrFailedReadPassword
	}
	defer zeroBytes(newPW1)
	defer zeroBytes(newPW2)

	if !bytes.Equal(newPW1, newPW2) {
		return ErrMismatchPassword
	}

	// Generate new salt
	newSalt, err := app.Vault.GenerateSalt()
	if err != nil {
		return err
	}
	defer zeroBytes(newSalt)

	// Derive new key
	newEncKey, newMacKey, err := app.Vault.DeriveKeys(newPW1, newSalt)
	if err != nil {
		return err
	}
	defer zeroBytes(newEncKey)
	defer zeroBytes(newMacKey)

	// Use saveVault to re-encrypt and save
	newSess := &Session{
		Vault:  v,
		EncKey: newEncKey,
		MacKey: newMacKey,
		File: &storage.VaultFile{
			Version: app.Storage.VaultFile.Version,
			Salt:    app.Vault.EncodeSalt(newSalt),
		},
	}

	if err := app.SaveSession(newSess); err != nil {
		return err
	}

	fmt.Println("✅ Master password changed successfully")
	app.Logger.Info("Master password changed")
	return nil
}

func Doctor(app *App) error {
	return doctor.Run(app.Config)
}

func Export(app *App) error {
	if len(os.Args) < 3 {
		return fmt.Errorf("Usage: secretd export <json|csv>")
	}

	format := strings.ToLower(os.Args[2])

	sess, err := app.OpenSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	// Create a read-only copy without passwords
	type ExportEntry struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Username string `json:"username"`
		URL      string `json:"url"`
		Notes    string `json:"notes"`
	}

	exportEntries := make([]ExportEntry, len(sess.Vault.Entries))
	for i, e := range sess.Vault.Entries {
		exportEntries[i] = ExportEntry{
			ID:       e.ID,
			Title:    e.Title,
			Username: e.Username,
			URL:      e.URL,
			Notes:    e.Notes,
		}
	}

	switch format {
	case "json":
		out, _ := json.MarshalIndent(exportEntries, "", "  ")
		os.WriteFile("vault_export.json", out, 0o600)
		fmt.Println("Vault exported to vault_export.json (read-only)")
		fmt.Println("⚠️  Exported file is NOT encrypted. Store it securely.")
	case "csv":
		var buf bytes.Buffer
		buf.WriteString("ID,Title,Username,URL,Notes\n")
		for _, e := range exportEntries {
			fmt.Fprintf(&buf, "%s,%q,%q,%q,%q\n", e.ID, e.Title, e.Username, e.URL, e.Notes)
		}
		os.WriteFile("vault_export.csv", buf.Bytes(), 0o600)
		fmt.Println("Vault exported to vault_export.csv (read-only)")
		fmt.Println("⚠️  Exported file is NOT encrypted. Store it securely.")
	default:
		return fmt.Errorf("Unsupported format: %v", format)
	}

	app.Logger.Info(fmt.Sprintf("Vault exported as %s", format))
	return nil
}

func Config(cfg config.Config) error {
	if len(os.Args) < 3 {
		return fmt.Errorf(`
			Usage:
				secretd config show   # print effective configuration
				secretd config edit   # edit config file in $EDITOR
		`)
	}
	switch os.Args[2] {
	case "show":
		cfg, _, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Failed to load config:", err)
			os.Exit(1)
		}

		out, err := cfg.MarshalPretty()
		if err != nil {
			fmt.Println("Failed to render config:", err)
			os.Exit(1)
		}

		fmt.Println(out)
	case "edit":
		path, err := config.ConfigPath()
		if err != nil {
			fmt.Println("Failed to locate config:", err)
			os.Exit(1)
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim"
		}

		cmd := exec.Command(editor, path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Println("Failed to open editor:", err)
			fmt.Println("Try setting $EDITOR")
			os.Exit(1)
		}
	default:
		return fmt.Errorf(`
			Usage:
				secretd config show   # print effective configuration
				secretd config edit   # edit config file in $EDITOR
		`)
	}

	return nil
}
