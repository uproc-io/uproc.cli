package cmd

import (
	"fmt"

	"bizzmod-cli/internal/config"
	"github.com/spf13/cobra"
)

func newProfileCmd() *cobra.Command {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage local login profiles",
	}

	profileCmd.AddCommand(newProfileListCmd())
	profileCmd.AddCommand(newProfileUseCmd())
	profileCmd.AddCommand(newProfileShowCmd())

	return profileCmd
}

func newProfileListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, active, err := config.ListProfiles()
			if err != nil {
				return err
			}

			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(no profiles)")
				return nil
			}

			for _, name := range names {
				marker := " "
				if name == active {
					marker = "*"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", marker, name)
			}

			return nil
		},
	}
}

func newProfileUseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use <profile_name>",
		Short: "Set active profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.SetActiveProfile(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "ok: active profile set to %q\n", args[0])
			return nil
		},
	}
}

func newProfileShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show active profile name",
		RunE: func(cmd *cobra.Command, args []string) error {
			active, err := config.GetActiveProfileName()
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), active)
			return nil
		},
	}
}
