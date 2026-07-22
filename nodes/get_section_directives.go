package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// GetSectionDirectives extracts every directive within one named section,
// in source order and with full multiplicity — the unfiltered view of a
// section.
func GetSectionDirectives(ctx context.Context, ax axiom.Context, input *gen.GetSectionDirectivesInput) (*gen.GetSectionDirectivesOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.GetSectionDirectivesOutput{Error: errStr}, nil
	}
	s, found := findSection(u, input.GetSectionName())
	if !found {
		return &gen.GetSectionDirectivesOutput{Found: false}, nil
	}
	out := &gen.GetSectionDirectivesOutput{Found: true}
	for _, d := range s.GetDirectives() {
		out.Directives = append(out.Directives, &gen.Directive{Key: d.GetKey(), Value: d.GetValue()})
	}
	return out, nil
}
