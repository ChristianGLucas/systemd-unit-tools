package nodes_test

import (
	"context"
	"reflect"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

func TestGetDependencies_Golden(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetDependencies(ctx, ax, &gen.GetDependenciesInput{Text: serviceFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// After= appears twice; both occurrences' whitespace-tokenized names are
	// concatenated in source order, with the duplicate network-online.target
	// preserved (it also appears standalone in Wants=).
	wantAfter := []string{"network.target", "remote-fs.target", "nss-lookup.target", "network-online.target"}
	if !reflect.DeepEqual(got.After, wantAfter) {
		t.Errorf("After = %v, want %v", got.After, wantAfter)
	}
	wantWants := []string{"network-online.target"}
	if !reflect.DeepEqual(got.Wants, wantWants) {
		t.Errorf("Wants = %v, want %v", got.Wants, wantWants)
	}
	for _, empty := range [][]string{got.Requires, got.Before, got.Conflicts, got.BindsTo, got.PartOf} {
		if len(empty) != 0 {
			t.Errorf("expected an empty list, got %v", empty)
		}
	}
}

func TestGetDependencies_MultiValueDirective(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	text := "[Unit]\nRequires=a.service b.service\nConflicts=c.service\n"
	got, err := nodes.GetDependencies(ctx, ax, &gen.GetDependenciesInput{Text: text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got.Requires, []string{"a.service", "b.service"}) {
		t.Errorf("Requires = %v", got.Requires)
	}
	if !reflect.DeepEqual(got.Conflicts, []string{"c.service"}) {
		t.Errorf("Conflicts = %v", got.Conflicts)
	}
}

func TestGetDependencies_NoUnitSection(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetDependencies(ctx, ax, &gen.GetDependenciesInput{Text: "[Service]\nExecStart=/bin/true\n"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.After) != 0 || len(got.Requires) != 0 {
		t.Errorf("expected all-empty dependency lists, got %+v", got)
	}
}
