// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package process waits on asynchronous Morpheus processes. Several 9.x
// mutations (for example the HVM virtual switch endpoints) return an async
// envelope carrying a processId rather than the mutated object; the real work
// fans out over the hosts in the background. Callers poll
// GET /api/processes/{id} until the process reaches a terminal status.
package process

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// Terminal process statuses, per Morpheus's ProcessService.
const (
	StatusComplete = "complete"
	StatusFailed   = "failed"
	StatusExpired  = "expired"
)

// Every async mutation returns the same envelope -- {success, processId,
// warnings} -- but the generated *200Response models are distinct Go types,
// so callers read the processId out of whichever one they have and pass it to
// Wait; there is no shared envelope type here.

// Wait polls GET /api/processes/{processID} until the process reaches a
// terminal status, using exponential backoff (starting at 2s, capped at 15s)
// and honouring ctx cancellation. It returns the terminal process on success
// (status "complete"); an error naming the process and its message/error/output
// on "failed" or "expired"; an error naming the process id on ctx cancellation
// or when the timeout elapses.
func Wait(
	ctx context.Context,
	client *sdk.APIClient,
	processID int64,
	timeout time.Duration,
) (*sdk.GetHistory200ResponseProcess, error) {
	deadlineCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	const (
		startDelay = 2 * time.Second
		maxDelay   = 15 * time.Second
	)

	delay := startDelay

	for {
		resp, hresp, err := client.HistoryAPI.GetHistory(deadlineCtx, processID).Execute()
		if hresp != nil {
			_ = hresp.Body.Close()
		}

		var proc *sdk.GetHistory200ResponseProcess

		switch {
		case err != nil && deadlineCtx.Err() != nil:
			// A cancelled or deadline-exceeded context is reported as a
			// timeout rather than a raw transport error.
			return nil, waitTimeoutErr(processID, deadlineCtx.Err())
		case err != nil && hresp != nil && hresp.StatusCode >= http.StatusBadRequest && hresp.StatusCode < http.StatusInternalServerError:
			// A 4xx is a definite answer (the process is gone, or we may not
			// read it); nothing a retry can change.
			return nil, fmt.Errorf("polling process %d failed: %w", processID, err)
		case err != nil:
			// A transport error or 5xx while the appliance is busy is
			// transient: keep polling until the deadline rather than abandon a
			// job that is still running.
			tflog.Warn(ctx, "process poll failed; retrying", map[string]any{
				"process": processID, "error": err.Error(),
			})
		case resp != nil:
			proc = resp.Process
		}

		if proc != nil && proc.Status != nil {
			switch *proc.Status {
			case StatusComplete:
				return proc, nil
			case StatusFailed, StatusExpired:
				return proc, processFailure(processID, proc)
			}
		}

		select {
		case <-deadlineCtx.Done():
			return nil, waitTimeoutErr(processID, deadlineCtx.Err())
		case <-time.After(delay):
		}

		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}
}

func waitTimeoutErr(processID int64, cause error) error {
	if cause != nil && cause == context.Canceled {
		return fmt.Errorf("waiting on process %d cancelled", processID)
	}

	return fmt.Errorf("timed out waiting on process %d to complete", processID)
}

// processFailure renders a failed/expired process into a single error carrying
// the display name and whichever of message, error, or output is present.
func processFailure(processID int64, proc *sdk.GetHistory200ResponseProcess) error {
	name := fmt.Sprintf("process %d", processID)
	if proc.DisplayName != nil && *proc.DisplayName != "" {
		name = fmt.Sprintf("%s (process %d)", *proc.DisplayName, processID)
	}

	status := ""
	if proc.Status != nil {
		status = *proc.Status
	}

	var detail []string
	if m := proc.Message.Get(); m != nil && *m != "" {
		detail = append(detail, "message: "+*m)
	}

	if e := proc.Error.Get(); e != nil && *e != "" {
		detail = append(detail, "error: "+*e)
	}

	if o := proc.Output.Get(); o != nil && *o != "" {
		detail = append(detail, "output: "+*o)
	}

	if len(detail) == 0 {
		return fmt.Errorf("%s %s with no detail", name, status)
	}

	return fmt.Errorf("%s %s: %s", name, status, strings.Join(detail, "; "))
}

// Warnings turns an async envelope's warnings slice into a diag.Diagnostics of
// warnings so callers can surface them (for example the noSpareNic warning).
func Warnings(source, name string, warnings []string) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, w := range warnings {
		if strings.TrimSpace(w) == "" {
			continue
		}

		diags.AddWarning(
			fmt.Sprintf("%s %s: %s", source, name, w),
			w,
		)
	}

	return diags
}
