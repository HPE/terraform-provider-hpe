// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP


//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_instance"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List instance resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListInstances200ResponseAllOfInstancesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.InstancesAPI.ListInstances(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetInstances(), hresp, err
		},
		// Is this a test instance?
		func(item sdk.ListInstances200ResponseAllOfInstancesInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test instance.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			instance sdk.ListInstances200ResponseAllOfInstancesInner,
		) (*http.Response, error) {
			serverIDs, err := getServerIDs(instance)
			if err != nil {
				return nil, err
			}

			for _, serverID := range serverIDs {
				stopID := sdk.UpdateHostIdParameter{Int64: &serverID}
				_, hresp, err := client.HostsAPI.StopHost(ctx, stopID).Execute()
				if err != nil || (hresp.StatusCode != http.StatusOK && hresp.StatusCode != http.StatusConflict) {
					return hresp, err
				}
			}

			for _, serverID := range serverIDs {
				deleteID := sdk.UpdateHostIdParameter{Int64: &serverID}
				_, hresp, err := client.HostsAPI.RemoveHost(ctx, deleteID).Force("on").
					RemoveResources("on").RemoveInstances("on").Execute()
				if err != nil || hresp.StatusCode != http.StatusOK {
					return hresp, err
				}
			}

			return &http.Response{StatusCode: http.StatusOK}, nil
		},
	)
}

func getServerIDs(instance sdk.ListInstances200ResponseAllOfInstancesInner) ([]int64, error) {
	containers, ok := instance.GetContainerDetailsOk()
	if !ok || containers == nil {
		return nil, fmt.Errorf("failed to get container details")
	}

	serverIDs := make([]int64, 0, len(containers))
	for _, container := range containers {
		server, ok := container.GetServerOk()
		if !ok || server == nil {
			return nil, fmt.Errorf("failed to get server details from container")
		}

		serverID, ok := server.GetIdOk()
		if !ok || serverID == nil {
			return nil, fmt.Errorf("failed to get server ID from container")
		}

		serverIDs = append(serverIDs, *serverID)
	}

	return serverIDs, nil
}
