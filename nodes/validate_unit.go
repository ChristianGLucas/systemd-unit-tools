package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// ValidateUnit checks basic structural correctness: unrecognized sections,
// a [Service] with no ExecStart (warning), a [Socket] with no Listen*
// (warning), a [Timer] with no schedule (warning), and an unrecognized
// Type= value (error). Static analysis only — never invokes systemd.
func ValidateUnit(ctx context.Context, ax axiom.Context, input *gen.ValidateUnitInput) (*gen.ValidateUnitOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.ValidateUnitOutput{Error: errStr}, nil
	}
	issues, valid, unitType := validateUnit(u)
	return &gen.ValidateUnitOutput{Valid: valid, Issues: issues, UnitType: unitType}, nil
}
