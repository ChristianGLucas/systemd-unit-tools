package nodes_test

import (
	"context"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

// TestGetDirectiveValue_RepeatedKeyReturnsAll is the repeated-key-safe
// lookup's own dedicated test: After= appears twice in [Unit]; both raw
// values (unsplit) must come back, never just the last.
func TestGetDirectiveValue_RepeatedKeyReturnsAll(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetDirectiveValue(ctx, ax, &gen.GetDirectiveValueInput{
		Text: serviceFixture, SectionName: "Unit", Key: "After",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Found {
		t.Fatalf("expected found=true")
	}
	want := []string{
		"network.target remote-fs.target nss-lookup.target",
		"network-online.target",
	}
	if len(got.Values) != len(want) {
		t.Fatalf("got %d values, want %d: %v", len(got.Values), len(want), got.Values)
	}
	for i, w := range want {
		if got.Values[i] != w {
			t.Errorf("Values[%d] = %q, want %q", i, got.Values[i], w)
		}
	}
}

func TestGetDirectiveValue_NotFound(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetDirectiveValue(ctx, ax, &gen.GetDirectiveValueInput{
		Text: serviceFixture, SectionName: "Service", Key: "DoesNotExist",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Found {
		t.Fatalf("expected found=false")
	}
	if len(got.Values) != 0 {
		t.Fatalf("expected no values, got %v", got.Values)
	}
}

func TestGetDirectiveValue_SectionNotFound(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetDirectiveValue(ctx, ax, &gen.GetDirectiveValueInput{
		Text: serviceFixture, SectionName: "Timer", Key: "OnCalendar",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Found {
		t.Fatalf("expected found=false for an absent section")
	}
}
