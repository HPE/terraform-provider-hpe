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

const sweeperName = "hpe_morpheus_certificate"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List certificate resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListCertificates200ResponseCertificatesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.SSLCertificatesAPI.ListCertificates(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.GetSafe(&resp.Certificates), hresp, err
		},
		// Is this a test certificate?
		func(item sdk.ListCertificates200ResponseCertificatesInner) bool {
			name, ok := getsafe.GetSafeOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test certificate.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListCertificates200ResponseCertificatesInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetSafeOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.SSLCertificatesAPI.DeleteCertificate(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListCertificates200ResponseCertificatesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
