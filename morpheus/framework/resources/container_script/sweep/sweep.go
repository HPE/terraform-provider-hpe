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

const sweeperName = "hpe_morpheus_container_script"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List library container script resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListScripts200ResponseAllOfContainerScriptsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.LibraryAPI.ListScripts(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetContainerScripts(), hresp, err
		},
		// Is this a test library container script?
		func(item sdk.ListScripts200ResponseAllOfContainerScriptsInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test library container script.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListScripts200ResponseAllOfContainerScriptsInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LibraryAPI.DeleteScript(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListScripts200ResponseAllOfContainerScriptsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
