package processes

import "testing"

func TestNewSigningCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newDocumentSigningCmd()
	if cmd == nil {
		t.Fatal("expected signing command")
	}

	expected := map[string]bool{
		"list":          false,
		"cancel":        false,
		"reopen":        false,
		"send-reminder": false,
		"sync-status":   false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected signing subcommand %s", name)
		}
	}
}
