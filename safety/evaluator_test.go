package safety

import (
	"context"
	"testing"
)

func TestParseRiskLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  RiskLevel
	}{
		{
			name:  "none",
			level: "none",
			want:  RiskNone,
		},
		{
			name:  "low",
			level: "low",
			want:  RiskLow,
		},
		{
			name:  "high",
			level: "high",
			want:  RiskHigh,
		},
		{
			name:  "invalid",
			level: "invalid",
			want:  RiskNone,
		},
		{
			name:  "empty",
			level: "",
			want:  RiskNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseRiskLevel(tt.level)
			if got != tt.want {
				t.Errorf("parseRiskLevel(%q) = %v, want %v", tt.level, got, tt.want)
			}
		})
	}
}

func TestRiskLevelString(t *testing.T) {
	tests := []struct {
		name  string
		level RiskLevel
		want  string
	}{
		{
			name:  "none",
			level: RiskNone,
			want:  "None",
		},
		{
			name:  "low",
			level: RiskLow,
			want:  "Low",
		},
		{
			name:  "high",
			level: RiskHigh,
			want:  "High",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.level.String()
			if got != tt.want {
				t.Errorf("RiskLevel.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildPrompt(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		contains []string
	}{
		{
			name:     "single command",
			commands: []string{"rm -rf /tmp/*"},
			contains: []string{"rm -rf /tmp/*", "1."},
		},
		{
			name:     "multiple commands",
			commands: []string{"ls -la", "git status", "rm -rf /"},
			contains: []string{"ls -la", "git status", "rm -rf /", "1.", "2.", "3."},
		},
		{
			name:     "empty commands",
			commands: []string{},
			contains: []string{"Evaluate these commands:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := buildPrompt(tt.commands)
			for _, substr := range tt.contains {
				if len(substr) > 0 && !contains(prompt, substr) {
					t.Errorf("buildPrompt() missing expected substring %q", substr)
				}
			}
		})
	}
}

func TestEvaluateEmptyCommands(t *testing.T) {
	// Test with nil client since we won't make API calls
	evaluator := &Evaluator{client: nil, model: "test"}

	results, err := evaluator.Evaluate(context.Background(), []string{})

	if err != nil {
		t.Errorf("Evaluate() with empty commands should not error, got: %v", err)
	}

	if results != nil {
		t.Errorf("Evaluate() with empty commands should return nil, got: %v", results)
	}
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
