package config

import "github.com/spf13/cobra"

func newOpenConfigCommand() *cobra.Command {
	return &cobra.Command{
		Use:           "open-config",
		Short:         "Open argus config file",
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return openConfigEditor()
		},
	}
}
