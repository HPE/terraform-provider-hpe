// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package networkdomain implements a data source for network_domain
package networkdomain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                     = "read network domain data source"
	ErrorNoValidSearchTerms     = `no valid search terms - an id or name is required`
	ErrorRunningPreApply        = `Error running pre-apply plan: exit status 1`
	ErrorNoNetworkDomainFound   = `no network domain found`
	ErrorMultipleNetworkDomains = `multiple network domains were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_network_domain"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkDomainDataSourceSchema(ctx)
}

// networkDomainListItem is a minimal struct used to decode the untyped list response.
type networkDomainListItem struct {
	Id   *int64  `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

func networkDomainAsState(
	ctx context.Context,
	domain *sdk.GetNetworkDomain200ResponseNetworkDomain,
) (NetworkDomainModel, error) {
	state := NetworkDomainModel{
		Id:               convert.Int64ToType(domain.Id),
		Name:             convert.StrToType(domain.Name),
		Active:           convert.BoolToType(domain.Active),
		Fqdn:             convert.StrToType(domain.Fqdn.Get()),
		Description:      convert.StrToType(domain.Description.Get()),
		Visibility:       convert.StrToType(domain.Visibility),
		DomainController: convert.BoolToType(domain.DomainController),
		PublicZone:       convert.BoolToType(domain.PublicZone),
		DomainUsername:   convert.StrToType(domain.DomainUsername.Get()),
		DomainPassword:   convert.StrToType(domain.DomainPassword.Get()),
		RefType:          convert.StrToType(domain.RefType.Get()),
		RefId:            convert.Int64ToType(domain.RefId.Get()),
		RefSource:        convert.StrToType(domain.RefSource.Get()),
		InternalId:       convert.StrToType(domain.InternalId.Get()),
		OuPath:           convert.StrToType(domain.OuPath.Get()),
		DcServer:         convert.StrToType(domain.DcServer.Get()),
		CloudType:        convert.StrToType(domain.ZoneType.Get()),
		Dnssec:           convert.StrToType(domain.Dnssec.Get()),
		DomainSerial:     convert.StrToType(domain.DomainSerial.Get()),
	}

	if domain.Account != nil {
		acct, diags := NewAccountValue(
			AccountValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(domain.Account.Id),
				"name": convert.StrToType(domain.Account.Name),
			},
		)
		if diags.HasError() {
			return NetworkDomainModel{}, fmt.Errorf("error creating account value")
		}

		state.Account = acct
	} else {
		state.Account = NewAccountValueNull()
	}

	if domain.Owner != nil {
		owner, diags := NewOwnerValue(
			OwnerValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(domain.Owner.Id),
				"name": convert.StrToType(domain.Owner.Name),
			},
		)
		if diags.HasError() {
			return NetworkDomainModel{}, fmt.Errorf("error creating owner value")
		}

		state.Owner = owner
	} else {
		state.Owner = NewOwnerValueNull()
	}

	return state, nil
}

func getNetworkDomainByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkDomain200ResponseNetworkDomain, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkDomain(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for network domain %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	domain := r.GetNetworkDomain()

	return &domain, nil
}

func getNetworkDomainByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkDomain200ResponseNetworkDomain, error) {
	rs, hresp, err := apiClient.NetworksAPI.GetNetworkDomains(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for network domain %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	raw, err := json.Marshal(rs.GetNetworkDomains())
	if err != nil {
		return nil, fmt.Errorf("error marshalling network domains list: %w", err)
	}

	var items []networkDomainListItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("error decoding network domains list: %w", err)
	}

	var matchedIDs []int64

	for _, item := range items {
		if item.Name != nil && *item.Name == name && item.Id != nil {
			matchedIDs = append(matchedIDs, *item.Id)
		}
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNoNetworkDomainFound)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleNetworkDomains)
	}

	return getNetworkDomainByID(ctx, matchedIDs[0], apiClient)
}

func getNetworkDomain(
	ctx context.Context,
	config *NetworkDomainModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkDomain200ResponseNetworkDomain, error) {
	if !config.Id.IsNull() {
		return getNetworkDomainByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getNetworkDomainByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkDomainModel

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

	domain, err := getNetworkDomain(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state, err := networkDomainAsState(ctx, domain)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
