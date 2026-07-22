package nodes_test

import (
	"context"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

// TestGetEnvironment_Golden hand-verifies the quoted multi-assignment
// tokenizer against the fixture's two Environment= lines: one plain
// multi-assignment, one double-quoted value containing a space.
func TestGetEnvironment_Golden(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetEnvironment(ctx, ax, &gen.GetEnvironmentInput{Text: serviceFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []struct{ key, value string }{
		{"NGINX_ENV", "production"},
		{"DEBUG", "0"},
		{"EXTRA", "hello world"},
	}
	if len(got.Environment) != len(want) {
		t.Fatalf("got %d env vars, want %d: %+v", len(got.Environment), len(want), got.Environment)
	}
	for i, w := range want {
		if got.Environment[i].GetKey() != w.key || got.Environment[i].GetValue() != w.value {
			t.Errorf("Environment[%d] = {%q,%q}, want {%q,%q}", i,
				got.Environment[i].GetKey(), got.Environment[i].GetValue(), w.key, w.value)
		}
	}
	if len(got.EnvironmentFiles) != 1 {
		t.Fatalf("got %d environment files, want 1", len(got.EnvironmentFiles))
	}
	ef := got.EnvironmentFiles[0]
	if ef.GetPath() != "/etc/default/nginx" || !ef.GetOptional() {
		t.Errorf("EnvironmentFiles[0] = {%q,optional=%v}, want {/etc/default/nginx,optional=true}", ef.GetPath(), ef.GetOptional())
	}
}

// TestGetEnvironment_EscapeSequences exercises the documented systemd.syntax
// escapes (\s \t \" \\) independent of the main fixture.
func TestGetEnvironment_EscapeSequences(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	text := `[Service]
Environment="A=one\stwo" B=plain C="quote\"inside"
`
	got, err := nodes.GetEnvironment(ctx, ax, &gen.GetEnvironmentInput{Text: text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]string{"A": "one two", "B": "plain", "C": `quote"inside`}
	if len(got.Environment) != len(want) {
		t.Fatalf("got %d env vars, want %d: %+v", len(got.Environment), len(want), got.Environment)
	}
	for _, e := range got.Environment {
		wantVal, ok := want[e.GetKey()]
		if !ok {
			t.Errorf("unexpected key %q", e.GetKey())
			continue
		}
		if e.GetValue() != wantVal {
			t.Errorf("%s = %q, want %q", e.GetKey(), e.GetValue(), wantVal)
		}
	}
}

func TestGetEnvironment_Absent(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetEnvironment(ctx, ax, &gen.GetEnvironmentInput{Text: "[Service]\nExecStart=/bin/true\n"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Environment) != 0 || len(got.EnvironmentFiles) != 0 {
		t.Errorf("expected no env vars/files, got %+v / %+v", got.Environment, got.EnvironmentFiles)
	}
}
