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

var outputMode = flag.String("output", "clipboard", "Output mode: clipboard, shell-function, stdout")

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Parse command-line flags, supporting flags anywhere in the arg list.
	// Go's flag package stops at the first non-flag argument, so
	// "1lm my query --output=shell-function" would leave --output unparsed.
	// Re-order args to put flags first so they're always processed.
	var flagArgs, queryArgs []string
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			flagArgs = append(flagArgs, arg)
		} else {
			queryArgs = append(queryArgs, arg)
		}
	}
	os.Args = append(append([]string{os.Args[0]}, flagArgs...), queryArgs...)
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.AnthropicAPIKey == "" {
		return fmt.Errorf("anthropic_api_key not set in config (~/.config/1lm/config.toml)")
	}

	client, err := llm.NewAnthropicClient(cfg.AnthropicAPIKey, cfg.Model)
	if err != nil {
		return fmt.Errorf("failed to create LLM client: %w", err)
	}

	// Separate Anthropic client needed for safety evaluation (different API surface)
	anthropicClient := anthropic.NewClient(
		option.WithAPIKey(cfg.AnthropicAPIKey),
	)

	generator := commands.NewGenerator(client, &anthropicClient, cfg.Model)

	var initialModel tea.Model
	if args := flag.Args(); len(args) > 0 {
		query := strings.Join(args, " ")
		initialModel = ui.NewLoadingModel(generator, query)
	} else {
		initialModel = ui.NewInputModel(generator)
	}

	// In shell-function mode, use /dev/tty so stdout stays clean for command output
	var p *tea.Program
	if *outputMode == "shell-function" {
		tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
		if err != nil {
			return fmt.Errorf("failed to open /dev/tty: %w", err)
		}
		defer func() { _ = tty.Close() }()

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

	if loadingModel, ok := finalModel.(ui.LoadingModel); ok {
		if err := loadingModel.Err(); err != nil {
			return fmt.Errorf("failed to generate options: %w", err)
		}
	}

	selectorModel, ok := finalModel.(ui.SelectorModel)
	if !ok {
		return nil
	}

	selected := selectorModel.Selected()
	if selected == nil {
		if *outputMode != "shell-function" {
			fmt.Println("No option selected")
		}
		return nil
	}

	handler := output.NewHandler(output.Mode(*outputMode))
	if err := handler.Output(selected); err != nil {
		return fmt.Errorf("failed to output command: %w", err)
	}

	return nil
}
