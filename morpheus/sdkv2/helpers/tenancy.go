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

// The tenancy determination is cached per client. In a Terraform run a provider
// process serves exactly one provider configuration (Terraform launches a
// separate plugin process per provider block, aliases included), so this still
// costs one whoami call per process. The in-process acceptance harness is
// different: it can hold several provider instances with different credentials
// in one process, and a process-wide answer would then be attributed to the
// wrong caller. Keying by client keeps each answer with the credentials that
// produced it. The mechanism otherwise mirrors morpheus/utils/tenancy (whoami
// behind a mutex; the two muxed stacks keep separate state so they stay
// uncoupled): only a successful determination is cached, and the mutex is held
// across the whoami call (single flight), so concurrent callers on a cold cache
// wait for the first determination and retry in turn on failure.
type caller struct {
	isMaster  bool
	accountID int64
}

var (
	cacheMu sync.Mutex
	callers = map[*morpheus.Client]caller{}
)

// resolveCaller returns the (cached) whoami determination for client. The
// caller must hold cacheMu.
func resolveCaller(client *morpheus.Client) (caller, error) {
	if client == nil {
		return caller{}, fmt.Errorf("no client available to determine tenancy")
	}

	if c, ok := callers[client]; ok {
		return c, nil
	}

	resp, err := client.Whoami()
	if err != nil {
		return caller{}, err
	}

	who, ok := resp.Result.(*morpheus.WhoamiResult)
	if !ok || who == nil {
		return caller{}, fmt.Errorf("whoami returned an unexpected response")
	}

	c := caller{isMaster: who.IsMasterAccount}
	if who.User != nil {
		c.accountID = who.User.AccountID
	}
	callers[client] = c

	return c, nil
}

// CallerIsMaster reports whether the legacy client's caller is the master
// tenant, via whoami. whoami introspects the current user, so it is reachable
// by any authenticated caller including a subtenant.
//
// The legacy WhoamiResult.IsMasterAccount is a non-pointer bool, so an omitted
// isMasterAccount decodes to false (non-master), which is correct.
func CallerIsMaster(client *morpheus.Client) (bool, error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	c, err := resolveCaller(client)
	if err != nil {
		return false, err
	}

	return c.isMaster, nil
}

// CallerAccountID returns the ID of the tenant (account) the legacy client's
// caller belongs to, via whoami. Resources use it to recognise their owning
// tenant in API responses that list it alongside explicitly shared tenants.
func CallerAccountID(client *morpheus.Client) (int64, error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	c, err := resolveCaller(client)
	if err != nil {
		return 0, err
	}

	if c.accountID == 0 {
		return 0, fmt.Errorf("whoami did not report the caller's account")
	}

	return c.accountID, nil
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
