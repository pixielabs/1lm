# Agent Development Guide

This file provides context and instructions for AI coding agents working on 1lm. For user-facing documentation, see README.md.

## Project Structure

```
1lm/
├── cmd/              # CLI entry point
├── config/           # TOML configuration handling
├── llm/             # LLM provider implementations
├── commands/        # Core generation logic
└── ui/              # Bubbletea TUI components
```

Each package has a single, focused responsibility. Files are kept small (< 200 lines typically).

## Build & Test

```bash
# Build
go build -v ./...

# Run tests
go test -v ./...

# Run linter (must match CI)
golangci-lint run

# Install locally
go install

# Run with args
1lm "find all go files modified in last week"

# Run interactive
1lm
```

## Code Conventions

### Documentation
- All public functions use TomDoc format
- Comments explain WHY, not WHAT
- Keep comments concise

Example:
```go
// Public: Generates shell command options from a natural language query.
//
// Sends the query to the configured LLM provider and returns 3 command
// options with explanations.
//
// ctx   - Context for cancellation and timeouts
// query - Natural language description of desired command
//
// Returns slice of CommandOption structs or error if generation fails.
func (g *Generator) Generate(ctx context.Context, query string) ([]Option, error)
```

### Style
- Use `gofmt` (enforced by CI)
- Prefer explicit error handling over panic
- Close resources with defer + error check:
  ```go
  defer func() {
      if err := file.Close(); err != nil {
          // handle error
      }
  }()
  ```

### Testing
- Mock external dependencies (LLM clients, clipboard)
- Test files live alongside implementation (`foo.go` → `foo_test.go`)
- Use table-driven tests for multiple cases
- Example mock pattern:
  ```go
  type MockLLMClient struct {
      GenerateFunc func(ctx context.Context, query string) ([]CommandOption, error)
  }
  ```

## Architecture Patterns

### LLM Provider Interface

To add a new provider:

1. Implement `llm.Client` interface in `llm/provider.go`:
   ```go
   type Client interface {
       GenerateOptions(ctx context.Context, query string) ([]CommandOption, error)
   }
   ```

2. Use structured outputs for reliable JSON parsing
3. Define JSON schema programmatically (see `AnthropicClient.GenerateOptions`)
4. Return exactly 3 options with title, command, and description

Current implementation: Anthropic Claude via `anthropic-sdk-go` with Beta structured outputs API.

### Bubbletea State Machine

UI flows through three models:

1. **InputModel** (`ui/input.go`): Text input for query
   - Transitions to LoadingModel on Enter

2. **LoadingModel** (`ui/loading.go`): Spinner during API call
   - Runs async: `fetchOptions()` in goroutine
   - Transitions to SelectorModel when complete

3. **SelectorModel** (`ui/selector.go`): Interactive option picker
   - Returns selected command on Enter
   - Quits on 'q'

Each model is self-contained with `Init()`, `Update()`, and `View()` methods.

### Configuration

Config file: `~/.config/1lm/config.toml`

```toml
provider = "anthropic"
anthropic_api_key = "sk-ant-..."
model = "claude-sonnet-4-5-20250929"
```

Loading: `config.Load()` uses `github.com/BurntSushi/toml`
- Creates default config if missing
- Returns error if API key not set
- Validates provider name

## Important Implementation Details

### Structured Outputs

We use Anthropic's Beta structured outputs API (header: `structured-outputs-2025-11-13`) to guarantee valid JSON. Do NOT use prompt engineering or try to parse markdown code blocks.

Schema is defined programmatically in `llm/provider.go`:
```go
schema := map[string]any{
    "type": "object",
    "properties": map[string]any{
        "options": map[string]any{
            "type": "array",
            // ...
        },
    },
    "required": []string{"options"},
    "additionalProperties": false,
}
```

Response comes as clean JSON in `message.Content[0].Text`.

### Text Wrapping

Terminal width is detected via `golang.org/x/term`:
```go
width := 80
if w, _, err := term.GetSize(0); err == nil && w > 0 {
    width = w
}
```

Lipgloss styles apply width constraints:
```go
// Don't use deprecated .Copy() - just call methods directly
description := DescriptionStyle.Width(contentWidth).Render(option.Description)
```

### Clipboard Integration

macOS: Uses `pbcopy` command
Linux: TODO - needs `xclip` or `wl-copy` detection (see PLAN.md)

Implementation in `main.go` after selection.

## CI/CD

### GitHub Actions

Two workflows:

1. **test.yml**: Runs on push/PR
   - Tests on Ubuntu + macOS matrix
   - Linting with `golangci-lint-action@v9`
   - Go 1.25

2. **release.yml**: Runs on version tags (`v*`)
   - Builds for macOS (arm64/amd64) and Linux (x86_64/arm64)
   - Uses GoReleaser
   - Creates GitHub release with binaries

### Linting

Must pass `golangci-lint run` locally before pushing.

Common issues:
- `errcheck`: Always check error returns (especially `defer file.Close()`)
- `staticcheck` SA1019: Don't use deprecated `.Copy()` on lipgloss styles

## Adding Features

### New LLM Provider

1. Add config field in `config/config.go`
2. Implement `llm.Client` interface
3. Update `cmd/1lm/main.go` to instantiate based on `cfg.Provider`
4. Add tests with mock
5. Update PLAN.md with provider status

### New UI Component

1. Create file in `ui/` package
2. Implement `bubbletea.Model` interface (Init, Update, View)
3. Define state transitions in `Update()`
4. Use shared styles from `ui/styles.go`
5. Handle terminal width for text wrapping

### Clipboard Support

1. Detect OS in `main.go`
2. Use appropriate command (`pbcopy`, `xclip`, `wl-copy`)
3. Fall back to printing if clipboard unavailable
4. Test on target platform

## Dependencies

Core:
- `github.com/anthropics/anthropic-sdk-go` - Official Anthropic API client
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/charmbracelet/bubbles` - UI components (spinner, textinput)
- `github.com/BurntSushi/toml` - TOML parsing
- `golang.org/x/term` - Terminal width detection

Keep dependencies minimal. Prefer stdlib when possible.

## Gotchas

1. **Lipgloss `.Copy()` is deprecated**: Just call methods directly (they return new styles)
2. **Anthropic client is value, not pointer**: Use `anthropic.Client`, not `*anthropic.Client`
3. **Beta API required**: Structured outputs need Beta header
4. **Terminal width detection can fail**: Always have fallback width (80)
5. **Go 1.25 required**: For latest syntax and stdlib features

## Future Work

See PLAN.md for roadmap. High priority items:
- Linux clipboard support
- Shell integration (context awareness)
- Safety checks before execution
- Multi-provider support
