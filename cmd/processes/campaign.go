package processes

import "github.com/spf13/cobra"

func newCampaignCmd() *cobra.Command {
	campaignCmd := &cobra.Command{
		Use:   "campaign",
		Short: "Business verbs for campaign automation workflows",
	}

	campaignCmd.AddCommand(newCampaignPreviewAudienceCmd())
	campaignCmd.AddCommand(newCollectionListCmd("list", "List campaign-automation campaigns", "campaign-automation", "campaigns"))
	campaignCmd.AddCommand(newCollectionListCmd("list-audiences", "List campaign-automation audiences", "campaign-automation", "audiences"))
	campaignCmd.AddCommand(newCampaignAddAudienceCmd())
	campaignCmd.AddCommand(newCampaignPauseCmd())
	campaignCmd.AddCommand(newCampaignActivateCmd())

	return campaignCmd
}

func newCampaignPreviewAudienceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "preview-audience <campaign_id> [limit]",
		Short: "Preview a campaign automation audience",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			campaignID, err := parsePositiveIntArg("campaign_id", args[0])
			if err != nil {
				return err
			}
			payload := map[string]any{"campaign_id": campaignID}
			if len(args) == 2 {
				limit, err := parsePositiveIntArg("limit", args[1])
				if err != nil {
					return err
				}
				payload["limit"] = limit
			}
			return runModuleAction(cmd, "campaign-automation", "preview_audience", payload)
		},
	}
}

func newCampaignAddAudienceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add-audience <campaign_id> [mode]",
		Short: "Add an audience definition to a campaign",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			campaignID, err := parsePositiveIntArg("campaign_id", args[0])
			if err != nil {
				return err
			}
			payload := map[string]any{"campaign_id": campaignID}
			if len(args) == 2 {
				payload["mode"] = args[1]
			}
			return runModuleAction(cmd, "campaign-automation", "add_audience", payload)
		},
	}
}

func newCampaignPauseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pause <campaign_id>",
		Short: "Pause a campaign automation flow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "campaign-automation", "pause", "campaign_id", args[0])
		},
	}
}

func newCampaignActivateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "activate <campaign_id>",
		Short: "Activate a campaign automation flow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "campaign-automation", "activate", "campaign_id", args[0])
		},
	}
}
