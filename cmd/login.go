package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"bizzmod-cli/internal/config"
	"github.com/spf13/cobra"
)

func newLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Save API credentials from args or environment",
		Args:  cobra.RangeArgs(0, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.FromEnv()

			if len(args) > 0 {
				if len(args) != 4 {
					return fmt.Errorf("when using arguments, provide all 4 values")
				}
				cfg = config.Config{
					APIURL:         args[0],
					CustomerAPIKey: args[1],
					CustomerDomain: args[2],
					UserEmail:      args[3],
				}
			}

			if len(args) == 0 && missingAny(cfg) {
				interactiveCfg, err := promptMissingFields(cfg, cmd)
				if err != nil {
					return err
				}
				cfg = interactiveCfg
				if err := config.SaveDotEnv(cfg, ".env"); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), "ok: .env created")
			}

			if err := config.Save(cfg); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "ok: credentials saved")
			return nil
		},
	}

	return cmd
}

func missingAny(cfg config.Config) bool {
	return strings.TrimSpace(cfg.APIURL) == "" ||
		strings.TrimSpace(cfg.CustomerAPIKey) == "" ||
		strings.TrimSpace(cfg.CustomerDomain) == "" ||
		strings.TrimSpace(cfg.UserEmail) == ""
}

func promptMissingFields(cfg config.Config, cmd *cobra.Command) (config.Config, error) {
	reader := bufio.NewReader(os.Stdin)

	if strings.TrimSpace(cfg.APIURL) == "" {
		value, err := prompt(reader, cmd, "BIZZMOD_API_URL")
		if err != nil {
			return cfg, err
		}
		cfg.APIURL = value
	}

	if strings.TrimSpace(cfg.CustomerDomain) == "" {
		for {
			value, err := prompt(reader, cmd, "CUSTOMER_DOMAIN")
			if err != nil {
				return cfg, err
			}
			if strings.Contains(value, "://") || strings.Contains(value, "/") {
				fmt.Fprintln(cmd.OutOrStdout(), "CUSTOMER_DOMAIN must be the customer domain value, not a URL")
				continue
			}
			cfg.CustomerDomain = value
			break
		}
	}

	if strings.TrimSpace(cfg.CustomerAPIKey) == "" {
		value, err := prompt(reader, cmd, "CUSTOMER_API_KEY")
		if err != nil {
			return cfg, err
		}
		cfg.CustomerAPIKey = value
	}

	if strings.TrimSpace(cfg.UserEmail) == "" {
		value, err := prompt(reader, cmd, "CUSTOMER_USER_EMAIL")
		if err != nil {
			return cfg, err
		}
		cfg.UserEmail = value
	}

	return cfg, nil
}

func prompt(reader *bufio.Reader, cmd *cobra.Command, field string) (string, error) {
	for {
		fmt.Fprintf(cmd.OutOrStdout(), "%s: ", field)
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		value := strings.TrimSpace(line)
		if value == "" {
			fmt.Fprintln(cmd.OutOrStdout(), "value is required")
			continue
		}
		return value, nil
	}
}
