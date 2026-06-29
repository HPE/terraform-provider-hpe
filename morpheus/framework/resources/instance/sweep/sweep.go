// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const sweeperName = "hpe_morpheus_instance"

// unmanagedServerSweeperName sweeps leftover discovered VMs (unmanaged servers)
// that acceptance tests create but that are not tied to a managed instance, so
// the instance sweeper above does not catch them.
const unmanagedServerSweeperName = "hpe_morpheus_instance_unmanaged_server"

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

			return getsafe.Get(&resp.Instances), hresp, err
		},
		// Is this a test instance?
		func(item sdk.ListInstances200ResponseAllOfInstancesInner) bool {
			name, ok := getsafe.GetOk(item.Name)
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

	testsweep.RegisterTypedAPISweeper(
		unmanagedServerSweeperName,
		// List unmanaged servers (discovered VMs not under Morpheus management).
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListHosts200ResponseAllOfServersInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.HostsAPI.ListHosts(ctx).
				Managed(false).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.Get(&resp.Servers), hresp, err
		},
		// Is this a test server? The API lowercases host names, so match the
		// lowercased test resource prefix ("testaccmorpheus") case-sensitively.
		func(item sdk.ListHosts200ResponseAllOfServersInner) bool {
			name, ok := getsafe.GetOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(
				*name, strings.ToLower(testsweep.TestResourcePrefix),
			)
		},
		// Stop and delete the unmanaged server.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			server sdk.ListHosts200ResponseAllOfServersInner,
		) (*http.Response, error) {
			if server.Id == nil {
				return nil, fmt.Errorf("failed to get server ID")
			}
			serverID := *server.Id

			stopID := sdk.UpdateHostIdParameter{Int64: &serverID}
			_, hresp, err := client.HostsAPI.StopHost(ctx, stopID).Execute()
			if err != nil || (hresp.StatusCode != http.StatusOK &&
				hresp.StatusCode != http.StatusConflict) {
				return hresp, err
			}

			deleteID := sdk.UpdateHostIdParameter{Int64: &serverID}
			_, hresp, err = client.HostsAPI.RemoveHost(ctx, deleteID).
				Force("on").RemoveResources("on").Execute()
			if err != nil || hresp.StatusCode != http.StatusOK {
				return hresp, err
			}

			return &http.Response{StatusCode: http.StatusOK}, nil
		},
	)
}

func getServerIDs(instance sdk.ListInstances200ResponseAllOfInstancesInner) ([]int64, error) {
	containers := instance.ContainerDetails
	if containers == nil {
		return nil, fmt.Errorf("failed to get container details")
	}

	serverIDs := make([]int64, 0, len(containers))
	for _, container := range containers {
		server := container.Server
		if server == nil {
			return nil, fmt.Errorf("failed to get server details from container")
		}

		serverID := server.Id
		if serverID == nil {
			return nil, fmt.Errorf("failed to get server ID from container")
		}

		serverIDs = append(serverIDs, *serverID)
	}

	return serverIDs, nil
}
