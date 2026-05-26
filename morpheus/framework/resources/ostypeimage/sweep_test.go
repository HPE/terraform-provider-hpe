// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package ostypeimage_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_os_type_image"

type osTypeImageSweepItem struct {
	id               int64
	virtualImageName string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List OS type image resources by iterating OS types.
		func(ctx context.Context, client *sdk.APIClient) (
			[]osTypeImageSweepItem,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.LibraryAPI.ListOsTypes(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			items := make([]osTypeImageSweepItem, 0)
			for _, osType := range resp.GetOsTypes() {
				osTypeID, ok := osType.GetIdOk()
				if !ok || osTypeID == nil {
					continue
				}

				osTypeResp, osTypeHresp, osTypeErr := client.LibraryAPI.GetOsType(ctx, *osTypeID).Execute()
				if osTypeErr != nil || osTypeHresp == nil || osTypeHresp.StatusCode != http.StatusOK {
					if osTypeHresp != nil && (osTypeHresp.StatusCode == http.StatusNotFound ||
						osTypeHresp.StatusCode == http.StatusForbidden) {
						continue
					}

					return nil, osTypeHresp, osTypeErr
				}

				osTypeDetail := osTypeResp.GetOsType()

				for _, image := range osTypeDetail.GetImages() {
					id, ok := image.GetIdOk()
					if !ok || id == nil {
						continue
					}

					name, ok := image.GetVirtualImageNameOk()
					if !ok || name == nil {
						continue
					}

					items = append(items, osTypeImageSweepItem{
						id:               *id,
						virtualImageName: *name,
					})
				}
			}

			return items, hresp, nil
		},
		// Is this a test OS type image?
		func(item osTypeImageSweepItem) bool {
			return strings.HasPrefix(item.virtualImageName, testsweep.TestResourcePrefix)
		},
		// Delete the test OS type image.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item osTypeImageSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LibraryAPI.DeleteOsTypeImage(ctx, item.id).Execute()

			return hresp, err
		},
	)
}
