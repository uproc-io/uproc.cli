package processes

import "testing"

func TestNewDocumentsCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newDocumentsCmd()
	if cmd == nil {
		t.Fatal("expected documents command")
	}

	expected := map[string]bool{
		"list":        false,
		"mark-ready":  false,
		"mark-review": false,
		"archive":     false,
		"restore":     false,
		"regenerate":  false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected documents subcommand %s", name)
		}
	}
}
