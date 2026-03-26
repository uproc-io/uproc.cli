package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"bizzmod-cli/internal/api"
	"bizzmod-cli/internal/config"
	"github.com/spf13/cobra"
)

func mustClient() (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if err := config.Validate(cfg); err != nil {
		return nil, err
	}
	return api.NewClient(cfg), nil
}

func printResponse(cmd *cobra.Command, body []byte, status int, err error) error {
	if err != nil {
		if len(body) > 0 {
			pretty := prettyJSON(body)
			if pretty != "" {
				return fmt.Errorf("%w:\n%s", err, pretty)
			}
			return fmt.Errorf("%w: %s", err, string(body))
		}
		return err
	}
	pretty := prettyJSON(body)
	if pretty != "" {
		fmt.Fprintln(cmd.OutOrStdout(), pretty)
		return nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), string(body))
	return nil
}

func prettyJSON(body []byte) string {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return ""
	}

	var parsed any
	if err := json.Unmarshal(trimmed, &parsed); err != nil {
		return ""
	}

	formatted, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return ""
	}

	return string(formatted)
}
