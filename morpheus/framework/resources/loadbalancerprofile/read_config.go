// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

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

// knownOr returns prior when it is known, otherwise the value resolved from the
// API response.
//
// Optional+Computed attributes that the configuration omits arrive as unknown
// in the plan. They must be resolved before state is written, or Terraform
// reports "Provider returned invalid result object after apply". Attributes the
// practitioner did configure are known and are always preserved, so the API's
// own defaults never overwrite them.
func knownOr[T attr.Value](prior, fromAPI T) T {
	if prior.IsUnknown() {
		return fromAPI
	}

	return prior
}

// mergeConfigBlocks writes the typed config_* blocks into state, preserving
// every value the practitioner configured while resolving any attribute that
// arrived unknown from the API response.
//
// Copying the prior block verbatim would leak the plan's unknown attributes
// into state; rebuilding it purely from the response would overwrite configured
// values with the API's own defaults (for example x_forwarded_for defaults to
// INSERT server-side) and produce "Provider produced inconsistent result after
// apply". Merging per attribute avoids both failure modes.
func mergeConfigBlocks(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	prior LoadBalancerProfileModel,
	cfg *sdk.GetLoadBalancerProfile200ResponseLoadBalancerProfileConfig,
) {
	var (
		httpCfg      *sdk.HTTPLoadBalancerProfileConfig3
		fastTCPCfg   *sdk.FastTCPLoadBalancerProfileConfig3
		fastUDPCfg   *sdk.FastUDPLoadBalancerProfileConfig3
		cookieCfg    *sdk.CookiePersistenceLoadBalancerProfileConfig3
		sourceIPCfg  *sdk.SourceIPPersistenceLoadBalancerProfileConfig3
		genericCfg   *sdk.GenericPersistenceLoadBalancerProfileConfig3
		clientSSLCfg *sdk.ClientSSLLoadBalancerProfileConfig3
		serverSSLCfg *sdk.ServerSSLLoadBalancerProfileConfig3
	)

	if cfg != nil {
		httpCfg = cfg.HTTPLoadBalancerProfileConfig3
		fastTCPCfg = cfg.FastTCPLoadBalancerProfileConfig3
		fastUDPCfg = cfg.FastUDPLoadBalancerProfileConfig3
		cookieCfg = cfg.CookiePersistenceLoadBalancerProfileConfig3
		sourceIPCfg = cfg.SourceIPPersistenceLoadBalancerProfileConfig3
		genericCfg = cfg.GenericPersistenceLoadBalancerProfileConfig3
		clientSSLCfg = cfg.ClientSSLLoadBalancerProfileConfig3
		serverSSLCfg = cfg.ServerSSLLoadBalancerProfileConfig3
	}

	state.ConfigHttp = mergeHTTPConfig(ctx, prior.ConfigHttp, httpCfg)
	state.ConfigFastTcp = mergeFastTCPConfig(prior.ConfigFastTcp, fastTCPCfg)
	state.ConfigFastUdp = mergeFastUDPConfig(prior.ConfigFastUdp, fastUDPCfg)
	state.ConfigCookiePersistence = mergeCookiePersistenceConfig(
		prior.ConfigCookiePersistence, cookieCfg,
	)
	state.ConfigSourceIpPersistence = mergeSourceIPPersistenceConfig(
		prior.ConfigSourceIpPersistence, sourceIPCfg,
	)
	state.ConfigGenericPersistence = mergeGenericPersistenceConfig(
		prior.ConfigGenericPersistence, genericCfg,
	)
	state.ConfigClientSsl = mergeClientSSLConfig(prior.ConfigClientSsl, clientSSLCfg)
	state.ConfigServerSsl = mergeServerSSLConfig(prior.ConfigServerSsl, serverSSLCfg)
}

func mergeHTTPConfig(
	ctx context.Context,
	prior ConfigHttpValue,
	cfg *sdk.HTTPLoadBalancerProfileConfig3,
) ConfigHttpValue {
	if prior.IsNull() {
		return prior
	}

	var fromAPI ConfigHttpValue
	if cfg != nil {
		fromAPI = httpConfigFromResponse(ctx, cfg)
	}

	if prior.IsUnknown() {
		if cfg == nil {
			return NewConfigHttpValueNull()
		}

		return fromAPI
	}

	return ConfigHttpValue{
		HttpIdleTimeout:    knownOr(prior.HttpIdleTimeout, fromAPI.HttpIdleTimeout),
		HttpsRedirect:      knownOr(prior.HttpsRedirect, fromAPI.HttpsRedirect),
		NtlmAuthentication: knownOr(prior.NtlmAuthentication, fromAPI.NtlmAuthentication),
		RedirectAddress:    knownOr(prior.RedirectAddress, fromAPI.RedirectAddress),
		RequestBodySize:    knownOr(prior.RequestBodySize, fromAPI.RequestBodySize),
		RequestHeaderSize:  knownOr(prior.RequestHeaderSize, fromAPI.RequestHeaderSize),
		ResponseHeaderSize: knownOr(prior.ResponseHeaderSize, fromAPI.ResponseHeaderSize),
		ResponseTimeout:    knownOr(prior.ResponseTimeout, fromAPI.ResponseTimeout),
		XForwardedFor:      knownOr(prior.XForwardedFor, fromAPI.XForwardedFor),
		state:              attr.ValueStateKnown,
	}
}

func mergeFastTCPConfig(
	prior ConfigFastTcpValue,
	cfg *sdk.FastTCPLoadBalancerProfileConfig3,
) ConfigFastTcpValue {
	if prior.IsNull() {
		return prior
	}

	var fromAPI ConfigFastTcpValue
	if cfg != nil {
		fromAPI = fastTCPConfigFromResponse(cfg)
	}

	if prior.IsUnknown() {
		if cfg == nil {
			return NewConfigFastTcpValueNull()
		}

		return fromAPI
	}

	return ConfigFastTcpValue{
		ConnectionCloseTimeout: knownOr(
			prior.ConnectionCloseTimeout, fromAPI.ConnectionCloseTimeout,
		),
		FastTcpIdleTimeout: knownOr(prior.FastTcpIdleTimeout, fromAPI.FastTcpIdleTimeout),
		HaFlowMirroring:    knownOr(prior.HaFlowMirroring, fromAPI.HaFlowMirroring),
		state:              attr.ValueStateKnown,
	}
}

func mergeFastUDPConfig(
	prior ConfigFastUdpValue,
	cfg *sdk.FastUDPLoadBalancerProfileConfig3,
) ConfigFastUdpValue {
	if prior.IsNull() {
		return prior
	}

	var fromAPI ConfigFastUdpValue
	if cfg != nil {
		fromAPI = fastUDPConfigFromResponse(cfg)
	}

	if prior.IsUnknown() {
		if cfg == nil {
			return NewConfigFastUdpValueNull()
		}

		return fromAPI
	}

	return ConfigFastUdpValue{
		FastUdpIdleTimeout: knownOr(prior.FastUdpIdleTimeout, fromAPI.FastUdpIdleTimeout),
		HaFlowMirroring:    knownOr(prior.HaFlowMirroring, fromAPI.HaFlowMirroring),
		state:              attr.ValueStateKnown,
	}
}

func mergeCookiePersistenceConfig(
	prior ConfigCookiePersistenceValue,
	cfg *sdk.CookiePersistenceLoadBalancerProfileConfig3,
) ConfigCookiePersistenceValue {
	if prior.IsNull() {
		return prior
	}

	var fromAPI ConfigCookiePersistenceValue
	if cfg != nil {
		fromAPI = cookiePersistenceConfigFromResponse(cfg)
	}

	if prior.IsUnknown() {
		if cfg == nil {
			return NewConfigCookiePersistenceValueNull()
		}

		return fromAPI
	}

	return ConfigCookiePersistenceValue{
		CookieDomain:     knownOr(prior.CookieDomain, fromAPI.CookieDomain),
		CookieFallback:   knownOr(prior.CookieFallback, fromAPI.CookieFallback),
		CookieGarbling:   knownOr(prior.CookieGarbling, fromAPI.CookieGarbling),
		CookieMode:       knownOr(prior.CookieMode, fromAPI.CookieMode),
		CookieName:       knownOr(prior.CookieName, fromAPI.CookieName),
		CookiePath:       knownOr(prior.CookiePath, fromAPI.CookiePath),
		CookieType:       knownOr(prior.CookieType, fromAPI.CookieType),
		MaxCookieAge:     knownOr(prior.MaxCookieAge, fromAPI.MaxCookieAge),
		MaxIdleTime:      knownOr(prior.MaxIdleTime, fromAPI.MaxIdleTime),
		SharePersistence: knownOr(prior.SharePersistence, fromAPI.SharePersistence),
		state:            attr.ValueStateKnown,
	}
}

func mergeSourceIPPersistenceConfig(
	prior ConfigSourceIpPersistenceValue,
	cfg *sdk.SourceIPPersistenceLoadBalancerProfileConfig3,
) ConfigSourceIpPersistenceValue {
	if prior.IsNull() {
		return prior
	}

	var fromAPI ConfigSourceIpPersistenceValue
	if cfg != nil {
		fromAPI = sourceIPPersistenceConfigFromResponse(cfg)
	}

	if prior.IsUnknown() {
		if cfg == nil {
			return NewConfigSourceIpPersistenceValueNull()
		}

		return fromAPI
	}

	return ConfigSourceIpPersistenceValue{
		HaPersistenceMirroring: knownOr(
			prior.HaPersistenceMirroring, fromAPI.HaPersistenceMirroring,
		),
		PersistenceEntryTimeout: knownOr(
			prior.PersistenceEntryTimeout, fromAPI.PersistenceEntryTimeout,
		),
		PurgeEntries:     knownOr(prior.PurgeEntries, fromAPI.PurgeEntries),
		SharePersistence: knownOr(prior.SharePersistence, fromAPI.SharePersistence),
		state:            attr.ValueStateKnown,
	}
}

func mergeGenericPersistenceConfig(
	prior ConfigGenericPersistenceValue,
	cfg *sdk.GenericPersistenceLoadBalancerProfileConfig3,
) ConfigGenericPersistenceValue {
	if prior.IsNull() {
		return prior
	}

	var fromAPI ConfigGenericPersistenceValue
	if cfg != nil {
		fromAPI = genericPersistenceConfigFromResponse(cfg)
	}

	if prior.IsUnknown() {
		if cfg == nil {
			return NewConfigGenericPersistenceValueNull()
		}

		return fromAPI
	}

	return ConfigGenericPersistenceValue{
		HaPersistenceMirroring: knownOr(
			prior.HaPersistenceMirroring, fromAPI.HaPersistenceMirroring,
		),
		PersistenceEntryTimeout: knownOr(
			prior.PersistenceEntryTimeout, fromAPI.PersistenceEntryTimeout,
		),
		SharePersistence: knownOr(prior.SharePersistence, fromAPI.SharePersistence),
		state:            attr.ValueStateKnown,
	}
}

func mergeClientSSLConfig(
	prior ConfigClientSslValue,
	cfg *sdk.ClientSSLLoadBalancerProfileConfig3,
) ConfigClientSslValue {
	if prior.IsNull() {
		return prior
	}

	var fromAPI ConfigClientSslValue
	if cfg != nil {
		fromAPI = clientSSLConfigFromResponse(cfg)
	}

	if prior.IsUnknown() {
		if cfg == nil {
			return NewConfigClientSslValueNull()
		}

		return fromAPI
	}

	return ConfigClientSslValue{
		PreferServerCipher:  knownOr(prior.PreferServerCipher, fromAPI.PreferServerCipher),
		SessionCache:        knownOr(prior.SessionCache, fromAPI.SessionCache),
		SessionCacheTimeout: knownOr(prior.SessionCacheTimeout, fromAPI.SessionCacheTimeout),
		SslSuite:            knownOr(prior.SslSuite, fromAPI.SslSuite),
		SupportedSslCiphers: knownOr(
			prior.SupportedSslCiphers, fromAPI.SupportedSslCiphers,
		),
		SupportedSslProtocols: knownOr(
			prior.SupportedSslProtocols, fromAPI.SupportedSslProtocols,
		),
		state: attr.ValueStateKnown,
	}
}

func mergeServerSSLConfig(
	prior ConfigServerSslValue,
	cfg *sdk.ServerSSLLoadBalancerProfileConfig3,
) ConfigServerSslValue {
	if prior.IsNull() {
		return prior
	}

	var fromAPI ConfigServerSslValue
	if cfg != nil {
		fromAPI = serverSSLConfigFromResponse(cfg)
	}

	if prior.IsUnknown() {
		if cfg == nil {
			return NewConfigServerSslValueNull()
		}

		return fromAPI
	}

	return ConfigServerSslValue{
		SessionCache: knownOr(prior.SessionCache, fromAPI.SessionCache),
		SslSuite:     knownOr(prior.SslSuite, fromAPI.SslSuite),
		SupportedSslCiphers: knownOr(
			prior.SupportedSslCiphers, fromAPI.SupportedSslCiphers,
		),
		SupportedSslProtocols: knownOr(
			prior.SupportedSslProtocols, fromAPI.SupportedSslProtocols,
		),
		state: attr.ValueStateKnown,
	}
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
