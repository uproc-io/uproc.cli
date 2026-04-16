package processes

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func runAdminListWithContract(
	cmd *cobra.Command,
	resource string,
	method string,
	path string,
) error {
	client, err := mustClient()
	if err != nil {
		return err
	}

	body, status, reqErr := client.Do(method, path, nil)
	if reqErr != nil {
		return printResponse(cmd, body, status, reqErr)
	}

	contract, contractErr := fetchAdminContract(client, resource, "list")
	if contractErr != nil || len(contract.VisibleFields) == 0 {
		return printResponse(cmd, body, status, nil)
	}

	filteredBody, filterErr := applyVisibleFieldsFromContract(body, contract.VisibleFields)
	if filterErr != nil {
		return printResponse(cmd, body, status, nil)
	}

	return printResponse(cmd, filteredBody, status, nil)
}

func applyVisibleFieldsFromContract(body []byte, visibleFields []string) ([]byte, error) {
	if len(visibleFields) == 0 {
		return body, nil
	}

	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, err
	}

	dataObj, ok := root["data"].(map[string]any)
	if !ok {
		return body, nil
	}

	rawRows, ok := dataObj["rows"].([]any)
	if !ok {
		return body, nil
	}

	filteredRows := make([]any, 0, len(rawRows))
	for _, item := range rawRows {
		rowMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		ordered := map[string]any{}
		for _, field := range visibleFields {
			if value, exists := rowMap[field]; exists {
				ordered[field] = value
			}
		}
		filteredRows = append(filteredRows, ordered)
	}

	dataObj["rows"] = filteredRows
	dataObj["columns"] = buildColumnsFromVisibleFields(visibleFields)

	root["data"] = dataObj
	encoded, err := json.Marshal(root)
	if err != nil {
		return nil, fmt.Errorf("failed to encode filtered payload: %w", err)
	}
	return encoded, nil
}

func buildColumnsFromVisibleFields(visibleFields []string) []any {
	columns := make([]any, 0, len(visibleFields))
	for _, field := range visibleFields {
		columns = append(columns, map[string]any{"key": field, "label": field})
	}
	return columns
}
