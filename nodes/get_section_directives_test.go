package nodes_test

import (
	"context"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

func TestGetSectionDirectives_Found(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetSectionDirectives(ctx, ax, &gen.GetSectionDirectivesInput{Text: socketFixture, SectionName: "Socket"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Found {
		t.Fatalf("expected found=true")
	}
	if len(got.Directives) != 5 {
		t.Fatalf("got %d directives, want 5", len(got.Directives))
	}
	first := got.Directives[0]
	if first.GetKey() != "ListenStream" || first.GetValue() != "/run/docker.sock" {
		t.Errorf("directives[0] = {%q,%q}, want {ListenStream,/run/docker.sock}", first.GetKey(), first.GetValue())
	}
	// Both repeated ListenStream= entries must survive, distinctly.
	second := got.Directives[1]
	if second.GetKey() != "ListenStream" || second.GetValue() != "127.0.0.1:2375" {
		t.Errorf("directives[1] = {%q,%q}, want {ListenStream,127.0.0.1:2375}", second.GetKey(), second.GetValue())
	}
}

func TestGetSectionDirectives_NotFound(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetSectionDirectives(ctx, ax, &gen.GetSectionDirectivesInput{Text: socketFixture, SectionName: "Timer"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Found {
		t.Fatalf("expected found=false for an absent section")
	}
	if len(got.Directives) != 0 {
		t.Fatalf("expected no directives, got %v", got.Directives)
	}
}
