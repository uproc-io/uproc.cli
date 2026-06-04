package processes

import "testing"

func TestNewSignalsCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newSignalsCmd()
	if cmd == nil {
		t.Fatal("expected signals command")
	}

	expected := map[string]bool{
		"list":                false,
		"list-executions":     false,
		"list-activations":    false,
		"approve":             false,
		"discard":             false,
		"mark-pending-review": false,
		"activate":            false,
		"close":               false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected signals subcommand %s", name)
		}
	}
}
