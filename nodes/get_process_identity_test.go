package nodes_test

import (
	"context"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

func TestGetProcessIdentity_Golden(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetProcessIdentity(ctx, ax, &gen.GetProcessIdentityInput{Text: serviceFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.User != "www-data" || got.Group != "www-data" || got.WorkingDirectory != "/var/www" {
		t.Errorf("got User=%q Group=%q WorkingDirectory=%q, want www-data/www-data//var/www",
			got.User, got.Group, got.WorkingDirectory)
	}
}

func TestGetProcessIdentity_Absent(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetProcessIdentity(ctx, ax, &gen.GetProcessIdentityInput{Text: "[Service]\nExecStart=/bin/true\n"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.User != "" || got.Group != "" || got.WorkingDirectory != "" {
		t.Errorf("expected all-empty identity, got %+v", got)
	}
}
