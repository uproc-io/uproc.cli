package processes

import "github.com/spf13/cobra"

func newEmailCmd() *cobra.Command {
	emailCmd := &cobra.Command{
		Use:   "email",
		Short: "Business verbs for email assistant workflows",
	}

	emailCmd.AddCommand(newEmailMarkProcessedCmd())
	emailCmd.AddCommand(newCollectionListCmd("list", "List email assistant emails", "email-assistant", "emails"))
	emailCmd.AddCommand(newEmailArchiveCmd())

	return emailCmd
}

func newEmailMarkProcessedCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mark-processed <email_id>",
		Short: "Mark an email assistant item as processed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "email-assistant", "mark_processed", "email_id", args[0])
		},
	}
}

func newEmailArchiveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "archive <email_id>",
		Short: "Archive an email assistant item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "email-assistant", "archive_email", "email_id", args[0])
		},
	}
}
