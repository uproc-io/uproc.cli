package processes

import "github.com/spf13/cobra"

func newSigningCmd() *cobra.Command {
	signingCmd := &cobra.Command{
		Use:   "signing",
		Short: "Business verbs for document signing workflows",
	}

	signingCmd.AddCommand(newSigningCancelCmd())
	signingCmd.AddCommand(newCollectionListCmd("list", "List document-signing requests", "document-signing", "requests"))
	signingCmd.AddCommand(newSigningReopenCmd())
	signingCmd.AddCommand(newSigningSendReminderCmd())
	signingCmd.AddCommand(newSigningSyncStatusCmd())

	return signingCmd
}

func newSigningCancelCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cancel <request_id>",
		Short: "Cancel a signing request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "document-signing", "cancel_request", "request_id", args[0])
		},
	}
}

func newSigningReopenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reopen <request_id>",
		Short: "Reopen a signing request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "document-signing", "reopen_request", "request_id", args[0])
		},
	}
}

func newSigningSendReminderCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "send-reminder <request_id>",
		Short: "Send a signing reminder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "document-signing", "send_reminder", "request_id", args[0])
		},
	}
}

func newSigningSyncStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync-status <request_id>",
		Short: "Sync a signing request status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "document-signing", "sync_status", "request_id", args[0])
		},
	}
}
