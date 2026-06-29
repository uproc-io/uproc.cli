package processes

import "github.com/spf13/cobra"

func newOrderTrackCmd() *cobra.Command {
	orderCmd := &cobra.Command{
		Use:   "order-track",
		Short: "Business verbs for tracked order workflows",
	}

	orderCmd.AddCommand(newOrderMarkReceivedCmd())
	orderCmd.AddCommand(newCollectionListCmd("list", "List tracked orders", "order-track", "orders"))
	orderCmd.AddCommand(newOrderCancelCmd())
	orderCmd.AddCommand(newOrderSendReminderCmd())

	return orderCmd
}

func newOrderCmd() *cobra.Command {
	cmd := newOrderTrackCmd()
	cmd.Use = "order"
	cmd.Hidden = true
	cmd.Long = "DEPRECATED: use 'order-track' instead"
	return cmd
}

func newOrderMarkReceivedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mark-received <order_id>",
		Short: "Mark a tracked order as received",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "order-track", "mark_received", "order_id", args[0])
		},
	}
}

func newOrderCancelCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cancel <order_id>",
		Short: "Cancel a tracked order",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "order-track", "cancel_order", "order_id", args[0])
		},
	}
}

func newOrderSendReminderCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "send-reminder <order_id>",
		Short: "Send a reminder on a tracked order",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "order-track", "send_reminder", "order_id", args[0])
		},
	}
}
