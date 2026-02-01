package clipboard

// copy data to clipboard and clears it after `seconds`.
// Uses detached sleep commands for reliability.
func Copy(data string, timeout int) error {
	err := copy(data, timeout)
	return err
}
