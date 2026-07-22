package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// ListSections lists every section present in a unit file, in
// first-appearance order, with each section's directive count.
func ListSections(ctx context.Context, ax axiom.Context, input *gen.ListSectionsInput) (*gen.ListSectionsOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.ListSectionsOutput{Error: errStr}, nil
	}
	out := &gen.ListSectionsOutput{}
	for _, s := range u.GetSections() {
		out.Sections = append(out.Sections, &gen.SectionSummary{
			Name:           s.GetName(),
			DirectiveCount: int32(len(s.GetDirectives())),
		})
	}
	return out, nil
}
