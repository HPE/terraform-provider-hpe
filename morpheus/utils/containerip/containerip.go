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
// On timeout the function returns warned=true rather than an error, because
// the instance provisioned successfully — only the address is not yet
// available. The real reasons an address never arrives are provisioning
// failure, a network with neither DHCP nor static assignment, or the
// accessor throwing and yielding an empty string.
func WaitAny(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
	timeout time.Duration,
) (warned bool, err error) {
	poll := func() (struct{}, error) {
		getResp, hresp, getErr := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
		if getErr != nil {
			return struct{}{}, backoff.Permanent(
				fmt.Errorf("failed to read instance %d: %s",
					instanceID, errfmt.ErrMsg(getErr, hresp)),
			)
		}

		if getResp.Instance == nil {
			return struct{}{}, backoff.Permanent(
				fmt.Errorf("instance %d: response is nil", instanceID),
			)
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
// On timeout the function returns a warning-level message rather than an
// error, because the node was provisioned successfully — only the address
// is missing. The real reasons an address never arrives are provisioning
// failure, a network with neither DHCP nor static assignment, or the
// accessor throwing and yielding an empty string.
func Wait(
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
				if cd.Ip != nil && Ready(*cd.Ip) {
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
