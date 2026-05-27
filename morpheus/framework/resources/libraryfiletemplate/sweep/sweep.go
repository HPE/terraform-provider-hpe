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
)

const sweeperName = "hpe_morpheus_library_file_template"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List library file template resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListFileTemplates200ResponseAllOfContainerTemplatesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.LibraryAPI.ListFileTemplates(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetContainerTemplates(), hresp, err
		},
		// Is this a test library file template?
		func(item sdk.ListFileTemplates200ResponseAllOfContainerTemplatesInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test library file template.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListFileTemplates200ResponseAllOfContainerTemplatesInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LibraryAPI.DeleteFileTemplate(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListFileTemplates200ResponseAllOfContainerTemplatesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
