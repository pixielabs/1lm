// Package commands handles command generation and option management.
package commands

import "github.com/pixielabs/1lm/safety"

// Option represents a generated command with metadata and optional risk info.
type Option struct {
	Title       string
	Command     string
	Description string
	Risk        *safety.RiskInfo // nil when no risk detected
}
