package processes

import "github.com/spf13/cobra"

func newLeadIntelligenceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lead-intelligence",
		Short: "Business verbs for lead-intelligence",
	}
	cmd.AddCommand(newCollectionListCmd("list-leads", "List lead-intelligence leads", "lead-intelligence", "leads"))
	cmd.AddCommand(newCollectionListCmd("list-executions", "List lead-intelligence executions", "lead-intelligence", "executions"))
	cmd.AddCommand(newCollectionListCmd("list-sync-logs", "List lead-intelligence sync logs", "lead-intelligence", "sync-logs"))
	cmd.AddCommand(newCollectionListCmd("list-duplicates", "List lead-intelligence duplicates", "lead-intelligence", "duplicates"))
	cmd.AddCommand(newLeadIntelligenceAnalyzePlanCmd())
	cmd.AddCommand(newLeadIntelligenceListPlansCmd())
	cmd.AddCommand(newLeadIntelligenceGetPlanCmd())
	cmd.AddCommand(newLeadIntelligenceExecutePlanCmd())
	cmd.AddCommand(newLeadIntelligenceRegeneratePlanCmd())
	return cmd
}

func newLeadIntelligenceAnalyzePlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "analyze-plan",
		Short: "Generate AI-powered lead intelligence plan",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runModuleAction(cmd, "lead-intelligence", "analyze_plan", map[string]any{})
		},
	}
}

func newLeadIntelligenceListPlansCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-plans",
		Short: "List lead intelligence plans",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runModuleAction(cmd, "lead-intelligence", "list_plans", map[string]any{})
		},
	}
}

func newLeadIntelligenceGetPlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-plan <plan_id>",
		Short: "Get a lead intelligence plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "lead-intelligence", "get_plan", "plan_id", args[0])
		},
	}
}

func newLeadIntelligenceExecutePlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "execute-plan <plan_id>",
		Short: "Execute an accepted lead intelligence plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "lead-intelligence", "execute_plan", "plan_id", args[0])
		},
	}
}

func newLeadIntelligenceRegeneratePlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "regenerate-plan <plan_id>",
		Short: "Regenerate a lead intelligence plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "lead-intelligence", "regenerate_plan", "plan_id", args[0])
		},
	}
}
