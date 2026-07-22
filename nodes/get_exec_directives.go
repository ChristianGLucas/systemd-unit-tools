package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// execDirectiveNames is the ordered-preservation filter GetExecDirectives
// applies to [Service]: every directive matching one of these names is
// returned, in source order, with full multiplicity.
var execDirectiveNames = map[string]bool{
	"ExecStartPre": true, "ExecStart": true, "ExecStartPost": true,
	"ExecCondition": true, "ExecStop": true, "ExecStopPost": true,
	"ExecReload": true,
}

// GetExecDirectives extracts the [Service] section's execution-lifecycle
// directives, in source order and with full multiplicity, verbatim
// (including any leading "-"/"@"/"+"/"!"/"!!" prefix character).
func GetExecDirectives(ctx context.Context, ax axiom.Context, input *gen.GetExecDirectivesInput) (*gen.GetExecDirectivesOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.GetExecDirectivesOutput{Error: errStr}, nil
	}
	out := &gen.GetExecDirectivesOutput{}
	svc, ok := findSection(u, "Service")
	if !ok {
		return out, nil
	}
	for _, d := range svc.GetDirectives() {
		if execDirectiveNames[d.GetKey()] {
			out.Entries = append(out.Entries, &gen.ExecEntry{Directive: d.GetKey(), Command: d.GetValue()})
		}
	}
	return out, nil
}
