// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package certificate_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
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

			return resp.GetCertificates(), hresp, err
		},
		// Is this a test certificate?
		func(item sdk.ListCertificates200ResponseCertificatesInner) bool {
			name, ok := item.GetNameOk()
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
			id, ok := item.GetIdOk()
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
