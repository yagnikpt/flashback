package contentloaders

import (
	"context"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// GetClipboardContent retrieves text from the clipboard across different operating systems.
// It uses platform-specific approaches to ensure maximum compatibility.
func GetClipboardContent() ([]string, error) {
	var text []string
	var err error

	switch runtime.GOOS {
	case "linux":
		text, err = getLinuxClipboard()
	case "windows":
		text, err = getWindowsClipboard()
	case "darwin":
		text, err = getMacClipboard()
	default:
		return nil, errors.New("unsupported operating system: " + runtime.GOOS)
	}

	if err != nil {
		return nil, err
	}

	var nonEmptyLines []string
	for _, line := range text {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	if len(nonEmptyLines) == 0 && len(text) > 0 {
		return text, nil
	}

	return nonEmptyLines, nil
}

// getLinuxClipboard tries different clipboard tools available on Linux
func getLinuxClipboard() ([]string, error) {
	var errMessages []string

	// Check if running in Wayland
	wayland := isWayland()

	if wayland {
		// Try Wayland's wl-paste first when in Wayland
		text, err := runCommandWithTimeout("wl-paste", []string{}, 500*time.Millisecond)
		if err == nil && text != "" {
			return strings.Split(text, "\n"), nil
		}
		if err != nil {
			errMessages = append(errMessages, "wl-paste error: "+err.Error())
		}

		// Try wl-paste with -n flag (no newline)
		text, err = runCommandWithTimeout("wl-paste", []string{"-n"}, 500*time.Millisecond)
		if err == nil && text != "" {
			return strings.Split(text, "\n"), nil
		}
		if err != nil {
			errMessages = append(errMessages, "wl-paste -n error: "+err.Error())
		}
	}

	// Try X11's xclip
	text, err := runCommandWithTimeout("xclip", []string{"-out", "-selection", "clipboard"}, 500*time.Millisecond)
	if err == nil && text != "" {
		return strings.Split(text, "\n"), nil
	}
	if err != nil {
		errMessages = append(errMessages, "xclip error: "+err.Error())
	}

	// Try xclip with primary selection as fallback
	text, err = runCommandWithTimeout("xclip", []string{"-out", "-selection", "primary"}, 500*time.Millisecond)
	if err == nil && text != "" {
		return strings.Split(text, "\n"), nil
	}

	// Try X11's xsel as another alternative
	text, err = runCommandWithTimeout("xsel", []string{"--clipboard", "--output"}, 500*time.Millisecond)
	if err == nil && text != "" {
		return strings.Split(text, "\n"), nil
	}
	if err != nil {
		errMessages = append(errMessages, "xsel error: "+err.Error())
	}

	// Try xsel with primary selection as fallback
	text, err = runCommandWithTimeout("xsel", []string{"--primary", "--output"}, 500*time.Millisecond)
	if err == nil && text != "" {
		return strings.Split(text, "\n"), nil
	}

	return nil, errors.New("no suitable clipboard tool found on Linux: " + strings.Join(errMessages, "; "))
}

// getWindowsClipboard uses PowerShell to get clipboard content on Windows
func getWindowsClipboard() ([]string, error) {
	// Try PowerShell first
	text, err := runCommandWithTimeout("powershell.exe", []string{"-command", "Get-Clipboard"}, time.Second)
	if err == nil && text != "" {
		return strings.Split(text, "\n"), nil
	}

	// Try alternative clipboard command (Windows cmd)
	text, err = runCommandWithTimeout("cmd.exe", []string{"/c", "echo off && powershell -command \"Add-Type -AssemblyName System.Windows.Forms;[System.Windows.Forms.Clipboard]::GetText()\""}, time.Second)
	if err == nil && text != "" {
		return strings.Split(text, "\n"), nil
	}

	return nil, errors.New("failed to get clipboard content on Windows")
}

// getMacClipboard uses pbpaste to get clipboard content on macOS
func getMacClipboard() ([]string, error) {
	// Try pbpaste (standard macOS clipboard tool)
	text, err := runCommandWithTimeout("pbpaste", []string{}, 500*time.Millisecond)
	if err == nil && text != "" {
		return strings.Split(text, "\n"), nil
	}

	// Try with osascript as fallback
	script := `tell application "System Events" to get the clipboard as text`
	text, err = runCommandWithTimeout("osascript", []string{"-e", script}, 500*time.Millisecond)
	if err == nil && text != "" {
		return strings.Split(text, "\n"), nil
	}

	return nil, errors.New("failed to get clipboard content on macOS")
}

// isWayland checks if the current session is running in Wayland
func isWayland() bool {
	// Check for WAYLAND_DISPLAY environment variable
	output, err := runCommandWithTimeout("sh", []string{"-c", "echo $WAYLAND_DISPLAY"}, 100*time.Millisecond)
	if err == nil && output != "" {
		return true
	}

	// Check if the XDG_SESSION_TYPE is wayland
	output, err = runCommandWithTimeout("sh", []string{"-c", "echo $XDG_SESSION_TYPE"}, 100*time.Millisecond)
	if err == nil && strings.Contains(strings.ToLower(output), "wayland") {
		return true
	}

	// Alternative check with loginctl
	output, err = runCommandWithTimeout("loginctl", []string{"show-session", "$(loginctl | grep $(whoami) | awk '{print $1}')", "-p", "Type"}, 200*time.Millisecond)
	if err == nil && strings.Contains(output, "wayland") {
		return true
	}

	// Check if wl-paste exists and xclip doesn't
	_, wlErr := runCommandWithTimeout("which", []string{"wl-paste"}, 100*time.Millisecond)
	_, xclipErr := runCommandWithTimeout("which", []string{"xclip"}, 100*time.Millisecond)
	if wlErr == nil && xclipErr != nil {
		return true
	}

	return false
}

// runCommandWithTimeout executes a command with a timeout
func runCommandWithTimeout(command string, args []string, timeout time.Duration) (string, error) {
	// First check if the command exists
	_, err := exec.LookPath(command)
	if err != nil {
		return "", errors.New("command not found: " + command)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)
	output, err := cmd.CombinedOutput()

	// Handle context deadline exceeded errors
	if ctx.Err() == context.DeadlineExceeded {
		return "", errors.New("command timed out: " + command)
	}

	if err != nil {
		// Return the error along with any output that might have been produced
		if len(output) > 0 {
			return "", errors.New(err.Error() + ": " + strings.TrimSpace(string(output)))
		}
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
