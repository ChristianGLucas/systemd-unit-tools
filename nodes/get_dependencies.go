package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// GetDependencies extracts the [Unit] section's ordering/dependency
// directives — Requires, Wants, After, Before, Conflicts, BindsTo, PartOf —
// as structured lists of unit names, the key node for understanding
// service ordering and relationships.
func GetDependencies(ctx context.Context, ax axiom.Context, input *gen.GetDependenciesInput) (*gen.GetDependenciesOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.GetDependenciesOutput{Error: errStr}, nil
	}
	get := func(key string) []string {
		return tokenizeAll(directiveValues(u, "Unit", key))
	}
	return &gen.GetDependenciesOutput{
		Requires:  get("Requires"),
		Wants:     get("Wants"),
		After:     get("After"),
		Before:    get("Before"),
		Conflicts: get("Conflicts"),
		BindsTo:   get("BindsTo"),
		PartOf:    get("PartOf"),
	}, nil
}
