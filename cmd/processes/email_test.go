package processes

import "testing"

func TestNewEmailCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newEmailCmd()
	if cmd == nil {
		t.Fatal("expected email command")
	}

	expected := map[string]bool{
		"list":           false,
		"mark-processed": false,
		"archive":        false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected email subcommand %s", name)
		}
	}
}
