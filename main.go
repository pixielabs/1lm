// 1lm generates CLI one-liners from natural language using LLMs.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pixielabs/1lm/commands"
	"github.com/pixielabs/1lm/config"
	"github.com/pixielabs/1lm/llm"
	"github.com/pixielabs/1lm/ui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
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

	// Check if query provided as command line args
	if len(os.Args) >= 2 {
		// Use command line args as query
		query := strings.Join(os.Args[1:], " ")
		initialModel = ui.NewLoadingModel(generator, query)
	} else {
		// No args - show text input prompt
		initialModel = ui.NewInputModel(generator)
	}

	// Run the program (will transition through input → loading → selector)
	p := tea.NewProgram(initialModel)
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
		fmt.Println("No option selected")
		return nil
	}

	// Copy to clipboard using pbcopy
	if err := copyToClipboard(selected.Command); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	fmt.Printf("\n✓ Copied to clipboard: %s\n", selected.Command)

	return nil
}

// copyToClipboard copies text to the macOS clipboard using pbcopy.
//
// text - The text to copy
//
// Returns any error encountered.
func copyToClipboard(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
