// Package commands handles command generation and option management.
package commands

// Option represents a command option with metadata.
type Option struct {
	// Brief title for this option
	Title string

	// The actual shell command
	Command string

	// Explanation of what this command does
	Description string
}
