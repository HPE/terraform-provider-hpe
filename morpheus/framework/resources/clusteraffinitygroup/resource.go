// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// gatedFeature names this resource in the appliance version gate diagnostic.
// Phrased as a plural noun so the message reads "Cluster affinity groups
// require ...".
const gatedFeature = "Cluster affinity groups"

var (
	_ resource.Resource                = &clusterAffinityGroupResource{}
	_ resource.ResourceWithConfigure   = &clusterAffinityGroupResource{}
	_ resource.ResourceWithImportState = &clusterAffinityGroupResource{}
)

type clusterAffinityGroupResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &clusterAffinityGroupResource{}
}

func (r *clusterAffinityGroupResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "cluster_affinity_group"
}

func (r *clusterAffinityGroupResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = ClusterAffinityGroupResourceSchema(ctx)
}

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

	// CRITICAL BEHAVIOUR 6: pool is COMPUTED ONLY for clusters — the API force-assigns it.
	// Never send it on create or update; only read it back.

	// Servers — send as []int32 (SDK type).
	if !plan.Servers.IsNull() && !plan.Servers.IsUnknown() {
		var serverIDs []int64
		resp.Diagnostics.Append(plan.Servers.ElementsAs(ctx, &serverIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		servers32 := make([]int32, 0, len(serverIDs))
		for _, id := range serverIDs {
			servers32 = append(servers32, int32(id))
		}
		ag.Servers = servers32
	}

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
			rp.Sites = buildSitesPayload(groups)
		}
		ag.ResourcePermissions = &rp
	}

	body := sdk.SaveClusterAffinityGroupRequest{
		AffinityGroup: &ag,
	}

	// Tenants are sent as the request-root `tenantPermissions` wrapper
	// (`{"accounts": [<id>, ...]}`), NOT as the nested `affinityGroup.tenants`
	// (`[{"id": <id>}]`) form. This is DELIBERATE — do not "harmonise" it with the
	// cloud affinity group resource, which legitimately uses the nested form.
	//
	// AffinityGroupService (morpheus-core) resolves tenant permissions as:
	//
	//	params.tenantPermissions ?: params.tenantPermission
	//	    ?: params.affinityGroup?.tenantPermissions ?: params.affinityGroup?.tenantPermission
	//	    ?: params.affinityGroup?.tenants
	//
	// so the request-root form is checked FIRST, at the highest precedence. Both forms
	// are accepted, but they have different shapes and are parsed by different branches
	// of permissionService.parseTenantPermissions. This resource has shipped sending the
	// request-root form, so switching would change the on-the-wire payload for existing
	// users with no benefit.
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.TenantPermissions = &sdk.SaveClusterAffinityGroupRequestTenantPermissions{
			Accounts: tenantIDs,
		}
	}

	result, httpResp, err := client.ClustersAPI.SaveClusterAffinityGroup(ctx, clusterID).
		SaveClusterAffinityGroupRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpCreate, "cluster_affinity_group",
			plan.Name.ValueString(), err, httpResp,
		)

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
	readResult, httpResp, err := client.ClustersAPI.GetClusterAffinityGroup(ctx, clusterID, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpRead, "cluster_affinity_group",
			plan.Name.ValueString(), err, httpResp,
		)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "cluster_affinity_group",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	readAg := readResult.AffinityGroup
	if readAg == nil {
		resp.Diagnostics.AddError("API returned nil", "AffinityGroup is nil in the response")

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

func (r *clusterAffinityGroupResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
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

	var plan, state ClusterAffinityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := plan.ClusterId.ValueInt64()
	id := plan.Id.ValueInt64()

	ag := sdk.UpdateClusterAffinityGroupRequestAffinityGroup{
		Name:   plan.Name.ValueStringPointer(),
		Active: plan.Active.ValueBoolPointer(),
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		ag.Visibility = plan.Visibility.ValueStringPointer()
	}

	// CRITICAL BEHAVIOUR 2: servers is WHOLESALE REPLACE on update.
	// resolveUpdateServers falls back to the membership recorded in STATE when the
	// planned value is unknown or null. Read its doc comment before changing this.
	servers, serverDiags := resolveUpdateServers(ctx, plan.Servers, state.Servers)
	resp.Diagnostics.Append(serverDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ag.Servers = servers

	// ResourcePermissions.
	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		rp := sdk.UpdateClusterAffinityGroupRequestAffinityGroupResourcePermissions{
			All: plan.ResourcePermissions.All.ValueBoolPointer(),
		}
		if !plan.ResourcePermissions.Groups.IsNull() &&
			!plan.ResourcePermissions.Groups.IsUnknown() {
			var groups []GroupsValue
			resp.Diagnostics.Append(
				plan.ResourcePermissions.Groups.ElementsAs(ctx, &groups, false)...,
			)
			if resp.Diagnostics.HasError() {
				return
			}
			rp.Sites = buildSitesPayload(groups)
		}
		ag.ResourcePermissions = &rp
	}

	body := sdk.UpdateClusterAffinityGroupRequest{
		AffinityGroup: &ag,
	}

	// Tenants: request-root `tenantPermissions` wrapper, matching Create. See the longer
	// note in Create — this form is checked first by AffinityGroupService and is the form
	// this shipped resource has always sent. Do not move it into `affinityGroup.tenants`.
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.TenantPermissions = &sdk.UpdateClusterAffinityGroupRequestTenantPermissions{
			Accounts: tenantIDs,
		}
	}

	_, httpResp, err := client.ClustersAPI.UpdateClusterAffinityGroup(ctx, clusterID, id).
		UpdateClusterAffinityGroupRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpUpdate, "cluster_affinity_group",
			plan.Name.ValueString(), err, httpResp,
		)

		return
	}

	// Read-back.
	readResult, httpResp, err := client.ClustersAPI.GetClusterAffinityGroup(ctx, clusterID, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpRead, "cluster_affinity_group",
			plan.Name.ValueString(), err, httpResp,
		)

		return
	}

	readAg := readResult.AffinityGroup
	if readAg == nil {
		resp.Diagnostics.AddError("API returned nil", "AffinityGroup is nil in the response")

		return
	}

	// CRITICAL BEHAVIOUR 5: Refuse to manage sync-discovered groups (on update read-back).
	if readAg.Source != nil && *readAg.Source == "sync" {
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

	resp.Diagnostics.Append(mapAndResolveResponse(ctx, &plan, readAg, clusterID)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *clusterAffinityGroupResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	// MORPH-15506 appliance version gate — see Create. Delete is gated too, so
	// the contract is simply "this provider does not manage affinity groups
	// below 8.0.10" rather than a half-state where destroy is the one operation
	// that still works. Read is gated anyway, so a refreshing destroy already
	// stops here; a practitioner holding an orphaned group on an older
	// appliance removes it with `terraform state rm` and the Morpheus UI.
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

	clusterID := state.ClusterId.ValueInt64()
	id := state.Id.ValueInt64()

	_, httpResp, err := client.ClustersAPI.DeleteClusterAffinityGroup(ctx, clusterID, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpDelete, "cluster_affinity_group", "", err, httpResp,
		)

		return
	}
}

func (r *clusterAffinityGroupResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parts := strings.Split(req.ID, ".")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid Import ID",
			fmt.Sprintf("Expected format: cluster_id.affinity_group_id, got: %s", req.ID))

		return
	}

	clusterID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not parse cluster_id %q: %s", parts[0], err),
		)

		return
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not parse affinity_group_id %q: %s", parts[1], err),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cluster_id"), clusterID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// Ensure unused imports are satisfied.
var _ *http.Response

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
	// We send []int32 of ComputeServer IDs; response returns [{id, name}].
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

// resolveUpdateServers decides what to send for affinityGroup.servers on update.
//
// WHY THIS EXISTS — do not collapse the unknown/null branch back into an empty array.
//
// The Morpheus API treats servers on update as a WHOLESALE REPLACE.
// AffinityGroupService.updateAffinityGroup diffs the supplied collection against the
// current members and evicts everything absent from it; its else-branch removes every
// member. Because Groovy truthiness makes both a missing key and an empty list falsy,
// OMITTING servers and SENDING [] are identical to the API: both wipe the group.
//
// servers is Optional+Computed in the generated schema (schema_gen.go) and carries no
// UseStateForUnknown plan modifier. So whenever a practitioner has not configured
// servers and changes any other attribute — a rename, for instance — Terraform marks
// servers UNKNOWN in the plan. Treating that unknown as "send []" silently destroyed
// the group's entire membership, with no error and no warning.
//
// Resolution:
//
//	plan known, non-empty  -> the planned set   (the intended wholesale replace)
//	plan known, EMPTY      -> []                (practitioner wrote `servers = []`;
//	                                            a genuine, explicit "remove all")
//	plan UNKNOWN or NULL   -> the set in STATE  (re-asserts current membership, so the
//	                                            wholesale replace is a no-op and the
//	                                            members survive)
//
// The known-empty case must stay distinguishable from null/unknown: that distinction is
// the entire point of this function.
//
// If neither plan nor state carries a known set, the current membership is unknowable.
// The API has no encoding for "leave unchanged", so nil is returned and the key is
// omitted rather than asserting a membership we cannot substantiate.
func resolveUpdateServers(
	ctx context.Context,
	planServers types.Set,
	stateServers types.Set,
) ([]int32, diag.Diagnostics) {
	var diags diag.Diagnostics

	source := planServers
	if planServers.IsNull() || planServers.IsUnknown() {
		source = stateServers
	}

	if source.IsNull() || source.IsUnknown() {
		return nil, diags
	}

	var serverIDs []int64
	diags.Append(source.ElementsAs(ctx, &serverIDs, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// make never returns nil, so a known-empty set yields []int32{} rather than a nil
	// slice. The SDK's ToMap only omits the key when the slice is nil, so this is what
	// makes an explicit `servers = []` serialise as "servers": [].
	servers32 := make([]int32, 0, len(serverIDs))
	for _, sid := range serverIDs {
		servers32 = append(servers32, int32(sid))
	}

	return servers32, diags
}

// buildSitesPayload converts GroupsValue slice to the SDK sites type.
func buildSitesPayload(
	groups []GroupsValue,
) []sdk.SaveClusterAffinityGroupRequestAffinityGroupResourcePermissionsSitesInner {
	sites := make(
		[]sdk.SaveClusterAffinityGroupRequestAffinityGroupResourcePermissionsSitesInner,
		0, len(groups),
	)
	for _, g := range groups {
		id := g.Id.ValueInt64()
		inner := sdk.SaveClusterAffinityGroupRequestAffinityGroupResourcePermissionsSitesInner{
			Id: &id,
		}
		if !g.Default.IsNull() && !g.Default.IsUnknown() {
			inner.Default = g.Default.ValueBoolPointer()
		}
		sites = append(sites, inner)
	}

	return sites
}
