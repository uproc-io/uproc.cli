package processes

import "testing"

func TestNewCasesCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newCasesCmd()
	if cmd == nil {
		t.Fatal("expected cases command")
	}

	expected := map[string]bool{
		"list":           false,
		"list-by-status": false,
		"list-by-type":   false,
		"add-note":       false,
		"close":          false,
		"reopen":         false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected cases subcommand %s", name)
		}
	}
}
