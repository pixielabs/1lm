package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/pixielabs/1lm/commands"
	"golang.org/x/term"
)

// SelectorModel represents the bubbletea model for option selection.
type SelectorModel struct {
	options  []commands.Option
	cursor   int
	selected *commands.Option
	quitting bool
	width    int
}

// NewSelector creates a new option selector.
//
// options - The command options to choose from
//
// Returns an initialized SelectorModel.
func NewSelector(options []commands.Option) SelectorModel {
	// Get terminal width, default to 80 if unable to detect
	width := 80
	if w, _, err := term.GetSize(0); err == nil && w > 0 {
		width = w
	}

	return SelectorModel{
		options: options,
		cursor:  0,
		width:   width,
	}
}

// Init initializes the model. Required by bubbletea.
func (m SelectorModel) Init() tea.Cmd {
	return nil
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
	}

	return m, nil
}

// View renders the UI. Required by bubbletea.
//
// Returns the rendered string.
func (m SelectorModel) View() string {
	if m.quitting {
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
		command := CommandStyle.Copy().Width(contentWidth).Render(option.Command)
		description := DescriptionStyle.Copy().Width(contentWidth).Render(option.Description)

		b.WriteString(fmt.Sprintf("%s %s\n", cursor, title))
		b.WriteString(fmt.Sprintf("  %s\n", command))
		b.WriteString(fmt.Sprintf("  %s\n\n", description))
	}

	b.WriteString(HelpStyle.Render("↑/k: up • ↓/j: down • enter: select • q: quit"))
	b.WriteString("\n")

	return b.String()
}

// Selected returns the selected option, if any.
//
// Returns a pointer to the selected Option, or nil if none selected.
func (m SelectorModel) Selected() *commands.Option {
	return m.selected
}
