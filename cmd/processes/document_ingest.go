package processes

import "github.com/spf13/cobra"

func newDocumentIngestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "document-ingest",
		Short: "Business verbs for document-ingest",
	}
	cmd.AddCommand(newCollectionListCmd("list-documents", "List document-ingest documents", "document-ingest", "documents"))
	cmd.AddCommand(newCollectionListCmd("list-uploads", "List document-ingest uploads", "document-ingest", "uploads"))
	cmd.AddCommand(newCollectionListCmd("list-emails", "List document-ingest emails", "document-ingest", "emails"))
	return cmd
}
