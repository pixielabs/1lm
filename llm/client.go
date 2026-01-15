// Package llm provides LLM client interfaces and implementations.
package llm

import "context"

// Client is the interface for interacting with LLM providers.
type Client interface {
	// GenerateOptions generates command options from a natural language query.
	//
	// ctx   - The context for the request
	// query - The natural language description of desired command
	//
	// Returns a slice of CommandOptions and any error encountered.
	GenerateOptions(ctx context.Context, query string) ([]CommandOption, error)
}

// CommandOption represents a single command option with explanation.
type CommandOption struct {
	// Brief title/summary of this option
	Title string `json:"title"`

	// The shell command to execute
	Command string `json:"command"`

	// Human-readable explanation of what the command does
	Description string `json:"description"`
}
