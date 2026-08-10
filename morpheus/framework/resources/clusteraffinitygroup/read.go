// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *clusterAffinityGroupResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	// MORPH-15506 appliance version gate — see Create.
	resp.Diagnostics.Append(versioncheck.Require(
		ctx, client, gatedFeature, constants.AffinityGroupMinVersion,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ClusterAffinityGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// On import the prior state is empty, so the API response is the only source
	// of truth for tenant_ids and resource_permissions. On a normal refresh the
	// prior state is preserved for those two — see mapAndResolveResponse.
	isImport := state.Name.IsNull()

	clusterID := state.ClusterId.ValueInt64()
	id := state.Id.ValueInt64()

	result, httpResp, err := client.ClustersAPI.GetClusterAffinityGroup(ctx, clusterID, id).Execute()

	// CRITICAL BEHAVIOUR 4: Treat HTTP 404 as resource-gone.
	// Deleting the parent resource pool or cluster cascades and hard-deletes affinity groups.
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}

	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpRead, "cluster_affinity_group", "", err, httpResp,
		)

		return
	}

	ag := result.AffinityGroup
	if ag == nil {
		resp.Diagnostics.AddError("API returned nil", "AffinityGroup is nil in the response")

		return
	}

	// CRITICAL BEHAVIOUR 5: Refuse to manage sync-discovered groups.
	// If source == "sync", the group is owned by cloud sync (VMware DRS rule discovery)
	// and will be overwritten on every sync cycle.
	if ag.Source != nil && *ag.Source == "sync" {
		resp.Diagnostics.AddError(
			"Cannot manage sync-discovered affinity group",
			fmt.Sprintf(
				"Affinity group %d (cluster %d) has source=\"sync\" and is managed by "+
					"cloud sync. Use the hpe_morpheus_cluster_affinity_group data source instead.",
				id, clusterID,
			),
		)

		return
	}

	if isImport {
		resp.Diagnostics.Append(mapGetResponseToModel(ctx, &state, ag, clusterID)...)
	} else {
		resp.Diagnostics.Append(mapAndResolveResponse(ctx, &state, ag, clusterID)...)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// mapAndResolveResponse maps the read-back into the model and then resolves the two
// attributes that are NOT taken from the response, so that neither can leave an
// unknown behind.
//
// MORPH-7311: tenant_ids and resource_permissions are Optional+Computed. When the
// practitioner does not configure them the framework marks them UNKNOWN in the plan,
// and this resource deliberately preserves the planned value rather than the mapped
// one (see resolveTenantIds for why). Preserving an UNKNOWN put an unknown into
// post-apply state, which Terraform rejects with "Provider returned invalid result
// object after apply". UseStateForUnknown cannot rescue this on Create because there
// is no prior state to borrow from.
//
// The model passed in must still carry the planned (Create/Update) or prior
// (Read) values for those two attributes; they are read off the model before the
// mapping overwrites them.
func mapAndResolveResponse(
	ctx context.Context,
	model *ClusterAffinityGroupModel,
	ag *sdk.GetClusterAffinityGroup200ResponseAffinityGroup,
	clusterID int64,
) diag.Diagnostics {
	plannedTenantIds := model.TenantIds
	plannedRP := model.ResourcePermissions

	diags := mapGetResponseToModel(ctx, model, ag, clusterID)

	model.TenantIds = resolveTenantIds(plannedTenantIds, model.TenantIds)

	rp, rpDiags := resolveResourcePermissions(ctx, plannedRP, model.ResourcePermissions)
	diags.Append(rpDiags...)
	model.ResourcePermissions = rp

	return diags
}

// resolveTenantIds decides what tenant_ids is worth in post-apply state.
//
// WHY THIS EXISTS — do not collapse it into "always take the API value".
//
// tenant_ids is Optional+Computed with a UseStateForUnknown plan modifier. The API
// normalises and reorders the tenant list it echoes back, so writing the API value
// over a set the practitioner configured produces a perpetual diff. Preserving the
// configured value is the deliberate fix for that, and the plan modifier assumes it.
//
// But the preserve-verbatim rule is only sound while the planned value is KNOWN.
// When the practitioner has not configured tenant_ids at all, the plan holds UNKNOWN,
// and preserving that wrote an unknown into post-apply state — always a provider bug
// as far as Terraform is concerned.
//
// Resolution:
//
//	planned KNOWN (incl. null)  -> the planned set   (no diff churn from normalisation)
//	planned UNKNOWN             -> the API set, or a typed NULL when the API reports no
//	                               tenants (`"tenants": []` on a fresh group)
//
// The result is never unknown.
func resolveTenantIds(planned, fromAPI types.Set) types.Set {
	if !planned.IsUnknown() {
		return planned
	}

	if fromAPI.IsNull() || fromAPI.IsUnknown() || len(fromAPI.Elements()) == 0 {
		// A typed null, not a zero value: the element type has to match the schema.
		return types.SetNull(types.Int64Type)
	}

	return fromAPI
}

// resolveResourcePermissions decides what resource_permissions is worth in post-apply
// state. Same rationale as resolveTenantIds: preserve what the practitioner configured
// so a normalised API echo cannot churn the diff, but never preserve an unknown.
//
// The nested members need the same treatment. all and groups are themselves
// Optional+Computed, so a partial configuration such as
//
//	resource_permissions = { all = true }
//
// leaves groups UNKNOWN inside an otherwise known object — the identical defect one
// level down. Each member is therefore resolved individually rather than the object
// being preserved wholesale.
//
// The result, including every nested value, is never unknown.
func resolveResourcePermissions(
	ctx context.Context,
	planned ResourcePermissionsValue,
	fromAPI ResourcePermissionsValue,
) (ResourcePermissionsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Not configured at all: the API response is the only source of truth, and a
	// typed null when it has nothing to say (Morpheus returns
	// `"resourcePermissions": null` for a group with no permissions).
	if planned.IsUnknown() {
		if fromAPI.IsUnknown() {
			return NewResourcePermissionsValueNull(), diags
		}

		return fromAPI, diags
	}

	// Null is already a perfectly good known value.
	if planned.IsNull() {
		return planned, diags
	}

	apiKnown := !fromAPI.IsNull() && !fromAPI.IsUnknown()

	all := planned.All
	if all.IsUnknown() {
		all = types.BoolNull()
		if apiKnown && !fromAPI.All.IsUnknown() {
			all = fromAPI.All
		}
	}

	apiGroups := types.SetNull(GroupsValue{}.Type(ctx))
	if apiKnown {
		apiGroups = fromAPI.Groups
	}

	groups, groupDiags := resolveGroups(ctx, planned.Groups, apiGroups)
	diags.Append(groupDiags...)

	return ResourcePermissionsValue{
		All:    all,
		Groups: groups,
		state:  attr.ValueStateKnown,
	}, diags
}

// resolveGroups resolves resource_permissions.groups, and the Optional+Computed `id`
// and `default` inside each element, down to known values.
//
// An element whose `default` was omitted from the configuration comes back UNKNOWN in
// the plan, so it is filled from the API group carrying the same id, falling back to
// null when the API says nothing about it. Elements are correlated by id because this
// is a set and position is meaningless.
func resolveGroups(ctx context.Context, planned, fromAPI types.Set) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	nullGroups := types.SetNull(GroupsValue{}.Type(ctx))
	apiKnown := !fromAPI.IsNull() && !fromAPI.IsUnknown()

	if planned.IsUnknown() {
		if !apiKnown {
			return nullGroups, diags
		}

		return fromAPI, diags
	}

	if planned.IsNull() {
		return planned, diags
	}

	apiDefaults := make(map[int64]types.Bool)
	if apiKnown {
		for _, elem := range fromAPI.Elements() {
			g, ok := elem.(GroupsValue)
			if !ok || g.Id.IsNull() || g.Id.IsUnknown() {
				continue
			}
			apiDefaults[g.Id.ValueInt64()] = g.Default
		}
	}

	resolved := make([]GroupsValue, 0, len(planned.Elements()))
	for _, elem := range planned.Elements() {
		g, ok := elem.(GroupsValue)
		if !ok || g.IsNull() || g.IsUnknown() {
			// A wholly unknown or unexpected element cannot be correlated, so defer
			// to the API rather than invent one.
			if apiKnown {
				return fromAPI, diags
			}

			return nullGroups, diags
		}

		id := g.Id
		if id.IsUnknown() {
			id = types.Int64Null()
		}

		def := g.Default
		if def.IsUnknown() {
			def = types.BoolNull()
			if !id.IsNull() {
				if apiDefault, found := apiDefaults[id.ValueInt64()]; found &&
					!apiDefault.IsUnknown() {
					def = apiDefault
				}
			}
		}

		resolved = append(resolved, GroupsValue{
			Id:      id,
			Default: def,
			state:   attr.ValueStateKnown,
		})
	}

	groupsSet, setDiags := types.SetValueFrom(ctx, GroupsValue{}.Type(ctx), resolved)
	diags.Append(setDiags...)
	if diags.HasError() {
		return nullGroups, diags
	}

	return groupsSet, diags
}

// mapGetResponseToModel maps the SDK Get response to the Terraform model.
func mapGetResponseToModel(
	ctx context.Context,
	model *ClusterAffinityGroupModel,
	ag *sdk.GetClusterAffinityGroup200ResponseAffinityGroup,
	clusterID int64,
) diag.Diagnostics {
	var diags diag.Diagnostics

	model.ClusterId = types.Int64Value(clusterID)

	if ag.Id != nil {
		model.Id = types.Int64Value(*ag.Id)
	}
	if ag.Name != nil {
		model.Name = types.StringValue(*ag.Name)
	}

	// active is Optional+Computed, so the plan holds UNKNOWN whenever the
	// practitioner left it out. Falling straight through the nil guard would carry
	// that unknown into post-apply state, which Terraform rejects, so resolve it.
	if ag.Active != nil {
		model.Active = types.BoolValue(*ag.Active)
	} else if model.Active.IsUnknown() {
		model.Active = types.BoolNull()
	}

	model.AffinityType = convert.StrToType(ag.AffinityType)
	model.Visibility = convert.StrToType(ag.Visibility)
	model.Source = convert.StrToType(ag.Source)

	// CRITICAL BEHAVIOUR 6: Pool is computed-only for clusters.
	if ag.Pool != nil && ag.Pool.Id != nil {
		model.PoolId = types.Int64Value(*ag.Pool.Id)
	} else {
		model.PoolId = types.Int64Null()
	}

	// CRITICAL BEHAVIOUR 3: servers write/read asymmetry.
	// We send []int64 of ComputeServer IDs; response returns [{id, name}].
	// Map response objects back to a Set of Int64 IDs.
	serverVals := make([]attr.Value, 0, len(ag.Servers))
	for _, s := range ag.Servers {
		if s.Id != nil {
			serverVals = append(serverVals, types.Int64Value(*s.Id))
		}
	}
	serverSet, setDiags := types.SetValue(types.Int64Type, serverVals)
	diags.Append(setDiags...)
	model.Servers = serverSet

	// Tenants.
	//
	// No tenants is represented as a typed NULL, not an empty set, so that the two
	// paths that reach this attribute agree. Create resolves an unconfigured
	// (UNKNOWN) tenant_ids to null — see resolveTenantIds — and import maps
	// straight from the response, so an empty list mapped to an empty set here
	// would make an imported resource differ from the one that created it and fail
	// ImportStateVerify.
	tenantVals := make([]attr.Value, 0, len(ag.Tenants))
	for _, t := range ag.Tenants {
		if t.Id != nil {
			tenantVals = append(tenantVals, types.Int64Value(*t.Id))
		}
	}
	if len(tenantVals) == 0 {
		model.TenantIds = types.SetNull(types.Int64Type)
	} else {
		tenantSet, setDiags := types.SetValue(types.Int64Type, tenantVals)
		diags.Append(setDiags...)
		model.TenantIds = tenantSet
	}

	// ResourcePermissions.
	if ag.ResourcePermissions != nil {
		groupVals := make([]GroupsValue, 0, len(ag.ResourcePermissions.Sites))
		for _, s := range ag.ResourcePermissions.Sites {
			if s.Id != nil {
				var defaultVal types.Bool
				if s.Default != nil {
					defaultVal = types.BoolValue(*s.Default)
				} else {
					defaultVal = types.BoolNull()
				}
				groupVals = append(groupVals, GroupsValue{
					Id:      types.Int64Value(*s.Id),
					Default: defaultVal,
					state:   attr.ValueStateKnown,
				})
			}
		}
		groupsSet, d := types.SetValueFrom(ctx, GroupsValue{}.Type(ctx), groupVals)
		diags.Append(d...)
		rp, rpDiags := NewResourcePermissionsValue(
			map[string]attr.Type{
				"all":    types.BoolType,
				"groups": types.SetType{ElemType: GroupsValue{}.Type(ctx)},
			},
			map[string]attr.Value{
				"all":    types.BoolPointerValue(ag.ResourcePermissions.All),
				"groups": groupsSet,
			},
		)
		diags.Append(rpDiags...)
		model.ResourcePermissions = rp
	} else {
		model.ResourcePermissions = NewResourcePermissionsValueNull()
	}

	return diags
}
