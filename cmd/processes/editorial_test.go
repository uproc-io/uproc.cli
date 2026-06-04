package processes

import "testing"

func TestNewEditorialCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newEditorialCmd()
	if cmd == nil {
		t.Fatal("expected editorial command")
	}

	expected := map[string]bool{
		"list-opportunities": false,
		"list-projects":      false,
		"list-articles":      false,
		"list-combined":      false,
		"generate-proposal":  false,
		"generate-article":   false,
		"publish":            false,
		"schedule":           false,
		"discard":            false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected editorial subcommand %s", name)
		}
	}
}
