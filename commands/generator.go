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

// Public: Creates a new Generator with the given LLM client and a safety
// evaluator backed by the Anthropic client.
func NewGenerator(client llm.Client, anthropicClient *anthropic.Client, model string) *Generator {
	return &Generator{
		client:    client,
		evaluator: safety.NewEvaluator(anthropicClient, model),
	}
}

// Public: Generates command options from a natural language query.
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

// Public: Evaluates commands for safety risks and returns updated options.
// Best-effort: returns (nil, err) on failure so callers can ignore silently.
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
