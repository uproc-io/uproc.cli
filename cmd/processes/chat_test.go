package processes

import "testing"

func TestNewChatCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newDataChatbotCmd()
	if cmd == nil {
		t.Fatal("expected chat command")
	}

	expected := map[string]bool{
		"list": false,
		"ask":  false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected chat subcommand %s", name)
		}
	}
}
