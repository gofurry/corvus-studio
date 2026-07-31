package directorypicker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

const cancelledMarker = "__CORVUS_DIRECTORY_PICKER_CANCELLED__"

var ErrUnavailable = errors.New("native directory picker unavailable")

type lookPathFunc func(string) (string, error)
type runCommandFunc func(context.Context, string, ...string) ([]byte, error)

type Picker struct {
	platform string
	lookPath lookPathFunc
	run      runCommandFunc
}

type commandSpec struct {
	name            string
	args            []string
	cancelExitCodes []int
}

func New() *Picker {
	return &Picker{
		platform: runtime.GOOS,
		lookPath: exec.LookPath,
		run: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return exec.CommandContext(ctx, name, args...).Output()
		},
	}
}

func (picker *Picker) Select(ctx context.Context) (string, bool, error) {
	command, err := picker.command()
	if err != nil {
		return "", false, err
	}
	output, err := picker.run(ctx, command.name, command.args...)
	if err != nil {
		if ctx.Err() != nil {
			return "", false, ctx.Err()
		}
		var exitError interface{ ExitCode() int }
		if errors.As(err, &exitError) && slices.Contains(command.cancelExitCodes, exitError.ExitCode()) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("%w: %s failed: %v", ErrUnavailable, filepath.Base(command.name), err)
	}

	selected := strings.TrimRight(string(output), "\r\n")
	if selected == "" || selected == cancelledMarker {
		return "", false, nil
	}
	selected = filepath.Clean(selected)
	if !filepath.IsAbs(selected) {
		return "", false, fmt.Errorf("%w: picker returned a non-absolute path", ErrUnavailable)
	}
	info, err := os.Stat(selected)
	if err != nil || !info.IsDir() {
		return "", false, fmt.Errorf("%w: picker returned an unavailable directory", ErrUnavailable)
	}
	return selected, true, nil
}

func (picker *Picker) command() (commandSpec, error) {
	switch picker.platform {
	case "windows":
		name, err := picker.firstExecutable("powershell.exe", "pwsh.exe", "pwsh")
		if err != nil {
			return commandSpec{}, err
		}
		return commandSpec{
			name: name,
			args: []string{
				"-NoProfile",
				"-NonInteractive",
				"-STA",
				"-Command",
				windowsScript,
			},
		}, nil
	case "darwin":
		name, err := picker.firstExecutable("osascript")
		if err != nil {
			return commandSpec{}, err
		}
		return commandSpec{
			name: name,
			args: []string{"-e", macOSScript},
		}, nil
	case "linux":
		if name, err := picker.lookPath("zenity"); err == nil {
			return commandSpec{
				name:            name,
				args:            []string{"--file-selection", "--directory", "--title=Select a project directory"},
				cancelExitCodes: []int{1},
			}, nil
		}
		if name, err := picker.lookPath("kdialog"); err == nil {
			startDirectory, homeErr := os.UserHomeDir()
			if homeErr != nil {
				startDirectory = "."
			}
			return commandSpec{
				name:            name,
				args:            []string{"--getexistingdirectory", startDirectory, "--title", "Select a project directory"},
				cancelExitCodes: []int{1},
			}, nil
		}
		return commandSpec{}, fmt.Errorf("%w: install zenity or kdialog", ErrUnavailable)
	default:
		return commandSpec{}, fmt.Errorf("%w: unsupported platform %s", ErrUnavailable, picker.platform)
	}
}

func (picker *Picker) firstExecutable(names ...string) (string, error) {
	for _, name := range names {
		if path, err := picker.lookPath(name); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("%w: no supported picker command found", ErrUnavailable)
}

const windowsScript = `
Add-Type -AssemblyName System.Windows.Forms
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$dialog = New-Object System.Windows.Forms.FolderBrowserDialog
$dialog.Description = 'Select a project directory'
$dialog.ShowNewFolderButton = $false
if ($dialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
    [Console]::Out.Write($dialog.SelectedPath)
} else {
    [Console]::Out.Write('__CORVUS_DIRECTORY_PICKER_CANCELLED__')
}
`

const macOSScript = `
try
    return POSIX path of (choose folder with prompt "Select a project directory")
on error number -128
    return "__CORVUS_DIRECTORY_PICKER_CANCELLED__"
end try
`
