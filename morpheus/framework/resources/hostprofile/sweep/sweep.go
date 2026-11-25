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

const sweeperName = "hpe_morpheus_host_profile"

// seedCode is the immutable, system-seeded HVM default profile which must
// never be swept.
const seedCode = "host-profile.hvm.default"

// sweepPageSize is the number of host profiles requested per list page.
const sweepPageSize int64 = 100

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List host profiles, paging through every result so a sweep is not
		// silently limited to the API's default first page.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListHostProfiles200ResponseHostProfilesInner,
			*http.Response,
			error,
		) {
			var (
				offset int64
				all    []sdk.ListHostProfiles200ResponseHostProfilesInner
			)

			for {
				resp, hresp, err := client.HostProfilesAPI.
					ListHostProfiles(ctx).
					Max(sweepPageSize).
					Offset(offset).
					Execute()
				// Surface list errors instead of silently
				// returning a partial (or empty) result.
				if err != nil {
					return all, hresp, err
				}
				if resp == nil {
					return all, hresp, nil
				}

				page := getsafe.Get(&resp.HostProfiles)
				all = append(all, page...)
				offset += int64(len(page))

				// A short page is the reliable end-of-results
				// signal: it holds whether or not the API
				// reports a total in Meta. The Meta.Total check
				// is only a secondary guard for a full final
				// page.
				if len(page) < int(sweepPageSize) {
					return all, hresp, nil
				}
				if resp.Meta != nil {
					total := int64(getsafe.Get(resp.Meta.Total))
					if total > 0 && offset >= total {
						return all, hresp, nil
					}
				}
			}
		},
		// Is this a sweepable test host profile?
		func(item sdk.ListHostProfiles200ResponseHostProfilesInner) bool {
			if code, ok := getsafe.GetOk(item.Code.Get()); ok &&
				code != nil && *code == seedCode {
				return false
			}

			name, ok := getsafe.GetOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test host profile.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListHostProfiles200ResponseHostProfilesInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.HostProfilesAPI.
				DeleteHostProfile(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListHostProfiles200ResponseHostProfilesInner](
			http.StatusNotFound,
		),
	)
}
