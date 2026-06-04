package processes

import "github.com/spf13/cobra"

func newInventoryCmd() *cobra.Command {
	inventoryCmd := &cobra.Command{
		Use:   "inventory",
		Short: "Business verbs for inventory planning workflows",
	}

	inventoryCmd.AddCommand(newInventoryMarkReceivedCmd())
	inventoryCmd.AddCommand(newCollectionListCmd("list", "List inventory planning orders", "inventory-planning", "orders"))
	inventoryCmd.AddCommand(newInventoryCancelCmd())
	inventoryCmd.AddCommand(newInventorySendReminderCmd())

	return inventoryCmd
}

func newInventoryMarkReceivedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mark-received <order_id>",
		Short: "Mark an inventory planning order as received",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "inventory-planning", "mark_received", "order_id", args[0])
		},
	}
}

func newInventoryCancelCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cancel <order_id>",
		Short: "Cancel an inventory planning order",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "inventory-planning", "cancel_order", "order_id", args[0])
		},
	}
}

func newInventorySendReminderCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "send-reminder <order_id>",
		Short: "Flag an inventory planning order for reminder follow-up",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "inventory-planning", "send_reminder", "order_id", args[0])
		},
	}
}
