package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "bizzmod",
		Short: "Bizzmod CLI for /api/v1/external",
	}

	rootCmd.AddCommand(newLoginCmd())
	rootCmd.AddCommand(newRequestCmd())
	rootCmd.AddCommand(newModuleCmd())
	rootCmd.AddCommand(newInteractiveCmd())

	return rootCmd
}
