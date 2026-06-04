package processes

import "testing"

func TestNewLeadsCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newLeadsCmd()
	if cmd == nil {
		t.Fatal("expected leads command")
	}

	expected := map[string]bool{
		"list":               false,
		"generate-proposal":  false,
		"send-proposal":      false,
		"rerun-intelligence": false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected leads subcommand %s", name)
		}
	}
}
