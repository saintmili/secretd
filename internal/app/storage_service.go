package app

import "github.com/saintmili/secretd/internal/storage"

type StorageService struct {
	Path      string
	VaultFile *storage.VaultFile
}

func NewStorageService(path string) *StorageService {
	return &StorageService{
		Path: path,
	}
}

func (s *StorageService) Load() error {
	vf, err := storage.Load(s.Path)
	if err != nil {
		return err
	}
	s.VaultFile = vf
	return nil
}

func (s *StorageService) Save(vf *storage.VaultFile) error {
	return storage.Save(s.VaultFile, s.Path)
}
