// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package helpers

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"
)

// masterTenantCache caches, per legacy client pointer, whether that client's
// caller is the master tenant. Only successful determinations are cached.
var masterTenantCache sync.Map // map[*morpheus.Client]bool

// CallerIsMaster reports whether the legacy client's caller is the master
// tenant, via whoami. whoami introspects the current user, so it is reachable
// by any authenticated caller including a subtenant. Successful results are
// cached keyed by the client pointer.
//
// The legacy WhoamiResult.IsMasterAccount is a non-pointer bool, so an omitted
// isMasterAccount decodes to false (non-master), which is correct.
func CallerIsMaster(client *morpheus.Client) (bool, error) {
	if v, ok := masterTenantCache.Load(client); ok {
		return v.(bool), nil
	}

	resp, err := client.Whoami()
	if err != nil {
		return false, err
	}

	who, ok := resp.Result.(*morpheus.WhoamiResult)
	if !ok || who == nil {
		return false, fmt.Errorf("whoami returned an unexpected response")
	}

	masterTenantCache.Store(client, who.IsMasterAccount)

	return who.IsMasterAccount, nil
}

// VisibilityCustomizeDiff rejects, at plan time, visibility = "public" set by a
// non-master caller. Morpheus silently coerces visibility to "private" for
// sub-tenant callers, so failing fast avoids state disagreeing with config.
//
// It only calls whoami when visibility is a known "public" value. On a whoami
// error it fails open (logs and returns nil) so an unreachable appliance does
// not block offline planning.
func VisibilityCustomizeDiff(
	ctx context.Context,
	d *schema.ResourceDiff,
	meta any,
) error {
	if !d.NewValueKnown("visibility") {
		return nil
	}
	if d.Get("visibility") != "public" {
		return nil
	}

	client, ok := meta.(*morpheus.Client)
	if !ok || client == nil {
		return nil
	}

	isMaster, err := CallerIsMaster(client)
	if err != nil {
		tflog.Warn(ctx, "could not determine tenancy for visibility validation",
			map[string]any{"error": err.Error()})

		return nil
	}

	if !isMaster {
		return fmt.Errorf(
			"visibility = \"public\" requires the master tenant: Morpheus silently " +
				"stores \"private\" for objects created or updated by sub-tenant " +
				"users; set visibility = \"private\" or run as a master-tenant user")
	}

	return nil
}
