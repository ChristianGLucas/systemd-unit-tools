package nodes_test

import (
	"context"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

func TestSummarizeUnit_Golden(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.SummarizeUnit(ctx, ax, &gen.SummarizeUnitInput{Text: serviceFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.SectionCount != 3 {
		t.Errorf("SectionCount = %d, want 3", got.SectionCount)
	}
	if got.DirectiveCount != 26 {
		t.Errorf("DirectiveCount = %d, want 26 (5 Unit + 20 Service + 1 Install)", got.DirectiveCount)
	}
	if got.UnitType != "service" {
		t.Errorf("UnitType = %q, want service", got.UnitType)
	}
	wantCounts := map[string]int32{"Unit": 5, "Service": 20, "Install": 1}
	for _, s := range got.Sections {
		want, ok := wantCounts[s.GetName()]
		if !ok || s.GetDirectiveCount() != want {
			t.Errorf("section %q has %d directives, want %d", s.GetName(), s.GetDirectiveCount(), want)
		}
	}
}

func TestSummarizeUnit_Timer(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.SummarizeUnit(ctx, ax, &gen.SummarizeUnitInput{Text: timerFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.SectionCount != 3 || got.DirectiveCount != 7 || got.UnitType != "timer" {
		t.Errorf("got {sections=%d, directives=%d, type=%q}, want {3, 7, timer}",
			got.SectionCount, got.DirectiveCount, got.UnitType)
	}
}

func TestSummarizeUnit_MalformedInput(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.SummarizeUnit(ctx, ax, &gen.SummarizeUnitInput{Text: malformedNoEqualsSign})
	if err != nil {
		t.Fatalf("expected structured error, not a Go error: %v", err)
	}
	if got.Error == "" {
		t.Fatalf("expected a structured error")
	}
}
