// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data LoadBalancerProfileModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	state, diags := getLoadBalancerProfileAsState(
		ctx, data.LoadBalancerId.ValueInt64(), data.Id.ValueInt64(), client, data,
	)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getLoadBalancerProfileAsState(
	ctx context.Context,
	loadBalancerID int64,
	id int64,
	client *sdk.APIClient,
	prior LoadBalancerProfileModel,
) (LoadBalancerProfileModel, diag.Diagnostics) {
	var state LoadBalancerProfileModel
	var diags diag.Diagnostics

	profileResp, httpResp, err := client.LoadBalancersAPI.
		GetLoadBalancerProfile(ctx, loadBalancerID, id).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			"error reading load balancer profile",
			fmt.Sprintf("load balancer profile %d GET failed: ", id)+
				errfmt.ErrMsg(err, httpResp),
		)

		return state, diags
	}

	p := profileResp.LoadBalancerProfile
	if p == nil {
		diags.AddError(
			"error reading load balancer profile",
			fmt.Sprintf("load balancer profile %d GET returned no profile", id),
		)

		return state, diags
	}

	// Core fields
	state.Id = convert.Int64ToType(p.Id)
	state.LoadBalancerId = types.Int64Value(loadBalancerID)
	state.Name = convert.StrToType(p.Name)
	state.Description = convert.StrToType(p.Description)
	state.ServiceType = convert.StrToType(p.ServiceType)
	state.ServiceTypeDisplay = convert.StrToType(p.ServiceTypeDisplay)
	state.Category = convert.StrToType(p.Category)
	state.Visibility = convert.StrToType(p.Visibility)
	state.InternalId = convert.StrToType(p.InternalId)
	state.ExternalId = convert.StrToType(p.ExternalId)
	state.Enabled = convert.BoolToType(p.Enabled)
	state.Editable = convert.BoolToType(p.Editable)
	state.InsertXforwardedFor = convert.BoolToType(p.InsertXforwardedFor)

	// NullableString fields
	if p.ProxyType.IsSet() {
		state.ProxyType = convert.StrToType(p.ProxyType.Get())
	} else {
		state.ProxyType = types.StringNull()
	}

	if p.RedirectRewrite.IsSet() {
		state.RedirectRewrite = convert.StrToType(p.RedirectRewrite.Get())
	} else {
		state.RedirectRewrite = types.StringNull()
	}

	if p.PersistenceType.IsSet() {
		state.PersistenceType = convert.StrToType(p.PersistenceType.Get())
	} else {
		state.PersistenceType = types.StringNull()
	}

	if p.SslEnabled.IsSet() {
		state.SslEnabled = convert.StrToType(p.SslEnabled.Get())
	} else {
		state.SslEnabled = types.StringNull()
	}

	if p.SslCert.IsSet() {
		state.SslCert = convert.StrToType(p.SslCert.Get())
	} else {
		state.SslCert = types.StringNull()
	}

	if p.AccountCertificate.IsSet() {
		state.AccountCertificate = convert.StrToType(p.AccountCertificate.Get())
	} else {
		state.AccountCertificate = types.StringNull()
	}

	if p.RedirectUrl.IsSet() {
		state.RedirectUrl = convert.StrToType(p.RedirectUrl.Get())
	} else {
		state.RedirectUrl = types.StringNull()
	}

	if p.PersistenceCookieName.IsSet() {
		state.PersistenceCookieName = convert.StrToType(p.PersistenceCookieName.Get())
	} else {
		state.PersistenceCookieName = types.StringNull()
	}

	if p.PersistenceExpiresIn.IsSet() {
		state.PersistenceExpiresIn = convert.StrToType(p.PersistenceExpiresIn.Get())
	} else {
		state.PersistenceExpiresIn = types.StringNull()
	}

	// Time fields
	if p.DateCreated != nil {
		state.DateCreated = types.StringValue(p.DateCreated.String())
	} else {
		state.DateCreated = types.StringNull()
	}

	if p.LastUpdated != nil {
		state.LastUpdated = types.StringValue(p.LastUpdated.String())
	} else {
		state.LastUpdated = types.StringNull()
	}

	// Nested LoadBalancer object
	lb, lbDiags := buildProfileLoadBalancerValue(ctx, p.LoadBalancer)
	if diags.Append(lbDiags...); diags.HasError() {
		return state, diags
	}

	state.LoadBalancer = lb

	// serviceType selects which config variant and tag list apply. It is read
	// straight from the API response, so it is available on import too.
	serviceType := state.ServiceType.ValueString()

	// Tags: read from config response
	state.Tags = readTagsFromConfig(ctx, serviceType, p.Config, prior.Tags)

	// Config blocks: preserve every value the practitioner configured, while
	// resolving any attribute that arrived unknown from the API response.
	//
	// Optional+Computed attributes omitted from the configuration are unknown in
	// the plan. Copying the prior block verbatim would write those unknowns into
	// state ("Provider returned invalid result object after apply"), whereas
	// rebuilding it purely from the response would overwrite configured values
	// with the API's own defaults (for example x_forwarded_for defaults to
	// INSERT server-side) and produce "Provider produced inconsistent result
	// after apply". Merging per attribute avoids both.
	mergeConfigBlocks(ctx, &state, prior, p.Config)

	// On import the prior state has no typed config block (ImportState only sets
	// id and load_balancer_id). Reconstruct the active block from the API
	// response so the imported resource is complete and produces a clean plan.
	if allConfigBlocksNull(prior) {
		reconstructConfigBlockFromResponse(ctx, &state, serviceType, p.Config)
	}

	return state, diags
}

func buildProfileLoadBalancerValue(
	ctx context.Context,
	lb *sdk.GetLoadBalancerProfile200ResponseLoadBalancerProfileLoadBalancer,
) (LoadBalancerValue, diag.Diagnostics) {
	if lb == nil {
		return NewLoadBalancerValueNull(), nil
	}

	// Build the nested type object
	typeVal := types.ObjectNull(TypeValue{}.AttributeTypes(ctx))

	if lbType := lb.Type; lbType != nil {
		tv, d := NewTypeValue(
			TypeValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"code": convert.StrToType(lbType.Code),
				"id":   convert.Int64ToType(lbType.Id),
				"name": convert.StrToType(lbType.Name),
			},
		)
		if d.HasError() {
			return NewLoadBalancerValueNull(), d
		}

		var tvDiags diag.Diagnostics
		typeVal, tvDiags = tv.ToObjectValue(ctx)
		if tvDiags.HasError() {
			return NewLoadBalancerValueNull(), tvDiags
		}
	}

	lbVal, d := NewLoadBalancerValue(
		LoadBalancerValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"id":   convert.Int64ToType(lb.Id),
			"ip":   convert.StrToType(lb.Ip),
			"name": convert.StrToType(lb.Name),
			"type": typeVal,
		},
	)
	if d.HasError() {
		return NewLoadBalancerValueNull(), d
	}

	return lbVal, nil
}
