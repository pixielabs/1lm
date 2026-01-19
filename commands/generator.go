package commands

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/pixielabs/1lm/llm"
	"github.com/pixielabs/1lm/safety"
)

// ProgressStage represents the current stage of generation.
type ProgressStage int

const (
	// StageGenerating indicates command generation is in progress.
	StageGenerating ProgressStage = iota
	// StageEvaluating indicates safety evaluation is in progress.
	StageEvaluating
)

// ProgressCallback is called when generation progresses to a new stage.
type ProgressCallback func(stage ProgressStage)

// Generator handles command generation from natural language queries.
type Generator struct {
	client    llm.Client
	evaluator *safety.Evaluator
}

// NewGenerator creates a new command Generator.
//
// client          - The LLM client to use for generation
// anthropicClient - The Anthropic client for safety evaluation
// model           - The model to use for safety evaluation
//
// Returns an initialized Generator.
func NewGenerator(client llm.Client, anthropicClient *anthropic.Client, model string) *Generator {
	return &Generator{
		client:    client,
		evaluator: safety.NewEvaluator(anthropicClient, model),
	}
}

// Generate creates command options from a natural language query.
//
// ctx   - The context for the request
// query - The natural language description
//
// Returns a slice of Options and any error encountered.
//
// Examples
//
//   gen := commands.NewGenerator(client, anthropicClient, model)
//   options, err := gen.Generate(ctx, "search git history for myFunction")
//   if err != nil {
//       log.Fatal(err)
//   }
func (g *Generator) Generate(ctx context.Context, query string) ([]Option, error) {
	return g.GenerateWithProgress(ctx, query, nil)
}

// GenerateWithProgress creates command options with progress callbacks.
//
// ctx      - The context for the request
// query    - The natural language description
// progress - Optional callback for progress updates
//
// Returns a slice of Options and any error encountered.
func (g *Generator) GenerateWithProgress(ctx context.Context, query string, progress ProgressCallback) ([]Option, error) {
	llmOptions, err := g.client.GenerateOptions(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate options: %w", err)
	}

	// Convert LLM options to command options
	options := make([]Option, len(llmOptions))
	for i, opt := range llmOptions {
		options[i] = Option{
			Title:       opt.Title,
			Command:     opt.Command,
			Description: opt.Description,
		}
	}

	// Notify progress: moving to safety evaluation stage
	if progress != nil {
		progress(StageEvaluating)
	}

	// Evaluate commands for safety risks
	commands := make([]string, len(options))
	for i, opt := range options {
		commands[i] = opt.Command
	}

	risks, err := g.evaluator.Evaluate(ctx, commands)
	if err != nil {
		// Log error but don't fail - safety check is best-effort
		// Safety evaluation failure shouldn't prevent command generation
	} else {
		for i, risk := range risks {
			if risk != nil && risk.Level != safety.RiskNone {
				options[i].Risk = risk
			}
		}
	}

	return options, nil
}
