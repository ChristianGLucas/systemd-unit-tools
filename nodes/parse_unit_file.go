package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// ParseUnitFile parses raw systemd unit-file text into the canonical
// UnitFile envelope of sections and directives, preserving repeated
// directives and source order. Returns a structured error, never a crash,
// on malformed input.
func ParseUnitFile(ctx context.Context, ax axiom.Context, input *gen.ParseUnitFileInput) (*gen.ParseUnitFileOutput, error) {
	if input.GetText() == "" {
		return &gen.ParseUnitFileOutput{Error: "no input provided: 'text' is empty"}, nil
	}
	u, err := parseUnitText(input.GetText())
	if err != nil {
		return &gen.ParseUnitFileOutput{Error: err.Error()}, nil
	}
	return &gen.ParseUnitFileOutput{Unit: u}, nil
}
