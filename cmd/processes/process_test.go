package processes

import "testing"

func TestNewProcessCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newProcessCmd()
	if cmd == nil {
		t.Fatal("expected process command")
	}

	expected := map[string]bool{
		"list":           false,
		"retry-step":     false,
		"reassign-owner": false,
		"cancel":         false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected process subcommand %s", name)
		}
	}
}
