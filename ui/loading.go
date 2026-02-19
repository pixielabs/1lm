package ui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pixielabs/1lm/commands"
)

// LoadingModel shows a spinner while generating command options.
type LoadingModel struct {
	spinner   spinner.Model
	generator *commands.Generator
	query     string
	err       error
}

// optionsMsg is sent when the generation API call completes.
type optionsMsg struct {
	options []commands.Option
	err     error
}

// NewLoadingModel creates a loading model that generates options for the query.
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

// Init starts the spinner and kicks off the API call.
func (m LoadingModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.loadOptions)
}

func (m LoadingModel) loadOptions() tea.Msg {
	options, err := m.generator.Generate(context.Background(), m.query)
	return optionsMsg{options: options, err: err}
}

// Update handles spinner ticks, API responses, and quit keys.
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

// View renders the spinner with a "Generating options..." message.
func (m LoadingModel) View() string {
	if m.err != nil {
		return ""
	}

	return fmt.Sprintf("\n%s Generating options...\n", m.spinner.View())
}

// Err returns any error encountered during loading.
func (m LoadingModel) Err() error {
	return m.err
}
