// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// allConfigBlocksNull reports whether the model has no typed config_* block
// set. This is the case immediately after import (where ImportState sets only
// id and load_balancer_id), and is used to decide whether the active config
// block should be reconstructed from the API response.
func allConfigBlocksNull(m LoadBalancerProfileModel) bool {
	return m.ConfigHttp.IsNull() &&
		m.ConfigFastTcp.IsNull() &&
		m.ConfigFastUdp.IsNull() &&
		m.ConfigCookiePersistence.IsNull() &&
		m.ConfigSourceIpPersistence.IsNull() &&
		m.ConfigGenericPersistence.IsNull() &&
		m.ConfigClientSsl.IsNull() &&
		m.ConfigServerSsl.IsNull()
}

// reconstructConfigBlockFromResponse populates the single typed config block
// matching the profile's serviceType from the API response.
//
// This is used on IMPORT ONLY. On a normal refresh the config block is
// preserved from prior state, because the API applies its own config defaults
// that Terraform did not send (for example x_forwarded_for defaults to INSERT
// server-side); reconstructing on every read would therefore produce a
// "Provider produced inconsistent result after apply" error. On import there
// is no prior config and no plan to be consistent with, so reconstructing the
// block yields a complete, usable imported resource.
//
// Tags are not handled here; the caller reconstructs the single top-level tags
// set from the response separately.
func reconstructConfigBlockFromResponse(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	cfg *sdk.GetLoadBalancerProfile200ResponseLoadBalancerProfileConfig,
) {
	if cfg == nil {
		return
	}

	switch {
	case cfg.HTTPLoadBalancerProfileConfig3 != nil:
		state.ConfigHttp = httpConfigFromResponse(ctx, cfg.HTTPLoadBalancerProfileConfig3)
	case cfg.FastTCPLoadBalancerProfileConfig3 != nil:
		state.ConfigFastTcp = fastTCPConfigFromResponse(cfg.FastTCPLoadBalancerProfileConfig3)
	case cfg.FastUDPLoadBalancerProfileConfig3 != nil:
		state.ConfigFastUdp = fastUDPConfigFromResponse(cfg.FastUDPLoadBalancerProfileConfig3)
	case cfg.CookiePersistenceLoadBalancerProfileConfig3 != nil:
		state.ConfigCookiePersistence = cookiePersistenceConfigFromResponse(
			cfg.CookiePersistenceLoadBalancerProfileConfig3,
		)
	case cfg.SourceIPPersistenceLoadBalancerProfileConfig3 != nil:
		state.ConfigSourceIpPersistence = sourceIPPersistenceConfigFromResponse(
			cfg.SourceIPPersistenceLoadBalancerProfileConfig3,
		)
	case cfg.GenericPersistenceLoadBalancerProfileConfig3 != nil:
		state.ConfigGenericPersistence = genericPersistenceConfigFromResponse(
			cfg.GenericPersistenceLoadBalancerProfileConfig3,
		)
	case cfg.ClientSSLLoadBalancerProfileConfig3 != nil:
		state.ConfigClientSsl = clientSSLConfigFromResponse(cfg.ClientSSLLoadBalancerProfileConfig3)
	case cfg.ServerSSLLoadBalancerProfileConfig3 != nil:
		state.ConfigServerSsl = serverSSLConfigFromResponse(cfg.ServerSSLLoadBalancerProfileConfig3)
	}
}

func nullableInt64ToType(n sdk.NullableInt64) types.Int64 {
	if n.IsSet() {
		return convert.Int64ToType(n.Get())
	}

	return types.Int64Null()
}

func httpConfigFromResponse(
	ctx context.Context,
	cfg *sdk.HTTPLoadBalancerProfileConfig3,
) ConfigHttpValue {
	httpsRedirect := types.BoolNull()
	if cfg.HttpsRedirect != nil {
		httpsRedirect = convert.StringToBool(ctx, *cfg.HttpsRedirect)
	}

	return ConfigHttpValue{
		HttpIdleTimeout:    nullableInt64ToType(cfg.HttpIdleTimeout),
		RequestHeaderSize:  nullableInt64ToType(cfg.RequestHeaderSize),
		ResponseHeaderSize: nullableInt64ToType(cfg.ResponseHeaderSize),
		HttpsRedirect:      httpsRedirect,
		RedirectAddress:    convert.StrToType(cfg.RedirectAddress),
		XForwardedFor:      convert.StrToType(cfg.XForwardedFor),
		RequestBodySize:    nullableInt64ToType(cfg.RequestBodySize),
		ResponseTimeout:    nullableInt64ToType(cfg.ResponseTimeout),
		NtlmAuthentication: convert.BoolToType(cfg.NtlmAuthentication),
		state:              attr.ValueStateKnown,
	}
}

func fastTCPConfigFromResponse(cfg *sdk.FastTCPLoadBalancerProfileConfig3) ConfigFastTcpValue {
	return ConfigFastTcpValue{
		FastTcpIdleTimeout:     nullableInt64ToType(cfg.FastTcpIdleTimeout),
		HaFlowMirroring:        convert.BoolToType(cfg.HaFlowMirroring),
		ConnectionCloseTimeout: nullableInt64ToType(cfg.ConnectionCloseTimeout),
		state:                  attr.ValueStateKnown,
	}
}

func fastUDPConfigFromResponse(cfg *sdk.FastUDPLoadBalancerProfileConfig3) ConfigFastUdpValue {
	return ConfigFastUdpValue{
		FastUdpIdleTimeout: nullableInt64ToType(cfg.FastUdpIdleTimeout),
		HaFlowMirroring:    convert.BoolToType(cfg.HaFlowMirroring),
		state:              attr.ValueStateKnown,
	}
}

func cookiePersistenceConfigFromResponse(
	cfg *sdk.CookiePersistenceLoadBalancerProfileConfig3,
) ConfigCookiePersistenceValue {
	cookieType := types.StringNull()
	if cfg.CookieType != nil {
		cookieType = convert.CookieTypeFromAPI(cfg.CookieType)
	}

	return ConfigCookiePersistenceValue{
		SharePersistence: convert.BoolToType(cfg.SharePersistence),
		CookieName:       convert.StrToType(cfg.CookieName),
		CookieFallback:   convert.BoolToType(cfg.CookieFallback),
		CookieGarbling:   convert.BoolToType(cfg.CookieGarbling),
		CookieMode:       convert.StrToType(cfg.CookieMode),
		CookieDomain:     convert.StrToType(cfg.CookieDomain),
		CookiePath:       convert.StrToType(cfg.CookiePath),
		CookieType:       cookieType,
		MaxIdleTime:      nullableInt64ToType(cfg.MaxIdleTime),
		MaxCookieAge:     nullableInt64ToType(cfg.MaxCookieAge),
		state:            attr.ValueStateKnown,
	}
}

func sourceIPPersistenceConfigFromResponse(
	cfg *sdk.SourceIPPersistenceLoadBalancerProfileConfig3,
) ConfigSourceIpPersistenceValue {
	return ConfigSourceIpPersistenceValue{
		SharePersistence:        convert.BoolToType(cfg.SharePersistence),
		PurgeEntries:            convert.BoolToType(cfg.PurgeEntries),
		HaPersistenceMirroring:  convert.BoolToType(cfg.HaPersistenceMirroring),
		PersistenceEntryTimeout: nullableInt64ToType(cfg.PersistenceEntryTimeout),
		state:                   attr.ValueStateKnown,
	}
}

func genericPersistenceConfigFromResponse(
	cfg *sdk.GenericPersistenceLoadBalancerProfileConfig3,
) ConfigGenericPersistenceValue {
	return ConfigGenericPersistenceValue{
		SharePersistence:        convert.BoolToType(cfg.SharePersistence),
		HaPersistenceMirroring:  convert.BoolToType(cfg.HaPersistenceMirroring),
		PersistenceEntryTimeout: nullableInt64ToType(cfg.PersistenceEntryTimeout),
		state:                   attr.ValueStateKnown,
	}
}

func clientSSLConfigFromResponse(
	cfg *sdk.ClientSSLLoadBalancerProfileConfig3,
) ConfigClientSslValue {
	return ConfigClientSslValue{
		SslSuite:              convert.StrToType(cfg.SslSuite),
		SessionCache:          convert.BoolToType(cfg.SessionCache),
		SessionCacheTimeout:   nullableInt64ToType(cfg.SessionCacheTimeout),
		PreferServerCipher:    convert.BoolToType(cfg.PreferServerCipher),
		SupportedSslCiphers:   convert.StrSliceToSet(cfg.SupportedSslCiphers),
		SupportedSslProtocols: convert.StrSliceToSet(cfg.SupportedSslProtocols),
		state:                 attr.ValueStateKnown,
	}
}

func serverSSLConfigFromResponse(
	cfg *sdk.ServerSSLLoadBalancerProfileConfig3,
) ConfigServerSslValue {
	return ConfigServerSslValue{
		SslSuite:              convert.StrToType(cfg.SslSuite),
		SessionCache:          convert.BoolToType(cfg.SessionCache),
		SupportedSslCiphers:   convert.StrSliceToSet(cfg.SupportedSslCiphers),
		SupportedSslProtocols: convert.StrSliceToSet(cfg.SupportedSslProtocols),
		state:                 attr.ValueStateKnown,
	}
}
