package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newModuleCmd() *cobra.Command {
	moduleCmd := &cobra.Command{
		Use:   "module",
		Short: "Module commands for external API",
	}

	moduleCmd.AddCommand(newModuleListCmd())
	moduleCmd.AddCommand(newModuleGetCmd())
	moduleCmd.AddCommand(newModuleOverviewCmd())
	moduleCmd.AddCommand(newModuleCollectionsCmd())
	moduleCmd.AddCommand(newModuleCollectionCmd())
	moduleCmd.AddCommand(newModuleDataCmd())
	moduleCmd.AddCommand(newModuleUploadCmd())
	moduleCmd.AddCommand(newModuleWebhookCmd())

	return moduleCmd
}

func newModuleListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}
			body, status, reqErr := client.Do("GET", "/api/v1/external/modules", nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newModuleGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <module_slug>",
		Short: "Get module details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/api/v1/external/modules/%s", args[0])
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newModuleOverviewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "overview <module_slug> [section]",
		Short: "Get module overview (kpis, charts, tables)",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			section := "all"
			if len(args) == 2 {
				section = args[1]
			}
			if !isValidOverviewSection(section) {
				return fmt.Errorf("invalid section %q. allowed: kpis, charts, tables", section)
			}

			client, err := mustClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/api/v1/external/modules/%s/kpis", args[0])
			body, status, reqErr := client.Do("GET", path, nil)
			if reqErr != nil {
				return printResponse(cmd, body, status, reqErr)
			}

			rendered, renderErr := formatModuleOverviewOutput(body, args[0], section)
			if renderErr != nil {
				return printResponse(cmd, body, status, nil)
			}

			if rendered == "" {
				return printResponse(cmd, body, status, nil)
			}

			fmt.Fprintln(cmd.OutOrStdout(), rendered)
			return nil
		},
	}
}

func newModuleCollectionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "collections <module_slug>",
		Short: "Get module collections",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}
			path := fmt.Sprintf("/api/v1/external/modules/%s/collections", args[0])
			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newModuleCollectionCmd() *cobra.Command {
	var page int
	var sortField string
	var sortOrder string
	var filterField string
	var filterValue string

	cmd := &cobra.Command{
		Use:   "collection <module_slug> <collection_name>",
		Short: "Get module collection table",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			query := url.Values{}
			query.Set("page", fmt.Sprintf("%d", page))
			if sortField != "" {
				query.Set("sort_field", sortField)
			}
			if sortOrder != "" {
				query.Set("sort_order", sortOrder)
			}
			if filterField != "" {
				query.Set("filter_field", filterField)
			}
			if filterValue != "" {
				query.Set("filter_value", filterValue)
			}

			path := fmt.Sprintf(
				"/api/v1/external/modules/%s/collections/%s?%s",
				args[0],
				args[1],
				query.Encode(),
			)

			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}

	cmd.Flags().IntVar(&page, "page", 1, "page number")
	cmd.Flags().StringVar(&sortField, "sort-field", "", "sort field")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort order: asc|desc")
	cmd.Flags().StringVar(&filterField, "filter-field", "", "filter field")
	cmd.Flags().StringVar(&filterValue, "filter-value", "", "filter value")

	return cmd
}

func newModuleDataCmd() *cobra.Command {
	var page int
	var sortField string
	var sortOrder string
	var filterField string
	var filterValue string

	cmd := &cobra.Command{
		Use:   "data <module_slug> <collection_name>",
		Short: "Get module data from /data endpoint",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			query := url.Values{}
			query.Set("page", fmt.Sprintf("%d", page))
			if sortField != "" {
				query.Set("sort_field", sortField)
			}
			if sortOrder != "" {
				query.Set("sort_order", sortOrder)
			}
			if filterField != "" {
				query.Set("filter_field", filterField)
			}
			if filterValue != "" {
				query.Set("filter_value", filterValue)
			}

			path := fmt.Sprintf(
				"/api/v1/external/modules/%s/data/%s?%s",
				args[0],
				args[1],
				query.Encode(),
			)

			body, status, reqErr := client.Do("GET", path, nil)
			return printResponse(cmd, body, status, reqErr)
		},
	}

	cmd.Flags().IntVar(&page, "page", 1, "page number")
	cmd.Flags().StringVar(&sortField, "sort-field", "", "sort field")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort order: asc|desc")
	cmd.Flags().StringVar(&filterField, "filter-field", "", "filter field")
	cmd.Flags().StringVar(&filterValue, "filter-value", "", "filter value")

	return cmd
}

func newModuleUploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload <module_slug> <collection_name> <file_path>",
		Short: "Upload file to a module collection",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			content, err := os.ReadFile(args[2])
			if err != nil {
				return err
			}

			payload := map[string]string{
				"file_name": filepath.Base(args[2]),
				"content":   base64.StdEncoding.EncodeToString(content),
			}

			b, err := json.Marshal(payload)
			if err != nil {
				return err
			}

			path := fmt.Sprintf(
				"/api/v1/external/modules/%s/collections/%s/inputs/file",
				args[0],
				args[1],
			)
			body, status, reqErr := client.Do("POST", path, b)
			return printResponse(cmd, body, status, reqErr)
		},
	}
}

func newModuleWebhookCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "webhook <module_slug> <collection_name> <payload_json>",
		Short: "Send webhook payload to a module collection",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf(
				"/api/v1/external/modules/%s/collections/%s/inputs/webhook",
				args[0],
				args[1],
			)
			body, status, reqErr := client.Do("POST", path, []byte(args[2]))
			return printResponse(cmd, body, status, reqErr)
		},
	}
}
