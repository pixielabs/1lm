package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jalada/1lm/commands"
)

// InputModel represents the text input prompt for queries.
type InputModel struct {
	textInput textinput.Model
	generator *commands.Generator
	submitted bool
	query     string
}

// NewInputModel creates a new input model.
//
// generator - The command generator to use
//
// Returns an initialized InputModel.
func NewInputModel(generator *commands.Generator) InputModel {
	ti := textinput.New()
	ti.Placeholder = "e.g., search git history for myFunction"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 80

	return InputModel{
		textInput: ti,
		generator: generator,
	}
}

// Init initializes the model.
func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model.
//
// msg - The message to process
//
// Returns the updated model and any command to run.
func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Submit query and transition to loading
			m.query = m.textInput.Value()
			if m.query != "" {
				m.submitted = true
				loadingModel := NewLoadingModel(m.generator, m.query)
				return loadingModel, loadingModel.Init()
			}
			return m, nil

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		// Adjust input width based on terminal size
		m.textInput.Width = msg.Width - 4
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the input UI.
//
// Returns the rendered string.
func (m InputModel) View() string {
	if m.submitted {
		return ""
	}

	return fmt.Sprintf(
		"\n%s\n\n%s\n\n%s\n",
		TitleStyle.Render("What command do you need?"),
		m.textInput.View(),
		HelpStyle.Render("Enter to submit • Esc/Ctrl+C to quit"),
	)
}
