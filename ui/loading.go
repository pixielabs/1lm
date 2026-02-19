package ui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pixielabs/1lm/commands"
)

// LoadingModel represents the loading state with a spinner.
type LoadingModel struct {
	spinner   spinner.Model
	generator *commands.Generator
	query     string
	err       error
}

// optionsMsg is sent when options are loaded.
type optionsMsg struct {
	options []commands.Option
	err     error
}

// NewLoadingModel creates a new loading model.
//
// generator - The command generator to use
// query     - The user's query
//
// Returns an initialized LoadingModel.
func NewLoadingModel(generator *commands.Generator, query string) LoadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = TitleStyle

	return LoadingModel{
		spinner:   s,
		generator: generator,
		query:     query,
	}
}

// Init initializes the model and starts the spinner + API call.
func (m LoadingModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.loadOptions)
}

// loadOptions performs the generation API call asynchronously.
func (m LoadingModel) loadOptions() tea.Msg {
	options, err := m.generator.Generate(context.Background(), m.query)
	return optionsMsg{options: options, err: err}
}

// Update handles messages and updates the model.
//
// msg - The message to process
//
// Returns the updated model and any command to run.
func (m LoadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case optionsMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}

		if len(msg.options) == 0 {
			m.err = fmt.Errorf("no options generated")
			return m, tea.Quit
		}

		selector := NewSelector(msg.options, m.generator)
		return selector, selector.Init()

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the loading UI.
//
// Returns the rendered string.
func (m LoadingModel) View() string {
	if m.err != nil {
		return ""
	}

	return fmt.Sprintf("\n%s Generating options...\n", m.spinner.View())
}

// Err returns any error encountered during loading.
//
// Returns the error or nil.
func (m LoadingModel) Err() error {
	return m.err
}
