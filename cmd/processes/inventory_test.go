package processes

import "testing"

func TestNewInventoryCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newInventoryCmd()
	if cmd == nil {
		t.Fatal("expected inventory command")
	}

	expected := map[string]bool{
		"list":          false,
		"mark-received": false,
		"cancel":        false,
		"send-reminder": false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected inventory subcommand %s", name)
		}
	}
}
