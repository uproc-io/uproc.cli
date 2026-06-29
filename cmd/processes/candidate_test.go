package processes

import "testing"

func TestNewCandidateCmdContainsExpectedVerbs(t *testing.T) {
	cmd := newCandidateEvaluationCmd()
	if cmd == nil {
		t.Fatal("expected candidate command")
	}

	expected := map[string]bool{
		"list-profiles":      false,
		"list-job-openings":  false,
		"list-applications":  false,
		"list-evaluations":   false,
		"list-stage-events":  false,
		"create-profile":     false,
		"create-job-opening": false,
		"create-application": false,
		"move-stage":         false,
		"update-status":      false,
		"create-evaluation":  false,
	}

	for _, child := range cmd.Commands() {
		if _, ok := expected[child.Name()]; ok {
			expected[child.Name()] = true
		}
	}

	for name, found := range expected {
		if !found {
			t.Fatalf("expected candidate subcommand %s", name)
		}
	}
}
