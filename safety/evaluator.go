package safety

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

// RiskLevel represents the severity of detected risk.
type RiskLevel int

const (
	// RiskNone indicates no detected risk.
	RiskNone RiskLevel = iota
	// RiskLow indicates low-severity risk (e.g., network operations).
	RiskLow
	// RiskHigh indicates high-severity risk (e.g., destructive operations).
	RiskHigh
)

// RiskInfo contains details about detected risks.
type RiskInfo struct {
	Level   RiskLevel
	Message string // Human-readable warning
}

// Evaluator uses LLM to evaluate command safety.
type Evaluator struct {
	client *anthropic.Client
	model  string
}

// CommandRisk represents the safety evaluation for a single command.
type CommandRisk struct {
	Command   string `json:"command"`
	RiskLevel string `json:"risk_level"` // "none", "low", "high"
	Reason    string `json:"reason"`     // Brief explanation
}

// SafetyResponse is the structured output from the LLM.
type SafetyResponse struct {
	Evaluations []CommandRisk `json:"evaluations"`
}

// NewEvaluator creates a new safety evaluator.
//
// client - The Anthropic API client
// model  - The model to use for evaluation
//
// Returns a new Evaluator instance.
func NewEvaluator(client *anthropic.Client, model string) *Evaluator {
	return &Evaluator{
		client: client,
		model:  model,
	}
}

// Evaluate evaluates multiple commands for safety risks in a single API call.
//
// ctx      - Context for the API call
// commands - List of commands to evaluate
//
// Returns a slice of RiskInfo pointers (nil for safe commands) and any error.
func (e *Evaluator) Evaluate(ctx context.Context, commands []string) ([]*RiskInfo, error) {
	if len(commands) == 0 {
		return nil, nil
	}

	// Return error if client is nil (e.g., in tests)
	if e.client == nil {
		return nil, fmt.Errorf("evaluator client is nil")
	}

	// Build the prompt
	prompt := buildPrompt(commands)

	// Define JSON schema for structured output
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"evaluations": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"command": map[string]any{
							"type": "string",
						},
						"risk_level": map[string]any{
							"type": "string",
							"enum": []string{"none", "low", "high"},
						},
						"reason": map[string]any{
							"type":      "string",
							"maxLength": 100,
						},
					},
					"required":             []string{"command", "risk_level", "reason"},
					"additionalProperties": false,
				},
			},
		},
		"required":             []string{"evaluations"},
		"additionalProperties": false,
	}

	// Build system message
	systemMessage := "You are a security expert evaluating shell commands for safety risks. Respond with structured JSON output following the provided schema."

	// Make API call with structured output using Beta API
	message, err := e.client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
		Model:     anthropic.Model(e.model),
		MaxTokens: 1024,
		Betas: []anthropic.AnthropicBeta{
			"structured-outputs-2025-11-13",
		},
		Messages: []anthropic.BetaMessageParam{{
			Content: []anthropic.BetaContentBlockParamUnion{{
				OfText: &anthropic.BetaTextBlockParam{
					Text: prompt,
				},
			}},
			Role: anthropic.BetaMessageParamRoleUser,
		}},
		System: []anthropic.BetaTextBlockParam{{
			Text: systemMessage,
		}},
		OutputFormat: anthropic.BetaJSONOutputFormatParam{
			Schema: schema,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	// Extract text from response
	if len(message.Content) == 0 {
		return nil, fmt.Errorf("empty response from API")
	}

	textContent := message.Content[0].Text
	if textContent == "" {
		return nil, fmt.Errorf("no text content in response")
	}

	// Parse JSON response
	var response SafetyResponse
	if err := json.Unmarshal([]byte(textContent), &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Validate we got the right number of evaluations
	if len(response.Evaluations) != len(commands) {
		return nil, fmt.Errorf("expected %d evaluations, got %d", len(commands), len(response.Evaluations))
	}

	// Convert to RiskInfo
	results := make([]*RiskInfo, len(commands))
	for i, eval := range response.Evaluations {
		level := parseRiskLevel(eval.RiskLevel)
		if level == RiskNone {
			results[i] = nil
			continue
		}

		results[i] = &RiskInfo{
			Level:   level,
			Message: eval.Reason,
		}
	}

	return results, nil
}

// buildPrompt builds the evaluation prompt for the LLM.
//
// commands - List of commands to evaluate
//
// Returns the formatted prompt string.
func buildPrompt(commands []string) string {
	prompt := `You are a security expert evaluating shell commands for safety risks.

For each command below, determine the risk level and provide a brief reason.

Risk levels:
- HIGH: Destructive operations that could cause data loss or system damage (rm -rf, dd, mkfs, formatting, permanent deletion)
- LOW: Operations that interact with external systems or require careful attention (network operations, downloads, system scans, privilege changes)
- NONE: Safe read-only operations (ls, grep, find, echo, cat, viewing files)

Commands to evaluate:
`

	for i, cmd := range commands {
		prompt += fmt.Sprintf("%d. %s\n", i+1, cmd)
	}

	prompt += "\nBe practical and context-aware. Flag commands that users should think twice about before running."

	return prompt
}

// parseRiskLevel converts a string risk level to RiskLevel enum.
//
// level - String risk level from API
//
// Returns the corresponding RiskLevel.
func parseRiskLevel(level string) RiskLevel {
	switch level {
	case "low":
		return RiskLow
	case "high":
		return RiskHigh
	default:
		return RiskNone
	}
}

// String returns the string representation of a RiskLevel.
//
// Returns "None", "Low", or "High".
func (r RiskLevel) String() string {
	switch r {
	case RiskLow:
		return "Low"
	case RiskHigh:
		return "High"
	default:
		return "None"
	}
}
