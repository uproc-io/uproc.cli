package processes

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newRequestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request <method> <path> [json_body]",
		Short: "Raw request to /api/v1/external endpoints",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := mustClient()
			if err != nil {
				return err
			}

			method := strings.ToUpper(args[0])
			path := args[1]
			body := []byte{}
			if len(args) == 3 {
				body = []byte(args[2])
			}

			if !strings.HasPrefix(path, "/api/v1/external/") {
				return fmt.Errorf("path must start with /api/v1/external/")
			}

			respBody, status, reqErr := client.Do(method, path, body)
			return printResponse(cmd, respBody, status, reqErr)
		},
	}

	return cmd
}
