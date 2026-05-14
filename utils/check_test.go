package utils_test

import (
	"testing"

	"github.com/ahmedsat/ebda-cli/utils"
)

func TestEvaluateChecks(t *testing.T) {
	checks := []utils.Check{
		{Name: "heavy", Ok: true, Weight: 3},
		{Name: "light", Ok: false, Weight: 1},
	}
	got := utils.EvaluateChecks(checks, &utils.Options{MinScore: 0.5})
	if !got.Passed {
		t.Fatalf("expected score %v to pass MinScore 0.5", got.Score)
	}
	if got.Score != 0.75 {
		t.Fatalf("Score = %v, want 0.75", got.Score)
	}
	if got.TotalWeight != 4 || got.FailedWeight != 1 {
		t.Fatalf("weights = total %v failed %v, want total 4 failed 1", got.TotalWeight, got.FailedWeight)
	}
	if len(got.Issues) != 1 || got.Issues[0] != "light" {
		t.Fatalf("Issues = %v, want [light]", got.Issues)
	}
}

func TestEvaluateChecksEmpty(t *testing.T) {
	got := utils.EvaluateChecks(nil, nil)
	if !got.Passed || got.Score != 1 {
		t.Fatalf("empty checks = passed %v score %v, want passed true score 1", got.Passed, got.Score)
	}
}
