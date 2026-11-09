package report

import (
	"encoding/json"
	"fmt"
	"time"
)

type Report struct {
	Product   string
	CreatedAt time.Time

	Branches Branches
	Repo     Repository

	Completeness Completeness
	Statistics   Statistics

	Analysis Analysis
}

type Branches struct {
	Base   string
	Target string
}

type Repository struct {
	Name      string
	BaseSHA   string
	TargetSHA string
}

type Completeness struct {
	Complete bool
	Reasons  []string
}

type Statistics struct {
	Available    bool
	FilesChanged int
	Insertions   int
	Deletions    int
	Excluded     int
}

type Analysis struct {
	Summary     string              `json:"summary"`
	Product     ProductAnalysis     `json:"product"`
	Development DevelopmentAnalysis `json:"development"`
	QA          []Finding           `json:"qa"`
}

// Finding describes an analysis conclusion and the files supporting it.
type Finding struct {
	Description string          `json:"description"`
	Files       []FileReference `json:"files"`
}

type FileReference struct {
	Repository string `json:"repository"`
	Path       string `json:"path"`
}

type ProductAnalysis struct {
	Scenarios []Finding        `json:"scenarios"`
	Effects   []BusinessEffect `json:"effects"`
}

// BusinessEffect is a hypothesis to be measured after release.
type BusinessEffect struct {
	Finding
	Metrics []string `json:"metrics"`
}

type DevelopmentAnalysis struct {
	Changes        []Finding `json:"changes"`
	Risks          []Finding `json:"risks"`
	Unclear        []Finding `json:"unclear"`
	Contradictions []Finding `json:"contradictions"`
}

func ParseLLMAnalysis(content string) (Analysis, error) {
	var analysis Analysis

	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		return Analysis{}, fmt.Errorf("failed to parse LLM analysis: %w", err)
	}

	return analysis, nil
}
