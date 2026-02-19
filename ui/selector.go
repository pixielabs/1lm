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

// riskResultMsg is sent when background safety evaluation completes.
type riskResultMsg struct {
	options []commands.Option
	err     error
}

// SelectorModel lets the user pick from generated command options.
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

// NewSelector creates a new option selector with background safety evaluation.
func NewSelector(options []commands.Option, generator *commands.Generator) SelectorModel {
	width := 80
	if w, _, err := term.GetSize(0); err == nil && w > 0 {
		width = w
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = CheckingStyle

	return SelectorModel{
		options:   options,
		width:     width,
		generator: generator,
		spinner:   s,
	}
}

// Init starts background safety evaluation and the spinner animation.
func (m SelectorModel) Init() tea.Cmd {
	return tea.Batch(m.evaluateSafety, m.spinner.Tick)
}

func (m SelectorModel) evaluateSafety() tea.Msg {
	options, err := m.generator.EvaluateSafety(context.Background(), m.options)
	return riskResultMsg{options: options, err: err}
}

// Update handles key presses, safety results, and spinner ticks.
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

// View renders the option list with safety indicators.
func (m SelectorModel) View() string {
	if m.quitting && m.selected == nil {
		return ""
	}

	var b strings.Builder

	b.WriteString("\n")
	b.WriteString("Select a command:\n\n")

	contentWidth := m.width - 4

	for i, option := range m.options {
		isSelected := m.cursor == i

		cursor := " "
		title := TitleStyle.Render(option.Title)
		if isSelected {
			cursor = SelectedStyle.Render("▸")
			title = SelectedStyle.Render(option.Title)
		}

		command := CommandStyle.Width(contentWidth).Render(option.Command)

		var riskWarning string
		if option.Risk != nil {
			riskWarning = formatRiskWarning(option.Risk, isSelected)
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

// formatRiskWarning returns a styled warning string for the given risk level.
func formatRiskWarning(risk *safety.RiskInfo, selected bool) string {
	var icon string
	var style lipgloss.Style

	switch risk.Level {
	case safety.RiskLow:
		icon, style = "⚠️", WarningLowStyle
	case safety.RiskHigh:
		icon, style = "🚨", WarningHighStyle
	default:
		return ""
	}

	if selected {
		style = style.Bold(true)
	}

	return style.Render(fmt.Sprintf("%s %s", icon, risk.Message))
}

// Selected returns the chosen option, or nil if the user quit.
func (m SelectorModel) Selected() *commands.Option {
	return m.selected
}
