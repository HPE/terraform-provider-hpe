// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package networkfloatingip implements a singular data source that looks up a
// single Morpheus network floating IP by id or by IP address.
package networkfloatingip

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
	"github.com/HPE/terraform-provider-hpe/utils/paging"
)

const (
	summary = "read network floating ip data source"

	// ErrorNoValidSearchTerms is returned when neither id nor ip_address is set.
	ErrorNoValidSearchTerms = `no valid search terms - an id or ip_address is required`
	// ErrorNoFloatingIpFound is returned when a lookup by IP address matches nothing.
	ErrorNoFloatingIpFound = `no network floating ip found`
	// ErrorMultipleFloatingIps is returned when a lookup by IP address is ambiguous.
	ErrorMultipleFloatingIps = `multiple network floating ips were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "network_floating_ip"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkFloatingIpDataSourceSchema(ctx)
}

// FloatingIp is the subset of fields both the singular GET response and the
// plural list item carry. Mapping is written once against this shape so the two
// data sources stay consistent; each caller adapts its own SDK model to it.
type FloatingIp struct {
	Id                  *int64
	ExternalId          *string
	IpStatus            *string
	IpAddress           *string
	PtrId               *string
	CloudId             *int64
	CloudName           *string
	ServerId            *int64
	ServerName          *string
	IpRangeStartAddress *string
	IpRangeEndAddress   *string
	NetworkDomainId     *int64
	NetworkDomainName   *string
	DateCreated         *time.Time
	LastUpdated         *time.Time
}

// FloatingIpFromGet normalises the singular GET response into FloatingIp.
func FloatingIpFromGet(fip *sdk.GetNetworkFloatingIp200ResponseNetworkFloatingIp) FloatingIp {
	out := FloatingIp{
		Id:          fip.Id,
		ExternalId:  fip.ExternalId.Get(),
		IpStatus:    fip.IpStatus,
		IpAddress:   fip.IpAddress,
		PtrId:       fip.PtrId.Get(),
		DateCreated: fip.DateCreated,
		LastUpdated: fip.LastUpdated,
	}

	if fip.Cloud != nil {
		out.CloudId = fip.Cloud.Id
		out.CloudName = fip.Cloud.Name
	}

	if fip.Server != nil {
		out.ServerId = fip.Server.Id
		out.ServerName = fip.Server.Name
	}

	if fip.IpRange != nil {
		out.IpRangeStartAddress = fip.IpRange.StartAddress
		out.IpRangeEndAddress = fip.IpRange.EndAddress
	}

	if fip.NetworkDomain != nil {
		out.NetworkDomainId = fip.NetworkDomain.Id
		out.NetworkDomainName = fip.NetworkDomain.Name
	}

	return out
}

// FloatingIpFromListInner normalises a plural list item into FloatingIp. It
// mirrors FloatingIpFromGet for the sibling model so both data sources can serve
// their results through the same mapping.
func FloatingIpFromListInner(
	item *sdk.GetAllNetworkFloatingIps200ResponseAllOfNetworkFloatingIpsInner,
) FloatingIp {
	out := FloatingIp{
		Id:          item.Id,
		ExternalId:  item.ExternalId.Get(),
		IpStatus:    item.IpStatus,
		IpAddress:   item.IpAddress,
		PtrId:       item.PtrId.Get(),
		DateCreated: item.DateCreated,
		LastUpdated: item.LastUpdated,
	}

	if item.Cloud != nil {
		out.CloudId = item.Cloud.Id
		out.CloudName = item.Cloud.Name
	}

	if item.Server != nil {
		out.ServerId = item.Server.Id
		out.ServerName = item.Server.Name
	}

	if item.IpRange != nil {
		out.IpRangeStartAddress = item.IpRange.StartAddress
		out.IpRangeEndAddress = item.IpRange.EndAddress
	}

	if item.NetworkDomain != nil {
		out.NetworkDomainId = item.NetworkDomain.Id
		out.NetworkDomainName = item.NetworkDomain.Name
	}

	return out
}

// FloatingIpToState maps the normalised floating IP into the generated model.
// This is a pure function so it can be unit-tested without an API client.
func FloatingIpToState(fip FloatingIp) NetworkFloatingIpModel {
	return NetworkFloatingIpModel{
		Id:                  convert.Int64ToType(fip.Id),
		ExternalId:          convert.StrToType(fip.ExternalId),
		IpStatus:            convert.StrToType(fip.IpStatus),
		IpAddress:           convert.StrToType(fip.IpAddress),
		PtrId:               convert.StrToType(fip.PtrId),
		CloudId:             convert.Int64ToType(fip.CloudId),
		CloudName:           convert.StrToType(fip.CloudName),
		ServerId:            convert.Int64ToType(fip.ServerId),
		ServerName:          convert.StrToType(fip.ServerName),
		IpRangeStartAddress: convert.StrToType(fip.IpRangeStartAddress),
		IpRangeEndAddress:   convert.StrToType(fip.IpRangeEndAddress),
		NetworkDomainId:     convert.Int64ToType(fip.NetworkDomainId),
		NetworkDomainName:   convert.StrToType(fip.NetworkDomainName),
		DateCreated:         convert.TimeToType(fip.DateCreated),
		LastUpdated:         convert.TimeToType(fip.LastUpdated),
	}
}

func getFloatingIpByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkFloatingIp200ResponseNetworkFloatingIp, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkFloatingIp(ctx, id).Execute()
	if r == nil || err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network floating ip %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	if r.NetworkFloatingIp == nil {
		return nil, fmt.Errorf(
			"GET failed for network floating ip %d: response missing networkFloatingIp", id)
	}

	return r.NetworkFloatingIp, nil
}

// getFloatingIpByIPAddress uses the list endpoint's IpAddress filter and then
// requires exactly one result: the IP address is expected to be unique, so zero
// or multiple matches are treated as errors rather than silently guessing.
func getFloatingIpByIPAddress(
	ctx context.Context,
	ipAddress string,
	apiClient *sdk.APIClient,
) (FloatingIp, error) {
	// The IpAddress query narrows server-side; every page of what it returns
	// is collected, so an exact match is found even when the filter is a
	// partial match that fills more than one page.
	items, err := paging.Collect(ctx,
		func(ctx context.Context, offset, max int64) (
			[]sdk.GetAllNetworkFloatingIps200ResponseAllOfNetworkFloatingIpsInner, int64, error,
		) {
			rs, hresp, err := apiClient.NetworksAPI.
				GetAllNetworkFloatingIps(ctx).
				IpAddress(ipAddress).
				Max(max).
				Offset(offset).
				Execute()
			if rs == nil || err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
				return nil, 0, fmt.Errorf(
					"LIST failed for network floating ip %s: %s",
					ipAddress, providererrors.ErrMsg(err, hresp))
			}

			var total int64
			if rs.Meta != nil && rs.Meta.Total != nil {
				total = *rs.Meta.Total
			}

			return rs.NetworkFloatingIps, total, nil
		})
	if err != nil {
		return FloatingIp{}, err
	}

	// Match exactly on the address so partial matches cannot inflate the count.
	var matches []*sdk.GetAllNetworkFloatingIps200ResponseAllOfNetworkFloatingIpsInner

	for i := range items {
		item := &items[i]
		if item.IpAddress != nil && *item.IpAddress == ipAddress {
			matches = append(matches, item)
		}
	}

	switch {
	case len(matches) == 0:
		return FloatingIp{}, errors.New(ErrorNoFloatingIpFound)
	case len(matches) > 1:
		return FloatingIp{}, errors.New(ErrorMultipleFloatingIps)
	}

	return FloatingIpFromListInner(matches[0]), nil
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkFloatingIpModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Exactly one of id or ip_address must be supplied. The schema already
	// declares them as mutually exclusive; this guards the "neither set" case.
	hasID := !config.Id.IsNull() && !config.Id.IsUnknown()
	hasIP := !config.IpAddress.IsNull() && !config.IpAddress.IsUnknown()

	if !hasID && !hasIP {
		resp.Diagnostics.AddError(summary, ErrorNoValidSearchTerms)

		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	var fip FloatingIp

	if hasID {
		got, err := getFloatingIpByID(ctx, config.Id.ValueInt64(), apiClient)
		if err != nil {
			resp.Diagnostics.AddError(summary, err.Error())

			return
		}

		fip = FloatingIpFromGet(got)
	} else {
		got, err := getFloatingIpByIPAddress(ctx, config.IpAddress.ValueString(), apiClient)
		if err != nil {
			resp.Diagnostics.AddError(summary, err.Error())

			return
		}

		fip = got
	}

	state := FloatingIpToState(fip)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
