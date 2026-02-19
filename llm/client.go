// Package llm provides LLM client interfaces and implementations.
package llm

import "context"

// Client is the interface for interacting with LLM providers.
type Client interface {
	GenerateOptions(ctx context.Context, query string) ([]CommandOption, error)
}

// CommandOption represents a single command suggestion with explanation.
type CommandOption struct {
	Title       string `json:"title"`
	Command     string `json:"command"`
	Description string `json:"description"`
}
