package processes

import "github.com/spf13/cobra"

func newReconciliationCmd() *cobra.Command {
	reconciliationCmd := &cobra.Command{
		Use:   "reconciliation",
		Short: "Business verbs for financial reconciliation workflows",
	}

	reconciliationCmd.AddCommand(newReconciliationRunCmd())
	reconciliationCmd.AddCommand(newCollectionListCmd("list-entries", "List reconciliation entries", "financial-reconciliation", "entries"))
	reconciliationCmd.AddCommand(newCollectionListCmd("list-extracts", "List reconciliation extracts", "financial-reconciliation", "extracts"))
	reconciliationCmd.AddCommand(newCollectionListCmd("list-exports", "List reconciliation exports", "financial-reconciliation", "exports"))
	reconciliationCmd.AddCommand(newCollectionListCmd("list-matches", "List reconciliation matches", "financial-reconciliation", "matches"))

	return reconciliationCmd
}

func newReconciliationRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reconcile [process_id]",
		Short: "Run financial reconciliation",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			payload := map[string]any{}
			if len(args) == 1 {
				processID, err := parsePositiveIntArg("process_id", args[0])
				if err != nil {
					return err
				}
				payload["process_id"] = processID
			}
			return runModuleAction(cmd, "financial-reconciliation", "reconcile", payload)
		},
	}
}
