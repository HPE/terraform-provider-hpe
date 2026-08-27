// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package containerip

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// Ready reports whether an IP address represents a ready container.
// Ready iff trim(ip) is not empty and not "0.0.0.0".
//
// The sentinel values mirror Container.getContainerExternalIp() in core,
// which resolves container.externalIp → server.externalIp/internalIp →
// primaryInterface.publicIpAddress/ipAddress. The "0.0.0.0" sentinel is
// the default when the platform has not yet assigned an address. No IPv4
// validation is performed so that IPv6 addresses pass correctly.
func Ready(ip string) bool {
	trimmed := strings.TrimSpace(ip)

	return trimmed != "" && trimmed != "0.0.0.0"
}

// WaitAny polls GET /api/instances/{id} until at least one container on the
// instance reports an IP that satisfies Ready, or until the timeout expires.
//
// The instance resource owns the container it provisions; containers added
// later by hpe_morpheus_instance_node belong to that resource and are waited
// on there. Requiring all containers would make an instance apply block on
// nodes it does not own.
//
// Three outcomes:
//  1. Permanent error (API failure, nil response) — returned as a real error
//     so the caller can fail the apply and, where appropriate, taint state.
//  2. Context cancelled — returned as an error (the apply is being torn down).
//  3. Genuine timeout (max elapsed time reached with no usable address) —
//     warned=true, err=nil. The instance provisioned successfully but the
//     address is not yet available.
func WaitAny(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
	timeout time.Duration,
) (warned bool, err error) {
	var hardErr error // captured inside the closure on permanent failures

	poll := func() (struct{}, error) {
		getResp, hresp, getErr := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
		if getErr != nil {
			hardErr = fmt.Errorf("failed to read instance %d: %s",
				instanceID, errfmt.ErrMsg(getErr, hresp))

			return struct{}{}, backoff.Permanent(hardErr)
		}

		if getResp.Instance == nil {
			hardErr = fmt.Errorf("instance %d: response is nil", instanceID)

			return struct{}{}, backoff.Permanent(hardErr)
		}

		for i := range getResp.Instance.ContainerDetails {
			cd := &getResp.Instance.ContainerDetails[i]
			if cd.Ip != nil && Ready(*cd.Ip) {
				return struct{}{}, nil
			}
		}

		return struct{}{}, fmt.Errorf("instance %d: no container has a ready IP", instanceID)
	}

	_, err = backoff.Retry(
		ctx,
		poll,
		backoff.WithBackOff(backoff.NewConstantBackOff(10*time.Second)),
		backoff.WithMaxElapsedTime(timeout),
	)
	if err != nil {
		// Permanent API error — propagate so the caller fails the apply.
		if hardErr != nil {
			return false, hardErr
		}

		// Context cancelled — propagate.
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		// Genuine timeout — warn only.
		tflog.Warn(ctx, "IP address wait timed out; instance provisioned but address not yet available",
			map[string]any{
				"instance_id": instanceID,
				"error":       err.Error(),
			},
		)

		return true, nil
	}

	return false, nil
}

// Wait polls GET /api/instances/{id} until the container identified by
// containerID has an IP that satisfies Ready, or until the timeout expires.
// Returns the IP on success.
//
// Three outcomes:
//  1. Permanent error (API failure, nil response, container not found) —
//     returned as a real error so the caller can fail the apply and, where
//     appropriate, taint state.
//  2. Context cancelled — returned as an error (the apply is being torn down).
//  3. Genuine timeout (max elapsed time reached with no usable address) —
//     warned=true, err=nil. The node provisioned successfully but the
//     address is not yet available.
func Wait(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
	containerID int64,
	timeout time.Duration,
) (ip string, warned bool, err error) {
	var hardErr error // captured inside the closure on permanent failures

	poll := func() (string, error) {
		getResp, hresp, getErr := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
		if getErr != nil {
			hardErr = fmt.Errorf("failed to read instance %d: %s",
				instanceID, errfmt.ErrMsg(getErr, hresp))

			return "", backoff.Permanent(hardErr)
		}

		if getResp.Instance == nil {
			hardErr = fmt.Errorf("instance %d: response is nil", instanceID)

			return "", backoff.Permanent(hardErr)
		}

		for i := range getResp.Instance.ContainerDetails {
			cd := &getResp.Instance.ContainerDetails[i]
			if cd.Id != nil && *cd.Id == containerID {
				if cd.Ip != nil && Ready(*cd.Ip) {
					return *cd.Ip, nil
				}

				// Container found but IP not ready — retry.
				return "", fmt.Errorf("container %d: ip not ready", containerID)
			}
		}

		// Container not found at all — permanent error.
		hardErr = fmt.Errorf("container %d not found on instance %d",
			containerID, instanceID)

		return "", backoff.Permanent(hardErr)
	}

	ip, err = backoff.Retry(
		ctx,
		poll,
		backoff.WithBackOff(backoff.NewConstantBackOff(10*time.Second)),
		backoff.WithMaxElapsedTime(timeout),
	)
	if err != nil {
		// Permanent API error — propagate so the caller fails the apply.
		if hardErr != nil {
			return "", false, hardErr
		}

		// Context cancelled — propagate.
		if ctx.Err() != nil {
			return "", false, ctx.Err()
		}

		// Genuine timeout — warn only.
		tflog.Warn(ctx, "IP address wait timed out; node provisioned but address not yet available",
			map[string]any{
				"instance_id":  instanceID,
				"container_id": containerID,
				"error":        err.Error(),
			},
		)

		return "", true, nil
	}

	return ip, false, nil
}
