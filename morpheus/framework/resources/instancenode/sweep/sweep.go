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

const sweeperName = "hpe_morpheus_instance_node"

// TestResourcePrefix is used by acceptance tests for resource naming so
// sweepers can identify and clean up test resources.
const TestResourcePrefix = "TestAccMorpheusInstanceNode"

func init() {
	// Instance nodes are identified by their parent instance's name.
	// We list instances, filter test ones, then remove their extra containers.
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List instances.
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

			return strings.HasPrefix(*name, TestResourcePrefix)
		},
		// Delete the test instance (which removes its nodes).
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListInstances200ResponseAllOfInstancesInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.InstancesAPI.DeleteInstance(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListInstances200ResponseAllOfInstancesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
