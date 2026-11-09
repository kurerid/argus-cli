package report

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
)

//go:embed templates/*
var templatesFS embed.FS

type Service struct {
	tmpl *template.Template
	css  template.CSS
	js   template.JS
}

type templateData struct {
	Report
	CSS template.CSS
	JS  template.JS
}

func NewService() (*Service, error) {
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
		tmpl: tmpl,
		css:  template.CSS(css),
		js:   template.JS(js),
	}, nil
}

func (s *Service) Generate(data Report) error {
	var buf bytes.Buffer

	view := templateData{
		Report: data,
		CSS:    s.css,
		JS:     s.js,
	}

	if err := s.tmpl.ExecuteTemplate(&buf, "report.html", view); err != nil {
		return fmt.Errorf("failed to execute report template: %w", err)
	}

	if err := os.WriteFile("diff-report.html", buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	return nil
}
