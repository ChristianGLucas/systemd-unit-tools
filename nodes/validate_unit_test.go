package nodes_test

import (
	"context"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

func TestValidateUnit_CleanFixturesAreValid(t *testing.T) {
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
			got, err := nodes.ValidateUnit(ctx, ax, &gen.ValidateUnitInput{Text: c.text})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Valid {
				t.Errorf("expected valid=true, got issues: %+v", got.Issues)
			}
			if got.UnitType != c.wantType {
				t.Errorf("UnitType = %q, want %q", got.UnitType, c.wantType)
			}
			for _, iss := range got.Issues {
				if iss.GetSeverity() == "error" {
					t.Errorf("unexpected error-severity issue on a clean fixture: %+v", iss)
				}
			}
		})
	}
}

// TestValidateUnit_MissingExecStartIsWarningOnly: a real, legitimate
// ordering-only oneshot pattern must not be rejected as invalid.
func TestValidateUnit_MissingExecStartIsWarningOnly(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	text := "[Unit]\nDescription=ordering only\n\n[Service]\nType=oneshot\nRemainAfterExit=yes\n"
	got, err := nodes.ValidateUnit(ctx, ax, &gen.ValidateUnitInput{Text: text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Valid {
		t.Fatalf("expected valid=true (missing ExecStart is a warning, not an error), got %+v", got.Issues)
	}
	foundWarning := false
	for _, iss := range got.Issues {
		if iss.GetSection() == "Service" && iss.GetSeverity() == "warning" {
			foundWarning = true
		}
	}
	if !foundWarning {
		t.Errorf("expected a warning about the missing ExecStart, got %+v", got.Issues)
	}
}

func TestValidateUnit_UnknownTypeIsError(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	text := "[Service]\nType=not-a-real-type\nExecStart=/bin/true\n"
	got, err := nodes.ValidateUnit(ctx, ax, &gen.ValidateUnitInput{Text: text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Valid {
		t.Fatalf("expected valid=false for an unrecognized Type=")
	}
	foundError := false
	for _, iss := range got.Issues {
		if iss.GetSeverity() == "error" {
			foundError = true
		}
	}
	if !foundError {
		t.Errorf("expected an error-severity issue, got %+v", got.Issues)
	}
}

func TestValidateUnit_UnrecognizedSectionIsWarning(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	text := "[Unit]\nDescription=x\n\n[NotARealSection]\nFoo=bar\n"
	got, err := nodes.ValidateUnit(ctx, ax, &gen.ValidateUnitInput{Text: text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Valid {
		t.Errorf("an unrecognized section should warn, not invalidate: %+v", got.Issues)
	}
	found := false
	for _, iss := range got.Issues {
		if iss.GetSection() == "NotARealSection" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected a warning naming the unrecognized section, got %+v", got.Issues)
	}
}

// TestValidateUnit_VendorExtensionSectionAllowed: an "X-"-prefixed section
// is a documented vendor-extension convention and must not be flagged.
func TestValidateUnit_VendorExtensionSectionAllowed(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	text := "[Unit]\nDescription=x\n\n[X-MyVendorExtension]\nFoo=bar\n"
	got, err := nodes.ValidateUnit(ctx, ax, &gen.ValidateUnitInput{Text: text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, iss := range got.Issues {
		if iss.GetSection() == "X-MyVendorExtension" {
			t.Errorf("an X-prefixed vendor section should not be flagged: %+v", iss)
		}
	}
}

func TestValidateUnit_SocketWithNoListener(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.ValidateUnit(ctx, ax, &gen.ValidateUnitInput{Text: "[Socket]\nSocketMode=0660\n"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, iss := range got.Issues {
		if iss.GetSection() == "Socket" && iss.GetSeverity() == "warning" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected a warning about no Listen* directive, got %+v", got.Issues)
	}
}

func TestValidateUnit_EmptyFileIsError(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.ValidateUnit(ctx, ax, &gen.ValidateUnitInput{Text: "# just a comment\n"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Valid {
		t.Errorf("expected valid=false for a unit file with no sections at all")
	}
}

func TestValidateUnit_MalformedInput(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.ValidateUnit(ctx, ax, &gen.ValidateUnitInput{Text: malformedNoClosingBracket})
	if err != nil {
		t.Fatalf("expected structured error, not a Go error: %v", err)
	}
	if got.Error == "" {
		t.Fatalf("expected a structured error")
	}
}
