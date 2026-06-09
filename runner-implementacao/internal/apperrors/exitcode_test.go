package apperrors

import "testing"

func TestExitCodesAreStable(t *testing.T) {
	cases := map[string]int{
		"OK":               OK,
		"UsageError":       UsageError,
		"DependencyError":  DependencyError,
		"IntegrationError": IntegrationError,
		"InternalError":    InternalError,
	}

	expected := map[string]int{
		"OK":               0,
		"UsageError":       1,
		"DependencyError":  2,
		"IntegrationError": 3,
		"InternalError":    4,
	}

	for name, got := range cases {
		if got != expected[name] {
			t.Fatalf("%s deveria ser %d, obtido %d", name, expected[name], got)
		}
	}
}
