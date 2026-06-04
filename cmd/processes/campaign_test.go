package processes

import "testing"

func TestNewCampaignCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newCampaignCmd()
	if cmd == nil {
		t.Fatal("expected campaign command")
	}

	expected := map[string]bool{
		"list":             false,
		"list-audiences":   false,
		"preview-audience": false,
		"add-audience":     false,
		"pause":            false,
		"activate":         false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected campaign subcommand %s", name)
		}
	}
}
