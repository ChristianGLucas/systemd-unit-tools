package nodes_test

import (
	"context"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

func TestGetServiceType_Declared(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetServiceType(ctx, ax, &gen.GetServiceTypeInput{Text: serviceFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Type != "forking" || !got.Explicit {
		t.Errorf("got type=%q explicit=%v, want type=forking explicit=true", got.Type, got.Explicit)
	}
}

func TestGetServiceType_Absent(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetServiceType(ctx, ax, &gen.GetServiceTypeInput{Text: "[Service]\nExecStart=/bin/true\n"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Type != "" || got.Explicit {
		t.Errorf("got type=%q explicit=%v, want type=\"\" explicit=false", got.Type, got.Explicit)
	}
}

// TestGetServiceType_LastDeclarationWins: systemd's own semantics for a
// repeated single-value directive is "last one wins".
func TestGetServiceType_LastDeclarationWins(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	text := "[Service]\nType=simple\nType=oneshot\nExecStart=/bin/true\n"
	got, err := nodes.GetServiceType(ctx, ax, &gen.GetServiceTypeInput{Text: text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Type != "oneshot" {
		t.Errorf("got type=%q, want oneshot (last declaration should win)", got.Type)
	}
}
