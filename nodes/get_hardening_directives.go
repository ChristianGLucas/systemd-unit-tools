package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// GetHardeningDirectives extracts every [Service] directive matching this
// package's known allowlist of resource-limit and sandboxing directives —
// the security/hardening-audit node.
func GetHardeningDirectives(ctx context.Context, ax axiom.Context, input *gen.GetHardeningDirectivesInput) (*gen.GetHardeningDirectivesOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.GetHardeningDirectivesOutput{Error: errStr}, nil
	}
	out := &gen.GetHardeningDirectivesOutput{}
	svc, ok := findSection(u, "Service")
	if !ok {
		return out, nil
	}
	for _, d := range svc.GetDirectives() {
		if hardeningAllowlist[d.GetKey()] {
			out.Directives = append(out.Directives, &gen.Directive{Key: d.GetKey(), Value: d.GetValue()})
		}
	}
	return out, nil
}
