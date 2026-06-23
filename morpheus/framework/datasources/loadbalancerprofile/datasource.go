// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read load balancer profile data source"

var _ datasource.DataSource = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_load_balancer_profile"
}

func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = LoadBalancerProfileDataSourceSchema(ctx)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config LoadBalancerProfileModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("failed to create client: %s", err.Error()))

		return
	}

	state, err := getLoadBalancerProfile(ctx, config, client)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func getLoadBalancerProfile(
	ctx context.Context,
	config LoadBalancerProfileModel,
	client *sdk.APIClient,
) (*LoadBalancerProfileModel, error) {
	loadBalancerID := config.LoadBalancerId.ValueInt64()

	if !config.Id.IsNull() {
		return getLoadBalancerProfileByID(ctx, loadBalancerID, config.Id.ValueInt64(), client)
	}

	if !config.Name.IsNull() {
		return getLoadBalancerProfileByName(ctx, loadBalancerID, config.Name.ValueString(), client)
	}

	return nil, fmt.Errorf("either id or name must be specified")
}

func getLoadBalancerProfileByID(
	ctx context.Context,
	loadBalancerID int64,
	id int64,
	client *sdk.APIClient,
) (*LoadBalancerProfileModel, error) {
	resp, hresp, err := client.LoadBalancersAPI.
		GetLoadBalancerProfile(ctx, loadBalancerID, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"load balancer %d profile %d GET failed: %s",
			loadBalancerID, id, errfmt.ErrMsg(err, hresp),
		)
	}

	p := resp.LoadBalancerProfile
	if p == nil {
		return nil, fmt.Errorf(
			"load balancer %d profile %d GET returned no profile",
			loadBalancerID, id,
		)
	}

	state := populateLoadBalancerProfileState(ctx, loadBalancerID, p)

	return state, nil
}

func getLoadBalancerProfileByName(
	ctx context.Context,
	loadBalancerID int64,
	name string,
	client *sdk.APIClient,
) (*LoadBalancerProfileModel, error) {
	profiles, hresp, err := client.LoadBalancersAPI.
		ListLoadBalancerProfiles(ctx, loadBalancerID).
		Name(name).Max(10000).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"load balancer %d profile list failed: %s",
			loadBalancerID, errfmt.ErrMsg(err, hresp),
		)
	}

	var matching []sdk.ListLoadBalancerProfiles200ResponseAllOfLoadBalancerProfilesInner
	for _, p := range profiles.LoadBalancerProfiles {
		if p.Name != nil && *p.Name == name {
			matching = append(matching, p)
		}
	}

	if len(matching) == 0 {
		return nil, fmt.Errorf(
			"load balancer %d profile with name %q not found",
			loadBalancerID, name,
		)
	}

	if len(matching) > 1 {
		var ids []string
		for _, p := range matching {
			if p.Id != nil {
				ids = append(ids, fmt.Sprintf("%d", *p.Id))
			}
		}

		return nil, fmt.Errorf(
			"multiple load balancer profiles found with name %q. IDs: %s. "+
				"Please specify an ID instead",
			name,
			strings.Join(ids, ", "),
		)
	}

	id := matching[0].Id
	if id == nil {
		return nil, fmt.Errorf(
			"load balancer %d profile with name %q has missing ID",
			loadBalancerID, name,
		)
	}

	return getLoadBalancerProfileByID(ctx, loadBalancerID, *id, client)
}

//nolint:funlen,cyclop // mapping all fields requires length
func populateLoadBalancerProfileState(
	ctx context.Context,
	loadBalancerID int64,
	p *sdk.GetLoadBalancerProfile200ResponseLoadBalancerProfile,
) *LoadBalancerProfileModel {
	state := &LoadBalancerProfileModel{}

	// Core fields
	state.Id = convert.Int64ToType(p.Id)
	state.LoadBalancerId = types.Int64Value(loadBalancerID)
	state.Name = convert.StrToType(p.Name)
	state.Category = convert.StrToType(p.Category)
	state.ServiceType = convert.StrToType(p.ServiceType)
	state.ServiceTypeDisplay = convert.StrToType(p.ServiceTypeDisplay)
	state.Visibility = convert.StrToType(p.Visibility)
	state.Description = convert.StrToType(p.Description)
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

	// *time.Time fields
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

	// LoadBalancer nested object
	populateLoadBalancerObj(ctx, state, p)

	// Config — dispatch by populated variant
	populateConfig(ctx, state, p)

	return state
}

func populateLoadBalancerObj(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	p *sdk.GetLoadBalancerProfile200ResponseLoadBalancerProfile,
) {
	if p.LoadBalancer != nil {
		lb := p.LoadBalancer

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
			if !d.HasError() {
				tvObj, tvDiags := tv.ToObjectValue(ctx)
				if !tvDiags.HasError() {
					typeVal = tvObj
				}
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
		if !d.HasError() {
			state.LoadBalancer = lbVal
		} else {
			state.LoadBalancer = NewLoadBalancerValueNull()
		}
	} else {
		state.LoadBalancer = NewLoadBalancerValueNull()
	}
}

//nolint:funlen,cyclop // config mapping requires length
func populateConfig(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	p *sdk.GetLoadBalancerProfile200ResponseLoadBalancerProfile,
) {
	// Initialize all config blocks to null
	state.ConfigHttp = NewConfigHttpValueNull()
	state.ConfigFastTcp = NewConfigFastTcpValueNull()
	state.ConfigFastUdp = NewConfigFastUdpValueNull()
	state.ConfigCookiePersistence = NewConfigCookiePersistenceValueNull()
	state.ConfigSourceIpPersistence = NewConfigSourceIpPersistenceValueNull()
	state.ConfigGenericPersistence = NewConfigGenericPersistenceValueNull()
	state.ConfigClientSsl = NewConfigClientSslValueNull()
	state.ConfigServerSsl = NewConfigServerSslValueNull()
	state.Tags = types.SetNull(TagsType{types.ObjectType{AttrTypes: TagsValue{}.AttributeTypes(ctx)}})

	if p.Config == nil {
		return
	}

	cfg := p.Config

	switch {
	case cfg.HTTPLoadBalancerProfileConfig3 != nil:
		populateHTTPConfig(ctx, state, cfg.HTTPLoadBalancerProfileConfig3)
	case cfg.FastTCPLoadBalancerProfileConfig3 != nil:
		populateFastTCPConfig(ctx, state, cfg.FastTCPLoadBalancerProfileConfig3)
	case cfg.FastUDPLoadBalancerProfileConfig3 != nil:
		populateFastUDPConfig(ctx, state, cfg.FastUDPLoadBalancerProfileConfig3)
	case cfg.CookiePersistenceLoadBalancerProfileConfig3 != nil:
		populateCookiePersistenceConfig(ctx, state, cfg.CookiePersistenceLoadBalancerProfileConfig3)
	case cfg.SourceIPPersistenceLoadBalancerProfileConfig3 != nil:
		populateSourceIPPersistenceConfig(ctx, state, cfg.SourceIPPersistenceLoadBalancerProfileConfig3)
	case cfg.GenericPersistenceLoadBalancerProfileConfig3 != nil:
		populateGenericPersistenceConfig(ctx, state, cfg.GenericPersistenceLoadBalancerProfileConfig3)
	case cfg.ClientSSLLoadBalancerProfileConfig3 != nil:
		populateClientSSLConfig(ctx, state, cfg.ClientSSLLoadBalancerProfileConfig3)
	case cfg.ServerSSLLoadBalancerProfileConfig3 != nil:
		populateServerSSLConfig(ctx, state, cfg.ServerSSLLoadBalancerProfileConfig3)
	}
}

func populateHTTPConfig(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	cfg *sdk.HTTPLoadBalancerProfileConfig3,
) {
	httpIdleTimeout := types.Int64Null()
	if cfg.HttpIdleTimeout.IsSet() {
		httpIdleTimeout = convert.Int64ToType(cfg.HttpIdleTimeout.Get())
	}

	requestHeaderSize := types.Int64Null()
	if cfg.RequestHeaderSize.IsSet() {
		requestHeaderSize = convert.Int64ToType(cfg.RequestHeaderSize.Get())
	}

	responseHeaderSize := types.Int64Null()
	if cfg.ResponseHeaderSize.IsSet() {
		responseHeaderSize = convert.Int64ToType(cfg.ResponseHeaderSize.Get())
	}

	requestBodySize := types.Int64Null()
	if cfg.RequestBodySize.IsSet() {
		requestBodySize = convert.Int64ToType(cfg.RequestBodySize.Get())
	}

	responseTimeout := types.Int64Null()
	if cfg.ResponseTimeout.IsSet() {
		responseTimeout = convert.Int64ToType(cfg.ResponseTimeout.Get())
	}

	httpsRedirect := types.BoolNull()
	if cfg.HttpsRedirect != nil {
		httpsRedirect = convert.StringToBool(ctx, *cfg.HttpsRedirect)
	}

	state.ConfigHttp = ConfigHttpValue{
		HttpIdleTimeout:    httpIdleTimeout,
		RequestHeaderSize:  requestHeaderSize,
		ResponseHeaderSize: responseHeaderSize,
		HttpsRedirect:      httpsRedirect,
		RedirectAddress:    convert.StrToType(cfg.RedirectAddress),
		XForwardedFor:      convert.StrToType(cfg.XForwardedFor),
		RequestBodySize:    requestBodySize,
		ResponseTimeout:    responseTimeout,
		NtlmAuthentication: convert.BoolToType(cfg.NtlmAuthentication),
		state:              attr.ValueStateKnown,
	}

	populateTagsFromHTTP(ctx, state, cfg.Tags)
}

func populateFastTCPConfig(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	cfg *sdk.FastTCPLoadBalancerProfileConfig3,
) {
	fastTCPIdleTimeout := types.Int64Null()
	if cfg.FastTcpIdleTimeout.IsSet() {
		fastTCPIdleTimeout = convert.Int64ToType(cfg.FastTcpIdleTimeout.Get())
	}

	connectionCloseTimeout := types.Int64Null()
	if cfg.ConnectionCloseTimeout.IsSet() {
		connectionCloseTimeout = convert.Int64ToType(cfg.ConnectionCloseTimeout.Get())
	}

	state.ConfigFastTcp = ConfigFastTcpValue{
		FastTcpIdleTimeout:     fastTCPIdleTimeout,
		HaFlowMirroring:        convert.BoolToType(cfg.HaFlowMirroring),
		ConnectionCloseTimeout: connectionCloseTimeout,
		state:                  attr.ValueStateKnown,
	}

	populateTagsFromFastTCP(ctx, state, cfg.Tags)
}

func populateFastUDPConfig(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	cfg *sdk.FastUDPLoadBalancerProfileConfig3,
) {
	fastUDPIdleTimeout := types.Int64Null()
	if cfg.FastUdpIdleTimeout.IsSet() {
		fastUDPIdleTimeout = convert.Int64ToType(cfg.FastUdpIdleTimeout.Get())
	}

	state.ConfigFastUdp = ConfigFastUdpValue{
		FastUdpIdleTimeout: fastUDPIdleTimeout,
		HaFlowMirroring:    convert.BoolToType(cfg.HaFlowMirroring),
		state:              attr.ValueStateKnown,
	}

	populateTagsFromFastUDP(ctx, state, cfg.Tags)
}

func populateCookiePersistenceConfig(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	cfg *sdk.CookiePersistenceLoadBalancerProfileConfig3,
) {
	maxIdleTime := types.Int64Null()
	if cfg.MaxIdleTime.IsSet() {
		maxIdleTime = convert.Int64ToType(cfg.MaxIdleTime.Get())
	}

	maxCookieAge := types.Int64Null()
	if cfg.MaxCookieAge.IsSet() {
		maxCookieAge = convert.Int64ToType(cfg.MaxCookieAge.Get())
	}

	cookieType := types.StringNull()
	if cfg.CookieType != nil {
		cookieType = convert.CookieTypeFromAPI(cfg.CookieType)
	}

	state.ConfigCookiePersistence = ConfigCookiePersistenceValue{
		SharePersistence: convert.BoolToType(cfg.SharePersistence),
		CookieName:       convert.StrToType(cfg.CookieName),
		CookieFallback:   convert.BoolToType(cfg.CookieFallback),
		CookieGarbling:   convert.BoolToType(cfg.CookieGarbling),
		CookieMode:       convert.StrToType(cfg.CookieMode),
		CookieDomain:     convert.StrToType(cfg.CookieDomain),
		CookiePath:       convert.StrToType(cfg.CookiePath),
		CookieType:       cookieType,
		MaxIdleTime:      maxIdleTime,
		MaxCookieAge:     maxCookieAge,
		state:            attr.ValueStateKnown,
	}

	populateTagsFromCookiePersistence(ctx, state, cfg.Tags)
}

func populateSourceIPPersistenceConfig(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	cfg *sdk.SourceIPPersistenceLoadBalancerProfileConfig3,
) {
	persistenceEntryTimeout := types.Int64Null()
	if cfg.PersistenceEntryTimeout.IsSet() {
		persistenceEntryTimeout = convert.Int64ToType(cfg.PersistenceEntryTimeout.Get())
	}

	state.ConfigSourceIpPersistence = ConfigSourceIpPersistenceValue{
		SharePersistence:        convert.BoolToType(cfg.SharePersistence),
		PurgeEntries:            convert.BoolToType(cfg.PurgeEntries),
		HaPersistenceMirroring:  convert.BoolToType(cfg.HaPersistenceMirroring),
		PersistenceEntryTimeout: persistenceEntryTimeout,
		state:                   attr.ValueStateKnown,
	}

	populateTagsFromSourceIPPersistence(ctx, state, cfg.Tags)
}

func populateGenericPersistenceConfig(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	cfg *sdk.GenericPersistenceLoadBalancerProfileConfig3,
) {
	persistenceEntryTimeout := types.Int64Null()
	if cfg.PersistenceEntryTimeout.IsSet() {
		persistenceEntryTimeout = convert.Int64ToType(cfg.PersistenceEntryTimeout.Get())
	}

	state.ConfigGenericPersistence = ConfigGenericPersistenceValue{
		SharePersistence:        convert.BoolToType(cfg.SharePersistence),
		HaPersistenceMirroring:  convert.BoolToType(cfg.HaPersistenceMirroring),
		PersistenceEntryTimeout: persistenceEntryTimeout,
		state:                   attr.ValueStateKnown,
	}

	populateTagsFromGenericPersistence(ctx, state, cfg.Tags)
}

func populateClientSSLConfig(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	cfg *sdk.ClientSSLLoadBalancerProfileConfig3,
) {
	sessionCacheTimeout := types.Int64Null()
	if cfg.SessionCacheTimeout.IsSet() {
		sessionCacheTimeout = convert.Int64ToType(cfg.SessionCacheTimeout.Get())
	}

	state.ConfigClientSsl = ConfigClientSslValue{
		SslSuite:              convert.StrToType(cfg.SslSuite),
		SessionCache:          convert.BoolToType(cfg.SessionCache),
		SessionCacheTimeout:   sessionCacheTimeout,
		PreferServerCipher:    convert.BoolToType(cfg.PreferServerCipher),
		SupportedSslCiphers:   convert.StrSliceToSet(cfg.SupportedSslCiphers),
		SupportedSslProtocols: convert.StrSliceToSet(cfg.SupportedSslProtocols),
		state:                 attr.ValueStateKnown,
	}

	populateTagsFromClientSSL(ctx, state, cfg.Tags)
}

func populateServerSSLConfig(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	cfg *sdk.ServerSSLLoadBalancerProfileConfig3,
) {
	state.ConfigServerSsl = ConfigServerSslValue{
		SslSuite:              convert.StrToType(cfg.SslSuite),
		SessionCache:          convert.BoolToType(cfg.SessionCache),
		SupportedSslCiphers:   convert.StrSliceToSet(cfg.SupportedSslCiphers),
		SupportedSslProtocols: convert.StrSliceToSet(cfg.SupportedSslProtocols),
		state:                 attr.ValueStateKnown,
	}

	populateTagsFromServerSSL(ctx, state, cfg.Tags)
}

// Tag population helpers — each config variant has its own tag type.

func populateTagsFromHTTP(ctx context.Context, state *LoadBalancerProfileModel, tags []sdk.LoadBalancerProfileTag24) {
	tagSet, diags := convert.ToSetType(ctx, tags, func(t sdk.LoadBalancerProfileTag24) TagsValue {
		return TagsValue{
			Name:  convert.StrToType(t.Name),
			Value: convert.StrToType(t.Value),
			state: attr.ValueStateKnown,
		}
	})
	if !diags.HasError() {
		state.Tags = tagSet
	}
}

func populateTagsFromFastTCP(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	tags []sdk.LoadBalancerProfileTag25,
) {
	tagSet, diags := convert.ToSetType(ctx, tags, func(t sdk.LoadBalancerProfileTag25) TagsValue {
		return TagsValue{
			Name:  convert.StrToType(t.Name),
			Value: convert.StrToType(t.Value),
			state: attr.ValueStateKnown,
		}
	})
	if !diags.HasError() {
		state.Tags = tagSet
	}
}

func populateTagsFromFastUDP(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	tags []sdk.LoadBalancerProfileTag26,
) {
	tagSet, diags := convert.ToSetType(ctx, tags, func(t sdk.LoadBalancerProfileTag26) TagsValue {
		return TagsValue{
			Name:  convert.StrToType(t.Name),
			Value: convert.StrToType(t.Value),
			state: attr.ValueStateKnown,
		}
	})
	if !diags.HasError() {
		state.Tags = tagSet
	}
}

func populateTagsFromCookiePersistence(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	tags []sdk.LoadBalancerProfileTag27,
) {
	tagSet, diags := convert.ToSetType(ctx, tags, func(t sdk.LoadBalancerProfileTag27) TagsValue {
		return TagsValue{
			Name:  convert.StrToType(t.Name),
			Value: convert.StrToType(t.Value),
			state: attr.ValueStateKnown,
		}
	})
	if !diags.HasError() {
		state.Tags = tagSet
	}
}

func populateTagsFromSourceIPPersistence(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	tags []sdk.LoadBalancerProfileTag28,
) {
	tagSet, diags := convert.ToSetType(ctx, tags, func(t sdk.LoadBalancerProfileTag28) TagsValue {
		return TagsValue{
			Name:  convert.StrToType(t.Name),
			Value: convert.StrToType(t.Value),
			state: attr.ValueStateKnown,
		}
	})
	if !diags.HasError() {
		state.Tags = tagSet
	}
}

func populateTagsFromGenericPersistence(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	tags []sdk.LoadBalancerProfileTag29,
) {
	tagSet, diags := convert.ToSetType(ctx, tags, func(t sdk.LoadBalancerProfileTag29) TagsValue {
		return TagsValue{
			Name:  convert.StrToType(t.Name),
			Value: convert.StrToType(t.Value),
			state: attr.ValueStateKnown,
		}
	})
	if !diags.HasError() {
		state.Tags = tagSet
	}
}

func populateTagsFromClientSSL(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	tags []sdk.LoadBalancerProfileTag30,
) {
	tagSet, diags := convert.ToSetType(ctx, tags, func(t sdk.LoadBalancerProfileTag30) TagsValue {
		return TagsValue{
			Name:  convert.StrToType(t.Name),
			Value: convert.StrToType(t.Value),
			state: attr.ValueStateKnown,
		}
	})
	if !diags.HasError() {
		state.Tags = tagSet
	}
}

func populateTagsFromServerSSL(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	tags []sdk.LoadBalancerProfileTag31,
) {
	tagSet, diags := convert.ToSetType(ctx, tags, func(t sdk.LoadBalancerProfileTag31) TagsValue {
		return TagsValue{
			Name:  convert.StrToType(t.Name),
			Value: convert.StrToType(t.Value),
			state: attr.ValueStateKnown,
		}
	})
	if !diags.HasError() {
		state.Tags = tagSet
	}
}
