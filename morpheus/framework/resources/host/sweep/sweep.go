// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/paging"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const sweeperName = "hpe_morpheus_host"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List hosts (compute servers) -- every page, since /api/servers lists
		// guests as well and the test hosts need not be on the first one. The
		// first page's response is what the harness judges (a 403/404 arrives
		// on it).
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListHosts200ResponseAllOfServersInner,
			*http.Response,
			error,
		) {
			var (
				firstResp *http.Response
				once      sync.Once
			)

			servers, err := paging.Collect(ctx, func(
				ctx context.Context,
				offset, max int64,
			) ([]sdk.ListHosts200ResponseAllOfServersInner, int64, error) {
				rs, hresp, err := client.HostsAPI.ListHosts(ctx).Max(max).Offset(offset).Execute()

				once.Do(func() { firstResp = hresp })

				if err != nil || rs == nil {
					return nil, 0, fmt.Errorf("ListHosts failed: %s", errfmt.ErrMsg(err, hresp))
				}

				total := int64(0)
				if rs.Meta != nil && rs.Meta.Total != nil {
					total = *rs.Meta.Total
				}

				return rs.Servers, total, nil
			})
			if err != nil {
				return nil, firstResp, err
			}

			// Only servers of a type this resource creates (the appliance's
			// creatable host types): a guest VM that happens to carry the test
			// prefix must not be removed.
			creatable, err := creatableCodes(ctx, client)
			if err != nil {
				return nil, firstResp, err
			}

			hosts := make([]sdk.ListHosts200ResponseAllOfServersInner, 0, len(servers))
			for _, srv := range servers {
				if srv.ComputeServerType != nil && srv.ComputeServerType.Code != nil && creatable[*srv.ComputeServerType.Code] {
					hosts = append(hosts, srv)
				}
			}

			return hosts, firstResp, nil
		},
		// Is this a test host?
		func(item sdk.ListHosts200ResponseAllOfServersInner) bool {
			name, ok := getsafe.GetOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test host. removeResources is omitted so Morpheus applies
		// its type-aware default (stub nodes carry nothing).
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListHosts200ResponseAllOfServersInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.HostsAPI.
				RemoveHost(ctx, sdk.UpdateHostIdParameter{Int64: id}).
				Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListHosts200ResponseAllOfServersInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}

// creatableCodes returns the codes of the server types the appliance lets a
// host be added as -- the only types hpe_morpheus_host creates.
func creatableCodes(ctx context.Context, client *sdk.APIClient) (map[string]bool, error) {
	types, err := paging.Collect(ctx, func(
		ctx context.Context,
		offset, max int64,
	) ([]sdk.ListHostTypes200ResponseAllOfServerTypesInner, int64, error) {
		rs, hresp, err := client.HostsAPI.ListHostTypes(ctx).Creatable(true).Max(max).Offset(offset).Execute()
		if err != nil || rs == nil {
			return nil, 0, fmt.Errorf("ListHostTypes failed: %s", errfmt.ErrMsg(err, hresp))
		}

		total := int64(0)
		if rs.Meta != nil && rs.Meta.Total != nil {
			total = *rs.Meta.Total
		}

		return rs.ServerTypes, total, nil
	})
	if err != nil {
		return nil, err
	}

	codes := map[string]bool{}
	for _, t := range types {
		if t.Code != nil {
			codes[*t.Code] = true
		}
	}

	return codes, nil
}
