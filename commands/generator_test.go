package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/pixielabs/1lm/llm"
)

func TestGeneratorGenerate(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		mockOptions []llm.CommandOption
		mockErr     error
		wantErr     bool
		wantCount   int
	}{
		{
			name:  "successful generation",
			query: "search git history",
			mockOptions: []llm.CommandOption{
				{Title: "Option 1", Command: "git log", Description: "Show git log"},
				{Title: "Option 2", Command: "git log -p", Description: "Show git log with patches"},
				{Title: "Option 3", Command: "git log --all", Description: "Show all git log"},
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "LLM error",
			query:     "test query",
			mockErr:   errors.New("API error"),
			wantErr:   true,
			wantCount: 0,
		},
		{
			name:        "empty response",
			query:       "test query",
			mockOptions: []llm.CommandOption{},
			wantErr:     false,
			wantCount:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &llm.MockClient{
				Response: tt.mockOptions,
				Err:      tt.mockErr,
			}

			gen := NewGenerator(mock, nil, "test-model")
			options, err := gen.Generate(context.Background(), tt.query)

			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(options) != tt.wantCount {
				t.Errorf("Generate() got %d options, want %d", len(options), tt.wantCount)
			}

			if mock.LastQuery != tt.query {
				t.Errorf("Generate() passed query %q, want %q", mock.LastQuery, tt.query)
			}

			if !tt.wantErr && len(tt.mockOptions) > 0 {
				if options[0].Title != tt.mockOptions[0].Title {
					t.Errorf("Generate() option title = %q, want %q", options[0].Title, tt.mockOptions[0].Title)
				}
				if options[0].Command != tt.mockOptions[0].Command {
					t.Errorf("Generate() option command = %q, want %q", options[0].Command, tt.mockOptions[0].Command)
				}
			}
		})
	}
}
