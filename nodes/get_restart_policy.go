package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// lastOrEmpty returns the last-declared value of a single-value directive
// (systemd's own semantics for a repeated non-list assignment: the last
// occurrence wins), or "" if it was never declared.
func lastOrEmpty(vals []string) string {
	if len(vals) == 0 {
		return ""
	}
	return vals[len(vals)-1]
}

// GetRestartPolicy extracts the [Service] section's Restart= and
// RestartSec= values verbatim, or "" when absent.
func GetRestartPolicy(ctx context.Context, ax axiom.Context, input *gen.GetRestartPolicyInput) (*gen.GetRestartPolicyOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.GetRestartPolicyOutput{Error: errStr}, nil
	}
	return &gen.GetRestartPolicyOutput{
		Restart:    lastOrEmpty(directiveValues(u, "Service", "Restart")),
		RestartSec: lastOrEmpty(directiveValues(u, "Service", "RestartSec")),
	}, nil
}
