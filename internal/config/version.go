package config

import "fmt"

const (
	CurrentVersion      = 1
	MinSupportedVersion = 1
)

func (c *Config) CheckVersion() error {
	if c.Version == 0 {
		// legacy config → assume v1
		c.Version = 1
		return nil
	}

	if c.Version < MinSupportedVersion {
		return fmt.Errorf(
			"config version %d is too old (min supported: %d)",
			c.Version, MinSupportedVersion,
		)
	}

	if c.Version > CurrentVersion {
		return fmt.Errorf(
			"config version %d is newer than this secretd version supports",
			c.Version,
		)
	}

	return nil
}
