package processes

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	processesCmd := &cobra.Command{
		Use:   "processes",
		Short: "Commands for Uproc Processes API",
	}

	processesCmd.AddCommand(newLoginCmd())
	processesCmd.AddCommand(newRequestCmd())
	processesCmd.AddCommand(newModuleCmd())
	processesCmd.AddCommand(newAdminCmd())
	processesCmd.AddCommand(newInstallCmd())
	processesCmd.AddCommand(newUpdateCmd())
	processesCmd.AddCommand(newInteractiveCmd())

	return processesCmd
}
