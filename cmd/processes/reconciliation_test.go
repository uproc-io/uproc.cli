package processes

import "testing"

func TestNewReconciliationCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newReconciliationCmd()
	if cmd == nil {
		t.Fatal("expected reconciliation command")
	}

	expected := map[string]bool{
		"list-entries":  false,
		"list-extracts": false,
		"list-exports":  false,
		"list-matches":  false,
		"reconcile":     false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected reconciliation subcommand %s", name)
		}
	}
}
