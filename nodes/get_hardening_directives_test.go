package nodes_test

import (
	"context"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

func TestGetHardeningDirectives_Golden(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetHardeningDirectives(ctx, ax, &gen.GetHardeningDirectivesInput{Text: serviceFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []struct{ key, value string }{
		{"NoNewPrivileges", "true"},
		{"ProtectSystem", "full"},
		{"PrivateTmp", "true"},
		{"LimitNOFILE", "65536"},
	}
	if len(got.Directives) != len(want) {
		t.Fatalf("got %d hardening directives, want %d: %+v", len(got.Directives), len(want), got.Directives)
	}
	for i, w := range want {
		if got.Directives[i].GetKey() != w.key || got.Directives[i].GetValue() != w.value {
			t.Errorf("Directives[%d] = {%q,%q}, want {%q,%q}", i,
				got.Directives[i].GetKey(), got.Directives[i].GetValue(), w.key, w.value)
		}
	}
}

// TestGetHardeningDirectives_ExcludesNonAllowlisted proves ordinary
// [Service] directives (Type, ExecStart, User, ...) are NOT reported here
// even though they are real directives — this node is a filtered view.
func TestGetHardeningDirectives_ExcludesNonAllowlisted(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetHardeningDirectives(ctx, ax, &gen.GetHardeningDirectivesInput{
		Text: "[Service]\nType=simple\nExecStart=/bin/true\nUser=nobody\n",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Directives) != 0 {
		t.Errorf("expected no hardening directives, got %+v", got.Directives)
	}
}
