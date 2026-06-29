package processes

import "github.com/spf13/cobra"

func newCaseLifecycleCmd() *cobra.Command {
	casesCmd := &cobra.Command{
		Use:   "case-lifecycle",
		Short: "Business verbs for case-lifecycle",
	}

	casesCmd.AddCommand(newCasesAddNoteCmd())
	casesCmd.AddCommand(newCollectionListCmd("list", "List cases", "case-lifecycle", "cases"))
	casesCmd.AddCommand(newCollectionListCmd("list-by-status", "List cases grouped by status", "case-lifecycle", "cases_by_status"))
	casesCmd.AddCommand(newCollectionListCmd("list-by-type", "List cases grouped by type", "case-lifecycle", "cases_by_type"))
	casesCmd.AddCommand(newCasesCloseCmd())
	casesCmd.AddCommand(newCasesReopenCmd())

	return casesCmd
}

func newCasesCmd() *cobra.Command {
	cmd := newCaseLifecycleCmd()
	cmd.Use = "cases"
	cmd.Hidden = true
	cmd.Long = "DEPRECATED: Use \"case-lifecycle\" instead.\n\n" + cmd.Long
	return cmd
}

func newCasesAddNoteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add-note <case_id> <content> [created_by]",
		Short: "Add an internal note to a case",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			caseID, err := parsePositiveIntArg("case_id", args[0])
			if err != nil {
				return err
			}
			payload := map[string]any{
				"case_id": caseID,
				"content": args[1],
			}
			if len(args) == 3 {
				payload["created_by"] = args[2]
			}
			return runModuleAction(cmd, "case-lifecycle", "create_case_note", payload)
		},
	}
}

func newCasesCloseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "close <case_id>",
		Short: "Close a case",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "case-lifecycle", "close_case", "case_id", args[0])
		},
	}
}

func newCasesReopenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reopen <case_id>",
		Short: "Reopen a case",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "case-lifecycle", "reopen_case", "case_id", args[0])
		},
	}
}
