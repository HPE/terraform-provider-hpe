// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package tenants implements a plural data source for tenants.
package tenants

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/dsfilter"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary    = "read tenants data source"
	maxResults = 10000
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
	resp.TypeName = req.ProviderTypeName + "_" + "tenants"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = TenantsDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config TenantsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filters := compileFilters(ctx, config.Filter, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	// Sort ascending by default. sort_ascending has no schema default (data
	// sources cannot declare one), so a null value means ascending -- matching
	// the previous sdkv2 Default:true behavior. Only an explicit false sorts
	// descending.
	direction := "asc"
	if !config.SortAscending.IsNull() && !config.SortAscending.ValueBool() {
		direction = "desc"
	}

	rs, hresp, err := apiClient.TenantsAPI.ListTenants(ctx).
		Max(maxResults).
		Sort("id").
		Direction(direction).
		Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("LIST failed for tenants: %s",
			providererrors.ErrMsg(err, hresp)))

		return
	}

	// Both lists must be empty (not null) when nothing matches.
	ids := make([]attr.Value, 0, len(rs.Accounts))
	objs := make([]attr.Value, 0, len(rs.Accounts))

	for i := range rs.Accounts {
		a := &rs.Accounts[i]

		if !dsfilter.Matches(a, filters, fieldValue) {
			continue
		}

		v, diags := tenantToValue(ctx, a)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		objs = append(objs, v)
		ids = append(ids, convert.Int64ToType(a.Id))
	}

	idsVal, diags := types.ListValue(types.Int64Type, ids)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantsVal, diags := types.ListValue(TenantsValue{}.Type(ctx), objs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Ids = idsVal
	config.Tenants = tenantsVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// fieldValue extracts the named filterable field off a tenant account for regex
// matching, reporting whether the field is present. A missing or nil field does
// not match.
//
// The filter field `enabled` maps to the API `Active` field.
func fieldValue(a *sdk.ListTenants200ResponseAllOfAccountsInner, field string) (string, bool) {
	switch field {
	case "name":
		if a.Name != nil {
			return *a.Name, true
		}
	case "subdomain":
		if a.Subdomain != nil {
			return *a.Subdomain, true
		}
	case "currency":
		if a.Currency != nil {
			return *a.Currency, true
		}
	case "external_id":
		if v := a.ExternalId.Get(); v != nil {
			return *v, true
		}
	case "enabled":
		if a.Active != nil {
			return strconv.FormatBool(*a.Active), true
		}
	case "master":
		if a.Master != nil {
			return strconv.FormatBool(*a.Master), true
		}
	case "account_name":
		if a.AccountName != nil {
			return *a.AccountName, true
		}
	case "account_number":
		if v := a.AccountNumber.Get(); v != nil {
			return *v, true
		}
	case "customer_number":
		if v := a.CustomerNumber.Get(); v != nil {
			return *v, true
		}
	}

	return "", false
}

// compileFilters converts the configured filter blocks into compiled filters.
func compileFilters(
	ctx context.Context,
	filterSet types.Set,
	diags *diag.Diagnostics,
) []dsfilter.Compiled {
	if filterSet.IsNull() || filterSet.IsUnknown() {
		return nil
	}

	var filterBlocks []FilterValue

	diags.Append(filterSet.ElementsAs(ctx, &filterBlocks, false)...)

	if diags.HasError() {
		return nil
	}

	blocks := make([]dsfilter.Block, 0, len(filterBlocks))

	for _, b := range filterBlocks {
		var values []string

		diags.Append(b.Values.ElementsAs(ctx, &values, false)...)

		if diags.HasError() {
			return nil
		}

		blocks = append(blocks, dsfilter.Block{Name: b.Name.ValueString(), Values: values})
	}

	return dsfilter.Compile(blocks, summary, diags)
}

// tenantToValue maps an API tenant account into the generated custom object
// value, using the same field mapping (including nil-checks) as the singular
// tenant data source.
func tenantToValue(
	ctx context.Context,
	a *sdk.ListTenants200ResponseAllOfAccountsInner,
) (TenantsValue, diag.Diagnostics) {
	baseRoleID := types.Int64Null()
	baseRoleName := types.StringNull()

	if a.Role != nil {
		baseRoleID = convert.Int64ToType(a.Role.Id)
		baseRoleName = convert.StrToType(a.Role.Authority)
	}

	parentID := types.Int64Null()
	parentName := types.StringNull()
	parentSubdomain := types.StringNull()

	if a.Parent != nil {
		parentID = convert.Int64ToType(a.Parent.Id)
		parentName = convert.StrToType(a.Parent.Name)
		parentSubdomain = convert.StrToType(a.Parent.Subdomain.Get())
	}

	instanceCount := types.Int64Null()
	userCount := types.Int64Null()

	if a.Stats != nil {
		instanceCount = convert.Int64ToType(a.Stats.InstanceCount)
		userCount = convert.Int64ToType(a.Stats.UserCount)
	}

	return NewTenantsValue(
		TenantsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"account_name":     convert.StrToType(a.AccountName),
			"account_number":   convert.StrToType(a.AccountNumber.Get()),
			"base_role_id":     baseRoleID,
			"base_role_name":   baseRoleName,
			"currency":         convert.StrToType(a.Currency),
			"customer_number":  convert.StrToType(a.CustomerNumber.Get()),
			"date_created":     convert.TimeToType(a.DateCreated),
			"description":      convert.StrToType(a.Description.Get()),
			"enabled":          convert.BoolToType(a.Active),
			"external_id":      convert.StrToType(a.ExternalId.Get()),
			"id":               convert.Int64ToType(a.Id),
			"instance_count":   instanceCount,
			"last_updated":     convert.TimeToType(a.LastUpdated),
			"master":           convert.BoolToType(a.Master),
			"name":             convert.StrToType(a.Name),
			"parent_id":        parentID,
			"parent_name":      parentName,
			"parent_subdomain": parentSubdomain,
			"subdomain":        convert.StrToType(a.Subdomain),
			"user_count":       userCount,
		},
	)
}
