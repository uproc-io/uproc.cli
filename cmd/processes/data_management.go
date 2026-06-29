package processes

import "github.com/spf13/cobra"

func newDataManagementCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data-management",
		Short: "Business verbs for data-management",
	}
	cmd.AddCommand(newDMListTablesCmd())
	cmd.AddCommand(newDMCreateTableCmd())
	cmd.AddCommand(newDMListColumnsCmd())
	cmd.AddCommand(newDMEnqueueImportCmd())
	cmd.AddCommand(newDMGetUploadStatusCmd())
	cmd.AddCommand(newDMGetLatestUploadCmd())
	cmd.AddCommand(newDMAnalyzeNormalizationCmd())
	cmd.AddCommand(newDMExecuteNormalizationCmd())
	cmd.AddCommand(newDMListPlansCmd())
	cmd.AddCommand(newDMGetPlanCmd())
	cmd.AddCommand(newDMRegeneratePlanCmd())
	cmd.AddCommand(newDMGetRowMappingsCmd())
	cmd.AddCommand(newDMSmartImportPreviewCmd())
	cmd.AddCommand(newDMSmartImportExecuteCmd())
	return cmd
}

func newDMListTablesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-tables",
		Short: "List data-management tables",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runModuleAction(cmd, "data-management", "list_tables", map[string]any{})
		},
	}
}

func newDMCreateTableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-table <item>",
		Short: "Create a data-management table",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			payload, err := parseJSONObjectArg("item", args[0])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "data-management", "create_table", payload)
		},
	}
}

func newDMListColumnsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-columns <entity_id>",
		Short: "List columns for a data-management table",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "data-management", "list_columns", "entity_id", args[0])
		},
	}
}

func newDMEnqueueImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enqueue-import <collection_name> <file_name> <content>",
		Short: "Queue a file import into a data-management table",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runModuleAction(cmd, "data-management", "enqueue_import", map[string]any{
				"collection_name": args[0],
				"file_name":       args[1],
				"content":         args[2],
			})
		},
	}
}

func newDMGetUploadStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-upload-status <upload_id>",
		Short: "Check the status of a queued data-management import",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "data-management", "get_upload_status", "upload_id", args[0])
		},
	}
}

func newDMGetLatestUploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-latest-upload",
		Short: "Get the latest data-management import status",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runModuleAction(cmd, "data-management", "get_latest_upload", map[string]any{})
		},
	}
}

func newDMAnalyzeNormalizationCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "analyze-normalization <entity_id> [entity_id...]",
		Short: "Analyze entities and generate AI normalization plan",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityIDs := make([]int, 0, len(args))
			for _, raw := range args {
				id, err := parsePositiveIntArg("entity_id", raw)
				if err != nil {
					return err
				}
				entityIDs = append(entityIDs, id)
			}
			return runModuleAction(cmd, "data-management", "analyze_normalization", map[string]any{
				"entity_ids": entityIDs,
			})
		},
	}
}

func newDMExecuteNormalizationCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "execute-normalization <plan_id>",
		Short: "Execute an accepted normalization plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "data-management", "execute_normalization", "plan_id", args[0])
		},
	}
}

func newDMListPlansCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-plans",
		Short: "List normalization plans",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runModuleAction(cmd, "data-management", "list_plans", map[string]any{})
		},
	}
}

func newDMGetPlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-plan <plan_id>",
		Short: "Get a normalization plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "data-management", "get_plan", "plan_id", args[0])
		},
	}
}

func newDMRegeneratePlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "regenerate-plan <plan_id>",
		Short: "Regenerate a normalization plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "data-management", "regenerate_plan", "plan_id", args[0])
		},
	}
}

func newDMGetRowMappingsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-row-mappings <plan_id>",
		Short: "Get row-level change mappings for a normalization plan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "data-management", "get_row_mappings", "plan_id", args[0])
		},
	}
}

func newDMSmartImportPreviewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "smart-import-preview <raw_entity_id>",
		Short: "Preview which normalized entity a raw upload matches",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSingleIDAction(cmd, "data-management", "smart_import_preview", "raw_entity_id", args[0])
		},
	}
}

func newDMSmartImportExecuteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "smart-import-execute <raw_entity_id> <target_entity_id>",
		Short: "Execute smart import from raw upload to normalized entity",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			rawID, err := parsePositiveIntArg("raw_entity_id", args[0])
			if err != nil {
				return err
			}
			targetID, err := parsePositiveIntArg("target_entity_id", args[1])
			if err != nil {
				return err
			}
			return runModuleAction(cmd, "data-management", "smart_import_execute", map[string]any{
				"raw_entity_id":    rawID,
				"target_entity_id": targetID,
			})
		},
	}
}
