package cmd

import (
	"errors"
	"os"
	"testing"

	"bizzmod-cli/internal/config"
)

func TestNewUpdateCmdContainsCheck(t *testing.T) {
	cmd := newUpdateCmd()
	if cmd == nil {
		t.Fatal("expected update command")
	}

	found := false
	for _, child := range cmd.Commands() {
		if child.Name() == "check" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected update check subcommand")
	}
}

func TestParseInstallRequirementsUsesPayloadValues(t *testing.T) {
	payload := []byte(`{
  "data": {
    "install": {
      "required_services": ["docker", "dokploy", "whisper-asr"],
      "required_env": [
        {"name": "WHISPER_SERVICE_URL"},
        {"name": "WHISPER_SERVICE_TIMEOUT"}
      ]
    }
  }
}`)

	services, env := parseInstallRequirements(payload)
	if len(services) != 3 {
		t.Fatalf("expected 3 services, got %d", len(services))
	}
	if len(env) != 2 {
		t.Fatalf("expected 2 env vars, got %d", len(env))
	}
}

func TestRunLocalUpdateChecksMarksMissingEnvAsFail(t *testing.T) {
	originalLookPath := lookPathFn
	originalCommandOutput := commandOutputFn
	defer func() {
		lookPathFn = originalLookPath
		commandOutputFn = originalCommandOutput
	}()

	lookPathFn = func(file string) (string, error) {
		if file == "docker" || file == "dokploy" {
			return "/usr/bin/" + file, nil
		}
		return "", errors.New("not found")
	}
	commandOutputFn = func(name string, args ...string) (string, error) {
		return "bizzmod-back\nbizzmod-front\nwhisper-asr", nil
	}

	_ = os.Unsetenv("WHISPER_SERVICE_URL")
	_ = os.Unsetenv("WHISPER_SERVICE_TIMEOUT")

	rows := runLocalUpdateChecks(
		config.Config{APIURL: ""},
		[]string{"docker", "dokploy", "bizzmod-back", "bizzmod-front", "whisper-asr"},
		[]string{"WHISPER_SERVICE_URL", "WHISPER_SERVICE_TIMEOUT"},
	)

	hasMissingEnv := false
	for _, row := range rows {
		if row.ID == "env:WHISPER_SERVICE_URL" && row.Status == "fail" {
			hasMissingEnv = true
			break
		}
	}

	if !hasMissingEnv {
		t.Fatalf("expected missing WHISPER_SERVICE_URL check failure")
	}
}
