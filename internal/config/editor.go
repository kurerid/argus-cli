package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var EditorService = editor{}

type editor struct{}

func (e editor) InitCommands(root *cobra.Command) {
	root.AddCommand(newOpenConfigCommand())
}

func openConfigEditor() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := ensureConfigExists(path); err != nil {
		return err
	}

	editor, err := findEditor()
	if err != nil {
		return err
	}

	// editor может содержать аргументы, например "code --wait"
	cmd := exec.Command(editor[0], append(editor[1:], path)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run editor %q: %w", editor[0], err)
	}
	return nil
}

// findEditor возвращает команду редактора и её аргументы.
func findEditor() ([]string, error) {
	for _, env := range []string{"ARGUS_EDITOR", "VISUAL", "EDITOR"} {
		if v := strings.TrimSpace(os.Getenv(env)); v != "" {
			return strings.Fields(v), nil
		}
	}

	var candidates []string
	if runtime.GOOS == "windows" {
		candidates = []string{"edit", "notepad"}
	} else {
		candidates = []string{"nano", "vim", "vi"}
	}

	for _, c := range candidates {
		if p, err := exec.LookPath(c); err == nil {
			return []string{p}, nil
		}
	}

	return nil, errors.New("no editor found: set the EDITOR environment variable")
}

func ensureConfigExists(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return save(path, defaultConfig())
}
