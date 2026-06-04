package processes

import "testing"

func TestNewTaxCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newTaxCmd()
	if cmd == nil {
		t.Fatal("expected tax command")
	}

	expected := map[string]bool{
		"list":        false,
		"generate":    false,
		"recalculate": false,
		"validate":    false,
		"export":      false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected tax subcommand %s", name)
		}
	}
}
