// Package config handles application configuration loading and management.
package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the application configuration.
type Config struct {
	AnthropicAPIKey string `toml:"anthropic_api_key"`
	Model           string `toml:"model"`
	Provider        string `toml:"provider"`
}

// Public: Reads and parses the configuration file from ~/.config/1lm/config.toml.
// Returns default config if the file doesn't exist.
func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Public: Writes the configuration to ~/.config/1lm/config.toml.
func Save(cfg *Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(
		path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600,
	)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	return toml.NewEncoder(file).Encode(cfg)
}

// Public: Returns the path to the configuration file.
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "1lm", "config.toml"), nil
}

// Public: Returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Provider: "anthropic",
		Model:    "claude-sonnet-4-5-20250929",
	}
}
