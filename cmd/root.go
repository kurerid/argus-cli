package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "argus",
	Short:         "Argus CLI",
	Version:       version,
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type CommandInitializers interface {
	InitCommands(root *cobra.Command)
}

func InitCommands(cmds ...CommandInitializers) {
	for _, c := range cmds {
		c.InitCommands(rootCmd)
	}
}
