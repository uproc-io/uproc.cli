package processes

import "github.com/spf13/cobra"

func newInvoiceIngestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoice-ingest",
		Short: "Business verbs for invoice-ingest",
	}
	cmd.AddCommand(newCollectionListCmd("list-invoices", "List invoice-ingest invoices", "invoice-ingest", "invoices"))
	cmd.AddCommand(newCollectionListCmd("list-emails", "List invoice-ingest emails", "invoice-ingest", "emails"))
	return cmd
}
