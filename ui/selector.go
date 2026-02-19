package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pixielabs/1lm/commands"
	"github.com/pixielabs/1lm/safety"
	"golang.org/x/term"
)

// SelectorModel represents the bubbletea model for option selection.
type SelectorModel struct {
	options    []commands.Option
	cursor     int
	selected   *commands.Option
	quitting   bool
	width      int
	generator  *commands.Generator
	safetyDone bool
	spinner    spinner.Model
}

// NewSelector creates a new option selector.
//
// options   - The command options to choose from
// generator - The generator used to run background safety evaluation
//
// Returns an initialized SelectorModel.
func NewSelector(options []commands.Option, generator *commands.Generator) SelectorModel {
	// Get terminal width, default to 80 if unable to detect
	width := 80
	if w, _, err := term.GetSize(0); err == nil && w > 0 {
		width = w
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = CheckingStyle

	return SelectorModel{
		options:   options,
		cursor:    0,
		width:     width,
		generator: generator,
		spinner:   s,
	}
}

// Init fires the background safety evaluation and starts the spinner.
func (m SelectorModel) Init() tea.Cmd {
	return tea.Batch(m.evaluateSafety, m.spinner.Tick)
}

// evaluateSafety runs safety evaluation and returns a riskResultMsg.
func (m SelectorModel) evaluateSafety() tea.Msg {
	options, err := m.generator.EvaluateSafety(context.Background(), m.options)
	return riskResultMsg{options: options, err: err}
}

// Update handles messages and updates the model. Required by bubbletea.
//
// msg - The message to process
//
// Returns the updated model and any command to run.
func (m SelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}

		case "enter":
			m.selected = &m.options[m.cursor]
			m.quitting = true
			return m, tea.Quit
		}

	case riskResultMsg:
		m.safetyDone = true
		if msg.err == nil {
			m.options = msg.options
		}
		return m, nil

	case spinner.TickMsg:
		if !m.safetyDone {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// View renders the UI. Required by bubbletea.
//
// Returns the rendered string.
func (m SelectorModel) View() string {
	if m.quitting && m.selected == nil {
		return ""
	}

	var b strings.Builder

	b.WriteString("\n")
	b.WriteString("Select a command:\n\n")

	// Reserve space for cursor and indentation
	contentWidth := m.width - 4

	for i, option := range m.options {
		cursor := " "
		if m.cursor == i {
			cursor = SelectedStyle.Render("▸")
		}

		title := TitleStyle.Render(option.Title)
		if m.cursor == i {
			title = SelectedStyle.Render(option.Title)
		}

		// Wrap command and description to terminal width
		command := CommandStyle.Width(contentWidth).Render(option.Command)

		// Add risk warning if present, or animated placeholder while checking
		var riskWarning string
		if option.Risk != nil {
			riskWarning = formatRiskWarning(option.Risk, m.cursor == i)
		} else if !m.safetyDone {
			riskWarning = m.spinner.View() + CheckingStyle.Render(" checking safety...")
		}

		description := DescriptionStyle.Width(contentWidth).Render(option.Description)

		b.WriteString(fmt.Sprintf("%s %s\n", cursor, title))
		b.WriteString(fmt.Sprintf("  %s\n", command))
		if riskWarning != "" {
			b.WriteString(fmt.Sprintf("  %s\n", riskWarning))
		}
		b.WriteString(fmt.Sprintf("  %s\n\n", description))
	}

	if m.selected == nil {
		b.WriteString(HelpStyle.Render("↑/k: up • ↓/j: down • enter: select • q: quit"))
		b.WriteString("\n")
	}

	return b.String()
}

// formatRiskWarning formats a risk warning with appropriate styling.
//
// risk     - The risk information
// selected - Whether this option is currently selected
//
// Returns a styled warning string.
func formatRiskWarning(risk *safety.RiskInfo, selected bool) string {
	var icon string
	var style lipgloss.Style

	switch risk.Level {
	case safety.RiskLow:
		icon = "⚠️"
		style = WarningLowStyle
	case safety.RiskHigh:
		icon = "🚨"
		style = WarningHighStyle
	default:
		return ""
	}

	message := fmt.Sprintf("%s %s", icon, risk.Message)

	if selected {
		style = style.Bold(true)
	}

	return style.Render(message)
}

// Selected returns the selected option, if any.
//
// Returns a pointer to the selected Option, or nil if none selected.
func (m SelectorModel) Selected() *commands.Option {
	return m.selected
}
