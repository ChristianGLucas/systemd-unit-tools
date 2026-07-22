package nodes_test

import (
	"context"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

func TestGetExecDirectives_Golden(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetExecDirectives(ctx, ax, &gen.GetExecDirectivesInput{Text: serviceFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []struct{ directive, command string }{
		{"ExecStartPre", `/usr/sbin/nginx -t -q -g "daemon on; master_process on;"`},
		{"ExecStart", `/usr/sbin/nginx -g "daemon on; master_process on;"`},
		{"ExecReload", `/usr/sbin/nginx -g "daemon on; master_process on;" -s reload`},
		{"ExecStop", `-/sbin/start-stop-daemon --quiet --stop --retry QUIT/5 --pidfile /run/nginx.pid`},
	}
	if len(got.Entries) != len(want) {
		t.Fatalf("got %d entries, want %d: %+v", len(got.Entries), len(want), got.Entries)
	}
	for i, w := range want {
		if got.Entries[i].GetDirective() != w.directive || got.Entries[i].GetCommand() != w.command {
			t.Errorf("entries[%d] = {%q,%q}, want {%q,%q}", i,
				got.Entries[i].GetDirective(), got.Entries[i].GetCommand(), w.directive, w.command)
		}
	}
}

// TestGetExecDirectives_RepeatedExecStartPre proves several ExecStartPre=
// occurrences are all kept, in order — a common real pattern.
func TestGetExecDirectives_RepeatedExecStartPre(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	text := "[Service]\nExecStartPre=/bin/first\nExecStartPre=/bin/second\nExecStart=/bin/main\n"
	got, err := nodes.GetExecDirectives(ctx, ax, &gen.GetExecDirectivesInput{Text: text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Entries) != 3 {
		t.Fatalf("got %d entries, want 3", len(got.Entries))
	}
	if got.Entries[0].GetCommand() != "/bin/first" || got.Entries[1].GetCommand() != "/bin/second" {
		t.Errorf("ExecStartPre entries out of order or lost: %+v", got.Entries[:2])
	}
}

func TestGetExecDirectives_NoServiceSection(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetExecDirectives(ctx, ax, &gen.GetExecDirectivesInput{Text: timerFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Entries) != 0 {
		t.Errorf("expected no entries for a unit with no [Service], got %+v", got.Entries)
	}
}
