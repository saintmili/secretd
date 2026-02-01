package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func LoadConfig() (*Config, []ValidationWarning, error) {
	cfg := DefaultConfig()

	path, err := ConfigPath()
	if err != nil {
		return cfg, nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		data, _ := toml.Marshal(cfg)
		if err := os.WriteFile(path, data, 0o600); err != nil {
			return cfg, nil, fmt.Errorf("failed to write default config: %w", err)
		}

		return cfg, nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := cfg.CheckVersion(); err != nil {
		return cfg, nil, err
	}

	warnings := cfg.Validate()
	return cfg, warnings, nil
}
