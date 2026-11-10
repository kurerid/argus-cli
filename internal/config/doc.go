// Package config
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const configPath = "./argus-config.json"

type Application struct {
	Git `json:"git"`
	LLM `json:"llm"`
}

type Git struct {
	UserEmail  string `json:"user_email"`
	UserName   string `json:"user_name"`
	SSHKeyPath string `json:"ssh_key_path"`
}

type LLM struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
	APIMode string `json:"api_mode"`
}

func Load() (*Application, error) {
	data, err := os.ReadFile(configPath)
	if errors.Is(err, os.ErrNotExist) {
		app := defaultConfig()
		if err := save(configPath, app); err != nil {
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
		Git: Git{},
		LLM: LLM{},
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
