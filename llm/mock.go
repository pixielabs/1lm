package llm

import "context"

// MockClient is a test double for the Client interface.
type MockClient struct {
	Response  []CommandOption
	Err       error
	LastQuery string
}

// GenerateOptions returns the pre-configured response and captures the query.
func (m *MockClient) GenerateOptions(_ context.Context, query string) ([]CommandOption, error) {
	m.LastQuery = query
	return m.Response, m.Err
}

// NewMockClient creates a MockClient with three sample options.
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
