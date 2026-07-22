package nodes_test

import (
	"context"
	"reflect"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

func TestGetInstallDirectives_Golden(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetInstallDirectives(ctx, ax, &gen.GetInstallDirectivesInput{Text: serviceFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got.WantedBy, []string{"multi-user.target"}) {
		t.Errorf("WantedBy = %v, want [multi-user.target]", got.WantedBy)
	}
	if len(got.RequiredBy) != 0 || len(got.Alias) != 0 {
		t.Errorf("expected empty RequiredBy/Alias, got %v / %v", got.RequiredBy, got.Alias)
	}
}

func TestGetInstallDirectives_MultiTarget(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	text := "[Install]\nWantedBy=multi-user.target graphical.target\nAlias=foo.service\n"
	got, err := nodes.GetInstallDirectives(ctx, ax, &gen.GetInstallDirectivesInput{Text: text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got.WantedBy, []string{"multi-user.target", "graphical.target"}) {
		t.Errorf("WantedBy = %v", got.WantedBy)
	}
	if !reflect.DeepEqual(got.Alias, []string{"foo.service"}) {
		t.Errorf("Alias = %v", got.Alias)
	}
}
