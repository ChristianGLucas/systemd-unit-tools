package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// GetMetadata extracts the [Unit] section's Description= and
// Documentation= values.
func GetMetadata(ctx context.Context, ax axiom.Context, input *gen.GetMetadataInput) (*gen.GetMetadataOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.GetMetadataOutput{Error: errStr}, nil
	}
	return &gen.GetMetadataOutput{
		Description:   lastOrEmpty(directiveValues(u, "Unit", "Description")),
		Documentation: tokenizeAll(directiveValues(u, "Unit", "Documentation")),
	}, nil
}
