# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2026-01-19

### Added
- Shell integration for bash, zsh, and fish shells
- Selected commands now appear directly in shell prompt ready to execute or edit
- Multiple output modes via `--output` flag:
  - `shell-function`: Outputs command for shell function to inject into prompt
  - `clipboard`: Copies to system clipboard (default)
  - `stdout`: Prints command to terminal
- Cross-platform clipboard support:
  - macOS via `pbcopy`
  - Linux X11 via `xclip`
  - Linux Wayland via `wl-copy`
- Automatic fallback to stdout when clipboard tools unavailable
- Shell function examples for all supported shells in README

### Technical
- New `output` package with handler abstraction for different output modes
- Uses `/dev/tty` for TUI rendering when stdout is captured by shell functions
- Manual color profile detection from TTY to preserve styling in shell-function mode
- Comprehensive test coverage for output handlers
- Flag parsing infrastructure for command-line options

### Changed
- Removed hardcoded pbcopy-only clipboard implementation
- Updated Prerequisites in README to reflect cross-platform support

### Fixed
- Color and styling now work correctly in shell-function mode

## [0.2.0] - 2025-01-19

### Added
- LLM-based safety checks for generated commands
- Real-time progress indicator showing "Generating options..." and "Evaluating safety..." stages
- Visual risk warnings in command selector:
  - 🚨 High risk (red) for destructive operations (rm -rf, dd, mkfs, data loss)
  - ⚠️ Low risk (yellow) for operations requiring attention (network scans, downloads, privilege changes)
- Context-aware safety evaluation using same Claude model as command generation
- Single batch API call evaluates all commands simultaneously
- Non-blocking warnings - users can still select and run flagged commands
- Graceful fallback if safety evaluation fails

### Technical
- New `safety` package with LLM-based evaluator using structured outputs
- Extended `commands.Option` struct with optional Risk field
- Multi-stage progress tracking in loading UI with channel-based communication
- Three-level risk system (none/low/high) with appropriate styling

## [0.1.0] - 2025-01-15

### Added
- Initial release of 1lm
- Natural language to CLI one-liner generation using Anthropic Claude
- Interactive terminal UI with bubbletea
- Three command options with explanations for each query
- Animated loading spinner while generating options
- Text input prompt when run without arguments
- Command line argument support for quick queries
- Clipboard integration via pbcopy (macOS)
- TOML configuration file support (~/.config/1lm/config.toml)
- Structured outputs API integration for reliable JSON responses
- Text wrapping for descriptions to fit terminal width
- Keyboard navigation (arrow keys and vim bindings)

### Technical
- Built with Go 1.25
- Uses official anthropic-sdk-go
- Structured outputs beta API with JSON Schema validation
- Charmbracelet Bubbletea for TUI
- Charmbracelet Lipgloss for styling
- Comprehensive test suite with mocks

### Known Limitations
- macOS only (clipboard requires pbcopy)
- Requires Anthropic API key

[Unreleased]: https://github.com/pixielabs/1lm/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/pixielabs/1lm/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/pixielabs/1lm/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/pixielabs/1lm/releases/tag/v0.1.0
