package nodes_test

import (
	"context"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

func TestDetectUnitType_FromFilename(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.DetectUnitType(ctx, ax, &gen.DetectUnitTypeInput{Text: serviceFixture, Filename: "nginx.service"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Error != "" {
		t.Fatalf("unexpected error: %s", got.Error)
	}
	if got.UnitType != "service" || got.DetectionBasis != "filename_extension" {
		t.Errorf("got type=%q basis=%q, want type=service basis=filename_extension", got.UnitType, got.DetectionBasis)
	}
}

// TestDetectUnitType_FilenameOnlyNoContent: a filename hint with NO
// text/unit content at all must still resolve — the extension check needs
// no content, so it must not be gated behind requiring a unit body.
func TestDetectUnitType_FilenameOnlyNoContent(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.DetectUnitType(ctx, ax, &gen.DetectUnitTypeInput{Filename: "backup.timer"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Error != "" {
		t.Fatalf("unexpected error: %s", got.Error)
	}
	if got.UnitType != "timer" || got.DetectionBasis != "filename_extension" {
		t.Errorf("got type=%q basis=%q, want type=timer basis=filename_extension", got.UnitType, got.DetectionBasis)
	}
}

// TestDetectUnitType_UnrecognizedFilenameFallsBackToContent: an
// unrecognized extension (or none) must fall through to content-based
// detection rather than erroring outright.
func TestDetectUnitType_UnrecognizedFilenameFallsBackToContent(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.DetectUnitType(ctx, ax, &gen.DetectUnitTypeInput{Text: socketFixture, Filename: "weird.conf"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.UnitType != "socket" || got.DetectionBasis != "section_presence" {
		t.Errorf("got type=%q basis=%q, want type=socket basis=section_presence (fallback)", got.UnitType, got.DetectionBasis)
	}
}

func TestDetectUnitType_FromSectionPresence(t *testing.T) {
	cases := map[string]struct {
		text     string
		wantType string
	}{
		"service": {serviceFixture, "service"},
		"timer":   {timerFixture, "timer"},
		"socket":  {socketFixture, "socket"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctx, ax := context.Background(), newTestContext(t)
			got, err := nodes.DetectUnitType(ctx, ax, &gen.DetectUnitTypeInput{Text: c.text})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.UnitType != c.wantType || got.DetectionBasis != "section_presence" {
				t.Errorf("got type=%q basis=%q, want type=%s basis=section_presence", got.UnitType, got.DetectionBasis, c.wantType)
			}
		})
	}
}

// TestDetectUnitType_Ambiguous: a [Unit]+[Install]-only file (the ordinary
// shape of a .target unit, but also possible for an empty .slice/.scope)
// cannot be resolved from content alone.
func TestDetectUnitType_Ambiguous(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	text := "[Unit]\nDescription=a generic grouping unit\n\n[Install]\nWantedBy=multi-user.target\n"
	got, err := nodes.DetectUnitType(ctx, ax, &gen.DetectUnitTypeInput{Text: text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.UnitType != "" || got.DetectionBasis != "ambiguous" {
		t.Errorf("got type=%q basis=%q, want type=\"\" basis=ambiguous", got.UnitType, got.DetectionBasis)
	}
}

func TestDetectUnitType_MalformedInput(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.DetectUnitType(ctx, ax, &gen.DetectUnitTypeInput{Text: malformedGarbageAfterHeader})
	if err != nil {
		t.Fatalf("expected structured error, not a Go error: %v", err)
	}
	if got.Error == "" {
		t.Fatalf("expected a structured error")
	}
}
