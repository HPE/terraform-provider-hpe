package security_group_rule

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

var (
	_ resource.Resource                = &securityGroupRuleResource{}
	_ resource.ResourceWithConfigure   = &securityGroupRuleResource{}
	_ resource.ResourceWithImportState = &securityGroupRuleResource{}
)

type securityGroupRuleResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &securityGroupRuleResource{}
}

func (r *securityGroupRuleResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_security_group_rule"
}

func (r *securityGroupRuleResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = SecurityGroupRuleSchema(ctx)
}

func (r *securityGroupRuleResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan securityGroupRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sgID := plan.SecurityGroupID.ValueInt64()

	body := sdk.AddSecurityGroupRulesRequestRule{
		Protocol: plan.Protocol.ValueString(),
		RuleType: plan.RuleType.ValueString(),
	}
	if !plan.Name.IsNull() {
		body.Name = plan.Name.ValueStringPointer()
	}
	if !plan.Direction.IsNull() {
		body.Direction = plan.Direction.ValueStringPointer()
	}
	if !plan.Source.IsNull() {
		body.Source = plan.Source.ValueStringPointer()
	}
	if !plan.SourceType.IsNull() {
		body.SourceType = plan.SourceType.ValueStringPointer()
	}
	if !plan.Destination.IsNull() {
		body.Destination = plan.Destination.ValueStringPointer()
	}
	if !plan.DestinationType.IsNull() {
		body.DestinationType = plan.DestinationType.ValueStringPointer()
	}
	if !plan.PortRange.IsNull() {
		body.PortRange = plan.PortRange.ValueStringPointer()
	}
	if !plan.Policy.IsNull() {
		body.Policy = plan.Policy.ValueStringPointer()
	}

	result, httpResp, err := client.SecurityGroupsAPI.AddSecurityGroupRules(ctx, sgID).
		AddSecurityGroupRulesRequest(sdk.AddSecurityGroupRulesRequest{
			Rule: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "security_group_rule", "", err, httpResp)

		return
	}

	var id int64
	if createRule := result.GetRule(); createRule.Id != nil {
		id = *createRule.Id
	}
	ruleIDParam := float32(id) //nolint:gosec // value range is safe

	readResult, httpResp, err := client.SecurityGroupsAPI.GetSecurityGroupRules(ctx, sgID, ruleIDParam).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "security_group_rule", "", err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "security_group_rule",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	readRule := readResult.GetRule()
	mapResponseToModel(&plan, &readRule)
	plan.SecurityGroupID = types.Int64Value(sgID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *securityGroupRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state securityGroupRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sgID := state.SecurityGroupID.ValueInt64()
	ruleID := float32(state.ID.ValueInt64())

	result, httpResp, err := client.SecurityGroupsAPI.GetSecurityGroupRules(ctx, sgID, ruleID).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "security_group_rule", "", err, httpResp)

		return
	}

	rule := result.Rule
	if rule == nil {
		resp.Diagnostics.AddError("API returned nil", "Rule is nil in the response")

		return
	}
	mapResponseToModel(&state, rule)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *securityGroupRuleResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan securityGroupRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sgID := plan.SecurityGroupID.ValueInt64()
	ruleID := float32(plan.ID.ValueInt64())

	body := sdk.UpdateSecurityGroupRulesRequestRule{
		Protocol: plan.Protocol.ValueString(),
		RuleType: plan.RuleType.ValueString(),
	}
	if !plan.Name.IsNull() {
		body.Name = plan.Name.ValueStringPointer()
	}
	if !plan.Direction.IsNull() {
		body.Direction = plan.Direction.ValueStringPointer()
	}
	if !plan.Source.IsNull() {
		body.Source = plan.Source.ValueStringPointer()
	}
	if !plan.SourceType.IsNull() {
		body.SourceType = plan.SourceType.ValueStringPointer()
	}
	if !plan.Destination.IsNull() {
		body.Destination = plan.Destination.ValueStringPointer()
	}
	if !plan.DestinationType.IsNull() {
		body.DestinationType = plan.DestinationType.ValueStringPointer()
	}
	if !plan.PortRange.IsNull() {
		body.PortRange = plan.PortRange.ValueStringPointer()
	}
	if !plan.Policy.IsNull() {
		body.Policy = plan.Policy.ValueStringPointer()
	}

	_, httpResp, err := client.SecurityGroupsAPI.UpdateSecurityGroupRules(ctx, sgID, ruleID).
		UpdateSecurityGroupRulesRequest(sdk.UpdateSecurityGroupRulesRequest{
			Rule: body,
		}).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "security_group_rule", "", err, httpResp)

		return
	}

	readResult, httpResp, err := client.SecurityGroupsAPI.GetSecurityGroupRules(ctx, sgID, ruleID).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "security_group_rule", "", err, httpResp)

		return
	}

	readRule := readResult.GetRule()
	mapResponseToModel(&plan, &readRule)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *securityGroupRuleResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state securityGroupRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sgID := state.SecurityGroupID.ValueInt64()
	ruleID := float32(state.ID.ValueInt64())

	_, httpResp, err := client.SecurityGroupsAPI.RemoveSecurityGroupRules(ctx, sgID, ruleID).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "security_group_rule", "", err, httpResp)

		return
	}
}

func (r *securityGroupRuleResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parts := strings.Split(req.ID, ".")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid ID",
			fmt.Sprintf("Expected import ID in format 'security_group_id.rule_id', got %q", req.ID))

		return
	}

	sgID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID",
			fmt.Sprintf("Could not parse security_group_id %q as integer: %s", parts[0], err))

		return
	}
	ruleID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse rule_id %q as integer: %s", parts[1], err))

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("security_group_id"), sgID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), ruleID)...)
}

func mapCreateResponseToModel(model *securityGroupRuleModel, rule *sdk.AddSecurityGroupRules200ResponseAllOfRule) {
	if rule.Id != nil {
		model.ID = types.Int64Value(*rule.Id)
	}
	if rule.Name.IsSet() && rule.Name.Get() != nil {
		model.Name = types.StringValue(*rule.Name.Get())
	} else {
		model.Name = types.StringNull()
	}
	if rule.RuleType != nil {
		model.RuleType = types.StringValue(*rule.RuleType)
	}
	if rule.Direction != nil {
		model.Direction = types.StringValue(*rule.Direction)
	}
	if rule.Policy != nil {
		model.Policy = types.StringValue(*rule.Policy)
	}
	if rule.Protocol != nil {
		model.Protocol = types.StringValue(*rule.Protocol)
	}
	if rule.SourceType != nil {
		model.SourceType = types.StringValue(*rule.SourceType)
	}
	if rule.Source.IsSet() && rule.Source.Get() != nil {
		model.Source = types.StringValue(*rule.Source.Get())
	} else {
		model.Source = types.StringNull()
	}
	if rule.DestinationType != nil {
		model.DestinationType = types.StringValue(*rule.DestinationType)
	}
	if rule.Destination.IsSet() && rule.Destination.Get() != nil {
		model.Destination = types.StringValue(*rule.Destination.Get())
	} else {
		model.Destination = types.StringNull()
	}
	if rule.PortRange.IsSet() && rule.PortRange.Get() != nil {
		model.PortRange = types.StringValue(*rule.PortRange.Get())
	} else {
		model.PortRange = types.StringNull()
	}
}

func mapResponseToModel(model *securityGroupRuleModel, rule *sdk.GetSecurityGroupRules200ResponseRule) {
	if rule.Id != nil {
		model.ID = types.Int64Value(*rule.Id)
	}
	if rule.Name.IsSet() && rule.Name.Get() != nil {
		model.Name = types.StringValue(*rule.Name.Get())
	} else {
		model.Name = types.StringNull()
	}
	if rule.RuleType != nil {
		model.RuleType = types.StringValue(*rule.RuleType)
	}
	if rule.Direction != nil {
		model.Direction = types.StringValue(*rule.Direction)
	}
	if rule.Policy != nil {
		model.Policy = types.StringValue(*rule.Policy)
	}
	if rule.Protocol != nil {
		model.Protocol = types.StringValue(*rule.Protocol)
	}
	if rule.SourceType != nil {
		model.SourceType = types.StringValue(*rule.SourceType)
	}
	if rule.Source.IsSet() && rule.Source.Get() != nil {
		model.Source = types.StringValue(*rule.Source.Get())
	} else {
		model.Source = types.StringNull()
	}
	if rule.DestinationType != nil {
		model.DestinationType = types.StringValue(*rule.DestinationType)
	}
	if rule.Destination.IsSet() && rule.Destination.Get() != nil {
		model.Destination = types.StringValue(*rule.Destination.Get())
	} else {
		model.Destination = types.StringNull()
	}
	if rule.PortRange.IsSet() && rule.PortRange.Get() != nil {
		model.PortRange = types.StringValue(*rule.PortRange.Get())
	} else {
		model.PortRange = types.StringNull()
	}
}

func mapUpdateResponseToModel(model *securityGroupRuleModel, rule *sdk.UpdateSecurityGroupRules200ResponseAllOfRule) {
	if rule.Id != nil {
		model.ID = types.Int64Value(*rule.Id)
	}
	if rule.Name.IsSet() && rule.Name.Get() != nil {
		model.Name = types.StringValue(*rule.Name.Get())
	} else {
		model.Name = types.StringNull()
	}
	if rule.RuleType != nil {
		model.RuleType = types.StringValue(*rule.RuleType)
	}
	if rule.Direction != nil {
		model.Direction = types.StringValue(*rule.Direction)
	}
	if rule.Policy != nil {
		model.Policy = types.StringValue(*rule.Policy)
	}
	if rule.Protocol != nil {
		model.Protocol = types.StringValue(*rule.Protocol)
	}
	if rule.SourceType != nil {
		model.SourceType = types.StringValue(*rule.SourceType)
	}
	if rule.Source.IsSet() && rule.Source.Get() != nil {
		model.Source = types.StringValue(*rule.Source.Get())
	} else {
		model.Source = types.StringNull()
	}
	if rule.DestinationType != nil {
		model.DestinationType = types.StringValue(*rule.DestinationType)
	}
	if rule.Destination.IsSet() && rule.Destination.Get() != nil {
		model.Destination = types.StringValue(*rule.Destination.Get())
	} else {
		model.Destination = types.StringNull()
	}
	if rule.PortRange.IsSet() && rule.PortRange.Get() != nil {
		model.PortRange = types.StringValue(*rule.PortRange.Get())
	} else {
		model.PortRange = types.StringNull()
	}
}
