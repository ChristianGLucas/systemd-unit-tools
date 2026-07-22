package nodes_test

import (
	"context"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

func TestListSections_Golden(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.ListSections(ctx, ax, &gen.ListSectionsInput{Text: serviceFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Error != "" {
		t.Fatalf("unexpected error: %s", got.Error)
	}
	want := []struct {
		Name           string
		DirectiveCount int32
	}{
		{Name: "Unit", DirectiveCount: 5},
		{Name: "Service", DirectiveCount: 20},
		{Name: "Install", DirectiveCount: 1},
	}
	if len(got.Sections) != len(want) {
		t.Fatalf("got %d sections, want %d", len(got.Sections), len(want))
	}
	for i, w := range want {
		if got.Sections[i].GetName() != w.Name || got.Sections[i].GetDirectiveCount() != w.DirectiveCount {
			t.Errorf("section[%d] = {%q,%d}, want {%q,%d}", i,
				got.Sections[i].GetName(), got.Sections[i].GetDirectiveCount(), w.Name, w.DirectiveCount)
		}
	}
}

// TestListSections_ComposedFromParsedUnit exercises the flow-composition
// path: passing a pre-parsed `unit` instead of raw `text`.
func TestListSections_ComposedFromParsedUnit(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	parsed, err := nodes.ParseUnitFile(ctx, ax, &gen.ParseUnitFileInput{Text: socketFixture})
	if err != nil || parsed.Error != "" {
		t.Fatalf("setup: unexpected parse failure: %v / %s", err, parsed.Error)
	}
	got, err := nodes.ListSections(ctx, ax, &gen.ListSectionsInput{Unit: parsed.Unit})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Sections) != 3 {
		t.Fatalf("got %d sections, want 3", len(got.Sections))
	}
}

func TestListSections_NoInput(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.ListSections(ctx, ax, &gen.ListSectionsInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Error == "" {
		t.Fatalf("expected a structured error for no input")
	}
}
