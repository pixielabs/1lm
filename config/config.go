// Package config handles application configuration loading and management.
package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the application configuration.
type Config struct {
	// The Anthropic API key for Claude access
	AnthropicAPIKey string `toml:"anthropic_api_key"`

	// The Claude model to use (e.g., "claude-sonnet-4-5-20250929")
	Model string `toml:"model"`

	// The LLM provider (currently only "anthropic" supported)
	Provider string `toml:"provider"`
}

// Load reads and parses the configuration file from the standard location.
//
// Returns the parsed Config and any error encountered.
//
// Examples
//
//   cfg, err := config.Load()
//   if err != nil {
//       log.Fatal(err)
//   }
func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	// Return default config if file doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save writes the configuration to the standard location.
//
// cfg - The Config to save
//
// Returns any error encountered.
func Save(cfg *Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(cfg)
}

// ConfigPath returns the path to the configuration file.
//
// Returns the config file path and any error encountered.
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "1lm", "config.toml"), nil
}

// DefaultConfig returns a Config with sensible defaults.
//
// Returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		Provider: "anthropic",
		Model:    "claude-sonnet-4-5-20250929",
	}
}
