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
	spinner    spinner.Model
	generator  *commands.Generator
	query      string
	stage      commands.ProgressStage
	progressCh chan commands.ProgressStage
	err        error
}

// optionsMsg is sent when options are loaded.
type optionsMsg struct {
	options []commands.Option
	err     error
}

// progressMsg is sent when generation progresses to a new stage.
type progressMsg struct {
	stage commands.ProgressStage
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
		spinner:    s,
		generator:  generator,
		query:      query,
		stage:      commands.StageGenerating,
		progressCh: make(chan commands.ProgressStage, 2),
	}
}

// Init initializes the model and starts the spinner + API call.
func (m LoadingModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.loadOptions,
		m.waitForProgress,
	)
}

// waitForProgress listens for progress updates from the channel.
func (m LoadingModel) waitForProgress() tea.Msg {
	stage, ok := <-m.progressCh
	if !ok {
		// Channel closed, no more progress updates
		return nil
	}
	return progressMsg{stage: stage}
}

// loadOptions performs the API call asynchronously with progress updates.
func (m LoadingModel) loadOptions() tea.Msg {
	// Call generator with progress callback
	options, err := m.generator.GenerateWithProgress(
		context.Background(),
		m.query,
		func(stage commands.ProgressStage) {
			// Send progress update to channel (non-blocking)
			select {
			case m.progressCh <- stage:
			default:
			}
		},
	)

	// Close progress channel when done
	close(m.progressCh)

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
		// Allow quit during loading
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

	case progressMsg:
		// Update stage and keep listening for more progress
		m.stage = msg.stage
		return m, m.waitForProgress

	case optionsMsg:
		// Options loaded - transition to selector or show error
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}

		if len(msg.options) == 0 {
			m.err = fmt.Errorf("no options generated")
			return m, tea.Quit
		}

		// Transition to selector
		selector := NewSelector(msg.options)
		return selector, selector.Init()

	default:
		// Update spinner
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

	var message string
	switch m.stage {
	case commands.StageGenerating:
		message = "Generating options..."
	case commands.StageEvaluating:
		message = "Evaluating safety..."
	default:
		message = "Processing..."
	}

	return fmt.Sprintf("\n%s %s\n", m.spinner.View(), message)
}

// Err returns any error encountered during loading.
//
// Returns the error or nil.
func (m LoadingModel) Err() error {
	return m.err
}
