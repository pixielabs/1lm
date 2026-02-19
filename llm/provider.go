package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// AnthropicClient implements Client using Anthropic's Claude models.
type AnthropicClient struct {
	client anthropic.Client
	model  anthropic.Model
}

// optionsSchema defines the structured output format for command generation.
var optionsSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"options": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title": map[string]any{
						"type":        "string",
						"description": "Brief title for this command option (2-5 words)",
					},
					"command": map[string]any{
						"type":        "string",
						"description": "The actual shell command to execute",
					},
					"description": map[string]any{
						"type":        "string",
						"description": "Clear explanation of what this command does and any important details",
					},
				},
				"required":             []string{"title", "command", "description"},
				"additionalProperties": false,
			},
		},
	},
	"required":             []string{"options"},
	"additionalProperties": false,
}

// Public: Creates a new Anthropic client for command generation.
func NewAnthropicClient(apiKey, model string) (Client, error) {
	return &AnthropicClient{
		client: anthropic.NewClient(option.WithAPIKey(apiKey)),
		model:  anthropic.Model(model),
	}, nil
}

// Public: Generates command options from a natural language query using
// Anthropic's structured outputs API.
func (c *AnthropicClient) GenerateOptions(ctx context.Context, query string) ([]CommandOption, error) {
	promptText := fmt.Sprintf(`Given this user request: "%s"

Generate exactly 3 different shell command options that accomplish the task.

Requirements:
- Provide exactly 3 different approaches when possible
- Commands should be safe and practical
- Prefer commonly available tools
- Include relevant flags and options
- Descriptions should explain the approach and any caveats`, query)

	message, err := c.client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
		Model:     c.model,
		MaxTokens: 2048,
		Betas: []anthropic.AnthropicBeta{
			"structured-outputs-2025-11-13",
		},
		Messages: []anthropic.BetaMessageParam{{
			Content: []anthropic.BetaContentBlockParamUnion{{
				OfText: &anthropic.BetaTextBlockParam{
					Text: promptText,
				},
			}},
			Role: anthropic.BetaMessageParamRoleUser,
		}},
		OutputFormat: anthropic.BetaJSONOutputFormatParam{
			Schema: optionsSchema,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if len(message.Content) == 0 {
		return nil, fmt.Errorf("empty response from API")
	}

	textContent := message.Content[0].Text
	if textContent == "" {
		return nil, fmt.Errorf("no text content in response")
	}

	var result struct {
		Options []CommandOption `json:"options"`
	}

	if err := json.Unmarshal([]byte(textContent), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	if len(result.Options) == 0 {
		return nil, fmt.Errorf("no options returned")
	}

	return result.Options, nil
}
