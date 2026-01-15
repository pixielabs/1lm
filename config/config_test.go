package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Provider != "anthropic" {
		t.Errorf("DefaultConfig() provider = %q, want %q", cfg.Provider, "anthropic")
	}

	if cfg.Model == "" {
		t.Error("DefaultConfig() model is empty")
	}

	if cfg.AnthropicAPIKey != "" {
		t.Error("DefaultConfig() should not set API key")
	}
}

func TestGetProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		wantOk   bool
	}{
		{
			name:     "anthropic exists",
			provider: "anthropic",
			wantOk:   true,
		},
		{
			name:     "unknown provider",
			provider: "unknown",
			wantOk:   false,
		},
		{
			name:     "empty provider",
			provider: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, ok := GetProvider(tt.provider)
			if ok != tt.wantOk {
				t.Errorf("GetProvider(%q) ok = %v, want %v", tt.provider, ok, tt.wantOk)
			}

			if ok && provider.Name != tt.provider {
				t.Errorf("GetProvider(%q) name = %q, want %q", tt.provider, provider.Name, tt.provider)
			}
		})
	}
}

func TestSupportedProviders(t *testing.T) {
	providers := SupportedProviders()

	if len(providers) == 0 {
		t.Error("SupportedProviders() returned empty list")
	}

	// Should at least have Anthropic
	found := false
	for _, p := range providers {
		if p.Name == "anthropic" {
			found = true
			if p.DefaultModel == "" {
				t.Error("Anthropic provider has no default model")
			}
			if !p.RequiresAPIKey {
				t.Error("Anthropic provider should require API key")
			}
		}
	}

	if !found {
		t.Error("SupportedProviders() missing 'anthropic'")
	}
}
