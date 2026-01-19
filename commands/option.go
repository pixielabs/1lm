// Package commands handles command generation and option management.
package commands

import "github.com/pixielabs/1lm/safety"

// Option represents a command option with metadata.
type Option struct {
	// Brief title for this option
	Title string

	// The actual shell command
	Command string

	// Explanation of what this command does
	Description string

	// Risk information if the command is potentially dangerous
	// nil if no risk detected
	Risk *safety.RiskInfo
}
