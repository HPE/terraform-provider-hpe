// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package library_spec_template_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_library_spec_template"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List library spec template resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListSpecTemplates200ResponseAllOfSpecTemplatesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.LibraryAPI.ListSpecTemplates(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetSpecTemplates(), hresp, err
		},
		// Is this a test library spec template?
		func(item sdk.ListSpecTemplates200ResponseAllOfSpecTemplatesInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test library spec template.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListSpecTemplates200ResponseAllOfSpecTemplatesInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LibraryAPI.DeleteSpecTemplate(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListSpecTemplates200ResponseAllOfSpecTemplatesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
