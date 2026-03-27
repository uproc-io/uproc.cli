package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"bizzmod-cli/internal/api"
	"bizzmod-cli/internal/config"
	"github.com/spf13/cobra"
)

func mustClient() (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return mustClientFromConfig(cfg)
}

func mustClientFromConfig(cfg config.Config) (*api.Client, error) {
	if err := config.Validate(cfg); err != nil {
		return nil, err
	}
	return api.NewClient(cfg), nil
}

func printResponse(cmd *cobra.Command, body []byte, status int, err error) error {
	_ = status
	formatted := formatStructuredOutput(body)

	if err != nil {
		if formatted != "" {
			return fmt.Errorf("%w:\n%s", err, formatted)
		}
		trimmed := strings.TrimSpace(string(body))
		if trimmed != "" {
			return fmt.Errorf("%w: %s", err, trimmed)
		}
		return fmt.Errorf("%w", err)
	}

	if formatted != "" {
		fmt.Fprintln(cmd.OutOrStdout(), formatted)
		return nil
	}

	trimmed := strings.TrimSpace(string(body))
	if trimmed != "" {
		fmt.Fprintln(cmd.OutOrStdout(), trimmed)
	}
	return nil
}

func formatStructuredOutput(body []byte) string {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return ""
	}

	var parsed any
	if err := json.Unmarshal(trimmed, &parsed); err != nil {
		return ""
	}

	return renderValue(parsed, 0)
}

func renderValue(v any, indent int) string {
	pad := strings.Repeat("  ", indent)

	switch value := v.(type) {
	case map[string]any:
		return renderMap(value, indent)
	case []any:
		if len(value) == 0 {
			return pad + "(empty)"
		}
		if isListOfObjects(value) {
			return renderObjectTable(value, indent)
		}
		lines := make([]string, 0, len(value))
		for _, item := range value {
			if isScalar(item) {
				lines = append(lines, fmt.Sprintf("%s- %v", pad, item))
				continue
			}
			rendered := renderValue(item, indent+1)
			lines = append(lines, fmt.Sprintf("%s-", pad))
			lines = append(lines, rendered)
		}
		return strings.Join(lines, "\n")
	default:
		return fmt.Sprintf("%s%v", pad, value)
	}
}

func renderMap(m map[string]any, indent int) string {
	if len(m) == 0 {
		return strings.Repeat("  ", indent) + "(empty)"
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pad := strings.Repeat("  ", indent)
	lines := make([]string, 0, len(keys)*2)
	for _, key := range keys {
		value := m[key]
		if isScalar(value) {
			lines = append(lines, fmt.Sprintf("%s- %s: %v", pad, key, value))
			continue
		}
		lines = append(lines, fmt.Sprintf("%s- %s:", pad, key))
		lines = append(lines, renderValue(value, indent+1))
	}

	return strings.Join(lines, "\n")
}

func renderObjectTable(items []any, indent int) string {
	columns := tableColumns(items)
	if len(columns) == 0 {
		return strings.Repeat("  ", indent) + "(empty)"
	}

	var b strings.Builder
	tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, strings.Join(columns, "\t"))
	for _, item := range items {
		rowMap, _ := item.(map[string]any)
		row := make([]string, 0, len(columns))
		for _, column := range columns {
			row = append(row, tableCell(rowMap[column]))
		}
		_, _ = fmt.Fprintln(tw, strings.Join(row, "\t"))
	}
	_ = tw.Flush()

	rendered := strings.TrimRight(b.String(), "\n")
	if indent == 0 {
		return rendered
	}

	pad := strings.Repeat("  ", indent)
	lines := strings.Split(rendered, "\n")
	for i := range lines {
		lines[i] = pad + lines[i]
	}
	return strings.Join(lines, "\n")
}

func tableColumns(items []any) []string {
	seen := make(map[string]struct{})
	columns := make([]string, 0)

	for _, item := range items {
		rowMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rowKeys := make([]string, 0, len(rowMap))
		for k := range rowMap {
			rowKeys = append(rowKeys, k)
		}
		sort.Strings(rowKeys)
		for _, key := range rowKeys {
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			columns = append(columns, key)
		}
	}

	return columns
}

func tableCell(v any) string {
	if v == nil {
		return ""
	}
	if isScalar(v) {
		return fmt.Sprintf("%v", v)
	}
	return "..."
}

func isListOfObjects(items []any) bool {
	if len(items) == 0 {
		return false
	}
	for _, item := range items {
		if _, ok := item.(map[string]any); !ok {
			return false
		}
	}
	return true
}

func isScalar(v any) bool {
	switch v.(type) {
	case nil, string, bool, float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}
