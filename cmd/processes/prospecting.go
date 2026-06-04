package processes

import "github.com/spf13/cobra"

func newProspectingCmd() *cobra.Command {
	prospectingCmd := &cobra.Command{
		Use:   "prospecting",
		Short: "Business verbs for lead prospecting workflows",
	}

	prospectingCmd.AddCommand(newProspectingRunDiscoveryCmd())
	prospectingCmd.AddCommand(newCollectionListCmd("list-strategies", "List prospecting strategies", "lead-prospecting", "strategies"))
	prospectingCmd.AddCommand(newCollectionListCmd("list-opportunities", "List prospecting opportunities", "lead-prospecting", "opportunities"))
	prospectingCmd.AddCommand(newCollectionListCmd("list-prospects", "List prospecting prospects", "lead-prospecting", "prospects"))
	prospectingCmd.AddCommand(newCollectionListCmd("list-executions", "List prospecting executions", "lead-prospecting", "executions"))
	prospectingCmd.AddCommand(newProspectingSendToLeadsCmd())

	return prospectingCmd
}

func newProspectingRunDiscoveryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run-discovery <strategy_id> [company] [domain]",
		Short: "Run prospecting discovery",
		Args:  cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			strategyID, err := parsePositiveIntArg("strategy_id", args[0])
			if err != nil {
				return err
			}
			payload := map[string]any{"strategy_id": strategyID}
			if len(args) >= 2 {
				payload["company"] = args[1]
			}
			if len(args) == 3 {
				payload["domain"] = args[2]
			}
			return runModuleAction(cmd, "lead-prospecting", "run_discovery", payload)
		},
	}
}

func newProspectingSendToLeadsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "send-to-leads <opportunity_id>",
		Short: "Send a prospecting opportunity into lead management",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "lead-prospecting", "send_opportunity_to_leads", "opportunity_id", args[0])
		},
	}
}
