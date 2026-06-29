package processes

import "github.com/spf13/cobra"

func newFormGeneratorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "form-generator",
		Short: "Business verbs for form-generator",
	}
	cmd.AddCommand(newFormsSubmitPublicCmd())
	cmd.AddCommand(newCollectionListCmd("list", "List form-generator forms", "form-generator", "forms"))
	cmd.AddCommand(newCollectionListCmd("list-fields", "List form-generator fields", "form-generator", "fields"))
	cmd.AddCommand(newCollectionListCmd("list-submissions", "List form-generator submissions", "form-generator", "submissions"))
	cmd.AddCommand(newFormsPublishCmd())
	cmd.AddCommand(newFormsArchiveCmd())
	cmd.AddCommand(newFormsRestoreCmd())
	cmd.AddCommand(newFormsMarkSubmissionProcessedCmd())
	cmd.AddCommand(newFormsArchiveSubmissionCmd())
	return cmd
}
