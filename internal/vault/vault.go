package vault

type Vault struct {
	Version int      `json:"version"`
	Entries []*Entry `json:"entries"`
}

func New() *Vault {
	return &Vault{
		Version: 1,
		Entries: []*Entry{},
	}
}
