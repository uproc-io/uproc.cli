package processes

import "github.com/spf13/cobra"

func newTaxReportingCmd() *cobra.Command {
	taxCmd := &cobra.Command{
		Use:   "tax-reporting",
		Short: "Business verbs for tax reporting workflows",
	}

	taxCmd.AddCommand(newTaxGenerateCmd())
	taxCmd.AddCommand(newCollectionListCmd("list", "List tax reports", "tax-reporting", "reports"))
	taxCmd.AddCommand(newTaxRecalculateCmd())
	taxCmd.AddCommand(newTaxValidateCmd())
	taxCmd.AddCommand(newTaxExportCmd())

	return taxCmd
}

func newTaxCmd() *cobra.Command {
	cmd := newTaxReportingCmd()
	cmd.Use = "tax"
	cmd.Hidden = true
	cmd.Long = "DEPRECATED: use 'tax-reporting' instead"
	return cmd
}

func newTaxGenerateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate <report_id>",
		Short: "Generate a draft tax report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "tax-reporting", "generate_report", "report_id", args[0])
		},
	}
}

func newTaxRecalculateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "recalculate <report_id>",
		Short: "Recalculate a tax report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "tax-reporting", "recalculate_report", "report_id", args[0])
		},
	}
}

func newTaxValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <report_id>",
		Short: "Validate a tax report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "tax-reporting", "validate_report", "report_id", args[0])
		},
	}
}

func newTaxExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export <report_id>",
		Short: "Mark a tax report as exported",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "tax-reporting", "export_report", "report_id", args[0])
		},
	}
}
