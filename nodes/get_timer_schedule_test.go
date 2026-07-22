package nodes_test

import (
	"context"
	"reflect"
	"testing"

	gen "christiangeorgelucas/systemd-unit-tools/gen"
	"christiangeorgelucas/systemd-unit-tools/nodes"
)

// TestGetTimerSchedule_Golden asserts both repeated OnCalendar= expressions
// survive, never whitespace-split (each is one calendar expression, not a
// token list — unlike GetDependencies/GetInstallDirectives).
func TestGetTimerSchedule_Golden(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetTimerSchedule(ctx, ax, &gen.GetTimerScheduleInput{Text: timerFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantCalendar := []string{"*-*-* 6,18:00", "*-*-* 12:00"}
	if !reflect.DeepEqual(got.OnCalendar, wantCalendar) {
		t.Errorf("OnCalendar = %v, want %v", got.OnCalendar, wantCalendar)
	}
	if !reflect.DeepEqual(got.OnBootSec, []string{"10min"}) {
		t.Errorf("OnBootSec = %v, want [10min]", got.OnBootSec)
	}
	if !reflect.DeepEqual(got.OnUnitActiveSec, []string{"1h"}) {
		t.Errorf("OnUnitActiveSec = %v, want [1h]", got.OnUnitActiveSec)
	}
}

func TestGetTimerSchedule_NoTimerSection(t *testing.T) {
	ctx, ax := context.Background(), newTestContext(t)
	got, err := nodes.GetTimerSchedule(ctx, ax, &gen.GetTimerScheduleInput{Text: serviceFixture})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.OnCalendar) != 0 || len(got.OnBootSec) != 0 || len(got.OnUnitActiveSec) != 0 {
		t.Errorf("expected all-empty schedule, got %+v", got)
	}
}
