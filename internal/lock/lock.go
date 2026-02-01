package lock

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type AppLock struct {
	file *os.File
}

func lockPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(homeDir, ".local", "share", "secretd")
	os.MkdirAll(dir, 0700)

	return filepath.Join(dir, "secretd.lock"), nil
}

func Aquire() (*AppLock, error) {
	path, err := lockPath()
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return nil, fmt.Errorf("another secretd instance is already running")
	}

	return &AppLock{file: f}, nil
}

func(al *AppLock) Release() {
	if al == nil || al.file == nil {
		return
	}

	syscall.Flock(int(al.file.Fd()), syscall.LOCK_UN)
	al.file.Close()
}
