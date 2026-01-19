// Package ui provides terminal user interface components.
package ui

import "github.com/charmbracelet/lipgloss"

var (
	// TitleStyle is used for option titles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	// CommandStyle is used for displaying commands
	CommandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Background(lipgloss.Color("235")).
			Padding(0, 1)

	// DescriptionStyle is used for option descriptions
	DescriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241"))

	// SelectedStyle is used for the currently selected option
	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true)

	// HelpStyle is used for help text
	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)

	// WarningLowStyle for low-risk operations (network, downloads, scans)
	WarningLowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")).
			Italic(true)

	// WarningHighStyle for high-risk operations (destructive, data loss)
	WarningHighStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")).
				Bold(true)
)
