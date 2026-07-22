package nodes_test

import (
	"context"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

func TestGetRestartPolicy_Declared(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetRestartPolicy(ctx, ax, &gen.GetRestartPolicyInput{Text: serviceFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Restart != "on-failure" || got.RestartSec != "5s" {
		t.Errorf("got Restart=%q RestartSec=%q, want on-failure/5s", got.Restart, got.RestartSec)
	}
}

func TestGetRestartPolicy_Absent(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetRestartPolicy(ctx, ax, &gen.GetRestartPolicyInput{Text: "[Service]\nExecStart=/bin/true\n"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Restart != "" || got.RestartSec != "" {
		t.Errorf("got Restart=%q RestartSec=%q, want both empty", got.Restart, got.RestartSec)
	}
}
