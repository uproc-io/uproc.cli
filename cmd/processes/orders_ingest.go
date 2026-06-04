package processes

import "github.com/spf13/cobra"

func newOrdersIngestCmd() *cobra.Command {
	ordersIngestCmd := &cobra.Command{
		Use:   "orders-ingest",
		Short: "Business verbs for order ingest workflows",
	}

	ordersIngestCmd.AddCommand(newOrdersIngestReprocessCmd())
	ordersIngestCmd.AddCommand(newCollectionListCmd("list", "List ingested orders", "order-ingest", "orders"))
	ordersIngestCmd.AddCommand(newCollectionListCmd("list-emails", "List ingest emails", "order-ingest", "emails"))
	ordersIngestCmd.AddCommand(newOrdersIngestValidateCmd())
	ordersIngestCmd.AddCommand(newOrdersIngestSendToERPCmd())

	return ordersIngestCmd
}

func newOrdersIngestReprocessCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reprocess <order_id>",
		Short: "Reprocess an ingested order",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "order-ingest", "reprocess_order", "order_id", args[0])
		},
	}
}

func newOrdersIngestValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <order_id>",
		Short: "Validate an ingested order",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "order-ingest", "validate_order", "order_id", args[0])
		},
	}
}

func newOrdersIngestSendToERPCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "send-to-erp <order_id>",
		Short: "Send an ingested order to ERP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "order-ingest", "send_order_to_erp", "order_id", args[0])
		},
	}
}
