// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const sweeperName = "hpe_morpheus_provisioning_license"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List provisioning license resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListProvisioningLicenses200ResponseAllOfLicensesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.ProvisioningLicensesAPI.ListProvisioningLicenses(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.GetSafe(&resp.Licenses), hresp, err
		},
		// Is this a test provisioning license?
		func(item sdk.ListProvisioningLicenses200ResponseAllOfLicensesInner) bool {
			name, ok := getsafe.GetSafeOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test provisioning license.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListProvisioningLicenses200ResponseAllOfLicensesInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetSafeOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.ProvisioningLicensesAPI.RemoveProvisioningLicense(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListProvisioningLicenses200ResponseAllOfLicensesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
