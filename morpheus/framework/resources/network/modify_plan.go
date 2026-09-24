// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/networkconfig"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/tenancy"
)

var _ resource.ResourceWithModifyPlan = &Resource{}

// ModifyPlan applies what the network type catalogue knows and what Terraform
// cannot see on its own.
//
// The type decides which changes Morpheus applies in place. Its flags and
// option types say which top-level fields an update would silently ignore;
// a change to one of those is planned as a replacement. The static modifiers
// this replaces treated every type alike.
//
// The write-only configuration never produces a diff: its prior and planned
// values are always null. The configuration is available here, though, and
// the digests of what was last applied are in private state. Comparing them:
//
//   - a change that touches a value Morpheus will not edit after create
//     requires replacement, since an update would be accepted and ignored;
//   - any other change marks the read-only config_from_api unknown, which is a
//     planned change, so Terraform runs Update and the new values are sent;
//   - no stored digests (a network created before configuration was write-only
//     whose refresh has not yet seeded them, an import, or -refresh=false) is
//     treated as a change to be applied once, never as grounds for replacement.
//
// Replacement is only ever planned when the configuration itself changed: a
// differing immutable digest with an unchanged overall digest means the
// catalogue's idea of what is immutable moved between provider versions, not
// that the user changed anything.
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	// Destroy: nothing to plan.
	if req.Plan.Raw.IsNull() {
		return
	}

	var config, state NetworkModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	creating := req.State.Raw.IsNull()

	if !creating {
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	typ := r.planTypeSpec(ctx, config, state)

	inputs, diags := resolveConfigInputs(ctx, config, typ)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if d := checkBlockMatchesType(inputs.block, typ); d != nil {
		resp.Diagnostics.Append(d)

		return
	}

	// Values the appliance's definition of the type does not declare. The
	// client is needed only to name the appliance's version in the message.
	client, _ := r.NewClient(ctx)

	resp.Diagnostics.Append(checkDeclared(ctx, client, inputs)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if creating {
		// Create stores the digests; nothing to compare yet. The visibility
		// rule, however, only bites on create.
		r.guardPublicVisibility(ctx, req, resp)

		return
	}

	if typ != nil {
		for _, attr := range networkconfig.ImmutableAttributes(*typ) {
			if knownAndDifferent(req.State.Raw, req.Plan.Raw, attr) {
				tflog.Debug(ctx, "network: attribute is not editable on this type; requiring replacement",
					map[string]any{"attribute": attr, "type": typ.Code})
				resp.RequiresReplace = append(resp.RequiresReplace, path.Root(attr))
			}
		}
	}

	stored, ok, ldiags := networkconfig.Load(ctx, req.Private)
	resp.Diagnostics.Append(ldiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	switch {
	case !ok:
		// Nothing recorded: apply the configuration once so it is, unless there
		// is nothing to send — then there is nothing to record either.
		if inputs.all == nil {
			return
		}

		tflog.Debug(ctx, "network: no configuration digests in private state; planning an update to apply and record the configuration")
		planConfigUpdate(ctx, req, resp)

	case stored.All == inputs.digests.All:
		// Unchanged. (Whatever the immutable digest says: see above.)

	case stored.Immutable != inputs.digests.Immutable:
		attr := "config"
		if inputs.block != nil {
			attr = inputs.block.Attr
		}

		tflog.Debug(ctx, "network: immutable configuration changed; requiring replacement", map[string]any{"attribute": attr})
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root(attr))

	default:
		tflog.Debug(ctx, "network: configuration changed; planning an update")
		planConfigUpdate(ctx, req, resp)
	}
}

// planConfigUpdate marks the attributes that mirror the config map unknown:
// config_from_api, which is the planned change that makes Terraform run
// Update, and the deprecated connected_gateway when it is not configured —
// it is read back from the config map, so a block that changes the gateway
// changes it too, and a plan that kept its prior value would be inconsistent
// with the result.
func planConfigUpdate(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("config_from_api"), types.DynamicUnknown())...)

	var configured types.String

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("connected_gateway"), &configured)...)
	if resp.Diagnostics.HasError() || !configured.IsNull() {
		return
	}

	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("connected_gateway"), types.StringUnknown())...)
}

// checkBlockMatchesType rejects a typed block set for a network of another
// type. With the type unknown at plan time there is nothing to check; Create
// checks again once the type is resolved.
func checkBlockMatchesType(block *networkconfig.Block, typ *networkconfig.TypeSpec) diag.Diagnostic {
	if block == nil || typ == nil || block.Serves(typ.Code) {
		return nil
	}

	detail := fmt.Sprintf("%s configures the network types %s; this network's type is %s (%s).",
		block.Attr, block.TypeNames(), typ.Code, typ.Name)

	if other, ok := networkconfig.BlockForType(typ.Code); ok {
		detail += fmt.Sprintf(" Use %s instead.", other.Attr)
	} else {
		detail += " Use the free-form config map for this type."
	}

	return diag.NewAttributeErrorDiagnostic(path.Root(block.Attr), "configuration block does not match the network type", detail)
}

// guardPublicVisibility applies, at plan time, the rule Morpheus enforces
// silently when a network is created: it stores "private" unless the
// network's tenant permissions include the master tenant. So visibility =
// "public" needs tenant_ids to include the caller's tenant — the sub-tenant
// case, where public is never possible, is the shared guard's. An unknown
// tenant_ids cannot be checked; the apply-time check reports it instead.
func (r *Resource) guardPublicVisibility(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var visibility types.String

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("visibility"), &visibility)...)
	if resp.Diagnostics.HasError() || visibility.IsNull() || visibility.IsUnknown() || visibility.ValueString() != "public" {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		return
	}

	tenancy.GuardVisibilityPublic(ctx, client, req.Plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// From the configuration, not the plan: unset is null there, whereas the
	// plan of an Optional+Computed attribute is unknown when it is unset.
	var tenantIDs types.Set

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("tenant_ids"), &tenantIDs)...)
	if resp.Diagnostics.HasError() || tenantIDs.IsUnknown() {
		return
	}

	caller, err := tenancy.CallerTenantID(ctx, client)
	if err != nil {
		tflog.Debug(ctx, "network: could not determine the caller's tenant; visibility not validated",
			map[string]any{"error": err.Error()})

		return
	}

	if !tenantIDs.IsNull() {
		var ids []types.Int64

		resp.Diagnostics.Append(tenantIDs.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, id := range ids {
			if id.IsUnknown() || (!id.IsNull() && id.ValueInt64() == caller) {
				return
			}
		}
	}

	resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(path.Root("visibility"),
		"visibility = \"public\" requires the master tenant in tenant_ids",
		fmt.Sprintf("Morpheus creates a network as \"private\" unless its tenant permissions include the master tenant. "+
			"Add the master tenant's id (%d) to tenant_ids, or set visibility = \"private\".", caller)))
}

// knownAndDifferent reports whether the named top-level attribute is known and
// non-null in both state and plan, and differs between them. Unknown or null
// on either side is not a change the provider can act on: an omitted
// computed attribute plans as unknown, and a value the API did not return is
// null in state.
func knownAndDifferent(state, plan tftypes.Value, attr string) bool {
	s, err := state.ApplyTerraform5AttributePathStep(tftypes.AttributeName(attr))
	if err != nil {
		return false
	}

	p, err := plan.ApplyTerraform5AttributePathStep(tftypes.AttributeName(attr))
	if err != nil {
		return false
	}

	sv, ok := s.(tftypes.Value)
	if !ok || !sv.IsKnown() || sv.IsNull() {
		return false
	}

	pv, ok := p.(tftypes.Value)
	if !ok || !pv.IsKnown() || pv.IsNull() {
		return false
	}

	return !sv.Equal(pv)
}
