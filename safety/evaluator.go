package safety

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

// RiskInfo contains details about a detected risk.
type RiskInfo struct {
	Level   RiskLevel
	Message string
}

// Evaluator uses an LLM to evaluate command safety.
type Evaluator struct {
	client *anthropic.Client
	model  string
}

// CommandRisk represents the safety evaluation for a single command.
type CommandRisk struct {
	Command   string `json:"command"`
	RiskLevel string `json:"risk_level"`
	Reason    string `json:"reason"`
}

// SafetyResponse is the structured output from the safety LLM call.
type SafetyResponse struct {
	Evaluations []CommandRisk `json:"evaluations"`
}

// Public: Creates a new safety evaluator.
func NewEvaluator(client *anthropic.Client, model string) *Evaluator {
	return &Evaluator{
		client: client,
		model:  model,
	}
}

// safetySchema defines the structured output format for safety evaluation.
var safetySchema = map[string]any{
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

// Public: Evaluates multiple commands for safety risks in a single API call.
// Returns a slice of RiskInfo pointers (nil elements indicate safe commands).
func (e *Evaluator) Evaluate(ctx context.Context, commands []string) ([]*RiskInfo, error) {
	if len(commands) == 0 {
		return nil, nil
	}

	if e.client == nil {
		return nil, fmt.Errorf("evaluator client is nil")
	}

	prompt := buildPrompt(commands)

	systemMessage := `You are a security expert evaluating shell commands for safety risks.

Risk levels:
- HIGH: Destructive operations that could cause data loss or system damage (rm -rf, dd, mkfs, formatting, permanent deletion)
- LOW: Operations that interact with external systems or require careful attention (network operations, downloads, system scans, privilege changes)
- NONE: Safe read-only operations (ls, grep, find, echo, cat, viewing files)

Be practical and context-aware. Flag commands that users should think twice about before running.`

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
			Schema: safetySchema,
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

	var response SafetyResponse
	if err := json.Unmarshal([]byte(textContent), &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(response.Evaluations) != len(commands) {
		return nil, fmt.Errorf("expected %d evaluations, got %d", len(commands), len(response.Evaluations))
	}

	results := make([]*RiskInfo, len(commands))
	for i, eval := range response.Evaluations {
		if level := parseRiskLevel(eval.RiskLevel); level != RiskNone {
			results[i] = &RiskInfo{
				Level:   level,
				Message: eval.Reason,
			}
		}
	}

	return results, nil
}

// buildPrompt formats the list of commands into an evaluation prompt.
func buildPrompt(commands []string) string {
	var b strings.Builder
	b.WriteString("Evaluate these commands:\n\n")

	for i, cmd := range commands {
		fmt.Fprintf(&b, "%d. %s\n", i+1, cmd)
	}

	return b.String()
}

// parseRiskLevel converts a string risk level to a RiskLevel enum value.
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

// String returns the human-readable name of a RiskLevel.
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
