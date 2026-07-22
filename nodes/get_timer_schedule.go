package nodes

import (
	"context"

	"christiangeorgelucas/systemd-unit-tools/axiom"
	gen "christiangeorgelucas/systemd-unit-tools/gen"
)

// GetTimerSchedule extracts a .timer unit's [Timer] scheduling
// directives — OnCalendar, OnBootSec, OnUnitActiveSec — each occurrence
// kept as one raw expression (never whitespace-split; OnCalendar= in
// particular is commonly repeated to declare a union of schedules).
func GetTimerSchedule(ctx context.Context, ax axiom.Context, input *gen.GetTimerScheduleInput) (*gen.GetTimerScheduleOutput, error) {
	u, errStr := resolveUnit(input.GetText(), input.GetUnit())
	if errStr != "" {
		return &gen.GetTimerScheduleOutput{Error: errStr}, nil
	}
	return &gen.GetTimerScheduleOutput{
		OnCalendar:      directiveValues(u, "Timer", "OnCalendar"),
		OnBootSec:       directiveValues(u, "Timer", "OnBootSec"),
		OnUnitActiveSec: directiveValues(u, "Timer", "OnUnitActiveSec"),
	}, nil
}
