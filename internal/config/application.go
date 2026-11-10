// Package config
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Application struct {
	Git `json:"git"`
	LLM `json:"llm"`
}

func Load() (*Application, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		app := defaultConfig()
		if err := save(path, app); err != nil {
			return nil, fmt.Errorf("create config %s: %w", configPath, err)
		}
		return app, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", configPath, err)
	}

	app := defaultConfig()
	if err := json.Unmarshal(data, app); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", configPath, err)
	}
	return app, nil
}

func defaultConfig() *Application {
	return &Application{
		Git: loadSystemGitSettings(),
		LLM: LLM{
			BaseURL: "api.openai.com",
		},
	}
}

// save записывает конфиг в новый файл. O_EXCL гарантирует,
// что существующий файл не будет случайно перезаписан.
func save(path string, app *Application) error {
	data, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	// 0600: в файле лежит API-ключ, читать его должен только владелец.
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir = filepath.Join(dir, "argus")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}
