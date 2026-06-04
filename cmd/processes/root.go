package processes

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	processesCmd := &cobra.Command{
		Use:   "processes",
		Short: "Commands for Uproc Processes API",
	}

	processesCmd.AddCommand(newLoginCmd())
	processesCmd.AddCommand(newRequestCmd())
	processesCmd.AddCommand(newModuleCmd())
	processesCmd.AddCommand(newChatCmd())
	processesCmd.AddCommand(newFormsCmd())
	processesCmd.AddCommand(newCandidateCmd())
	processesCmd.AddCommand(newSupportCmd())
	processesCmd.AddCommand(newApprovalCmd())
	processesCmd.AddCommand(newCampaignCmd())
	processesCmd.AddCommand(newContractCmd())
	processesCmd.AddCommand(newOrderCmd())
	processesCmd.AddCommand(newEmailCmd())
	processesCmd.AddCommand(newProcessCmd())
	processesCmd.AddCommand(newSignalsCmd())
	processesCmd.AddCommand(newEditorialCmd())
	processesCmd.AddCommand(newSigningCmd())
	processesCmd.AddCommand(newTaxCmd())
	processesCmd.AddCommand(newDocumentsCmd())
	processesCmd.AddCommand(newInventoryCmd())
	processesCmd.AddCommand(newOrdersIngestCmd())
	processesCmd.AddCommand(newCasesCmd())
	processesCmd.AddCommand(newInvoiceCmd())
	processesCmd.AddCommand(newInvoiceLinesCmd())
	processesCmd.AddCommand(newSyncCmd())
	processesCmd.AddCommand(newLeadsCmd())
	processesCmd.AddCommand(newProspectingCmd())
	processesCmd.AddCommand(newReconciliationCmd())
	processesCmd.AddCommand(newAdminCmd())
	processesCmd.AddCommand(newInstallCmd())
	processesCmd.AddCommand(newUpdateCmd())
	processesCmd.AddCommand(newInteractiveCmd())

	return processesCmd
}
