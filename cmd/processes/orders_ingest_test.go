package processes

import "testing"

func TestNewOrdersIngestCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newOrdersIngestCmd()
	if cmd == nil {
		t.Fatal("expected orders-ingest command")
	}

	expected := map[string]bool{
		"list":        false,
		"list-emails": false,
		"reprocess":   false,
		"validate":    false,
		"send-to-erp": false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected orders-ingest subcommand %s", name)
		}
	}
}
