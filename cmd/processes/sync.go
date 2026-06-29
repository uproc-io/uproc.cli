package processes

import "github.com/spf13/cobra"

func newDataSyncCmd() *cobra.Command {
	syncCmd := &cobra.Command{
		Use:   "data-sync",
		Short: "Business verbs for data-sync",
	}

	syncCmd.AddCommand(newSyncRunCmd())
	syncCmd.AddCommand(newCollectionListCmd("list-workflows", "List data sync workflows", "data-sync", "workflows"))
	syncCmd.AddCommand(newCollectionListCmd("list-runs", "List data sync runs", "data-sync", "runs"))
	syncCmd.AddCommand(newCollectionListCmd("list-records", "List data sync records", "data-sync", "records"))
	syncCmd.AddCommand(newSyncPreviewCmd())
	syncCmd.AddCommand(newSyncDryRunCmd())

	return syncCmd
}

func newSyncCmd() *cobra.Command {
	cmd := newDataSyncCmd()
	cmd.Use = "sync"
	cmd.Hidden = true
	cmd.Long = "DEPRECATED: Use \"data-sync\" instead.\n\n" + cmd.Long
	return cmd
}

func newSyncRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run <workflow_id>",
		Short: "Run a data sync workflow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "data-sync", "run_workflow", "workflow_id", args[0])
		},
	}
}

func newSyncPreviewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "preview <workflow_id> [limit]",
		Short: "Preview a data sync workflow",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID, err := parsePositiveIntArg("workflow_id", args[0])
			if err != nil {
				return err
			}
			payload := map[string]any{"workflow_id": workflowID}
			if len(args) == 2 {
				limit, err := parsePositiveIntArg("limit", args[1])
				if err != nil {
					return err
				}
				payload["limit"] = limit
			}
			return runModuleAction(cmd, "data-sync", "preview_workflow", payload)
		},
	}
}

func newSyncDryRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dry-run <workflow_id> [limit]",
		Short: "Dry-run a data sync workflow",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID, err := parsePositiveIntArg("workflow_id", args[0])
			if err != nil {
				return err
			}
			payload := map[string]any{"workflow_id": workflowID}
			if len(args) == 2 {
				limit, err := parsePositiveIntArg("limit", args[1])
				if err != nil {
					return err
				}
				payload["limit"] = limit
			}
			return runModuleAction(cmd, "data-sync", "dry_run_workflow", payload)
		},
	}
}
