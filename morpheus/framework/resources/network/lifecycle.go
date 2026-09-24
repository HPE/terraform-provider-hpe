// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// Morpheus answers a network create as soon as the record exists, with status
// "provisioning" while the network's provider (a cloud, or a plugin such as
// the HVM network server configuring host interfaces) finishes the work — an
// update passes through the same status. Create waits for that to end so a
// failed provisioning is reported at apply, and state is read from a settled
// network. Delete is synchronous in Morpheus, but its answer has to be read:
// a refusal is HTTP 200 with success:false.

// The cadence is adjustable so the unit tests need not wait it out.
var (
	// settleInterval is how often a network's status is read while waiting.
	settleInterval = 3 * time.Second
	// deleteRetryAfter is the pause before a refused delete is sent again.
	deleteRetryAfter = 10 * time.Second
)

const (
	// settleTimeout bounds the wait for a new network to leave "provisioning".
	settleTimeout = 5 * time.Minute
	// deleteAttempts bounds the number of delete requests.
	deleteAttempts = 3
)

// networkStatus reads a network's status and status message. gone is true
// when the network no longer exists.
func networkStatus(ctx context.Context, client *sdk.APIClient, id int64) (status, message string, gone bool, err error) {
	r, hresp, err := client.NetworksAPI.GetNetwork(ctx, id).Execute()
	if hresp != nil && hresp.StatusCode == http.StatusNotFound {
		return "", "", true, nil
	}

	if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK || r == nil || r.Network == nil {
		return "", "", false, fmt.Errorf("network %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	if r.Network.Status.IsSet() && r.Network.Status.Get() != nil {
		status = *r.Network.Status.Get()
	}

	// statusMessage is not in the response model.
	if m, ok := r.Network.AdditionalProperties["statusMessage"].(string); ok {
		message = m
	}

	return status, message, false, nil
}

// waitForNetworkSettled waits for a newly created network to leave
// "provisioning". A failed provisioning is an error; a network that still
// reports provisioning when the wait runs out is a warning, since the record
// exists and some providers are simply slow.
func waitForNetworkSettled(ctx context.Context, client *sdk.APIClient, id int64) diag.Diagnostics {
	var diags diag.Diagnostics

	type outcome struct{ status, message string }

	poll := func() (outcome, error) {
		status, message, gone, err := networkStatus(ctx, client, id)
		if err != nil {
			// A failed read (the appliance answering 502 for a moment) is
			// retried until the wait runs out, not treated as the network's
			// fate.
			tflog.Debug(ctx, "network: status read failed while waiting; retrying", map[string]any{"error": err.Error()})

			return outcome{}, fmt.Errorf("status read failed: %w", err)
		}

		if gone {
			return outcome{}, backoff.Permanent(fmt.Errorf("network %d disappeared while provisioning", id))
		}

		switch status {
		case "provisioning":
			return outcome{}, fmt.Errorf("still provisioning")
		case "error", "failed":
			return outcome{}, backoff.Permanent(fmt.Errorf("network %d provisioning failed: status %q%s", id, status, detail(message)))
		default:
			return outcome{status, message}, nil
		}
	}

	_, err := backoff.Retry(ctx, poll,
		backoff.WithBackOff(backoff.NewConstantBackOff(settleInterval)),
		backoff.WithMaxElapsedTime(settleTimeout))

	switch {
	case err == nil:
		return diags
	case err.Error() == "still provisioning":
		diags.AddWarning("network still provisioning",
			fmt.Sprintf("Network %d still reports status \"provisioning\" after %s; continuing. "+
				"Its provider may not have finished configuring it.", id, settleTimeout))
	case strings.HasPrefix(err.Error(), "status read failed:"):
		diags.AddWarning("network status not confirmed",
			fmt.Sprintf("Network %d was created but its status could not be read for %s; continuing. "+
				"Last error: %s", id, settleTimeout, err.Error()))
	default:
		diags.AddError("create network resource", err.Error())
	}

	return diags
}

// deleteRefused is a delete Morpheus accepted as a request and refused as an
// operation. The API renders every failed delete as HTTP 200 with
// success:false and a message, so the status code alone says nothing; the
// body has to be read.
type deleteRefused struct {
	id  int64
	msg string
}

func (e deleteRefused) Error() string {
	return fmt.Sprintf("network %d: Morpheus refused the delete: %s", e.id, e.msg)
}

// deleteNetworkAndWait deletes a network and waits until it is gone. A
// refusal is retried a few times — the network's provider may be doing
// host-side work that fails transiently — and then reported with Morpheus's
// own message. Morpheus deletes synchronously, so once it answers success the
// network is gone or about to be; the wait is a guard, not an expectation.
func deleteNetworkAndWait(ctx context.Context, client *sdk.APIClient, id int64, timeout time.Duration) error {
	deleteOnce := func() error {
		r, hresp, err := client.NetworksAPI.DeleteNetwork(ctx, id).Execute()
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			return nil
		}

		if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
			return fmt.Errorf("network %d DELETE failed: %s", id, errfmt.ErrMsg(err, hresp))
		}

		if r != nil && r.Success != nil && !*r.Success {
			return deleteRefused{id: id, msg: refusalMessage(r.AdditionalProperties)}
		}

		return nil
	}

	for attempt := 1; ; attempt++ {
		err := deleteOnce()

		var refused deleteRefused
		if !errors.As(err, &refused) {
			if err != nil {
				return err
			}

			break
		}

		if attempt >= deleteAttempts {
			return refused
		}

		tflog.Warn(ctx, "network: Morpheus refused the delete; retrying",
			map[string]any{"id": id, "attempt": attempt, "msg": refused.msg})

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(deleteRetryAfter):
		}
	}

	deadline := time.Now().Add(timeout)

	var lastErr error

	for {
		status, message, gone, err := networkStatus(ctx, client, id)
		if err != nil {
			// Retried until the deadline: a momentary 502 is not a failed delete.
			tflog.Debug(ctx, "network: status read failed while waiting for deletion; retrying", map[string]any{"error": err.Error()})
			lastErr = err
		} else {
			lastErr = nil

			if gone {
				return nil
			}
		}

		if time.Now().After(deadline) {
			if lastErr != nil {
				return fmt.Errorf("network %d: could not confirm the deletion within %s: %w", id, timeout, lastErr)
			}

			return fmt.Errorf("network %d still exists %s after Morpheus reported it deleted (status %q%s)",
				id, timeout, status, detail(message))
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(settleInterval):
		}
	}
}

// refusalMessage assembles the message of a refused delete from the body's
// msg and errors, which the response model keeps as additional properties.
func refusalMessage(body map[string]any) string {
	var parts []string

	if m, ok := body["msg"].(string); ok && m != "" {
		parts = append(parts, m)
	}

	if errs, ok := body["errors"].(map[string]any); ok {
		keys := make([]string, 0, len(errs))
		for k := range errs {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			parts = append(parts, fmt.Sprintf("%s: %v", k, errs[k]))
		}
	}

	if len(parts) == 0 {
		return "no reason given"
	}

	return strings.Join(parts, "; ")
}

func detail(message string) string {
	if message == "" {
		return ""
	}

	return ": " + message
}
