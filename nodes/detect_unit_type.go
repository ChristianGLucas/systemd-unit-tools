package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// DetectUnitType determines the systemd unit type from a supplied
// filename's extension, or, failing that, from which type-distinguishing
// section is present in the parsed content. A filename alone (no text/unit
// content at all) is a fully valid call — the extension check needs no
// content, so it is tried BEFORE requiring any unit body; content is only
// required when the filename doesn't resolve it.
func DetectUnitType(ctx context.Context, ax axiom.Context, input *gen.DetectUnitTypeInput) (*gen.DetectUnitTypeOutput, error) {
	if filename := input.GetFilename(); filename != "" {
		if t, ok := extensionToType[lowerExt(filename)]; ok {
			return &gen.DetectUnitTypeOutput{UnitType: t, DetectionBasis: "filename_extension"}, nil
		}
	}
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.DetectUnitTypeOutput{Error: errStr}, nil
	}
	t, basis := detectUnitType(u, "")
	return &gen.DetectUnitTypeOutput{UnitType: t, DetectionBasis: basis}, nil
}
