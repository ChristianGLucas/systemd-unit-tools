package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// GetProcessIdentity extracts the [Service] section's User=, Group=, and
// WorkingDirectory= values verbatim, or "" when absent.
func GetProcessIdentity(ctx context.Context, ax axiom.Context, input *gen.GetProcessIdentityInput) (*gen.GetProcessIdentityOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.GetProcessIdentityOutput{Error: errStr}, nil
	}
	return &gen.GetProcessIdentityOutput{
		User:             lastOrEmpty(directiveValues(u, "Service", "User")),
		Group:            lastOrEmpty(directiveValues(u, "Service", "Group")),
		WorkingDirectory: lastOrEmpty(directiveValues(u, "Service", "WorkingDirectory")),
	}, nil
}
