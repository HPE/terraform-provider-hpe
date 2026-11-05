// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolip

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
)

const (
	summary = "read network pool IP data source"
	// ErrorNoValidSearchTerms is returned when none of the mutually exclusive
	// lookup keys (id, ip_address, hostname) are set in the configuration.
	ErrorNoValidSearchTerms = `no valid search terms - exactly one of id, ` +
		`ip_address, or hostname is required`
	// ErrorNoNetworkPoolIPFound is returned when a lookup by ip_address or
	// hostname matches no reservation in the pool.
	ErrorNoNetworkPoolIPFound = `no network pool IP reservation found`
	// ErrorMultipleNetworkPoolIPs is returned when a lookup by ip_address or
	// hostname matches more than one reservation in the pool.
	ErrorMultipleNetworkPoolIPs = `multiple network pool IP reservations were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "network_pool_ip"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkPoolIpDataSourceSchema(ctx)
}

// networkPoolIPFields is the common subset of fields shared by the single-get
// and list-item response models. Mapping from a plain struct keeps the state
// conversion in one pure, unit-testable function regardless of which API call
// produced the data.
type networkPoolIPFields struct {
	Id             *int64
	NetworkPoolId  *int64
	IpType         *string
	IpAddress      *string
	GatewayAddress *string
	SubnetMask     *string
	DNSServer      *string
	InterfaceName  *string
	Description    *string
	Active         *bool
	StaticIp       *bool
	Fqdn           *string
	DomainName     *string
	Hostname       *string
	ExternalId     *string
	StartDate      *time.Time
	EndDate        *time.Time
	DateCreated    *time.Time
	LastUpdated    *time.Time
	RefType        *string
	RefId          *int64
}

// networkPoolIPAsState converts the common field set into Terraform state. It
// is pure so nullable and nil handling can be exercised by unit tests.
func networkPoolIPAsState(f networkPoolIPFields) NetworkPoolIpModel {
	return NetworkPoolIpModel{
		Id:             convert.Int64ToType(f.Id),
		NetworkPoolId:  convert.Int64ToType(f.NetworkPoolId),
		IpType:         convert.StrToType(f.IpType),
		IpAddress:      convert.StrToType(f.IpAddress),
		GatewayAddress: convert.StrToType(f.GatewayAddress),
		SubnetMask:     convert.StrToType(f.SubnetMask),
		DnsServer:      convert.StrToType(f.DNSServer),
		InterfaceName:  convert.StrToType(f.InterfaceName),
		Description:    convert.StrToType(f.Description),
		Active:         convert.BoolToType(f.Active),
		StaticIp:       convert.BoolToType(f.StaticIp),
		Fqdn:           convert.StrToType(f.Fqdn),
		DomainName:     convert.StrToType(f.DomainName),
		Hostname:       convert.StrToType(f.Hostname),
		ExternalId:     convert.StrToType(f.ExternalId),
		StartDate:      convert.TimeToType(f.StartDate),
		EndDate:        convert.TimeToType(f.EndDate),
		DateCreated:    convert.TimeToType(f.DateCreated),
		LastUpdated:    convert.TimeToType(f.LastUpdated),
		RefType:        convert.StrToType(f.RefType),
		RefId:          convert.Int64ToType(f.RefId),
	}
}

// fromSingle flattens the single-get response model into the common field set.
func fromSingle(
	ip *sdk.GetNetworkPoolIp200ResponseNetworkPoolIp,
) networkPoolIPFields {
	return networkPoolIPFields{
		Id:             ip.Id,
		NetworkPoolId:  ip.NetworkPoolId,
		IpType:         ip.IpType,
		IpAddress:      ip.IpAddress,
		GatewayAddress: ip.GatewayAddress.Get(),
		SubnetMask:     ip.SubnetMask.Get(),
		DNSServer:      ip.DnsServer.Get(),
		InterfaceName:  ip.InterfaceName.Get(),
		Description:    ip.Description.Get(),
		Active:         ip.Active,
		StaticIp:       ip.StaticIp,
		Fqdn:           ip.Fqdn,
		DomainName:     ip.DomainName.Get(),
		Hostname:       ip.Hostname,
		ExternalId:     ip.ExternalId.Get(),
		StartDate:      ip.StartDate,
		EndDate:        ip.EndDate.Get(),
		DateCreated:    ip.DateCreated,
		LastUpdated:    ip.LastUpdated,
		RefType:        ip.RefType.Get(),
		RefId:          ip.RefId.Get(),
	}
}

// fromListItem flattens a list-response item into the common field set.
func fromListItem(
	ip *sdk.GetNetworkPoolIps200ResponseAllOfNetworkPoolIpsInner,
) networkPoolIPFields {
	return networkPoolIPFields{
		Id:             ip.Id,
		NetworkPoolId:  ip.NetworkPoolId,
		IpType:         ip.IpType,
		IpAddress:      ip.IpAddress,
		GatewayAddress: ip.GatewayAddress.Get(),
		SubnetMask:     ip.SubnetMask.Get(),
		DNSServer:      ip.DnsServer.Get(),
		InterfaceName:  ip.InterfaceName.Get(),
		Description:    ip.Description.Get(),
		Active:         ip.Active,
		StaticIp:       ip.StaticIp,
		Fqdn:           ip.Fqdn,
		DomainName:     ip.DomainName.Get(),
		Hostname:       ip.Hostname,
		ExternalId:     ip.ExternalId.Get(),
		StartDate:      ip.StartDate,
		EndDate:        ip.EndDate.Get(),
		DateCreated:    ip.DateCreated,
		LastUpdated:    ip.LastUpdated,
		RefType:        ip.RefType.Get(),
		RefId:          ip.RefId.Get(),
	}
}

func getByID(
	ctx context.Context,
	networkPoolID int64,
	id int64,
	apiClient *sdk.APIClient,
) (networkPoolIPFields, error) {
	r, hresp, err := apiClient.NetworksAPI.
		GetNetworkPoolIp(ctx, networkPoolID, id).Execute()
	if hresp != nil && hresp.StatusCode == http.StatusNotFound {
		return networkPoolIPFields{}, fmt.Errorf(
			"%s: id %d in pool %d", ErrorNoNetworkPoolIPFound, id, networkPoolID)
	}

	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return networkPoolIPFields{}, fmt.Errorf(
			"GET failed for network pool IP %d in pool %d: %s",
			id, networkPoolID, providererrors.ErrMsg(err, hresp))
	}

	if r.NetworkPoolIp == nil {
		return networkPoolIPFields{}, fmt.Errorf(
			"GET failed for network pool IP %d: response missing networkPoolIp", id)
	}

	return fromSingle(r.NetworkPoolIp), nil
}

// getByList performs an exact-match lookup by ip_address or hostname. The
// server applies the exact match; this function enforces that exactly one
// reservation is returned.
func getByList(
	networkPoolID int64,
	req sdk.ApiGetNetworkPoolIpsRequest,
	label string,
) (networkPoolIPFields, error) {
	rs, hresp, err := req.Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return networkPoolIPFields{}, fmt.Errorf(
			"GET failed for network pool IP %s in pool %d: %s",
			label, networkPoolID, providererrors.ErrMsg(err, hresp))
	}

	switch len(rs.NetworkPoolIps) {
	case 0:
		return networkPoolIPFields{}, errors.New(ErrorNoNetworkPoolIPFound)
	case 1:
		item := rs.NetworkPoolIps[0]

		return fromListItem(&item), nil
	default:
		return networkPoolIPFields{}, errors.New(ErrorMultipleNetworkPoolIPs)
	}
}

func getNetworkPoolIP(
	ctx context.Context,
	config *NetworkPoolIpModel,
	apiClient *sdk.APIClient,
) (networkPoolIPFields, error) {
	poolID := config.NetworkPoolId.ValueInt64()

	switch {
	case !config.Id.IsNull():
		return getByID(ctx, poolID, config.Id.ValueInt64(), apiClient)
	case !config.IpAddress.IsNull():
		req := apiClient.NetworksAPI.
			GetNetworkPoolIps(ctx, poolID).
			IpAddress(config.IpAddress.ValueString())

		return getByList(poolID, req, config.IpAddress.ValueString())
	case !config.Hostname.IsNull():
		req := apiClient.NetworksAPI.
			GetNetworkPoolIps(ctx, poolID).
			Hostname(config.Hostname.ValueString())

		return getByList(poolID, req, config.Hostname.ValueString())
	default:
		return networkPoolIPFields{}, errors.New(ErrorNoValidSearchTerms)
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkPoolIpModel

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

	fields, err := getNetworkPoolIP(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := networkPoolIPAsState(fields)
	// Preserve the configured parent pool id, which is the required lookup
	// scope and not always echoed back by the list-item path.
	state.NetworkPoolId = config.NetworkPoolId

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
