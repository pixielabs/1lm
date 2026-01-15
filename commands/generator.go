package commands

import (
	"context"
	"fmt"

	"github.com/pixielabs/1lm/llm"
)

// Generator handles command generation from natural language queries.
type Generator struct {
	client llm.Client
}

// NewGenerator creates a new command Generator.
//
// client - The LLM client to use for generation
//
// Returns an initialized Generator.
func NewGenerator(client llm.Client) *Generator {
	return &Generator{
		client: client,
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
//   gen := commands.NewGenerator(client)
//   options, err := gen.Generate(ctx, "search git history for myFunction")
//   if err != nil {
//       log.Fatal(err)
//   }
func (g *Generator) Generate(ctx context.Context, query string) ([]Option, error) {
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

	return options, nil
}
