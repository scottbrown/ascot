package ascot

import (
	"testing"
)

type MockDescribeRegionsClient struct {
}

func TestActiveRegionsRunnerRun(t *testing.T) {
}

func TestActiveRegionsRunnerRequiredPermissions(t *testing.T) {
	var runner ActiveRegionsRunner

	expectedAtLeast := 1
	privs := runner.RequiredPermissions()
	if len(privs) < expectedAtLeast {
		t.Fatalf("Expected %d but got %d", expectedAtLeast, len(privs))
	}
}

func TestActiveRegionsRunnerHowItWorks(t *testing.T) {
	var runner ActiveRegionsRunner

	expectedCallsAtLeast := 1
	expectedNotesAtLeast := 0
	calls, notes := runner.HowItWorks()
	if len(calls) < expectedCallsAtLeast {
		t.Fatalf("Expected at least %d but got %d", expectedCallsAtLeast, len(calls))
	}
	if len(notes) < expectedNotesAtLeast {
		t.Fatalf("Expected at least %d but got %d", expectedNotesAtLeast, len(notes))
	}
}
