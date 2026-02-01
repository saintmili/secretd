package config

type ValidationWarning struct {
	Field   string
	Message string
}

func (c *Config) Validate() []ValidationWarning {
	var warnings []ValidationWarning

	if c.Clipboard.ClearAfterSeconds < 5 || c.Clipboard.ClearAfterSeconds > 300 {
		warnings = append(warnings, ValidationWarning{
			Field:   "timeout",
			Message: "must be between 5 and 300 seconds; using default",
		})
		c.Clipboard.ClearAfterSeconds = DefaultConfig().Clipboard.ClearAfterSeconds
	}

	if c.Security.Argon2Memory < 32*1024 {
		warnings = append(warnings, ValidationWarning{
			Field:   "argon2_memory",
			Message: "too low for secure key derivation; using default",
		})
		c.Security.Argon2Memory = DefaultConfig().Security.Argon2Memory
	}

	if c.Security.Argon2Time < 2 {
		warnings = append(warnings, ValidationWarning{
			Field:   "argon2_time",
			Message: "must be >= 2; using default",
		})
		c.Security.Argon2Time = DefaultConfig().Security.Argon2Time
	}

	if c.Security.Argon2Threads < 1 {
		warnings = append(warnings, ValidationWarning{
			Field:   "argon2_threads",
			Message: "must be >= 1; using default",
		})
		c.Security.Argon2Threads = DefaultConfig().Security.Argon2Threads
	}

	return warnings
}
