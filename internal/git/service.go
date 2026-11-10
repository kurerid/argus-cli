package git

import (
	"argus-agent/internal/config"
	"argus-agent/internal/llm"
	"argus-agent/internal/report"

	"github.com/spf13/cobra"
)

type GitService interface {
	CloneRepo(url string) error
	InitCommands(root *cobra.Command)
}

type Service struct {
	cfg           config.Git
	AI            llm.Provider
	reportService reportGenerator
}

type reportGenerator interface {
	Generate(data report.Report) (string, error)
}

func NewService(
	cfg config.Git,
	ai llm.Provider,
	reportService reportGenerator,
) (GitService, error) {
	return &Service{
		cfg:           cfg,
		AI:            ai,
		reportService: reportService,
	}, nil
}

func (s *Service) CloneRepo(url string) error {
	return cloneRepo(s.cfg, url)
}

func (s *Service) CompareDiff(
	pathToRepo, left, right, aiModel string,
	temperature float64,
	maxOutputTokens int64,
) error {
	return compareDiff(
		s.AI,
		s.reportService,
		pathToRepo,
		left,
		right,
		aiModel,
		temperature,
		maxOutputTokens,
	)
}

func (s *Service) InitCommands(root *cobra.Command) {
	root.AddCommand(
		s.newCloneRepoCommand(),
		s.newCompareCommand(),
	)
}
