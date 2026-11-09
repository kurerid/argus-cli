package git

import (
	"argus-agent/internal/llm"

	"github.com/spf13/cobra"
)

func (s *Service) newCloneRepoCommand() *cobra.Command {
	return &cobra.Command{
		Use:           "clone-repo <url>",
		Short:         "Clone git repository",
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.CloneRepo(args[0])
		},
		Args: cobra.ExactArgs(1),
	}
}

func (s *Service) newCompareCommand() *cobra.Command {
	var (
		model           string
		temperature     float64
		maxOutputTokens int64
	)

	cmd := &cobra.Command{
		Use:           "compare-branches <repo> <left> <right>",
		Short:         "Compare branches diff",
		SilenceErrors: true,
		Args:          cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.CompareDiff(
				args[0],
				args[1],
				args[2],
				model,
				temperature,
				maxOutputTokens,
			)
		},
	}

	cmd.Flags().StringVar(
		&model,
		"model",
		llm.DEFAULT_MODEL,
		"LLM model",
	)

	cmd.Flags().Float64Var(
		&temperature,
		"temperature",
		llm.DEFAULT_TEMPERATURE,
		"LLM temperature",
	)

	cmd.Flags().Int64Var(
		&maxOutputTokens,
		"max-output-tokens",
		llm.DEFAULT_MAX_OUTPUT_TOKENS,
		"Maximum output tokens",
	)

	return cmd
}
