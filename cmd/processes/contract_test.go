package processes

import "testing"

func TestNewContractCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newContractCmd()
	if cmd == nil {
		t.Fatal("expected contract command")
	}

	expected := map[string]bool{
		"list":                 false,
		"list-expiring":        false,
		"list-by-counterparty": false,
		"renew":                false,
		"terminate":            false,
		"update":               false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected contract subcommand %s", name)
		}
	}
}
