// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

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

// IPReady reports whether an IP address represents a ready container.
// Ready iff trim(ip) is not empty and not "0.0.0.0".
//
// This mirrors Container.getContainerExternalIp() in core, which returns
// the raw address string. No IPv4 validation is performed so that IPv6
// addresses pass correctly.
func IPReady(ip string) bool {
	trimmed := strings.TrimSpace(ip)

	return trimmed != "" && trimmed != "0.0.0.0"
}

// WaitForContainerIP polls GET /api/instances/{id} until the container
// identified by containerID has an IP that satisfies IPReady, or until the
// timeout expires. Returns the IP on success.
//
// On timeout the function returns a warning-level message rather than an error,
// because the node was provisioned successfully — only the address is missing.
func WaitForContainerIP(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
	containerID int64,
	timeout time.Duration,
) (ip string, warned bool, err error) {
	poll := func() (string, error) {
		getResp, hresp, getErr := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
		if getErr != nil {
			return "", backoff.Permanent(
				fmt.Errorf("failed to read instance %d: %s",
					instanceID, errfmt.ErrMsg(getErr, hresp)),
			)
		}

		if getResp.Instance == nil {
			return "", backoff.Permanent(
				fmt.Errorf("instance %d: response is nil", instanceID),
			)
		}

		for i := range getResp.Instance.ContainerDetails {
			cd := &getResp.Instance.ContainerDetails[i]
			if cd.Id != nil && *cd.Id == containerID {
				if cd.Ip != nil && IPReady(*cd.Ip) {
					return *cd.Ip, nil
				}

				// Container found but IP not ready — retry.
				return "", fmt.Errorf("container %d: ip not ready", containerID)
			}
		}

		// Container not found at all — permanent error.
		return "", backoff.Permanent(
			fmt.Errorf("container %d not found on instance %d",
				containerID, instanceID),
		)
	}

	ip, err = backoff.Retry(
		ctx,
		poll,
		backoff.WithBackOff(backoff.NewConstantBackOff(10*time.Second)),
		backoff.WithMaxElapsedTime(timeout),
	)
	if err != nil {
		// On timeout, warn rather than fail — the node is provisioned.
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
