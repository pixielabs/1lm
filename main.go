// 1lm generates CLI one-liners from natural language using LLMs.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/pixielabs/1lm/commands"
	"github.com/pixielabs/1lm/config"
	"github.com/pixielabs/1lm/llm"
	"github.com/pixielabs/1lm/output"
	"github.com/pixielabs/1lm/ui"
)

var (
	outputMode = flag.String("output", "clipboard", "Output mode: clipboard, shell-function, stdout")
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Parse command-line flags
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate API key
	if cfg.AnthropicAPIKey == "" {
		return fmt.Errorf("anthropic_api_key not set in config (~/.config/1lm/config.toml)")
	}

	// Initialize LLM client
	client, err := llm.NewAnthropicClient(cfg.AnthropicAPIKey, cfg.Model)
	if err != nil {
		return fmt.Errorf("failed to create LLM client: %w", err)
	}

	// Create Anthropic client for safety evaluation
	anthropicClient := anthropic.NewClient(
		option.WithAPIKey(cfg.AnthropicAPIKey),
	)

	// Create generator
	generator := commands.NewGenerator(client, &anthropicClient, cfg.Model)

	var initialModel tea.Model

	// Check if query provided as command line args (after flags)
	args := flag.Args()
	if len(args) >= 1 {
		// Use command line args as query
		query := strings.Join(args, " ")
		initialModel = ui.NewLoadingModel(generator, query)
	} else {
		// No args - show text input prompt
		initialModel = ui.NewInputModel(generator)
	}

	// Run the program (will transition through input → loading → selector)
	// In shell-function mode, use /dev/tty for TUI so stdout is clean for command output
	var p *tea.Program
	if *outputMode == "shell-function" {
		// Open TTY for direct terminal access
		tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
		if err != nil {
			return fmt.Errorf("failed to open /dev/tty: %w", err)
		}
		defer func() { _ = tty.Close() }()

		// Detect color profile from TTY (not stdout, which is captured)
		output := termenv.NewOutput(tty)
		lipgloss.SetColorProfile(output.ColorProfile())

		p = tea.NewProgram(initialModel, tea.WithInput(tty), tea.WithOutput(tty))
	} else {
		p = tea.NewProgram(initialModel)
	}

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running UI: %w", err)
	}

	// Check if loading failed
	if loadingModel, ok := finalModel.(ui.LoadingModel); ok {
		if err := loadingModel.Err(); err != nil {
			return fmt.Errorf("failed to generate options: %w", err)
		}
	}

	// Get selected option from selector
	selectorModel, ok := finalModel.(ui.SelectorModel)
	if !ok {
		// User quit before selecting (from input or loading)
		return nil
	}

	selected := selectorModel.Selected()
	if selected == nil {
		// Only print message in non-shell-function mode
		if *outputMode != "shell-function" {
			fmt.Println("No option selected")
		}
		return nil
	}

	// Create output handler based on mode
	handler := output.NewHandler(output.Mode(*outputMode))

	// Output the selected command
	if err := handler.Output(selected); err != nil {
		return fmt.Errorf("failed to output command: %w", err)
	}

	return nil
}
