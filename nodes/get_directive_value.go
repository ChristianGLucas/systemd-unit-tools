package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// GetDirectiveValue looks up one directive by section + key and returns
// every value it was declared with, in source order — correctly handling
// directives that legitimately repeat.
func GetDirectiveValue(ctx context.Context, ax axiom.Context, input *gen.GetDirectiveValueInput) (*gen.GetDirectiveValueOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.GetDirectiveValueOutput{Error: errStr}, nil
	}
	if !hasSection(u, input.GetSectionName()) {
		return &gen.GetDirectiveValueOutput{Found: false}, nil
	}
	vals := directiveValues(u, input.GetSectionName(), input.GetKey())
	return &gen.GetDirectiveValueOutput{Values: vals, Found: len(vals) > 0}, nil
}
