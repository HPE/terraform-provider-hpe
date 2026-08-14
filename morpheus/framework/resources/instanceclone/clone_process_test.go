// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instanceclone

import (
	"testing"
	"time"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// These drive the production selection and failure rules directly. Observed
// behaviour they encode: on one appliance a clone failed server-side in 27
// seconds while the provider kept polling for the full 15-minute timeout and
// then reported a bare "<nil>"; on another the clone was still running at the
// timeout. Those two cases need different advice, so they must not collapse
// into one another.

func proc(
	code, status string, start time.Time,
) sdk.GetInstanceHistory200ResponseAllOfProcessesInner {
	p := sdk.GetInstanceHistory200ResponseAllOfProcessesInner{
		ProcessType: &sdk.GetInstanceHistory200ResponseAllOfProcessesInnerProcessType{
			Code: &code,
		},
		StartDate: &start,
	}
	if status != "" {
		p.Status = &status
	}

	return p
}

func TestUnitPickLatestCloneProcess_IgnoresOtherTypes(t *testing.T) {
	t.Parallel()

	base := time.Date(2026, 8, 13, 12, 0, 0, 0, time.UTC)
	got := pickLatestCloneProcess([]sdk.GetInstanceHistory200ResponseAllOfProcessesInner{
		proc("provision", "failed", base.Add(2*time.Hour)),
		proc("cloning", "failed", base),
		proc("startup", "complete", base.Add(3*time.Hour)),
	})

	if got == nil {
		t.Fatal("expected the cloning process to be selected")
	}

	if *got.ProcessType.Code != "cloning" {
		t.Fatalf("selected %q, want cloning", *got.ProcessType.Code)
	}

	// A later provision/startup must not win: only cloning is relevant, and
	// picking a provision failure would abort a healthy clone.
	if !processStart(got).Equal(base) {
		t.Fatalf("selected the wrong process, start=%v", processStart(got))
	}
}

func TestUnitPickLatestCloneProcess_PicksMostRecent(t *testing.T) {
	t.Parallel()

	base := time.Date(2026, 8, 13, 12, 0, 0, 0, time.UTC)
	newest := base.Add(30 * time.Minute)
	got := pickLatestCloneProcess([]sdk.GetInstanceHistory200ResponseAllOfProcessesInner{
		proc("cloning", "failed", base),
		proc("cloning", "cloning", newest),
		proc("cloning", "failed", base.Add(10*time.Minute)),
	})

	if got == nil || !processStart(got).Equal(newest) {
		t.Fatalf("expected the newest cloning process, got %v", processStart(got))
	}

	// Stale failures from earlier attempts must not abort the current one.
	if reason := cloneFailureReason(got); reason != "" {
		t.Fatalf("running clone reported as failed: %q", reason)
	}
}

func TestUnitPickLatestCloneProcess_NoneWhenAbsent(t *testing.T) {
	t.Parallel()

	if got := pickLatestCloneProcess(nil); got != nil {
		t.Fatal("expected nil for empty history")
	}

	got := pickLatestCloneProcess([]sdk.GetInstanceHistory200ResponseAllOfProcessesInner{
		proc("provision", "failed", time.Now()),
	})
	if got != nil {
		t.Fatal("expected nil when no cloning process is present")
	}
}

func TestUnitCloneFailureReason_RunningIsNotAFailure(t *testing.T) {
	t.Parallel()

	p := proc("cloning", "cloning", time.Now())
	if reason := cloneFailureReason(&p); reason != "" {
		t.Fatalf("a running clone must not abort the poll, got %q", reason)
	}
}

func TestUnitCloneFailureReason_PrefersErrorText(t *testing.T) {
	t.Parallel()

	p := proc("cloning", "failed", time.Now())
	errText := "no addresses available in pool"
	msgText := "less specific message"
	p.Error.Set(&errText)
	p.Message.Set(&msgText)

	if got := cloneFailureReason(&p); got != errText {
		t.Fatalf("reason = %q, want %q", got, errText)
	}
}

func TestUnitCloneFailureReason_FallsBackWhenSilent(t *testing.T) {
	t.Parallel()

	p := proc("cloning", "failed", time.Now())

	got := cloneFailureReason(&p)
	if got == "" {
		t.Fatal("a failed clone must always produce a reason, got empty")
	}

	// The bug being guarded: the old code surfaced errors.Unwrap(err), which
	// was nil, so the diagnostic read "<nil>" and told the operator nothing.
	if got == "<nil>" || got == "nil" {
		t.Fatalf("reason must never be a bare nil rendering, got %q", got)
	}
}

func TestUnitCloneFailureReason_NilProcess(t *testing.T) {
	t.Parallel()

	if got := cloneFailureReason(nil); got != "" {
		t.Fatalf("nil process must not be treated as a failure, got %q", got)
	}
}
