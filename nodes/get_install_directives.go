package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// GetInstallDirectives extracts the [Install] section's directives —
// WantedBy, RequiredBy, Alias — as whitespace-tokenized lists, the key
// node for understanding how a unit gets enabled.
func GetInstallDirectives(ctx context.Context, ax axiom.Context, input *gen.GetInstallDirectivesInput) (*gen.GetInstallDirectivesOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.GetInstallDirectivesOutput{Error: errStr}, nil
	}
	return &gen.GetInstallDirectivesOutput{
		WantedBy:   tokenizeAll(directiveValues(u, "Install", "WantedBy")),
		RequiredBy: tokenizeAll(directiveValues(u, "Install", "RequiredBy")),
		Alias:      tokenizeAll(directiveValues(u, "Install", "Alias")),
	}, nil
}
