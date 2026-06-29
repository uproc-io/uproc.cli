package processes

import "github.com/spf13/cobra"

func newContractLifecycleCmd() *cobra.Command {
	contractCmd := &cobra.Command{
		Use:   "contract-lifecycle",
		Short: "Business verbs for contract-lifecycle",
	}

	contractCmd.AddCommand(newContractRenewCmd())
	contractCmd.AddCommand(newCollectionListCmd("list", "List contract-lifecycle contracts", "contract-lifecycle", "contracts"))
	contractCmd.AddCommand(newCollectionListCmd("list-expiring", "List expiring contracts", "contract-lifecycle", "expiring_contracts"))
	contractCmd.AddCommand(newCollectionListCmd("list-by-counterparty", "List contracts by counterparty", "contract-lifecycle", "contracts_by_counterparty"))
	contractCmd.AddCommand(newContractTerminateCmd())
	contractCmd.AddCommand(newContractUpdateCmd())

	return contractCmd
}

func newContractCmd() *cobra.Command {
	cmd := newContractLifecycleCmd()
	cmd.Use = "contract"
	cmd.Hidden = true
	cmd.Long = "DEPRECATED: Use \"contract-lifecycle\" instead.\n\n" + cmd.Long
	return cmd
}

func newContractRenewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "renew <contract_id>",
		Short: "Renew a contract lifecycle contract",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "contract-lifecycle", "renew", "contract_id", args[0])
		},
	}
}

func newContractTerminateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "terminate <contract_id>",
		Short: "Terminate a contract lifecycle contract",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "contract-lifecycle", "terminate", "contract_id", args[0])
		},
	}
}

func newContractUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update <contract_id> <data_json>",
		Short: "Update contract fields for an existing contract",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			contractID, err := parsePositiveIntArg("contract_id", args[0])
			if err != nil {
				return err
			}
			data, err := parseJSONObjectArg("data_json", args[1])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "contract-lifecycle", "update_contract", map[string]any{
				"contract_id": contractID,
				"data":        data,
			})
		},
	}
}
