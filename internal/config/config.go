package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	APIURL         string `json:"api_url"`
	CustomerAPIKey string `json:"customer_api_key"`
	CustomerDomain string `json:"customer_domain"`
	UserEmail      string `json:"user_email"`
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	baseDir := filepath.Join(dir, "bizzmod-cli")
	if err := os.MkdirAll(baseDir, 0o700); err != nil {
		return "", err
	}

	return filepath.Join(baseDir, "config.json"), nil
}

func Load() (Config, error) {
	loadDotEnvIfPresent(".env")
	cfg := FromEnv()

	cp, err := configPath()
	if err != nil {
		return cfg, err
	}

	b, err := os.ReadFile(cp)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	var saved Config
	if err := json.Unmarshal(b, &saved); err != nil {
		return cfg, nil
	}

	if cfg.APIURL == "" {
		cfg.APIURL = strings.TrimRight(strings.TrimSpace(saved.APIURL), "/")
	}
	if cfg.CustomerAPIKey == "" {
		cfg.CustomerAPIKey = strings.TrimSpace(saved.CustomerAPIKey)
	}
	if cfg.CustomerDomain == "" {
		cfg.CustomerDomain = strings.TrimSpace(saved.CustomerDomain)
	}
	if cfg.UserEmail == "" {
		cfg.UserEmail = strings.TrimSpace(saved.UserEmail)
	}

	return cfg, nil
}

func FromEnv() Config {
	loadDotEnvIfPresent(".env")

	return Config{
		APIURL:         strings.TrimRight(strings.TrimSpace(os.Getenv("BIZZMOD_API_URL")), "/"),
		CustomerAPIKey: strings.TrimSpace(os.Getenv("CUSTOMER_API_KEY")),
		CustomerDomain: strings.TrimSpace(os.Getenv("CUSTOMER_DOMAIN")),
		UserEmail:      strings.TrimSpace(os.Getenv("CUSTOMER_USER_EMAIL")),
	}
}

func SaveDotEnv(cfg Config, path string) error {
	if strings.TrimSpace(path) == "" {
		path = ".env"
	}

	lines := []string{
		"BIZZMOD_API_URL=" + strings.TrimRight(strings.TrimSpace(cfg.APIURL), "/"),
		"CUSTOMER_DOMAIN=" + strings.TrimSpace(cfg.CustomerDomain),
		"CUSTOMER_API_KEY=" + strings.TrimSpace(cfg.CustomerAPIKey),
		"CUSTOMER_USER_EMAIL=" + strings.TrimSpace(cfg.UserEmail),
	}

	content := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(path, []byte(content), 0o600)
}

func loadDotEnvIfPresent(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		entry := strings.TrimSpace(line)
		if entry == "" || strings.HasPrefix(entry, "#") {
			continue
		}

		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		if key == "" || os.Getenv(key) != "" {
			continue
		}

		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"")
		value = strings.Trim(value, "'")
		_ = os.Setenv(key, value)
	}
}

func Save(cfg Config) error {
	if strings.TrimSpace(cfg.APIURL) == "" || strings.TrimSpace(cfg.CustomerAPIKey) == "" ||
		strings.TrimSpace(cfg.CustomerDomain) == "" || strings.TrimSpace(cfg.UserEmail) == "" {
		return fmt.Errorf("all login fields are required")
	}

	cfg.APIURL = strings.TrimRight(strings.TrimSpace(cfg.APIURL), "/")
	cfg.CustomerAPIKey = strings.TrimSpace(cfg.CustomerAPIKey)
	cfg.CustomerDomain = strings.TrimSpace(cfg.CustomerDomain)
	cfg.UserEmail = strings.TrimSpace(cfg.UserEmail)

	cp, err := configPath()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cp, b, 0o600)
}

func Validate(cfg Config) error {
	if strings.TrimSpace(cfg.APIURL) == "" {
		return fmt.Errorf("missing BIZZMOD_API_URL or login config")
	}
	if strings.TrimSpace(cfg.CustomerAPIKey) == "" {
		return fmt.Errorf("missing CUSTOMER_API_KEY or login config")
	}
	if strings.TrimSpace(cfg.CustomerDomain) == "" {
		return fmt.Errorf("missing CUSTOMER_DOMAIN or login config")
	}

	if strings.Contains(cfg.CustomerDomain, "://") || strings.Contains(cfg.CustomerDomain, "/") {
		return fmt.Errorf("CUSTOMER_DOMAIN must be the customer domain value, not a URL")
	}

	if strings.TrimSpace(cfg.UserEmail) == "" {
		return fmt.Errorf("missing CUSTOMER_USER_EMAIL or login config")
	}
	return nil
}
