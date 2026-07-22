package nodes_test

import (
	"context"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

// TestParseUnitFile_RepeatedDirectivesPreserved is the load-bearing test for
// this whole package: a plain INI parser would collapse the two "After="
// lines into one, silently losing "network-online.target". This asserts
// both survive as separate Directive entries, in source order.
func TestParseUnitFile_RepeatedDirectivesPreserved(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.ParseUnitFile(ctx, ax, &gen.ParseUnitFileInput{Text: serviceFixture})
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if got.Error != "" {
		t.Fatalf("unexpected parse error: %s", got.Error)
	}
	if len(got.Unit.GetSections()) != 3 {
		t.Fatalf("expected 3 sections, got %d", len(got.Unit.GetSections()))
	}

	unitSec := got.Unit.GetSections()[0]
	if unitSec.GetName() != "Unit" {
		t.Fatalf("expected first section 'Unit', got %q", unitSec.GetName())
	}
	var afterValues []string
	for _, d := range unitSec.GetDirectives() {
		if d.GetKey() == "After" {
			afterValues = append(afterValues, d.GetValue())
		}
	}
	want := []string{
		"network.target remote-fs.target nss-lookup.target",
		"network-online.target",
	}
	if len(afterValues) != 2 {
		t.Fatalf("expected 2 separate After= directives preserved, got %d: %v", len(afterValues), afterValues)
	}
	for i, w := range want {
		if afterValues[i] != w {
			t.Errorf("After[%d] = %q, want %q", i, afterValues[i], w)
		}
	}
}

// TestParseUnitFile_GoldenStructure hand-verifies the full section/directive
// count against the fixture text (an independent count, not derived from
// running the parser).
func TestParseUnitFile_GoldenStructure(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.ParseUnitFile(ctx, ax, &gen.ParseUnitFileInput{Text: timerFixture})
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if got.Error != "" {
		t.Fatalf("unexpected parse error: %s", got.Error)
	}
	wantCounts := map[string]int{"Unit": 1, "Timer": 5, "Install": 1}
	if len(got.Unit.GetSections()) != len(wantCounts) {
		t.Fatalf("expected %d sections, got %d", len(wantCounts), len(got.Unit.GetSections()))
	}
	for _, s := range got.Unit.GetSections() {
		want, ok := wantCounts[s.GetName()]
		if !ok {
			t.Errorf("unexpected section %q", s.GetName())
			continue
		}
		if len(s.GetDirectives()) != want {
			t.Errorf("section %q: got %d directives, want %d", s.GetName(), len(s.GetDirectives()), want)
		}
	}
}

// TestParseUnitFile_MalformedInputs proves every malformed case returns a
// structured error (never a Go error / panic) and never a populated Unit.
func TestParseUnitFile_MalformedInputs(t *testing.T) {
	cases := map[string]string{
		"unterminated section header":        malformedNoClosingBracket,
		"garbage after section header":       malformedGarbageAfterHeader,
		"no equals sign in a directive line": malformedNoEqualsSign,
	}
	for name, text := range cases {
		t.Run(name, func(t *testing.T) {
			ctx, ax := context.Background(), newTestContext(t)
			got, err := nodes.ParseUnitFile(ctx, ax, &gen.ParseUnitFileInput{Text: text})
			if err != nil {
				t.Fatalf("expected a structured error, not a Go error: %v", err)
			}
			if got.Error == "" {
				t.Fatalf("expected a non-empty structured error for input %q", text)
			}
			if got.Unit != nil && len(got.Unit.GetSections()) != 0 {
				t.Fatalf("expected no partial Unit on error, got %+v", got.Unit)
			}
		})
	}
}

// TestParseUnitFile_EmptyText: no input at all is a structured error, not a
// panic on a nil/empty string.
func TestParseUnitFile_EmptyText(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.ParseUnitFile(ctx, ax, &gen.ParseUnitFileInput{})
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if got.Error == "" {
		t.Fatalf("expected a structured error for empty text")
	}
}

// TestParseUnitFile_OversizedInput proves the 1 MB text cap fires on raw
// input before any parsing work happens.
func TestParseUnitFile_OversizedInput(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	huge := make([]byte, (1<<20)+1)
	for i := range huge {
		huge[i] = 'a'
	}
	got, err := nodes.ParseUnitFile(ctx, ax, &gen.ParseUnitFileInput{Text: string(huge)})
	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}
	if got.Error == "" {
		t.Fatalf("expected the size cap to reject oversized input")
	}
}
