package nodes_test

import (
	"context"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

func TestGetSocketListeners_Golden(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetSocketListeners(ctx, ax, &gen.GetSocketListenersInput{Text: socketFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []struct{ typ, addr string }{
		{"Stream", "/run/docker.sock"},
		{"Stream", "127.0.0.1:2375"},
	}
	if len(got.Listeners) != len(want) {
		t.Fatalf("got %d listeners, want %d: %+v", len(got.Listeners), len(want), got.Listeners)
	}
	for i, w := range want {
		if got.Listeners[i].GetType() != w.typ || got.Listeners[i].GetAddress() != w.addr {
			t.Errorf("Listeners[%d] = {%q,%q}, want {%q,%q}", i,
				got.Listeners[i].GetType(), got.Listeners[i].GetAddress(), w.typ, w.addr)
		}
	}
}

// TestGetSocketListeners_AllListenTypes proves the prefix-strip works for
// every documented Listen* directive, not just ListenStream.
func TestGetSocketListeners_AllListenTypes(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	text := "[Socket]\nListenDatagram=/run/foo.sock\nListenFIFO=/run/foo.fifo\nListenNetlink=route 0\n"
	got, err := nodes.GetSocketListeners(ctx, ax, &gen.GetSocketListenersInput{Text: text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantTypes := []string{"Datagram", "FIFO", "Netlink"}
	if len(got.Listeners) != len(wantTypes) {
		t.Fatalf("got %d listeners, want %d", len(got.Listeners), len(wantTypes))
	}
	for i, wt := range wantTypes {
		if got.Listeners[i].GetType() != wt {
			t.Errorf("Listeners[%d].Type = %q, want %q", i, got.Listeners[i].GetType(), wt)
		}
	}
}

func TestGetSocketListeners_NoSocketSection(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetSocketListeners(ctx, ax, &gen.GetSocketListenersInput{Text: serviceFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Listeners) != 0 {
		t.Errorf("expected no listeners, got %+v", got.Listeners)
	}
}
