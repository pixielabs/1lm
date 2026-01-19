// Package output handles command output in different modes.
package output

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pixielabs/1lm/commands"
)

// Mode represents the output mode.
type Mode string

const (
	// ModeClipboard copies to system clipboard (default).
	ModeClipboard Mode = "clipboard"
	// ModeShellFunction outputs for shell function integration.
	ModeShellFunction Mode = "shell-function"
	// ModeStdout prints to stdout only.
	ModeStdout Mode = "stdout"
)

// Handler manages command output.
type Handler struct {
	mode Mode
}

// NewHandler creates a new output handler.
//
// mode - The output mode to use
//
// Returns an initialized Handler.
func NewHandler(mode Mode) *Handler {
	return &Handler{mode: mode}
}

// Output handles the selected command based on output mode.
//
// cmd - The selected command option
//
// Returns any error encountered.
func (h *Handler) Output(cmd *commands.Option) error {
	switch h.mode {
	case ModeShellFunction:
		return h.outputShellFunction(cmd)
	case ModeStdout:
		return h.outputStdout(cmd)
	case ModeClipboard:
		fallthrough
	default:
		return h.outputClipboard(cmd)
	}
}

// outputShellFunction outputs for shell function consumption.
func (h *Handler) outputShellFunction(cmd *commands.Option) error {
	// Print command to stdout (shell wrapper will read it)
	fmt.Println(cmd.Command)
	return nil
}

// outputStdout prints to stdout with confirmation.
func (h *Handler) outputStdout(cmd *commands.Option) error {
	fmt.Printf("\n✓ Selected command:\n%s\n", cmd.Command)
	return nil
}

// outputClipboard copies to system clipboard (current behavior).
func (h *Handler) outputClipboard(cmd *commands.Option) error {
	// Try pbcopy (macOS)
	if err := copyViaPbcopy(cmd.Command); err == nil {
		fmt.Printf("\n✓ Copied to clipboard: %s\n", cmd.Command)
		return nil
	}

	// Try xclip (Linux)
	if err := copyViaXclip(cmd.Command); err == nil {
		fmt.Printf("\n✓ Copied to clipboard: %s\n", cmd.Command)
		return nil
	}

	// Try wl-copy (Wayland)
	if err := copyViaWlCopy(cmd.Command); err == nil {
		fmt.Printf("\n✓ Copied to clipboard: %s\n", cmd.Command)
		return nil
	}

	// Fallback: print to stdout
	fmt.Printf("\n⚠ Clipboard not available\n")
	return h.outputStdout(cmd)
}

// copyViaPbcopy uses macOS pbcopy.
func copyViaPbcopy(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// copyViaXclip uses Linux xclip.
func copyViaXclip(text string) error {
	cmd := exec.Command("xclip", "-selection", "clipboard")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// copyViaWlCopy uses Wayland wl-copy.
func copyViaWlCopy(text string) error {
	cmd := exec.Command("wl-copy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
