// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkinterfacetype

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

const (
	summary                            = "read network interface type data source"
	ErrorNoNetworkInterfaceTypeFound   = `no network interface type found`
	ErrorMultipleNetworkInterfaceTypes = `multiple network interface types were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "network_interface_type"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkInterfaceTypeDataSourceSchema(ctx)
}

// matchedNetworkInterfaceType is the resolved id and code of a NIC type.
type matchedNetworkInterfaceType struct {
	id   int64
	code string
}

// matchNetworkInterfaceType finds the single NIC type whose name equals the
// requested name. It errors when zero or more than one match, so the data
// source fails clearly rather than silently picking an arbitrary one. It is a
// pure function so the match/error logic is unit testable without an appliance.
func matchNetworkInterfaceType(
	networkTypes []sdk.ZoneNetworkOptionsResponseNetworkTypesInner,
	name string,
) (matchedNetworkInterfaceType, error) {
	var matches []matchedNetworkInterfaceType

	for _, nt := range networkTypes {
		if nt.Name != nil && *nt.Name == name {
			m := matchedNetworkInterfaceType{}
			if nt.Id != nil {
				m.id = *nt.Id
			}

			if nt.Code != nil {
				m.code = *nt.Code
			}

			matches = append(matches, m)
		}
	}

	switch len(matches) {
	case 0:
		return matchedNetworkInterfaceType{}, errors.New(ErrorNoNetworkInterfaceTypeFound)
	case 1:
		return matches[0], nil
	default:
		return matchedNetworkInterfaceType{}, errors.New(ErrorMultipleNetworkInterfaceTypes)
	}
}

// resolveProvisionTypeID resolves a provision type code (e.g. "vmware") to its
// numeric id, requiring exactly one match.
func resolveProvisionTypeID(
	ctx context.Context,
	code string,
	apiClient *sdk.APIClient,
) (int64, error) {
	pTypes, hresp, err := apiClient.ProvisioningAPI.ListProvisionTypes(ctx).Code(code).Execute()
	if pTypes == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("GET failed for provision type code %s: %s",
			code, providererrors.ErrMsg(err, hresp))
	}

	var matching []sdk.ListProvisionTypes200ResponseAllOfProvisionTypesInner

	for _, pt := range pTypes.ProvisionTypes {
		if pt.Code != nil && *pt.Code == code {
			matching = append(matching, pt)
		}
	}

	switch {
	case len(matching) == 0:
		return 0, fmt.Errorf("provision type with code %s not found", code)
	case len(matching) > 1:
		return 0, fmt.Errorf("multiple provision types with code %s found", code)
	case matching[0].Id == nil:
		return 0, fmt.Errorf("id not found for provision type with code %s", code)
	default:
		return *matching[0].Id, nil
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkInterfaceTypeModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	// Data source schemas cannot declare attribute defaults, so apply hpegl's
	// vmware default here when the argument is omitted.
	provisionTypeCode := defaultProvisionTypeCode
	if !config.ProvisionTypeCode.IsNull() && !config.ProvisionTypeCode.IsUnknown() {
		provisionTypeCode = config.ProvisionTypeCode.ValueString()
	}

	ptID, err := resolveProvisionTypeID(ctx, provisionTypeCode, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	// For a cloud + provision-type pair, zoneNetworkOptions returns the
	// available NIC types under the generic name "networkTypes".
	opts, hresp, err := apiClient.OptionsAPI.ListOptionNetworkOptions(ctx).
		ZoneId(config.CloudId.ValueInt64()).
		ProvisionTypeId(ptID).
		Execute()
	if opts == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(summary,
			fmt.Sprintf("GET failed for network interface types: %s",
				providererrors.ErrMsg(err, hresp)))

		return
	}

	var networkTypes []sdk.ZoneNetworkOptionsResponseNetworkTypesInner
	if opts.Data != nil {
		networkTypes = opts.Data.NetworkTypes
	}

	match, err := matchNetworkInterfaceType(networkTypes, config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state := NetworkInterfaceTypeModel{
		Name:              config.Name,
		CloudId:           config.CloudId,
		ProvisionTypeCode: types.StringValue(provisionTypeCode),
		Id:                types.Int64Value(match.id),
		Code:              types.StringValue(match.code),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
