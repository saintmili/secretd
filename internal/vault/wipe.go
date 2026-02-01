package vault

func (e *Entry) Wipe() {
	if e.Password != nil {
		for i := range e.Password {
			e.Password[i] = 0
		}
	}
}

func (v *Vault) Wipe() {
	for i := range v.Entries {
		v.Entries[i].Wipe()
	}
}
