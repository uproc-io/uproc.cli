package processes

import "github.com/spf13/cobra"

func newDocumentGeneratorCmd() *cobra.Command {
	documentsCmd := &cobra.Command{
		Use:   "document-generator",
		Short: "Business verbs for document generator workflows",
	}

	documentsCmd.AddCommand(newDocumentsMarkReadyCmd())
	documentsCmd.AddCommand(newCollectionListCmd("list", "List generated documents", "document-generator", "documents"))
	documentsCmd.AddCommand(newDocumentsMarkReviewCmd())
	documentsCmd.AddCommand(newDocumentsArchiveCmd())
	documentsCmd.AddCommand(newDocumentsRestoreCmd())
	documentsCmd.AddCommand(newDocumentsRegenerateCmd())

	return documentsCmd
}

func newDocumentsCmd() *cobra.Command {
	cmd := newDocumentGeneratorCmd()
	cmd.Use = "documents"
	cmd.Hidden = true
	cmd.Long = "DEPRECATED: use 'document-generator' instead"
	return cmd
}

func newDocumentsMarkReadyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mark-ready <document_id>",
		Short: "Mark a generated document as ready",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "document-generator", "mark_ready", "document_id", args[0])
		},
	}
}

func newDocumentsMarkReviewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mark-review <document_id>",
		Short: "Send a generated document back to review",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "document-generator", "mark_review", "document_id", args[0])
		},
	}
}

func newDocumentsArchiveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "archive <document_id>",
		Short: "Archive a generated document",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "document-generator", "archive_document", "document_id", args[0])
		},
	}
}

func newDocumentsRestoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore <document_id>",
		Short: "Restore an archived generated document",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "document-generator", "restore_document", "document_id", args[0])
		},
	}
}

func newDocumentsRegenerateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "regenerate <document_id>",
		Short: "Queue document regeneration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "document-generator", "regenerate_document", "document_id", args[0])
		},
	}
}
