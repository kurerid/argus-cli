package report

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

//go:embed templates/*
var templatesFS embed.FS

type Service struct {
	tmpl      *template.Template
	css       template.CSS
	js        template.JS
	outputDir string
}

type templateData struct {
	Report
	CSS template.CSS
	JS  template.JS
}

// DefaultDir возвращает директорию для отчётов внутри пользовательского кеша:
// Linux: ~/.cache/arugs/reports, macOS: ~/Library/Caches/arugs/reports,
// Windows: %LocalAppData%\arugs\reports.
func DefaultDir() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve cache dir: %w", err)
	}
	return filepath.Join(dir, "arugs", "reports"), nil
}

func NewService(outputDir string) (*Service, error) {
	tmpl, err := template.New("report.html").
		ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse report template: %w", err)
	}

	css, err := templatesFS.ReadFile("templates/styles.css")
	if err != nil {
		return nil, fmt.Errorf("failed to read report styles: %w", err)
	}

	js, err := templatesFS.ReadFile("templates/scripts.js")
	if err != nil {
		return nil, fmt.Errorf("failed to read report scripts: %w", err)
	}

	return &Service{
		tmpl:      tmpl,
		css:       template.CSS(css),
		js:        template.JS(js),
		outputDir: outputDir,
	}, nil
}

// Generate рендерит отчёт и возвращает путь к созданному файлу.
func (s *Service) Generate(data Report) (string, error) {
	var buf bytes.Buffer

	view := templateData{
		Report: data,
		CSS:    s.css,
		JS:     s.js,
	}

	if err := s.tmpl.ExecuteTemplate(&buf, "report.html", view); err != nil {
		return "", fmt.Errorf("failed to execute report template: %w", err)
	}

	if err := os.MkdirAll(s.outputDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create reports dir: %w", err)
	}

	name := fmt.Sprintf("diff-report-%s.html", time.Now().Format("20060102-150405"))
	path := filepath.Join(s.outputDir, name)

	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return "", fmt.Errorf("failed to write report: %w", err)
	}

	return path, nil
}
