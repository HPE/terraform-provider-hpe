// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile

import (
	"context"
	"encoding/json"

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

// anyOfRawJSON marshals whichever anyOf variant the SDK populated to raw JSON and
// returns it. This is necessary because the SDK's UnmarshalJSON uses a
// "first-match-wins" strategy for anyOf: it tries each variant in alphabetical
// order and stops at the first one that produces a non-empty JSON object. Since
// all 8 profile Config3 variants share a tags field, any profile that has tags
// causes ClientSSLLoadBalancerProfileConfig3 (first alphabetically) to "win",
// leaving the semantically correct variant nil.
//
// The round-trip is safe because all SDK profile config structs carry an
// AdditionalProperties field tagged `json:",remain"`, which captures every key
// that did not match a declared struct field. When the wrong variant is selected,
// the type-specific fields (e.g. httpIdleTimeout for an HTTP profile) land in
// AdditionalProperties and are re-emitted verbatim on Marshal, so the raw JSON
// faithfully represents all original fields regardless of which variant was used.
func anyOfRawJSON(cfg *sdk.GetLoadBalancerProfile200ResponseLoadBalancerProfileConfig) []byte {
	if cfg == nil {
		return nil
	}

	switch {
	case cfg.ClientSSLLoadBalancerProfileConfig3 != nil:
		b, _ := json.Marshal(cfg.ClientSSLLoadBalancerProfileConfig3)
		return b
	case cfg.CookiePersistenceLoadBalancerProfileConfig3 != nil:
		b, _ := json.Marshal(cfg.CookiePersistenceLoadBalancerProfileConfig3)
		return b
	case cfg.FastTCPLoadBalancerProfileConfig3 != nil:
		b, _ := json.Marshal(cfg.FastTCPLoadBalancerProfileConfig3)
		return b
	case cfg.FastUDPLoadBalancerProfileConfig3 != nil:
		b, _ := json.Marshal(cfg.FastUDPLoadBalancerProfileConfig3)
		return b
	case cfg.GenericPersistenceLoadBalancerProfileConfig3 != nil:
		b, _ := json.Marshal(cfg.GenericPersistenceLoadBalancerProfileConfig3)
		return b
	case cfg.HTTPLoadBalancerProfileConfig3 != nil:
		b, _ := json.Marshal(cfg.HTTPLoadBalancerProfileConfig3)
		return b
	case cfg.ServerSSLLoadBalancerProfileConfig3 != nil:
		b, _ := json.Marshal(cfg.ServerSSLLoadBalancerProfileConfig3)
		return b
	case cfg.SourceIPPersistenceLoadBalancerProfileConfig3 != nil:
		b, _ := json.Marshal(cfg.SourceIPPersistenceLoadBalancerProfileConfig3)
		return b
	}

	return nil
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
// The serviceType parameter is the profile's top-level serviceType string
// (e.g. "LBHttpProfile"). It is used as the discriminator instead of checking
// which anyOf variant is non-nil, because the SDK's UnmarshalJSON may have
// populated the wrong variant (see anyOfRawJSON for details).
//
// Tags are not handled here; the caller reconstructs the single top-level tags
// set from the response separately.
func reconstructConfigBlockFromResponse(
	ctx context.Context,
	state *LoadBalancerProfileModel,
	cfg *sdk.GetLoadBalancerProfile200ResponseLoadBalancerProfileConfig,
	serviceType string,
) {
	rawJSON := anyOfRawJSON(cfg)
	if rawJSON == nil {
		return
	}

	// Decode the raw JSON into the correct typed struct for this serviceType.
	// Because anyOfRawJSON preserves all fields via AdditionalProperties, the
	// round-trip unmarshal into the specific struct type yields fully populated
	// fields regardless of which variant the SDK originally selected.
	switch serviceType {
	case "LBHttpProfile":
		var c sdk.HTTPLoadBalancerProfileConfig3
		if json.Unmarshal(rawJSON, &c) == nil {
			state.ConfigHttp = httpConfigFromResponse(ctx, &c)
		}
	case "LBFastTcpProfile":
		var c sdk.FastTCPLoadBalancerProfileConfig3
		if json.Unmarshal(rawJSON, &c) == nil {
			state.ConfigFastTcp = fastTCPConfigFromResponse(&c)
		}
	case "LBFastUdpProfile":
		var c sdk.FastUDPLoadBalancerProfileConfig3
		if json.Unmarshal(rawJSON, &c) == nil {
			state.ConfigFastUdp = fastUDPConfigFromResponse(&c)
		}
	case "LBCookiePersistenceProfile":
		var c sdk.CookiePersistenceLoadBalancerProfileConfig3
		if json.Unmarshal(rawJSON, &c) == nil {
			state.ConfigCookiePersistence = cookiePersistenceConfigFromResponse(&c)
		}
	case "LBSourceIpPersistenceProfile":
		var c sdk.SourceIPPersistenceLoadBalancerProfileConfig3
		if json.Unmarshal(rawJSON, &c) == nil {
			state.ConfigSourceIpPersistence = sourceIPPersistenceConfigFromResponse(&c)
		}
	case "LBGenericPersistenceProfile":
		var c sdk.GenericPersistenceLoadBalancerProfileConfig3
		if json.Unmarshal(rawJSON, &c) == nil {
			state.ConfigGenericPersistence = genericPersistenceConfigFromResponse(&c)
		}
	case "LBClientSslProfile":
		var c sdk.ClientSSLLoadBalancerProfileConfig3
		if json.Unmarshal(rawJSON, &c) == nil {
			state.ConfigClientSsl = clientSSLConfigFromResponse(&c)
		}
	case "LBServerSslProfile":
		var c sdk.ServerSSLLoadBalancerProfileConfig3
		if json.Unmarshal(rawJSON, &c) == nil {
			state.ConfigServerSsl = serverSSLConfigFromResponse(&c)
		}
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
