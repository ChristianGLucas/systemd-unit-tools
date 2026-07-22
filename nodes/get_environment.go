package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// GetEnvironment extracts the [Service] section's Environment= and
// EnvironmentFile= directives, tokenizing Environment='s quoted
// multi-assignment syntax into individual KEY=VALUE pairs.
func GetEnvironment(ctx context.Context, ax axiom.Context, input *gen.GetEnvironmentInput) (*gen.GetEnvironmentOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.GetEnvironmentOutput{Error: errStr}, nil
	}
	return &gen.GetEnvironmentOutput{
		Environment:      parseEnvironmentEntries(directiveValues(u, "Service", "Environment")),
		EnvironmentFiles: parseEnvironmentFiles(directiveValues(u, "Service", "EnvironmentFile")),
	}, nil
}
