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
	default:
		return h.outputClipboard(cmd)
	}
}

// outputShellFunction outputs for shell function consumption.
func (h *Handler) outputShellFunction(cmd *commands.Option) error {
	fmt.Println(cmd.Command)
	return nil
}

// outputStdout prints to stdout with confirmation.
func (h *Handler) outputStdout(cmd *commands.Option) error {
	fmt.Printf("\n✓ Selected command:\n%s\n", cmd.Command)
	return nil
}

// clipboardCmd describes one clipboard tool and how to invoke it.
type clipboardCmd struct {
	name string
	args []string
}

// clipboardTools lists the clipboard tools to try, in order of preference.
var clipboardTools = []clipboardCmd{
	{name: "pbcopy"},                                    // macOS
	{name: "xclip", args: []string{"-selection", "clipboard"}}, // Linux X11
	{name: "wl-copy"},                                   // Wayland
}

// outputClipboard copies to system clipboard (current behavior).
func (h *Handler) outputClipboard(cmd *commands.Option) error {
	for _, tool := range clipboardTools {
		c := exec.Command(tool.name, tool.args...)
		c.Stdin = strings.NewReader(cmd.Command)
		if c.Run() == nil {
			fmt.Printf("\n✓ Copied to clipboard: %s\n", cmd.Command)
			return nil
		}
	}

	fmt.Printf("\n⚠ Clipboard not available\n")
	return h.outputStdout(cmd)
}
