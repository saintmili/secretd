package app

import (
	"encoding/json"
	"fmt"

	"github.com/saintmili/secretd/internal/config"
	"github.com/saintmili/secretd/internal/logger"
	"github.com/saintmili/secretd/internal/security"
	"github.com/saintmili/secretd/internal/vault"
)

type App struct {
	Config *config.Config
	Logger *logger.Logger

	Vault     *VaultService
	Clipboard *ClipboardService
	Storage   *StorageService
}

func New(cfg *config.Config) (*App, error) {
	logger, err := logger.New(cfg.Logging)
	if err != nil {
		return nil, err
	}
	return &App{
		Config: cfg,
		Logger: logger,

		Vault: NewVaultService(cfg),
		Clipboard: NewClipboardService(
			cfg.Clipboard.ClearAfterSeconds,
		),
		Storage: NewStorageService(cfg.Vault.Path),
	}, nil
}

func (a *App) OpenSession() (*Session, error) {
	if err := security.CheckLocked(a.Config.Security.MaxFailedUnlocks, a.Config.Security.LockoutSeconds); err != nil {
		a.Logger.Error("Unlock blocked by rate limiter")
		return nil, err
	}
	password, err := a.Vault.ReadPassword("Enter master password: ")
	if err != nil {
		a.Logger.Error(fmt.Sprintf("Vault unlock aborted: %s", ErrFailedReadPassword))
		return nil, ErrFailedReadPassword
	}
	defer zeroBytes(password)

	if err := a.Storage.Load(); err != nil {
		a.Logger.Error(fmt.Sprintf("Vault unlock failed: %s", err))
		return nil, err
	}

	salt, err := a.Vault.DecodeSalt(a.Storage.VaultFile.Salt)
	if err != nil {
		a.Logger.Error(fmt.Sprintf("Vault unlock failed: %s", err))
		return nil, err
	}
	defer zeroBytes(salt)

	encKey, macKey, err := a.Vault.DeriveKeys(password, salt)
	if err != nil {
		a.Logger.Error(fmt.Sprintf("Vault unlock failed: %s", err))
		return nil, err
	}

	// verify MAC if present
	if len(a.Storage.VaultFile.MAC) > 0 && !a.Vault.VerifyMAC(macKey, a.Storage.VaultFile.MAC, a.Storage.VaultFile.Nonce, a.Storage.VaultFile.Ciphertext) {
		zeroBytes(encKey)
		zeroBytes(macKey)
		a.Logger.Warn("Vault integrity check failed")
		security.RecordFailure(
			a.Config.Security.MaxFailedUnlocks,
			a.Config.Security.LockoutSeconds,
		)
		return nil, ErrFailedVaultIntegrity
	}

	plain, err := a.Vault.Decrypt(encKey, a.Storage.VaultFile.Nonce, a.Storage.VaultFile.Ciphertext)
	if err != nil {
		zeroBytes(encKey)
		zeroBytes(macKey)
		a.Logger.Warn(fmt.Sprintf("Vault unlock failed: %s", ErrWrongMasterPassword))
		security.RecordFailure(
			a.Config.Security.MaxFailedUnlocks,
			a.Config.Security.LockoutSeconds,
		)
		return nil, ErrWrongMasterPassword
	}
	defer zeroBytes(plain)

	v := &vault.Vault{}
	if err := json.Unmarshal(plain, &v); err != nil {
		zeroBytes(encKey)
		zeroBytes(macKey)
		a.Logger.Error(fmt.Sprintf("Failed to parse vault JSON: %s", err))
		return nil, fmt.Errorf("Failed to parse vault JSON: %w", err)
	}

	// upgrade MAC if missing
	if len(a.Storage.VaultFile.MAC) == 0 {
		a.Storage.VaultFile.MAC = a.Vault.ComputeMAC(macKey, a.Storage.VaultFile.Nonce, a.Storage.VaultFile.Ciphertext)
		if err := a.Storage.Save(a.Storage.VaultFile); err != nil {
			a.Logger.Error("Failed to upgrade vault MAC")
			return nil, fmt.Errorf("Failed to upgrade vault MAC")
		}
		a.Logger.Info("Vault MAC added")
		fmt.Println("✅ Vault MAC added")
	}

	security.RecordSuccess()

	return &Session{
		Vault:  v,
		EncKey: encKey,
		MacKey: macKey,
		File:   a.Storage.VaultFile,
	}, nil
}

func (a *App) SaveSession(sess *Session) error {
	plain, err := json.Marshal(sess.Vault)
	if err != nil {
		return ErrFailedSerializeVault
	}
	defer zeroBytes(plain)

	nonce, ciphertext, err := a.Vault.Encrypt(sess.EncKey, plain)
	if err != nil {
		return err
	}
	defer zeroBytes(nonce)
	defer zeroBytes(ciphertext)

	mac := a.Vault.ComputeMAC(sess.MacKey, nonce, ciphertext)
	defer zeroBytes(mac)

	sess.File.Nonce = nonce
	sess.File.Ciphertext = ciphertext
	sess.File.MAC = mac

	if err := a.Storage.Save(sess.File); err != nil {
		return err
	}
	return nil
}

func (a *App) CloseSession(sess *Session) {
	sess.Close()
}
