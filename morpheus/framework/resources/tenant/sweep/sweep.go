// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const sweeperName = "hpe_morpheus_tenant"

const (
	// deleteTimeout bounds how long the sweeper waits for one tenant's
	// asynchronous deletion to complete before moving on.
	deleteTimeout = 3 * time.Minute
	// deletePollInterval is how often the sweeper polls for the tenant to be
	// gone.
	deletePollInterval = 5 * time.Second
)

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List tenants -- all of them, deepest first.
		//
		// Test tenants form trees (a tenant's parent_id may be another test
		// tenant), and Morpheus cannot delete a tenant that still has children:
		// the final delete fails on the parent reference and the tenant is left
		// in a deleting_failed state. Ordering the sweep deepest-first, together
		// with waiting for each deletion to finish below, removes children
		// before their parents.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListTenants200ResponseAllOfAccountsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.TenantsAPI.ListTenants(ctx).
				Max(testsweep.ListPageSize).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			accounts := getsafe.Get(&resp.Accounts)
			sortDeepestFirst(accounts)

			return accounts, hresp, err
		},
		// Is this a test tenant?
		func(item sdk.ListTenants200ResponseAllOfAccountsInner) bool {
			name, ok := getsafe.GetOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test tenant, de-provisioning any managed resources, and
		// wait for the asynchronous deletion to complete so that a parent
		// swept next does not race its children's removal.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListTenants200ResponseAllOfAccountsInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.TenantsAPI.RemoveTenant(ctx, *id).
				RemoveResources(true).Execute()
			if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
				return hresp, err
			}

			return hresp, waitForTenantGone(ctx, client, *id)
		},
		testsweep.WithIgnoreListStatuses[sdk.ListTenants200ResponseAllOfAccountsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}

// sortDeepestFirst orders tenants by depth in the tenant hierarchy, deepest
// first, so children are swept before their parents. Depth is the number of
// parent hops to a tenant with no parent in the list (the master tenant, or a
// parent not returned). The sort is stable, so tenants at the same depth keep
// their listing order.
func sortDeepestFirst(accounts []sdk.ListTenants200ResponseAllOfAccountsInner) {
	parentOf := make(map[int64]int64, len(accounts))
	for _, a := range accounts {
		if a.Id == nil || a.Parent == nil || a.Parent.Id == nil {
			continue
		}
		parentOf[*a.Id] = *a.Parent.Id
	}

	depth := func(a sdk.ListTenants200ResponseAllOfAccountsInner) int {
		if a.Id == nil {
			return 0
		}

		d := 0
		for id, ok := parentOf[*a.Id]; ok && d < len(accounts); id, ok = parentOf[id] {
			d++
		}

		return d
	}

	sort.SliceStable(accounts, func(i, j int) bool {
		return depth(accounts[i]) > depth(accounts[j])
	})
}

// waitForTenantGone polls the tenant until it returns 404 or deleteTimeout
// elapses. Tenant deletion is a background job that the API acknowledges
// immediately, so without waiting the next sweep target -- possibly this
// tenant's parent -- would be deleted while this one still exists. Any error
// other than the 404 that signals completion ends the wait at once, so a
// failing appliance is reported rather than polled until the deadline.
func waitForTenantGone(ctx context.Context, client *sdk.APIClient, id int64) error {
	deadline := time.Now().Add(deleteTimeout)

	for {
		_, hresp, err := client.TenantsAPI.GetTenant(ctx, id).Execute()
		if errfmt.IsNotFound(hresp) {
			return nil
		}
		if err := errfmt.CheckResponse(err, hresp); err != nil {
			return fmt.Errorf("waiting for tenant %d to be deleted: %w", id, err)
		}

		if time.Now().After(deadline) {
			return fmt.Errorf(
				"timed out after %s waiting for tenant %d to be deleted; the background "+
					"deletion job may still be running (or have failed) on the appliance",
				deleteTimeout, id)
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting for tenant %d to be deleted: %w", id, ctx.Err())
		case <-time.After(deletePollInterval):
		}
	}
}
