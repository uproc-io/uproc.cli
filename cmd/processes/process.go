package processes

import "github.com/spf13/cobra"

func newProcessVisibilityCmd() *cobra.Command {
	processCmd := &cobra.Command{
		Use:   "process-visibility",
		Short: "Business verbs for process visibility workflows",
	}

	processCmd.AddCommand(newProcessRetryStepCmd())
	processCmd.AddCommand(newCollectionListCmd("list", "List tracked processes", "process-visibility", "processes"))
	processCmd.AddCommand(newProcessReassignOwnerCmd())
	processCmd.AddCommand(newProcessCancelCmd())

	return processCmd
}

func newProcessCmd() *cobra.Command {
	cmd := newProcessVisibilityCmd()
	cmd.Use = "process"
	cmd.Hidden = true
	cmd.Long = "DEPRECATED: use 'process-visibility' instead"
	return cmd
}

func newProcessRetryStepCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "retry-step <process_id>",
		Short: "Retry the current step of a tracked process",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "process-visibility", "retry_step", "process_id", args[0])
		},
	}
}

func newProcessReassignOwnerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reassign-owner <process_id>",
		Short: "Reassign a tracked process to the fallback owner flow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "process-visibility", "reassign_owner", "process_id", args[0])
		},
	}
}

func newProcessCancelCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cancel <process_id>",
		Short: "Cancel a tracked process",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "process-visibility", "cancel_process", "process_id", args[0])
		},
	}
}
