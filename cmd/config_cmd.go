package cmd

import (
	"fmt"

	"bizzmod-cli/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Config file helpers",
	}

	configCmd.AddCommand(newConfigPathCmd())
	return configCmd
}

func newConfigPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Show resolved config.yml path",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := config.ResolvedConfigPath()
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), path)
			return nil
		},
	}
}
