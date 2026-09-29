// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package tenant implements the singular tenant data source.
package tenant

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read tenant data source"

const (
	// ErrorNoTenantFound is returned when a lookup matches no tenant.
	ErrorNoTenantFound = `no tenant found`
	// ErrorNoValidSearchTerms is returned when neither id nor name is set.
	ErrorNoValidSearchTerms = `no valid search terms - an id or name is required`
	// ErrorMultipleTenants is returned when a by-name lookup is ambiguous.
	ErrorMultipleTenants = `multiple tenants were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "tenant"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = TenantDataSourceSchema(ctx)
}

// getTenantByID fetches a single tenant account by its ID. A 404/not-found is
// reported as the clear "no tenant found" error.
func getTenantByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetTenant200ResponseAccount, error) {
	resp, hresp, err := apiClient.TenantsAPI.GetTenant(ctx, id).Execute()
	if errfmt.IsNotFound(hresp) {
		return nil, errors.New(ErrorNoTenantFound)
	}
	if err := errfmt.CheckResponse(err, hresp); err != nil {
		return nil, fmt.Errorf("GET failed for tenant %d: %s", id, errfmt.ErrMsg(err, hresp))
	}

	if resp == nil || resp.Account == nil {
		return nil, errors.New(ErrorNoTenantFound)
	}

	return resp.Account, nil
}

// getTenantByName resolves a tenant by its exact name. The LIST name parameter
// is an exact, case-sensitive match server-side, but names are not guaranteed
// unique, so the results are filtered client-side and 0 or >1 matches are
// treated as errors. The single match is re-fetched by id so the state is
// populated from the same account shape as the by-id path.
func getTenantByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetTenant200ResponseAccount, error) {
	rs, hresp, err := apiClient.TenantsAPI.ListTenants(ctx).Name(name).Execute()
	if err := errfmt.CheckResponse(err, hresp); err != nil {
		return nil, fmt.Errorf("LIST failed for tenant %s: %s", name, errfmt.ErrMsg(err, hresp))
	}

	if rs == nil {
		return nil, errors.New(ErrorNoTenantFound)
	}

	var matches []sdk.ListTenants200ResponseAllOfAccountsInner

	for _, a := range rs.Accounts {
		if a.Name != nil && *a.Name == name {
			matches = append(matches, a)
		}
	}

	switch {
	case len(matches) == 1:
		if matches[0].Id == nil {
			return nil, errors.New(ErrorNoTenantFound)
		}

		return getTenantByID(ctx, *matches[0].Id, apiClient)
	case len(matches) > 1:
		return nil, errors.New(ErrorMultipleTenants)
	default:
		return nil, errors.New(ErrorNoTenantFound)
	}
}

// getTenant dispatches on the configured search terms.
func getTenant(
	ctx context.Context,
	data TenantModel,
	apiClient *sdk.APIClient,
) (*sdk.GetTenant200ResponseAccount, error) {
	if !data.Id.IsNull() {
		return getTenantByID(ctx, data.Id.ValueInt64(), apiClient)
	} else if !data.Name.IsNull() {
		return getTenantByName(ctx, data.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data TenantModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	a, err := getTenant(ctx, data, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	data.Id = convert.Int64ToType(a.Id)
	data.Name = convert.StrToType(a.Name)
	data.Description = convert.StrToType(a.Description.Get())
	data.Enabled = convert.BoolToType(a.Active)
	data.Subdomain = convert.StrToType(a.Subdomain)
	data.Currency = convert.StrToType(a.Currency)
	data.AccountNumber = convert.StrToType(a.AccountNumber.Get())
	data.AccountName = convert.StrToType(a.AccountName)
	data.CustomerNumber = convert.StrToType(a.CustomerNumber.Get())
	data.ExternalId = convert.StrToType(a.ExternalId.Get())
	data.Master = convert.BoolToType(a.Master)
	data.DateCreated = convert.TimeToType(a.DateCreated)
	data.LastUpdated = convert.TimeToType(a.LastUpdated)

	if a.Role != nil {
		data.BaseRoleId = convert.Int64ToType(a.Role.Id)
		data.BaseRoleName = convert.StrToType(a.Role.Authority)
	} else {
		// The master tenant has no base role.
		data.BaseRoleId = convert.Int64ToType(nil)
		data.BaseRoleName = convert.StrToType(nil)
	}

	if a.Parent != nil {
		data.ParentId = convert.Int64ToType(a.Parent.Id)
		data.ParentName = convert.StrToType(a.Parent.Name)
		data.ParentSubdomain = convert.StrToType(a.Parent.Subdomain.Get())
	} else {
		// The master tenant has no parent.
		data.ParentId = convert.Int64ToType(nil)
		data.ParentName = convert.StrToType(nil)
		data.ParentSubdomain = convert.StrToType(nil)
	}

	if a.Stats != nil {
		data.InstanceCount = convert.Int64ToType(a.Stats.InstanceCount)
		data.UserCount = convert.Int64ToType(a.Stats.UserCount)
	} else {
		data.InstanceCount = convert.Int64ToType(nil)
		data.UserCount = convert.Int64ToType(nil)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
