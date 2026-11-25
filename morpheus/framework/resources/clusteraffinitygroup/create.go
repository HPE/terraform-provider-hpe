// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

func (r *clusterAffinityGroupResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	// MORPH-15506: refuse to operate against an appliance older than the first
	// release with stable affinity group semantics, so the practitioner gets a
	// diagnostic naming the required version instead of an opaque API error.
	//
	// The check sits at the top of each CRUD method rather than in Configure:
	// the framework calls Configure on every RPC for the type, including
	// ValidateResourceConfig and UpgradeResourceState, and neither should have
	// to reach the network. See versioncheck.Require for the full rationale,
	// including why an unreadable version fails open.
	resp.Diagnostics.Append(versioncheck.Require(
		ctx, client, gatedFeature, constants.AffinityGroupMinVersion,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan ClusterAffinityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := plan.ClusterId.ValueInt64()

	ag := sdk.SaveClusterAffinityGroupRequestAffinityGroup{
		Name: plan.Name.ValueStringPointer(),
	}

	// CRITICAL BEHAVIOUR 1: active MUST ALWAYS be sent on create.
	// The API does `active = (params.active == 'on' || params.active == true)` unconditionally,
	// which overrides the domain default of true. If omitted, the group is created INACTIVE.
	if plan.Active.IsNull() || plan.Active.IsUnknown() {
		active := true
		ag.Active = &active
	} else {
		ag.Active = plan.Active.ValueBoolPointer()
	}

	// CRITICAL BEHAVIOUR 7: affinity_type is create-only (absent from update model).
	if !plan.AffinityType.IsNull() && !plan.AffinityType.IsUnknown() {
		ag.AffinityType = plan.AffinityType.ValueStringPointer()
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		ag.Visibility = plan.Visibility.ValueStringPointer()
	}

	// CRITICAL BEHAVIOUR 8: servers is sent on create when it is configured.
	//
	// Configured, the practitioner owns the membership and the group is created
	// with exactly that set. Unconfigured, the attribute is unknown in the plan,
	// nothing is sent, and membership is left to
	// hpe_morpheus_cluster_affinity_group_member or to provisioning.
	//
	// Morpheus 9.1 refuses to create a group with no members ("An affinity rule
	// requires either a Servers list ... or both a VM Group and a Host Group"),
	// so against 9.1 a configuration that omits servers cannot create a group
	// at all; the hint added below says so where the API's own answer does not.
	// An empty set is sent as given -- 9.0 accepts it and 9.1 rejects it with
	// that same message -- rather than being second-guessed here.
	serversConfigured := !plan.Servers.IsNull() && !plan.Servers.IsUnknown()
	if serversConfigured {
		resp.Diagnostics.Append(plan.Servers.ElementsAs(ctx, &ag.Servers, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// CRITICAL BEHAVIOUR 6: pool is COMPUTED ONLY for clusters — the API force-assigns it.
	// Never send it on create or update; only read it back.

	// ResourcePermissions.
	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		rp := sdk.SaveClusterAffinityGroupRequestAffinityGroupResourcePermissions{
			All: plan.ResourcePermissions.All.ValueBoolPointer(),
		}
		if !plan.ResourcePermissions.Groups.IsNull() && !plan.ResourcePermissions.Groups.IsUnknown() {
			var groups []GroupsValue
			resp.Diagnostics.Append(plan.ResourcePermissions.Groups.ElementsAs(ctx, &groups, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			rp.Sites = buildSitesPayload(groups,
				func(id *int64, dflt *bool) sdk.UpdateAffinityGroupRequestAffinityGroupAllOfResourcePermissionsSitesInner {
					return sdk.UpdateAffinityGroupRequestAffinityGroupAllOfResourcePermissionsSitesInner{Id: id, Default: dflt}
				})
		}
		ag.ResourcePermissions = &rp
	}

	body := sdk.SaveClusterAffinityGroupRequest{
		AffinityGroup: &ag,
	}

	// tenant_ids is deliberately NOT sent. See MORPH-15806.
	//
	// The API accepts two shapes for this -- request-root `tenantPermissions`
	// ({"accounts": [<id>]}) and nested `affinityGroup.tenants` ([{"id": <id>}])
	// -- and parses both correctly. Neither works. Verified against 9.0.1:
	// supplying tenants on create answers 403 while still creating the group
	// with no tenants applied, supplying them on update answers 500, and either
	// one leaves that group's single-item GET returning 500 permanently. A
	// control group never given tenants reads back 200, so the damage is caused
	// by the request rather than being a property of the endpoint.
	//
	// Sending nothing keeps groups readable. The attribute stays in the schema,
	// carrying a deprecation message, so existing configurations continue to
	// plan; the value is state-only and does not reflect the appliance.

	result, httpResp, err := client.ClustersAPI.SaveClusterAffinityGroup(ctx, clusterID).
		SaveClusterAffinityGroupRequest(body).Execute()

	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpCreate, "cluster_affinity_group",
			plan.Name.ValueString(), err, httpResp,
		)

		// A second diagnostic, beside the API's own answer, when that answer
		// is Morpheus 9.1 refusing a group with no members and the
		// configuration sent none.
		if !serversConfigured && refusedForWantOfMembers(ctx, client, err, httpResp) {
			resp.Diagnostics.AddError(noMembersSummary, noMembersHint)
		}

		return
	}

	if result.AffinityGroup == nil || result.AffinityGroup.Id == nil {
		resp.Diagnostics.AddError(
			"API returned nil ID", "AffinityGroup ID is nil in the create response",
		)

		return
	}

	id := *result.AffinityGroup.Id

	// Read-back to populate full state.
	readAg, found := fetchAffinityGroup(ctx, client, clusterID, id, &resp.Diagnostics)
	if !found {
		if !resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError(
				"affinity group not found after create",
				"The affinity group was created but could not be read back.",
			)
		}

		// The group exists in Morpheus but never reached state. Taint so the
		// next apply replaces it rather than leaving it orphaned and invisible
		// to Terraform.
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "cluster_affinity_group",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	resp.Diagnostics.Append(mapAndResolveResponse(ctx, &plan, readAg, clusterID)...)

	// The create response is authoritative for the ID, so fall back to it if the
	// read-back somehow omitted it. id is Computed and therefore UNKNOWN in the
	// plan, and an unknown left in post-apply state is rejected outright.
	if plan.Id.IsUnknown() {
		plan.Id = types.Int64Value(id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// membersRequiredSince is the first Morpheus release that refuses to create an
// affinity group with no members (AffinityGroup.isValidRuleShape).
const membersRequiredSince = ">= 9.1.0"

// noMembersSummary and noMembersHint explain a create that Morpheus refused
// for want of members, in the terms the practitioner can act on.
const (
	noMembersSummary = "The affinity group needs its first members"
	noMembersHint    = "Morpheus 9.1 and later refuse to create an affinity group " +
		"with no members, and some releases report that refusal as a 403 rather than " +
		"as the rule-shape error. If these credentials can otherwise manage this " +
		"cluster's affinity groups, set servers so the group is created with its " +
		"first members; membership can then be managed there, or one server at a " +
		"time with hpe_morpheus_cluster_affinity_group_member."
)

// refusalKind classifies a failed create: Morpheus's rule-shape refusal
// (its errors map carries affinityGroup.ruleShape.invalid), a bare 403, or
// neither.
type refusalKind int

const (
	refusalOther refusalKind = iota
	refusalRuleShape
	refusalForbidden
)

func classifyRefusal(err error, httpResp *http.Response) refusalKind {
	if err != nil && strings.Contains(err.Error(), "ruleShape.invalid") {
		return refusalRuleShape
	}

	if httpResp != nil && httpResp.StatusCode == http.StatusForbidden {
		return refusalForbidden
	}

	return refusalOther
}

// refusedForWantOfMembers reports whether a failed create looks like Morpheus
// 9.1's refusal of a group with no members. The rule-shape error is that
// refusal wherever it appears. A bare 403 is it only on 9.1 and later: with
// servers absent from the request, 9.1.0 fails to render the refusal and
// answers "You do not have permissions to access this api endpoint" instead
// -- but before 9.1 there is no such refusal, so a 403 there is a genuine
// permission denial and gets no hint. An appliance whose version cannot be
// read gets none either.
func refusedForWantOfMembers(ctx context.Context, client *sdk.APIClient, err error, httpResp *http.Response) bool {
	switch classifyRefusal(err, httpResp) {
	case refusalRuleShape:
		return true
	case refusalForbidden:
		atLeast, verr := versioncheck.AtLeast(ctx, client, membersRequiredSince)

		return verr == nil && atLeast
	default:
		return false
	}
}
