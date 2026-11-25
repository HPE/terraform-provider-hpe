// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const sweeperName = "hpe_morpheus_network_server"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network server resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListNetworkServers200ResponseAllOfNetworkServersInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.NetworksAPI.ListNetworkServers(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.Get(&resp.NetworkServers), hresp, err
		},
		// Is this a test network server?
		func(item sdk.ListNetworkServers200ResponseAllOfNetworkServersInner) bool {
			name, ok := getsafe.GetOk(item.Name)
			if !ok || name == nil {
				return false
			}

			if !strings.HasPrefix(*name, testsweep.TestResourcePrefix) {
				return false
			}
			// Never take the shared SDN fixture, whatever it is called: the
			// networks, routers and load balancers of its cloud depend on it.
			if id, ok := getsafe.GetOk(item.Id); ok && id != nil &&
				strconv.FormatInt(*id, 10) == os.Getenv(testhelpers.EnvSdnNetworkServerID) {
				return false
			}

			return true
		},
		// Delete the test network server.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListNetworkServers200ResponseAllOfNetworkServersInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			result, hresp, err := client.NetworksAPI.DeleteNetworkServer(ctx, *id).Execute()
			// A 200 with success=false is a refusal, not a deletion.
			if err := errfmt.CheckDeleted(result, err, hresp); err != nil {
				return hresp, err
			}

			return hresp, nil
		},
		testsweep.WithIgnoreListStatuses[sdk.ListNetworkServers200ResponseAllOfNetworkServersInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
