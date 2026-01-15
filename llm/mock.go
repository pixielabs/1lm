package llm

import "context"

// MockClient is a mock LLM client for testing.
type MockClient struct {
	// Response to return from GenerateOptions
	Response []CommandOption

	// Error to return from GenerateOptions
	Err error

	// Captures the last query passed to GenerateOptions
	LastQuery string
}

// GenerateOptions returns the configured response or error.
//
// ctx   - The context for the request
// query - The natural language description
//
// Returns the configured response and error.
func (m *MockClient) GenerateOptions(ctx context.Context, query string) ([]CommandOption, error) {
	m.LastQuery = query
	return m.Response, m.Err
}

// NewMockClient creates a new mock client with default successful response.
//
// Returns a MockClient configured with sample options.
func NewMockClient() *MockClient {
	return &MockClient{
		Response: []CommandOption{
			{
				Title:       "Option 1",
				Command:     "echo 'test'",
				Description: "Test command 1",
			},
			{
				Title:       "Option 2",
				Command:     "echo 'test2'",
				Description: "Test command 2",
			},
			{
				Title:       "Option 3",
				Command:     "echo 'test3'",
				Description: "Test command 3",
			},
		},
	}
}
