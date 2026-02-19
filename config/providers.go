package config

// Provider represents an LLM provider configuration.
type Provider struct {
	Name           string
	DefaultModel   string
	RequiresAPIKey bool
}

// Public: Returns all supported LLM providers.
func SupportedProviders() []Provider {
	return []Provider{
		{
			Name:           "anthropic",
			DefaultModel:   "claude-sonnet-4-5-20250929",
			RequiresAPIKey: true,
		},
	}
}

// Public: Returns the provider configuration for a given name.
func GetProvider(name string) (Provider, bool) {
	for _, p := range SupportedProviders() {
		if p.Name == name {
			return p, true
		}
	}
	return Provider{}, false
}
