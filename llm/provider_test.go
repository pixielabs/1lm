package llm

import (
	"encoding/json"
	"testing"
)

func TestCommandOptionJSONSerialization(t *testing.T) {
	option := CommandOption{
		Title:       "Test Command",
		Command:     "echo test",
		Description: "Test description",
	}

	// Test JSON marshaling
	data, err := json.Marshal(option)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Test JSON unmarshaling
	var decoded CommandOption
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	// Verify fields match
	if decoded.Title != option.Title {
		t.Errorf("Title = %q, want %q", decoded.Title, option.Title)
	}
	if decoded.Command != option.Command {
		t.Errorf("Command = %q, want %q", decoded.Command, option.Command)
	}
	if decoded.Description != option.Description {
		t.Errorf("Description = %q, want %q", decoded.Description, option.Description)
	}
}

func TestCommandOptionResponseParsing(t *testing.T) {
	// Test parsing a typical API response
	jsonResponse := `{
		"options": [
			{
				"title": "Git log with search",
				"command": "git log -p -S myFunction",
				"description": "Search git history for modifications"
			},
			{
				"title": "Git grep",
				"command": "git log -G myFunction",
				"description": "Show commits with pattern in diff"
			}
		]
	}`

	var result struct {
		Options []CommandOption `json:"options"`
	}

	if err := json.Unmarshal([]byte(jsonResponse), &result); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if len(result.Options) != 2 {
		t.Errorf("got %d options, want 2", len(result.Options))
	}

	if result.Options[0].Title == "" {
		t.Error("first option title is empty")
	}
	if result.Options[0].Command == "" {
		t.Error("first option command is empty")
	}
	if result.Options[0].Description == "" {
		t.Error("first option description is empty")
	}
}
