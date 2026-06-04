package processes

import "testing"

func TestNewInvoiceLinesCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newInvoiceLinesCmd()
	if cmd == nil {
		t.Fatal("expected invoice-lines command")
	}

	expected := map[string]bool{
		"list":   false,
		"add":    false,
		"update": false,
		"delete": false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected invoice-lines subcommand %s", name)
		}
	}
}
