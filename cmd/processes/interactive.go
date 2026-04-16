package processes

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newInteractiveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "interactive",
		Short: "Interactive shell for processes commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			scanner := bufio.NewScanner(cmd.InOrStdin())
			for {
				fmt.Fprint(cmd.OutOrStdout(), "uproc> ")
				if !scanner.Scan() {
					break
				}

				line := strings.TrimSpace(scanner.Text())
				if line == "" {
					continue
				}

				switch strings.ToLower(line) {
				case "exit", "quit":
					return nil
				case "help":
					fmt.Fprintln(cmd.OutOrStdout(), "Type any processes command without the binary name.")
					fmt.Fprintln(cmd.OutOrStdout(), "Examples: module list, module get order-track, request GET /api/v1/external/modules")
					fmt.Fprintln(cmd.OutOrStdout(), "Use 'exit' or 'quit' to leave interactive mode.")
					continue
				}

				args, parseErr := parseInteractiveArgs(line)
				if parseErr != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "parse error: %v\n", parseErr)
					continue
				}

				runner := NewCmd()
				runner.SetOut(cmd.OutOrStdout())
				runner.SetErr(cmd.ErrOrStderr())
				runner.SetIn(cmd.InOrStdin())
				runner.SetArgs(args)
				if execErr := runner.Execute(); execErr != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), execErr)
				}
			}

			if err := scanner.Err(); err != nil {
				return err
			}

			return nil
		},
	}
}

func parseInteractiveArgs(input string) ([]string, error) {
	args := make([]string, 0)
	var current strings.Builder
	inQuote := false
	quoteChar := byte(0)
	escaped := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if escaped {
			current.WriteByte(ch)
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			continue
		}

		if inQuote {
			if ch == quoteChar {
				inQuote = false
				continue
			}
			current.WriteByte(ch)
			continue
		}

		if ch == '"' || ch == '\'' {
			inQuote = true
			quoteChar = ch
			continue
		}

		if ch == ' ' || ch == '\t' {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteByte(ch)
	}

	if escaped {
		current.WriteByte('\\')
	}

	if inQuote {
		return nil, fmt.Errorf("unterminated quote")
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args, nil
}
