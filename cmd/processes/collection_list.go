package processes

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func newCollectionListCmd(use, short, moduleSlug, collectionName string) *cobra.Command {
	var page int
	var sortField string
	var sortOrder string
	var filterField string
	var filterValue string

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
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
				moduleSlug,
				collectionName,
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
