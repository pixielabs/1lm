# 1lm

> LLM-powered CLI one-liner generator

Describe what you want in natural language, get shell commands. No more forgetting git flags or searching Stack Overflow for that perfect one-liner.

```bash
$ 1lm "search git history for modifications mentioning myFunction"

Select a command:

▸ Git log with pickaxe
  git log -p -S myFunction
  Search commit history for additions/deletions of the exact string "myFunction"

  Git log with regex
  git log -G myFunction
  Search commits where "myFunction" appears in the diff (supports regex)

  Git log all branches
  git log --all -p -S myFunction
  Search across all branches for changes mentioning "myFunction"

↑/k: up • ↓/j: down • enter: select • q: quit
```

## Features

- **Multiple options**: Get 3 different approaches to choose from
- **Interactive selection**: Arrow keys or vim bindings to navigate
- **Safety warnings**: LLM-powered risk evaluation with visual indicators
  - 🚨 High risk warnings for destructive operations (rm -rf, dd, etc.)
  - ⚠️ Low risk warnings for network operations, scans, and privilege changes
- **Clipboard ready**: Selected command automatically copied via `pbcopy`
- **Context-aware**: Descriptions explain what each command does and any caveats
- **Reliable**: Uses Anthropic's structured outputs API for guaranteed valid responses
- **Real-time progress**: See "Generating options..." and "Evaluating safety..." as it works

## Installation

### Prerequisites

- Go 1.25 or later
- macOS (for clipboard support via `pbcopy`) - Linux builds available but clipboard doesn't work yet
- [Anthropic API key](https://console.anthropic.com/)

### Build from source

```bash
git clone https://github.com/jalada/1lm.git
cd 1lm
go build -o 1lm
```

Optionally, move the binary to your PATH:

```bash
sudo mv 1lm /usr/local/bin/
```

## Configuration

Create `~/.config/1lm/config.toml`:

```toml
provider = "anthropic"
model = "claude-sonnet-4-5-20250929"
anthropic_api_key = "sk-ant-your-api-key-here"
```

### Getting an API key

1. Sign up at [console.anthropic.com](https://console.anthropic.com/)
2. Navigate to API Keys
3. Create a new key
4. Add it to your config file

## Usage

### Basic usage

```bash
1lm "your natural language query"
```

### Examples

Find files:
```bash
1lm "find all python files modified in the last week"
```

Process data:
```bash
1lm "count unique IP addresses in access.log"
```

Git operations:
```bash
1lm "show commits from last month by author alice"
```

System info:
```bash
1lm "check disk usage sorted by size"
```

### Keyboard controls

- `↑` or `k` - Move selection up
- `↓` or `j` - Move selection down
- `Enter` - Select command and copy to clipboard
- `q` or `Ctrl+C` - Quit without selecting

## How it works

1. **Query**: You describe what you want in natural language
2. **Generate**: Claude generates 3 command options using structured outputs API
3. **Evaluate**: Claude assesses each command for safety risks (destructive ops, network activity)
4. **Select**: Interactive TUI shows options with explanations and safety warnings
5. **Copy**: Selected command is copied to clipboard, ready to paste and run

Under the hood:
- Uses [Anthropic SDK for Go](https://github.com/anthropics/anthropic-sdk-go) with structured outputs beta
- JSON Schema enforces reliable response format
- [Bubbletea](https://github.com/charmbracelet/bubbletea) powers the interactive TUI
- [Lipgloss](https://github.com/charmbracelet/lipgloss) handles styling

## Troubleshooting

### "anthropic_api_key not set in config"

Make sure `~/.config/1lm/config.toml` exists and contains your API key:

```toml
anthropic_api_key = "sk-ant-..."
```

### Text wrapping issues

The UI automatically detects terminal width. If descriptions still overflow, try resizing your terminal or updating to the latest version.

### "API call failed"

Check:
- Your API key is valid
- You have API credits remaining
- Your network connection is working

## Development

### Project structure

```
1lm/
├── main.go          # Entry point and wiring
├── config/          # TOML configuration management
├── llm/             # Anthropic SDK integration
│   ├── client.go    # LLM client interface
│   ├── provider.go  # Structured outputs implementation
│   └── mock.go      # Mock for testing
├── commands/        # Command generation logic
├── safety/          # LLM-based safety evaluation
├── ui/              # Bubbletea interactive selector
└── tests/           # Unit tests
```

### Running tests

```bash
go test ./...
```

### Code style

- TomDoc format for public functions
- Clear interfaces for testability
- Small, focused files

See [PLAN.md](PLAN.md) for detailed architecture decisions and roadmap.

## Roadmap

### v0.2.0 (Complete ✓)
- ✓ Anthropic Claude integration
- ✓ Structured output with JSON schema
- ✓ Interactive TUI
- ✓ Clipboard support
- ✓ TOML configuration
- ✓ LLM-based safety warnings with visual indicators
- ✓ Multi-stage progress indicator

### Fast follows
- Shell integration (command insertion vs clipboard)
- Context awareness (cwd, OS, installed tools)

### Future
- Multiple LLM provider support
- Response caching
- Command history
- Cross-platform clipboard support

## Why?

Tools like GitHub Copilot CLI and Warp AI exist, but they're tied to specific editors or terminals. `1lm` is a standalone tool that:
- Works in any terminal
- Gives you multiple options with explanations
- Lets you pick the approach that fits your needs
- Is easy to understand and modify

## Contributing

This is a learning project built collaboratively with Claude Code. Contributions welcome!

## License

MIT

---

Built with Go, [Anthropic Claude](https://www.anthropic.com/), and [Charmbracelet](https://charm.sh/) tools.