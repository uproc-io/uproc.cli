package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newDataCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "data",
		Short: "Data workflows (under construction)",
		Long:  "Data command group is under construction. Dedicated data workflows will be added here. See https://uproc.io for updates.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "data is under construction")
			fmt.Fprintln(cmd.OutOrStdout(), "Use `uproc data --help` to check updates as new commands are added.")
			fmt.Fprintln(cmd.OutOrStdout(), "More info: https://uproc.io")
		},
	}
}
