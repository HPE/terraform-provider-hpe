// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package securitygrouprule implements a data source for security_group_rule
package securitygrouprule

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                         = "read security group rule data source"
	ErrorNoValidSearchTerms         = `no valid search terms - an id or name is required`
	ErrorNoSecurityGroupRuleFound   = `no security group rule found`
	ErrorMultipleSecurityGroupRules = `multiple security group rules were returned`
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "security_group_rule"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = SecurityGroupRuleDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config SecurityGroupRuleModel

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			"could not create sdk client",
		)

		return
	}

	rule, err := getRule(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	securityGroupID := config.SecurityGroupId.ValueInt64()
	state := ruleAsState(rule, securityGroupID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func ruleAsState(
	rule *sdk.GetSecurityGroupRules200ResponseRule,
	securityGroupID int64,
) SecurityGroupRuleModel {
	state := SecurityGroupRuleModel{
		Id:              convert.Int64ToType(rule.Id),
		SecurityGroupId: types.Int64Value(securityGroupID),
		RuleType:        convert.StrToType(rule.RuleType),
		CustomRule:      convert.BoolToType(rule.CustomRule),
		Direction:       convert.StrToType(rule.Direction),
		Policy:          convert.StrToType(rule.Policy),
		Protocol:        convert.StrToType(rule.Protocol),
		SourceType:      convert.StrToType(rule.SourceType),
		DestinationType: convert.StrToType(rule.DestinationType),
	}

	state.Name = convert.StrToType(rule.Name.Get())

	state.PortRange = convert.StrToType(rule.PortRange.Get())

	state.Source = convert.StrToType(rule.Source.Get())

	state.Destination = convert.StrToType(rule.Destination.Get())

	state.DestinationPortRange = convert.StrToType(rule.DestinationPortRange.Get())

	state.SourcePortRange = convert.StrToType(rule.SourcePortRange.Get())

	state.Enabled = convert.StrToType(rule.Enabled.Get())

	state.ExternalId = convert.StrToType(rule.ExternalId.Get())

	state.InstanceTypeId = convert.StrToType(rule.InstanceTypeId.Get())

	// Nested objects — nil-check parent pointer
	state.SourceGroup = sourceGroupValue(rule.SourceGroup)
	state.DestinationGroup = destinationGroupValue(rule.DestinationGroup)
	state.SourceTier = sourceTierValue(rule.SourceTier)
	state.DestinationTier = destinationTierValue(rule.DestinationTier)

	return state
}

func sourceGroupValue(in *sdk.GetSecurityGroupRules200ResponseRuleSourceGroup) SourceGroupValue {
	if in == nil {
		return NewSourceGroupValueNull()
	}

	return SourceGroupValue{
		Id:    convert.Int64ToType(in.Id),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}

func destinationGroupValue(in *sdk.GetSecurityGroupRules200ResponseRuleDestinationGroup) DestinationGroupValue {
	if in == nil {
		return NewDestinationGroupValueNull()
	}

	return DestinationGroupValue{
		Id:    convert.Int64ToType(in.Id),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}

func sourceTierValue(in *sdk.GetSecurityGroupRules200ResponseRuleSourceTier) SourceTierValue {
	if in == nil {
		return NewSourceTierValueNull()
	}

	return SourceTierValue{
		Id:    convert.Int64ToType(in.Id),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}

func destinationTierValue(in *sdk.GetSecurityGroupRules200ResponseRuleDestinationTier) DestinationTierValue {
	if in == nil {
		return NewDestinationTierValueNull()
	}

	return DestinationTierValue{
		Id:    convert.Int64ToType(in.Id),
		Name:  convert.StrToType(in.Name),
		state: attr.ValueStateKnown,
	}
}

func getRuleByID(
	ctx context.Context,
	ruleID int64,
	securityGroupID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetSecurityGroupRules200ResponseRule, error) {
	// GetSecurityGroupRules(ctx, id, sgId) where id=securityGroupId (parent), sgId=ruleId (child, float32)
	r, hresp, err := apiClient.SecurityGroupsAPI.GetSecurityGroupRules(
		ctx, securityGroupID, float32(ruleID),
	).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for security group rule %d: %s",
			ruleID, providererrors.ErrMsg(err, hresp),
		)
	}

	return r.Rule, nil
}

func getRuleByName(
	ctx context.Context,
	name string,
	securityGroupID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetSecurityGroupRules200ResponseRule, error) {
	// ListSecurityGroupRules(ctx, id) where id=securityGroupId (parent)
	rs, hresp, err := apiClient.SecurityGroupsAPI.ListSecurityGroupRules(
		ctx, securityGroupID,
	).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for security group rules with name %s: %s",
			name, providererrors.ErrMsg(err, hresp),
		)
	}

	items := rs.Rules
	if len(items) == 0 {
		return nil, errors.New(ErrorNoSecurityGroupRuleFound)
	}

	var matchedIDs []int64

	for i := range items {
		itemName := items[i].Name.Get()
		if itemName == nil || *itemName != name {
			continue
		}
		if items[i].Id == nil {
			continue
		}

		matchedIDs = append(matchedIDs, *items[i].Id)
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNoSecurityGroupRuleFound)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleSecurityGroupRules)
	}

	return getRuleByID(ctx, matchedIDs[0], securityGroupID, apiClient)
}

func getRule(
	ctx context.Context,
	config *SecurityGroupRuleModel,
	apiClient *sdk.APIClient,
) (*sdk.GetSecurityGroupRules200ResponseRule, error) {
	securityGroupID := config.SecurityGroupId.ValueInt64()

	if !config.Id.IsNull() {
		return getRuleByID(ctx, config.Id.ValueInt64(), securityGroupID, apiClient)
	} else if !config.Name.IsNull() {
		return getRuleByName(ctx, config.Name.ValueString(), securityGroupID, apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}
