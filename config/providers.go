package config

// Provider represents an LLM provider configuration.
type Provider struct {
	// The provider name (e.g., "anthropic", "openai")
	Name string

	// Default model for this provider
	DefaultModel string

	// Whether API key is required
	RequiresAPIKey bool
}

// SupportedProviders returns a list of all supported LLM providers.
//
// Returns a slice of Provider definitions.
func SupportedProviders() []Provider {
	return []Provider{
		{
			Name:           "anthropic",
			DefaultModel:   "claude-sonnet-4-5-20250929",
			RequiresAPIKey: true,
		},
		// Future providers can be added here
	}
}

// GetProvider returns the provider configuration for a given name.
//
// name - The provider name to look up
//
// Returns the Provider and a boolean indicating if found.
func GetProvider(name string) (Provider, bool) {
	for _, p := range SupportedProviders() {
		if p.Name == name {
			return p, true
		}
	}
	return Provider{}, false
}
