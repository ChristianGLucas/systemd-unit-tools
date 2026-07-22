package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// SummarizeUnit produces a directive-count summary of a unit file: total
// section/directive counts, per-section directive counts, and the detected
// unit type.
func SummarizeUnit(ctx context.Context, ax axiom.Context, input *gen.SummarizeUnitInput) (*gen.SummarizeUnitOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.SummarizeUnitOutput{Error: errStr}, nil
	}
	out := &gen.SummarizeUnitOutput{SectionCount: int32(len(u.GetSections()))}
	for _, s := range u.GetSections() {
		n := int32(len(s.GetDirectives()))
		out.DirectiveCount += n
		out.Sections = append(out.Sections, &gen.SectionSummary{Name: s.GetName(), DirectiveCount: n})
	}
	unitType, _ := detectUnitType(u, "")
	out.UnitType = unitType
	return out, nil
}
