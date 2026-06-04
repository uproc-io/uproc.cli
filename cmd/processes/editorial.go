package processes

import "github.com/spf13/cobra"

func newEditorialCmd() *cobra.Command {
	editorialCmd := &cobra.Command{
		Use:   "editorial",
		Short: "Business verbs for editorial engine workflows",
	}

	editorialCmd.AddCommand(newEditorialGenerateProposalCmd())
	editorialCmd.AddCommand(newCollectionListCmd("list-opportunities", "List editorial opportunities", "editorial-engine", "opportunities"))
	editorialCmd.AddCommand(newCollectionListCmd("list-projects", "List editorial projects", "editorial-engine", "projects"))
	editorialCmd.AddCommand(newCollectionListCmd("list-articles", "List editorial articles", "editorial-engine", "articles"))
	editorialCmd.AddCommand(newCollectionListCmd("list-combined", "List combined editorial view", "editorial-engine", "combined"))
	editorialCmd.AddCommand(newEditorialGenerateArticleCmd())
	editorialCmd.AddCommand(newEditorialPublishCmd())
	editorialCmd.AddCommand(newEditorialScheduleCmd())
	editorialCmd.AddCommand(newEditorialDiscardCmd())

	return editorialCmd
}

func newEditorialGenerateProposalCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate-proposal <opportunity_id>",
		Short: "Generate an editorial proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "editorial-engine", "generate_proposal", "opportunity_id", args[0])
		},
	}
}

func newEditorialGenerateArticleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate-article <opportunity_id>",
		Short: "Generate an editorial article draft",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "editorial-engine", "generate_article", "opportunity_id", args[0])
		},
	}
}

func newEditorialPublishCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "publish <opportunity_id>",
		Short: "Mark an editorial opportunity as published",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "editorial-engine", "publish", "opportunity_id", args[0])
		},
	}
}

func newEditorialScheduleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "schedule <opportunity_id>",
		Short: "Schedule an editorial opportunity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "editorial-engine", "schedule", "opportunity_id", args[0])
		},
	}
}

func newEditorialDiscardCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "discard <opportunity_id>",
		Short: "Discard an editorial opportunity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "editorial-engine", "discard", "opportunity_id", args[0])
		},
	}
}
