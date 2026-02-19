package commands

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/pixielabs/1lm/llm"
	"github.com/pixielabs/1lm/safety"
)

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
func (g *Generator) Generate(ctx context.Context, query string) ([]Option, error) {
	llmOptions, err := g.client.GenerateOptions(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate options: %w", err)
	}

	options := make([]Option, len(llmOptions))
	for i, opt := range llmOptions {
		options[i] = Option{
			Title:       opt.Title,
			Command:     opt.Command,
			Description: opt.Description,
		}
	}

	return options, nil
}

// EvaluateSafety evaluates commands for safety risks and returns updated options.
// This is best-effort: returns nil, err on failure so callers can ignore silently.
//
// ctx     - The context for the request
// options - The options to evaluate
//
// Returns updated options with Risk fields populated, or nil and an error.
func (g *Generator) EvaluateSafety(ctx context.Context, options []Option) ([]Option, error) {
	cmds := make([]string, len(options))
	for i, opt := range options {
		cmds[i] = opt.Command
	}

	risks, err := g.evaluator.Evaluate(ctx, cmds)
	if err != nil {
		return nil, err
	}

	result := make([]Option, len(options))
	copy(result, options)
	for i, risk := range risks {
		if risk != nil && risk.Level != safety.RiskNone {
			result[i].Risk = risk
		}
	}

	return result, nil
}
