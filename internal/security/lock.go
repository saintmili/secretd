package security

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LockState struct {
	FailedAttempts int       `json:"failed_attempts"`
	LockedUntil    time.Time `json:"locked_until"`
}

func lockFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(homeDir, ".local", "share", "secretd")
	os.MkdirAll(dir, 0700)
	return filepath.Join(dir, "lock.json"), nil
}

func loadState() (LockState, error) {
	var state LockState

	path, err := lockFilePath()
	if err != nil {
		return state, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return state, nil
	}

	json.Unmarshal(data, &state)
	return state, nil
}

func saveState(state LockState) error {
	path, err := lockFilePath()
	if err != nil {
		return err
	}

	data, _ := json.MarshalIndent(state, "", " ")

	return os.WriteFile(path, data, 0600)
}

func CheckLocked(maxTry int, lockSeconds int) error {
	state, err := loadState()
	if err != nil {
		return err
	}

	now := time.Now()
	if now.Before(state.LockedUntil) {
		remaining := int(time.Until(state.LockedUntil).Seconds())
		return fmt.Errorf("vault locked. Try again in %d seconds", remaining)
	}

	if state.FailedAttempts > 0 {
		fmt.Printf("%d tries left\n", maxTry-state.FailedAttempts)

		// we use 1<< instead of 2^
		// it performs better. its a `left bit shift` operation
		delay := time.Duration(1<<uint(state.FailedAttempts-1)) * time.Second
		time.Sleep(delay)
	}

	return nil
}

// adds one to failedAttemps, if it reach maxTry, user will be locked
func RecordFailure(maxTry int, lockSeconds int) error {
	state, err := loadState()
	if err != nil {
		return err
	}

	state.FailedAttempts++

	if state.FailedAttempts >= maxTry {
		state.LockedUntil = time.Now().Add(time.Duration(lockSeconds) * time.Second)
		state.FailedAttempts = 0
	}

	return saveState(state)
}

func RecordSuccess() {
	state := LockState{}
	saveState(state)
}
