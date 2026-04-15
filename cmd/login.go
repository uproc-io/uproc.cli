package cmd

import (
	"bufio"
	"fmt"
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
			currentCfg := config.FromEnv()
			cfg := currentCfg

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
			} else {
				interactiveCfg, err := promptReviewFields(cfg, cmd)
				if err != nil {
					return err
				}
				cfg = interactiveCfg
			}

			changed := hasConfigChanged(currentCfg, cfg)
			if changed {
				if err := validateExternalCredentials(cfg); err != nil {
					return err
				}
			}

			if err := config.SaveDotEnv(cfg, ".env"); err != nil {
				return err
			}

			if err := config.Save(cfg); err != nil {
				return err
			}

			if changed {
				fmt.Fprintln(cmd.OutOrStdout(), "ok: credentials saved and validated")
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "ok: credentials unchanged")
			}
			return nil
		},
	}

	return cmd
}

func promptReviewFields(cfg config.Config, cmd *cobra.Command) (config.Config, error) {
	reader := bufio.NewReader(cmd.InOrStdin())

	apiURL, err := promptWithDefault(reader, cmd, "BIZZMOD_API_URL", cfg.APIURL)
	if err != nil {
		return cfg, err
	}
	cfg.APIURL = apiURL

	for {
		domain, promptErr := promptWithDefault(reader, cmd, "CUSTOMER_DOMAIN", cfg.CustomerDomain)
		if promptErr != nil {
			return cfg, promptErr
		}
		if strings.Contains(domain, "://") || strings.Contains(domain, "/") {
			fmt.Fprintln(cmd.OutOrStdout(), "CUSTOMER_DOMAIN must be the customer domain value, not a URL")
			continue
		}
		cfg.CustomerDomain = domain
		break
	}

	apiKey, err := promptWithDefault(reader, cmd, "CUSTOMER_API_KEY", cfg.CustomerAPIKey)
	if err != nil {
		return cfg, err
	}
	cfg.CustomerAPIKey = apiKey

	userEmail, err := promptWithDefault(reader, cmd, "CUSTOMER_USER_EMAIL", cfg.UserEmail)
	if err != nil {
		return cfg, err
	}
	cfg.UserEmail = userEmail

	return cfg, nil
}

func promptWithDefault(reader *bufio.Reader, cmd *cobra.Command, field string, current string) (string, error) {
	current = strings.TrimSpace(current)
	for {
		if current == "" {
			fmt.Fprintf(cmd.OutOrStdout(), "%s: ", field)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "%s [%s]: ", field, current)
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		value := strings.TrimSpace(line)
		if value == "" {
			if current != "" {
				return current, nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), "value is required")
			continue
		}
		return value, nil
	}
}

func hasConfigChanged(before config.Config, after config.Config) bool {
	return strings.TrimRight(strings.TrimSpace(before.APIURL), "/") != strings.TrimRight(strings.TrimSpace(after.APIURL), "/") ||
		strings.TrimSpace(before.CustomerDomain) != strings.TrimSpace(after.CustomerDomain) ||
		strings.TrimSpace(before.CustomerAPIKey) != strings.TrimSpace(after.CustomerAPIKey) ||
		strings.TrimSpace(before.UserEmail) != strings.TrimSpace(after.UserEmail)
}

func validateExternalCredentials(cfg config.Config) error {
	client, err := mustClientFromConfig(cfg)
	if err != nil {
		return err
	}

	body, status, reqErr := client.Do("GET", "/api/v1/external/modules", nil)
	if reqErr != nil {
		details := strings.TrimSpace(formatStructuredOutput(body))
		if details != "" {
			return fmt.Errorf("credentials validation failed: %w\n%s", reqErr, details)
		}
		return fmt.Errorf("credentials validation failed: %w", reqErr)
	}

	if status < 200 || status >= 300 {
		trimmed := strings.TrimSpace(string(body))
		if trimmed != "" {
			return fmt.Errorf("credentials validation failed with http %d: %s", status, trimmed)
		}
		return fmt.Errorf("credentials validation failed with http %d", status)
	}

	return nil
}
