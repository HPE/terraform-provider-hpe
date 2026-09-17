// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// callerIsMaster reports whether the authenticated tenant is the master tenant,
// via whoami. whoami introspects the current user (not an arbitrary account), so
// it is reachable by any authenticated caller including a subtenant. It returns
// an error when the status could not be determined so the caller can surface it
// as a diagnostic.
//
// The whoami response omits isMasterAccount unless it is true (see
// views/morpheus_api/whoami/index.gson), so a nil value means non-master, not
// "unknown".
// callerIsMaster reports whether the authenticated tenant is the master tenant.
// It delegates to the shared tenancy helper on the embedded configure struct so
// the whoami determination is cached and shared across the provider.
func (r *Resource) callerIsMaster(ctx context.Context) (bool, error) {
	return r.CallerIsMaster(ctx)
}

// validateParentIdRequiresMaster warns, at plan time, when parent_id is set by a
// non-master caller. The tenant create endpoint only binds a nominated parent
// for master callers; for everyone else it silently forces the parent to the
// caller's own tenant, so a parent_id set by a subtenant is a no-op.
func validateParentIdRequiresMaster(
	parentID int64,
	isMaster bool,
	diags *diag.Diagnostics,
) {
	if isMaster {
		return
	}

	diags.AddAttributeWarning(
		path.Root("parent_id"),
		"parent_id is ignored for non-master tenants",
		fmt.Sprintf(
			"parent_id is only honoured when authenticated as the master tenant. "+
				"This provider is authenticated as a non-master tenant, so parent_id "+
				"(%d) is ignored and the tenant will be created under your own tenant.",
			parentID),
	)
}

// validateAvailableParentTenantRoles warns, at plan time, when base_role_id will
// not resolve to a usable role under the effective parent tenant. It mirrors the
// Morpheus tenant create/update role logic:
//
//   - base_role_id set: the role must be one of the effective parent's available
//     account roles, otherwise the API rejects with "Invalid role id".
//   - base_role_id omitted: Morpheus assigns the master tenant's default account
//     role when the effective parent is the master tenant (which may not be the
//     role the user expects), and otherwise rejects with "Role is required".
//
// The effective parent depends on the caller: only the master tenant may
// nominate a parent, so for a non-master caller the effective parent is always
// the caller's own tenant (parent_id is ignored). The assignable-roles lookup is
// therefore routed accordingly:
//
//   - master  -> GET /api/accounts/{parent_id}/available-roles (any tenant)
//   - non-master -> GET /api/accounts/available-roles (the caller's own tenant)
//
// Routing the non-master case to the no-id endpoint avoids the 404 the tenant-
// scoped endpoint returns for a subtenant querying its own account.
//
// Best-effort: on any client or lookup failure it does not block the plan.
func (r *Resource) validateAvailableParentTenantRoles(
	ctx context.Context,
	parentID int64,
	baseRole types.Int64,
	isMaster bool,
	diags *diag.Diagnostics,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		return
	}

	// base_role_id set: check it against the effective parent's base roles.
	if !baseRole.IsNull() && !baseRole.IsUnknown() {
		want := baseRole.ValueInt64()

		roles, ok := r.availableParentRoles(ctx, client, parentID, isMaster, diags)
		if !ok {
			return
		}

		for _, role := range roles {
			if role.id == want {
				return
			}
		}

		available := formatAvailableRoles(roles)

		// parent_id is only honoured for the master tenant; for a non-master
		// caller the roles belong to the caller's own tenant (the nominated
		// parent_id is ignored), so word the message accordingly.
		subject := "your tenant"
		if isMaster {
			subject = fmt.Sprintf("parent tenant %d", parentID)
		}

		diags.AddAttributeWarning(
			path.Root("base_role_id"),
			"base_role_id is not assignable for "+subject,
			fmt.Sprintf(
				"role %d is not an assignable base role for %s; applying this "+
					"configuration will fail with \"Invalid role id\". The roles "+
					"assignable for %s are: %s.",
				want, subject, subject, available),
		)

		return
	}

	// base_role_id omitted: the outcome depends on whether the effective parent
	// is the master tenant.
	if !isMaster {
		// The effective parent is the caller's own (non-master) tenant.
		diags.AddAttributeWarning(
			path.Root("base_role_id"),
			"base_role_id omitted; apply will fail",
			"base_role_id is omitted and your tenant is not the master tenant, so "+
				"applying this configuration will fail with \"Role is required\". "+
				"Set base_role_id to an assignable base role.",
		)

		return
	}

	tresp, hresp, err := client.TenantsAPI.GetTenant(ctx, parentID).Execute()
	if err := errfmt.CheckResponse(err, hresp); err != nil {
		return
	}
	if tresp == nil || tresp.Account == nil {
		return
	}

	if tresp.Account.Master != nil && *tresp.Account.Master {
		diags.AddAttributeWarning(
			path.Root("base_role_id"),
			"base_role_id omitted; the master default role will be assigned",
			fmt.Sprintf(
				"base_role_id is omitted and parent tenant %d is the master tenant, "+
					"so Morpheus will assign the master tenant's default account role. "+
					"This may not be the role you expect; set base_role_id explicitly "+
					"to choose the tenant's base role.",
				parentID),
		)

		return
	}

	diags.AddAttributeWarning(
		path.Root("base_role_id"),
		"base_role_id omitted; apply will fail",
		fmt.Sprintf(
			"base_role_id is omitted and parent tenant %d is not the master tenant, "+
				"so applying this configuration will fail with \"Role is required\". "+
				"Set base_role_id to an assignable base role for the parent tenant.",
			parentID),
	)
}

// availableRole is an account role assignable as a tenant base role. It carries
// the name (authority) alongside the id so the validator can both match on id
// and tell the user which roles are available.
type availableRole struct {
	id   int64
	name string
}

// formatAvailableRoles renders roles as "name" (1), "name2" (2). A role with no
// name falls back to (N); an empty list renders as "(none)".
func formatAvailableRoles(roles []availableRole) string {
	if len(roles) == 0 {
		return "(none)"
	}

	parts := make([]string, 0, len(roles))
	for _, role := range roles {
		if role.name != "" {
			parts = append(parts, fmt.Sprintf("%q (%d)", role.name, role.id))
		} else {
			parts = append(parts, fmt.Sprintf("(%d)", role.id))
		}
	}

	return strings.Join(parts, ", ")
}

// availableParentRoles returns the account roles assignable under the effective
// parent tenant, choosing the endpoint the caller can reach. The second return
// is false when the lookup failed (a warning is emitted and the caller should
// stop).
func (r *Resource) availableParentRoles(
	ctx context.Context,
	client *sdk.APIClient,
	parentID int64,
	isMaster bool,
	diags *diag.Diagnostics,
) ([]availableRole, bool) {
	if isMaster {
		resp, hresp, err := client.TenantsAPI.ListTenantAvailableRoles(ctx, parentID).Execute()
		if err := errfmt.CheckResponse(err, hresp); err != nil {
			diags.AddAttributeWarning(
				path.Root("base_role_id"),
				"Could not validate base_role_id at plan time",
				fmt.Sprintf(
					"Unable to list assignable base roles for parent tenant %d; "+
						"Morpheus will enforce the role on apply. Error: %s",
					parentID, errfmt.ErrMsg(err, hresp)),
			)

			return nil, false
		}
		if resp == nil {
			return nil, false
		}

		roles := make([]availableRole, 0, len(resp.Roles))
		for _, role := range resp.Roles {
			if role.Id == nil {
				continue
			}
			name := ""
			if role.Authority != nil {
				name = *role.Authority
			}
			roles = append(roles, availableRole{id: *role.Id, name: name})
		}

		return roles, true
	}

	// Non-master caller: the effective parent is the caller's own tenant, which
	// is exactly what the no-id available-roles endpoint returns.
	resp, hresp, err := client.TenantsAPI.ListTenantsAvailableRoles(ctx).Execute()
	if err := errfmt.CheckResponse(err, hresp); err != nil {
		diags.AddAttributeWarning(
			path.Root("base_role_id"),
			"Could not validate base_role_id at plan time",
			"Unable to list assignable base roles for your tenant; Morpheus will "+
				"enforce the role on apply. Error: "+errfmt.ErrMsg(err, hresp),
		)

		return nil, false
	}
	if resp == nil {
		return nil, false
	}

	roles := make([]availableRole, 0, len(resp.Roles))
	for _, role := range resp.Roles {
		if role.Id == nil {
			continue
		}
		name := ""
		if role.Authority != nil {
			name = *role.Authority
		}
		roles = append(roles, availableRole{id: *role.Id, name: name})
	}

	return roles, true
}
