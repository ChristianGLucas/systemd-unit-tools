package nodes

import (
	"context"
	"strings"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// GetSocketListeners extracts a .socket unit's [Socket] Listen* directives
// as a list of {type, address} pairs, in source order and with full
// multiplicity.
func GetSocketListeners(ctx context.Context, ax axiom.Context, input *gen.GetSocketListenersInput) (*gen.GetSocketListenersOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.GetSocketListenersOutput{Error: errStr}, nil
	}
	out := &gen.GetSocketListenersOutput{}
	sock, ok := findSection(u, "Socket")
	if !ok {
		return out, nil
	}
	for _, d := range sock.GetDirectives() {
		if strings.HasPrefix(d.GetKey(), "Listen") {
			out.Listeners = append(out.Listeners, &gen.ListenEntry{
				Type:    strings.TrimPrefix(d.GetKey(), "Listen"),
				Address: d.GetValue(),
			})
		}
	}
	return out, nil
}
