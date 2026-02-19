package output

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/pixielabs/1lm/commands"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name string
		mode Mode
		want Mode
	}{
		{
			name: "clipboard mode",
			mode: ModeClipboard,
			want: ModeClipboard,
		},
		{
			name: "shell-function mode",
			mode: ModeShellFunction,
			want: ModeShellFunction,
		},
		{
			name: "stdout mode",
			mode: ModeStdout,
			want: ModeStdout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.mode)
			if handler.mode != tt.want {
				t.Errorf("NewHandler() mode = %v, want %v", handler.mode, tt.want)
			}
		})
	}
}

func TestShellFunctionOutput(t *testing.T) {
	handler := NewHandler(ModeShellFunction)
	cmd := &commands.Option{
		Title:       "List files",
		Command:     "ls -la",
		Description: "List all files",
	}

	output := captureOutput(func() {
		err := handler.Output(cmd)
		if err != nil {
			t.Errorf("Output() error = %v", err)
		}
	})

	expected := "ls -la\n"
	if output != expected {
		t.Errorf("Output() = %q, want %q", output, expected)
	}
}

func TestStdoutOutput(t *testing.T) {
	handler := NewHandler(ModeStdout)
	cmd := &commands.Option{
		Title:       "List files",
		Command:     "ls -la",
		Description: "List all files",
	}

	output := captureOutput(func() {
		err := handler.Output(cmd)
		if err != nil {
			t.Errorf("Output() error = %v", err)
		}
	})

	if !strings.Contains(output, "✓ Selected command:") {
		t.Errorf("Output() missing '✓ Selected command:', got %q", output)
	}
	if !strings.Contains(output, "ls -la") {
		t.Errorf("Output() missing command 'ls -la', got %q", output)
	}
}

func TestClipboardFallback(t *testing.T) {
	handler := NewHandler(ModeClipboard)
	cmd := &commands.Option{
		Title:       "List files",
		Command:     "ls -la",
		Description: "List all files",
	}

	output := captureOutput(func() {
		// Error is acceptable if clipboard tools are missing
		_ = handler.Output(cmd)
	})

	hasSuccess := strings.Contains(output, "✓ Copied to clipboard:")
	hasFallback := strings.Contains(output, "⚠ Clipboard not available")
	hasCommand := strings.Contains(output, "ls -la")

	if !hasCommand {
		t.Errorf("Output() missing command, got %q", output)
	}
	if !hasSuccess && !hasFallback {
		t.Errorf("Output() missing expected message, got %q", output)
	}
}

func TestModeSelection(t *testing.T) {
	tests := []struct {
		name     string
		mode     Mode
		contains string
	}{
		{
			name:     "shell-function outputs command only",
			mode:     ModeShellFunction,
			contains: "ls -la\n",
		},
		{
			name:     "stdout outputs with formatting",
			mode:     ModeStdout,
			contains: "✓ Selected command:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.mode)
			cmd := &commands.Option{
				Title:       "List files",
				Command:     "ls -la",
				Description: "List all files",
			}

			output := captureOutput(func() {
				err := handler.Output(cmd)
				if err != nil {
					t.Errorf("Output() error = %v", err)
				}
			})

			if !strings.Contains(output, tt.contains) {
				t.Errorf("Output() = %q, want to contain %q", output, tt.contains)
			}
		})
	}
}

func TestOutputShellFunction(t *testing.T) {
	handler := &Handler{mode: ModeShellFunction}
	cmd := &commands.Option{
		Command: "git status",
	}

	output := captureOutput(func() {
		err := handler.outputShellFunction(cmd)
		if err != nil {
			t.Errorf("outputShellFunction() error = %v", err)
		}
	})

	expected := "git status\n"
	if output != expected {
		t.Errorf("outputShellFunction() = %q, want %q", output, expected)
	}
}

func TestOutputStdoutFormatting(t *testing.T) {
	handler := &Handler{mode: ModeStdout}
	cmd := &commands.Option{
		Command: "docker ps -a",
	}

	output := captureOutput(func() {
		err := handler.outputStdout(cmd)
		if err != nil {
			t.Errorf("outputStdout() error = %v", err)
		}
	})

	if !strings.HasPrefix(output, "\n✓") {
		t.Errorf("outputStdout() should start with newline and checkmark, got %q", output)
	}
	if !strings.Contains(output, "docker ps -a") {
		t.Errorf("outputStdout() missing command, got %q", output)
	}
}

func TestDefaultModeIsClipboard(t *testing.T) {
	handler := &Handler{mode: "invalid"}
	cmd := &commands.Option{
		Command: "echo test",
	}

	output := captureOutput(func() {
		// Error is acceptable if clipboard tools are missing
		_ = handler.Output(cmd)
	})

	if !strings.Contains(output, "echo test") {
		t.Errorf("Output() with invalid mode missing command, got %q", output)
	}
}
