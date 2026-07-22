package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// GetServiceType extracts the [Service] section's declared Type= value
// verbatim. Reports only what is explicitly declared — it does not apply
// systemd's own implicit default.
func GetServiceType(ctx context.Context, ax axiom.Context, input *gen.GetServiceTypeInput) (*gen.GetServiceTypeOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.GetServiceTypeOutput{Error: errStr}, nil
	}
	vals := directiveValues(u, "Service", "Type")
	if len(vals) == 0 {
		return &gen.GetServiceTypeOutput{Explicit: false}, nil
	}
	return &gen.GetServiceTypeOutput{Type: lastOrEmpty(vals), Explicit: true}, nil
}
