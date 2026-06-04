package processes

import "testing"

func TestNewProspectingCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newProspectingCmd()
	if cmd == nil {
		t.Fatal("expected prospecting command")
	}

	expected := map[string]bool{
		"list-strategies":    false,
		"list-opportunities": false,
		"list-prospects":     false,
		"list-executions":    false,
		"run-discovery":      false,
		"send-to-leads":      false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected prospecting subcommand %s", name)
		}
	}
}
