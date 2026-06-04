package processes

import "testing"

func TestNewSyncCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newSyncCmd()
	if cmd == nil {
		t.Fatal("expected sync command")
	}

	expected := map[string]bool{
		"list-workflows": false,
		"list-runs":      false,
		"list-records":   false,
		"run":            false,
		"preview":        false,
		"dry-run":        false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected sync subcommand %s", name)
		}
	}
}
