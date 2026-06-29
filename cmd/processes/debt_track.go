package processes

import "github.com/spf13/cobra"

func newDebtTrackCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debt-track",
		Short: "Business verbs for debt-track",
	}
	cmd.AddCommand(newCollectionListCmd("list-invoices", "List debt-track invoices", "debt-track", "invoices"))
	cmd.AddCommand(newCollectionListCmd("list-payments", "List debt-track payments", "debt-track", "payments"))
	cmd.AddCommand(newCollectionListCmd("list-reminders", "List debt-track reminders", "debt-track", "reminders"))
	cmd.AddCommand(newCollectionListCmd("list-communications", "List debt-track communications", "debt-track", "communications"))
	cmd.AddCommand(newDebtTrackAnalyzePlanCmd())
	cmd.AddCommand(newDebtTrackListPlansCmd())
	cmd.AddCommand(newDebtTrackGetPlanCmd())
	cmd.AddCommand(newDebtTrackExecutePlanCmd())
	cmd.AddCommand(newDebtTrackRegeneratePlanCmd())
	return cmd
}

func newDebtTrackAnalyzePlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "analyze-plan",
		Short: "Generate AI-powered collection strategy plan",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runModuleAction(cmd, "debt-track", "analyze_plan", map[string]any{})
		},
	}
}

func newDebtTrackListPlansCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-plans",
		Short: "List debt collection plans",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runModuleAction(cmd, "debt-track", "list_plans", map[string]any{})
		},
	}
}

func newDebtTrackGetPlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-plan <plan_id>",
		Short: "Get a debt collection plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "debt-track", "get_plan", "plan_id", args[0])
		},
	}
}

func newDebtTrackExecutePlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "execute-plan <plan_id>",
		Short: "Execute an accepted debt collection plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "debt-track", "execute_plan", "plan_id", args[0])
		},
	}
}

func newDebtTrackRegeneratePlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "regenerate-plan <plan_id>",
		Short: "Regenerate a debt collection plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "debt-track", "regenerate_plan", "plan_id", args[0])
		},
	}
}
