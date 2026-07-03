// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// profileTypeForServiceType returns the profileType API value for a given serviceType.
func profileTypeForServiceType(serviceType string) string {
	switch serviceType {
	case "LBHttpProfile", "LBFastTcpProfile", "LBFastUdpProfile":
		return "application-profile"
	case "LBCookiePersistenceProfile", "LBSourceIpPersistenceProfile", "LBGenericPersistenceProfile":
		return "persistence-profile"
	case "LBClientSslProfile", "LBServerSslProfile":
		return "ssl-profile"
	default:
		return ""
	}
}

// buildCreateConfig constructs the create config wrapper from the plan.
func buildCreateConfig(
	ctx context.Context,
	plan LoadBalancerProfileModel,
) (*sdk.CreateLoadBalancerProfileRequestLoadBalancerProfileConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	serviceType := plan.ServiceType.ValueString()
	profileType := profileTypeForServiceType(serviceType)

	cfg := &sdk.CreateLoadBalancerProfileRequestLoadBalancerProfileConfig{}

	switch serviceType {
	case "LBHttpProfile":
		variant := &sdk.HTTPLoadBalancerProfileConfig1{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigHttp
		if !c.IsNull() && !c.IsUnknown() {
			setNullableInt64(&variant.HttpIdleTimeout, c.HttpIdleTimeout)
			setNullableInt64(&variant.RequestHeaderSize, c.RequestHeaderSize)
			setNullableInt64(&variant.ResponseHeaderSize, c.ResponseHeaderSize)
			setNullableInt64(&variant.RequestBodySize, c.RequestBodySize)
			setNullableInt64(&variant.ResponseTimeout, c.ResponseTimeout)
			if !c.HttpsRedirect.IsNull() && !c.HttpsRedirect.IsUnknown() {
				variant.HttpsRedirect = convert.BoolTypeToStringPointerOnOff(c.HttpsRedirect)
			}
			setStringPtr(&variant.RedirectAddress, c.RedirectAddress)
			setStringPtr(&variant.XForwardedFor, c.XForwardedFor)
			setBoolPtr(&variant.NtlmAuthentication, c.NtlmAuthentication)
			variant.Tags = buildTagsForHTTPCreate(ctx, plan.Tags)
		}
		cfg.HTTPLoadBalancerProfileConfig1 = variant

	case "LBFastTcpProfile":
		variant := &sdk.FastTCPLoadBalancerProfileConfig1{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigFastTcp
		if !c.IsNull() && !c.IsUnknown() {
			setNullableInt64(&variant.FastTcpIdleTimeout, c.FastTcpIdleTimeout)
			setNullableInt64(&variant.ConnectionCloseTimeout, c.ConnectionCloseTimeout)
			setBoolPtr(&variant.HaFlowMirroring, c.HaFlowMirroring)
			variant.Tags = buildTagsForFastTCPCreate(ctx, plan.Tags)
		}
		cfg.FastTCPLoadBalancerProfileConfig1 = variant

	case "LBFastUdpProfile":
		variant := &sdk.FastUDPLoadBalancerProfileConfig1{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigFastUdp
		if !c.IsNull() && !c.IsUnknown() {
			setNullableInt64(&variant.FastUdpIdleTimeout, c.FastUdpIdleTimeout)
			setBoolPtr(&variant.HaFlowMirroring, c.HaFlowMirroring)
			variant.Tags = buildTagsForFastUDPCreate(ctx, plan.Tags)
		}
		cfg.FastUDPLoadBalancerProfileConfig1 = variant

	case "LBCookiePersistenceProfile":
		variant := &sdk.CookiePersistenceLoadBalancerProfileConfig1{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigCookiePersistence
		if !c.IsNull() && !c.IsUnknown() {
			setBoolPtr(&variant.SharePersistence, c.SharePersistence)
			setStringPtr(&variant.CookieName, c.CookieName)
			setBoolPtr(&variant.CookieFallback, c.CookieFallback)
			setBoolPtr(&variant.CookieGarbling, c.CookieGarbling)
			setStringPtr(&variant.CookieMode, c.CookieMode)
			setStringPtr(&variant.CookieDomain, c.CookieDomain)
			setStringPtr(&variant.CookiePath, c.CookiePath)
			if !c.CookieType.IsNull() && !c.CookieType.IsUnknown() {
				variant.CookieType = convert.CookieTypeToAPI(c.CookieType)
			}
			setNullableInt64(&variant.MaxIdleTime, c.MaxIdleTime)
			setNullableInt64(&variant.MaxCookieAge, c.MaxCookieAge)
			variant.Tags = buildTagsForCookieCreate(ctx, plan.Tags)
		}
		cfg.CookiePersistenceLoadBalancerProfileConfig1 = variant

	case "LBSourceIpPersistenceProfile":
		variant := &sdk.SourceIPPersistenceLoadBalancerProfileConfig1{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigSourceIpPersistence
		if !c.IsNull() && !c.IsUnknown() {
			setBoolPtr(&variant.SharePersistence, c.SharePersistence)
			setBoolPtr(&variant.PurgeEntries, c.PurgeEntries)
			setBoolPtr(&variant.HaPersistenceMirroring, c.HaPersistenceMirroring)
			setNullableInt64(&variant.PersistenceEntryTimeout, c.PersistenceEntryTimeout)
			variant.Tags = buildTagsForSourceIPCreate(ctx, plan.Tags)
		}
		cfg.SourceIPPersistenceLoadBalancerProfileConfig1 = variant

	case "LBGenericPersistenceProfile":
		variant := &sdk.GenericPersistenceLoadBalancerProfileConfig1{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigGenericPersistence
		if !c.IsNull() && !c.IsUnknown() {
			setBoolPtr(&variant.SharePersistence, c.SharePersistence)
			setBoolPtr(&variant.HaPersistenceMirroring, c.HaPersistenceMirroring)
			setNullableInt64(&variant.PersistenceEntryTimeout, c.PersistenceEntryTimeout)
			variant.Tags = buildTagsForGenericCreate(ctx, plan.Tags)
		}
		cfg.GenericPersistenceLoadBalancerProfileConfig1 = variant

	case "LBClientSslProfile":
		variant := &sdk.ClientSSLLoadBalancerProfileConfig1{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigClientSsl
		if !c.IsNull() && !c.IsUnknown() {
			setStringPtr(&variant.SslSuite, c.SslSuite)
			setBoolPtr(&variant.SessionCache, c.SessionCache)
			setNullableInt64(&variant.SessionCacheTimeout, c.SessionCacheTimeout)
			setBoolPtr(&variant.PreferServerCipher, c.PreferServerCipher)
			if !c.SupportedSslCiphers.IsNull() && !c.SupportedSslCiphers.IsUnknown() {
				ciphers, err := convert.SetToStrSlice(c.SupportedSslCiphers)
				if err == nil {
					variant.SupportedSslCiphers = ciphers
				}
			}
			if !c.SupportedSslProtocols.IsNull() && !c.SupportedSslProtocols.IsUnknown() {
				protocols, err := convert.SetToStrSlice(c.SupportedSslProtocols)
				if err == nil {
					variant.SupportedSslProtocols = protocols
				}
			}
			variant.Tags = buildTagsForClientSSLCreate(ctx, plan.Tags)
		}
		cfg.ClientSSLLoadBalancerProfileConfig1 = variant

	case "LBServerSslProfile":
		variant := &sdk.ServerSSLLoadBalancerProfileConfig1{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigServerSsl
		if !c.IsNull() && !c.IsUnknown() {
			setStringPtr(&variant.SslSuite, c.SslSuite)
			setBoolPtr(&variant.SessionCache, c.SessionCache)
			if !c.SupportedSslCiphers.IsNull() && !c.SupportedSslCiphers.IsUnknown() {
				ciphers, err := convert.SetToStrSlice(c.SupportedSslCiphers)
				if err == nil {
					variant.SupportedSslCiphers = ciphers
				}
			}
			if !c.SupportedSslProtocols.IsNull() && !c.SupportedSslProtocols.IsUnknown() {
				protocols, err := convert.SetToStrSlice(c.SupportedSslProtocols)
				if err == nil {
					variant.SupportedSslProtocols = protocols
				}
			}
			variant.Tags = buildTagsForServerSSLCreate(ctx, plan.Tags)
		}
		cfg.ServerSSLLoadBalancerProfileConfig1 = variant
	}

	return cfg, diags
}

// buildUpdateConfig constructs the update config wrapper from the plan.
func buildUpdateConfig(
	ctx context.Context,
	plan LoadBalancerProfileModel,
) (*sdk.UpdateLoadBalancerProfileRequestLoadBalancerProfileConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	serviceType := plan.ServiceType.ValueString()
	profileType := profileTypeForServiceType(serviceType)

	cfg := &sdk.UpdateLoadBalancerProfileRequestLoadBalancerProfileConfig{}

	switch serviceType {
	case "LBHttpProfile":
		variant := &sdk.HTTPLoadBalancerProfileConfig4{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigHttp
		if !c.IsNull() && !c.IsUnknown() {
			setNullableInt64(&variant.HttpIdleTimeout, c.HttpIdleTimeout)
			setNullableInt64(&variant.RequestHeaderSize, c.RequestHeaderSize)
			setNullableInt64(&variant.ResponseHeaderSize, c.ResponseHeaderSize)
			setNullableInt64(&variant.RequestBodySize, c.RequestBodySize)
			setNullableInt64(&variant.ResponseTimeout, c.ResponseTimeout)
			if !c.HttpsRedirect.IsNull() && !c.HttpsRedirect.IsUnknown() {
				variant.HttpsRedirect = convert.BoolTypeToStringPointerOnOff(c.HttpsRedirect)
			}
			setStringPtr(&variant.RedirectAddress, c.RedirectAddress)
			setStringPtr(&variant.XForwardedFor, c.XForwardedFor)
			setBoolPtr(&variant.NtlmAuthentication, c.NtlmAuthentication)
			variant.Tags = buildTagsForHTTPUpdate(ctx, plan.Tags)
		}
		cfg.HTTPLoadBalancerProfileConfig4 = variant

	case "LBFastTcpProfile":
		variant := &sdk.FastTCPLoadBalancerProfileConfig4{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigFastTcp
		if !c.IsNull() && !c.IsUnknown() {
			setNullableInt64(&variant.FastTcpIdleTimeout, c.FastTcpIdleTimeout)
			setNullableInt64(&variant.ConnectionCloseTimeout, c.ConnectionCloseTimeout)
			setBoolPtr(&variant.HaFlowMirroring, c.HaFlowMirroring)
			variant.Tags = buildTagsForFastTCPUpdate(ctx, plan.Tags)
		}
		cfg.FastTCPLoadBalancerProfileConfig4 = variant

	case "LBFastUdpProfile":
		variant := &sdk.FastUDPLoadBalancerProfileConfig4{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigFastUdp
		if !c.IsNull() && !c.IsUnknown() {
			setNullableInt64(&variant.FastUdpIdleTimeout, c.FastUdpIdleTimeout)
			setBoolPtr(&variant.HaFlowMirroring, c.HaFlowMirroring)
			variant.Tags = buildTagsForFastUDPUpdate(ctx, plan.Tags)
		}
		cfg.FastUDPLoadBalancerProfileConfig4 = variant

	case "LBCookiePersistenceProfile":
		variant := &sdk.CookiePersistenceLoadBalancerProfileConfig4{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigCookiePersistence
		if !c.IsNull() && !c.IsUnknown() {
			setBoolPtr(&variant.SharePersistence, c.SharePersistence)
			setStringPtr(&variant.CookieName, c.CookieName)
			setBoolPtr(&variant.CookieFallback, c.CookieFallback)
			setBoolPtr(&variant.CookieGarbling, c.CookieGarbling)
			setStringPtr(&variant.CookieMode, c.CookieMode)
			setStringPtr(&variant.CookieDomain, c.CookieDomain)
			setStringPtr(&variant.CookiePath, c.CookiePath)
			if !c.CookieType.IsNull() && !c.CookieType.IsUnknown() {
				variant.CookieType = convert.CookieTypeToAPI(c.CookieType)
			}
			setNullableInt64(&variant.MaxIdleTime, c.MaxIdleTime)
			setNullableInt64(&variant.MaxCookieAge, c.MaxCookieAge)
			variant.Tags = buildTagsForCookieUpdate(ctx, plan.Tags)
		}
		cfg.CookiePersistenceLoadBalancerProfileConfig4 = variant

	case "LBSourceIpPersistenceProfile":
		variant := &sdk.SourceIPPersistenceLoadBalancerProfileConfig4{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigSourceIpPersistence
		if !c.IsNull() && !c.IsUnknown() {
			setBoolPtr(&variant.SharePersistence, c.SharePersistence)
			setBoolPtr(&variant.PurgeEntries, c.PurgeEntries)
			setBoolPtr(&variant.HaPersistenceMirroring, c.HaPersistenceMirroring)
			setNullableInt64(&variant.PersistenceEntryTimeout, c.PersistenceEntryTimeout)
			variant.Tags = buildTagsForSourceIPUpdate(ctx, plan.Tags)
		}
		cfg.SourceIPPersistenceLoadBalancerProfileConfig4 = variant

	case "LBGenericPersistenceProfile":
		variant := &sdk.GenericPersistenceLoadBalancerProfileConfig4{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigGenericPersistence
		if !c.IsNull() && !c.IsUnknown() {
			setBoolPtr(&variant.SharePersistence, c.SharePersistence)
			setBoolPtr(&variant.HaPersistenceMirroring, c.HaPersistenceMirroring)
			setNullableInt64(&variant.PersistenceEntryTimeout, c.PersistenceEntryTimeout)
			variant.Tags = buildTagsForGenericUpdate(ctx, plan.Tags)
		}
		cfg.GenericPersistenceLoadBalancerProfileConfig4 = variant

	case "LBClientSslProfile":
		variant := &sdk.ClientSSLLoadBalancerProfileConfig4{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigClientSsl
		if !c.IsNull() && !c.IsUnknown() {
			setStringPtr(&variant.SslSuite, c.SslSuite)
			setBoolPtr(&variant.SessionCache, c.SessionCache)
			setNullableInt64(&variant.SessionCacheTimeout, c.SessionCacheTimeout)
			setBoolPtr(&variant.PreferServerCipher, c.PreferServerCipher)
			if !c.SupportedSslCiphers.IsNull() && !c.SupportedSslCiphers.IsUnknown() {
				ciphers, err := convert.SetToStrSlice(c.SupportedSslCiphers)
				if err == nil {
					variant.SupportedSslCiphers = ciphers
				}
			}
			if !c.SupportedSslProtocols.IsNull() && !c.SupportedSslProtocols.IsUnknown() {
				protocols, err := convert.SetToStrSlice(c.SupportedSslProtocols)
				if err == nil {
					variant.SupportedSslProtocols = protocols
				}
			}
			variant.Tags = buildTagsForClientSSLUpdate(ctx, plan.Tags)
		}
		cfg.ClientSSLLoadBalancerProfileConfig4 = variant

	case "LBServerSslProfile":
		variant := &sdk.ServerSSLLoadBalancerProfileConfig4{
			ProfileType: sdk.PtrString(profileType),
		}
		c := plan.ConfigServerSsl
		if !c.IsNull() && !c.IsUnknown() {
			setStringPtr(&variant.SslSuite, c.SslSuite)
			setBoolPtr(&variant.SessionCache, c.SessionCache)
			if !c.SupportedSslCiphers.IsNull() && !c.SupportedSslCiphers.IsUnknown() {
				ciphers, err := convert.SetToStrSlice(c.SupportedSslCiphers)
				if err == nil {
					variant.SupportedSslCiphers = ciphers
				}
			}
			if !c.SupportedSslProtocols.IsNull() && !c.SupportedSslProtocols.IsUnknown() {
				protocols, err := convert.SetToStrSlice(c.SupportedSslProtocols)
				if err == nil {
					variant.SupportedSslProtocols = protocols
				}
			}
			variant.Tags = buildTagsForServerSSLUpdate(ctx, plan.Tags)
		}
		cfg.ServerSSLLoadBalancerProfileConfig4 = variant
	}

	return cfg, diags
}

// readTagsFromConfig extracts tags from the read response config and returns them
// as a Terraform Set, preserving user-specified name casing from the plan/state.
func readTagsFromConfig(
	ctx context.Context,
	cfg *sdk.GetLoadBalancerProfile200ResponseLoadBalancerProfileConfig,
	priorTags types.Set,
) types.Set {
	if cfg == nil {
		return priorTags
	}

	// Collect API tags as name→value pairs
	type tagPair struct {
		name  string
		value string
	}

	var apiTags []tagPair

	switch {
	case cfg.HTTPLoadBalancerProfileConfig3 != nil:
		for _, t := range cfg.HTTPLoadBalancerProfileConfig3.Tags {
			apiTags = append(apiTags, tagPair{
				name:  ptrStr(t.Tag),
				value: ptrStr(t.Scope),
			})
		}
	case cfg.FastTCPLoadBalancerProfileConfig3 != nil:
		for _, t := range cfg.FastTCPLoadBalancerProfileConfig3.Tags {
			apiTags = append(apiTags, tagPair{
				name:  ptrStr(t.Tag),
				value: ptrStr(t.Scope),
			})
		}
	case cfg.FastUDPLoadBalancerProfileConfig3 != nil:
		for _, t := range cfg.FastUDPLoadBalancerProfileConfig3.Tags {
			apiTags = append(apiTags, tagPair{
				name:  ptrStr(t.Tag),
				value: ptrStr(t.Scope),
			})
		}
	case cfg.CookiePersistenceLoadBalancerProfileConfig3 != nil:
		for _, t := range cfg.CookiePersistenceLoadBalancerProfileConfig3.Tags {
			apiTags = append(apiTags, tagPair{
				name:  ptrStr(t.Tag),
				value: ptrStr(t.Scope),
			})
		}
	case cfg.SourceIPPersistenceLoadBalancerProfileConfig3 != nil:
		for _, t := range cfg.SourceIPPersistenceLoadBalancerProfileConfig3.Tags {
			apiTags = append(apiTags, tagPair{
				name:  ptrStr(t.Tag),
				value: ptrStr(t.Scope),
			})
		}
	case cfg.GenericPersistenceLoadBalancerProfileConfig3 != nil:
		for _, t := range cfg.GenericPersistenceLoadBalancerProfileConfig3.Tags {
			apiTags = append(apiTags, tagPair{
				name:  ptrStr(t.Tag),
				value: ptrStr(t.Scope),
			})
		}
	case cfg.ClientSSLLoadBalancerProfileConfig3 != nil:
		for _, t := range cfg.ClientSSLLoadBalancerProfileConfig3.Tags {
			apiTags = append(apiTags, tagPair{
				name:  ptrStr(t.Tag),
				value: ptrStr(t.Scope),
			})
		}
	case cfg.ServerSSLLoadBalancerProfileConfig3 != nil:
		for _, t := range cfg.ServerSSLLoadBalancerProfileConfig3.Tags {
			apiTags = append(apiTags, tagPair{
				name:  ptrStr(t.Tag),
				value: ptrStr(t.Scope),
			})
		}
	}

	if len(apiTags) == 0 {
		return types.SetNull(TagsValue{}.Type(ctx))
	}

	vals := make([]attr.Value, 0, len(apiTags))
	for _, t := range apiTags {
		tv := TagsValue{
			Name:  types.StringValue(t.name),
			Value: types.StringValue(t.value),
			state: attr.ValueStateKnown,
		}
		vals = append(vals, tv)
	}

	setVal, diags := types.SetValue(TagsValue{}.Type(ctx), vals)
	if diags.HasError() {
		return priorTags
	}

	return setVal
}

// --- Helper functions ---

func setStringPtr(dst **string, src types.String) {
	if !src.IsNull() && !src.IsUnknown() {
		*dst = sdk.PtrString(src.ValueString())
	}
}

func setBoolPtr(dst **bool, src types.Bool) {
	if !src.IsNull() && !src.IsUnknown() {
		*dst = sdk.PtrBool(src.ValueBool())
	}
}

func setNullableInt64(dst *sdk.NullableInt64, src types.Int64) {
	if !src.IsNull() && !src.IsUnknown() {
		dst.Set(sdk.PtrInt64(src.ValueInt64()))
	}
}

func ptrStr(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

// --- Tag builders per variant/context ---

func tagsFromPlan(ctx context.Context, tags types.Set) []struct {
	name  string
	value string
} {
	if tags.IsNull() || tags.IsUnknown() {
		return nil
	}

	var elems []TagsValue
	if d := tags.ElementsAs(ctx, &elems, false); d.HasError() {
		return nil
	}

	result := make([]struct {
		name  string
		value string
	}, 0, len(elems))
	for _, e := range elems {
		result = append(result, struct {
			name  string
			value string
		}{
			name:  e.Name.ValueString(),
			value: e.Value.ValueString(),
		})
	}

	return result
}

func buildTagsForHTTPCreate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag8 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag8, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag8{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}

func buildTagsForFastTCPCreate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag9 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag9, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag9{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}

func buildTagsForFastUDPCreate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag10 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag10, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag10{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}

func buildTagsForCookieCreate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag11 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag11, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag11{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}

func buildTagsForSourceIPCreate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag12 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag12, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag12{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}

func buildTagsForGenericCreate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag13 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag13, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag13{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}

func buildTagsForClientSSLCreate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag14 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag14, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag14{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}

func buildTagsForServerSSLCreate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag15 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag15, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag15{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}

// Update tag builders use Config4 tag types (Tag32-39).
func buildTagsForHTTPUpdate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag32 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag32, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag32{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}

func buildTagsForFastTCPUpdate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag33 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag33, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag33{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}

func buildTagsForFastUDPUpdate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag34 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag34, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag34{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}

func buildTagsForCookieUpdate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag35 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag35, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag35{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}

func buildTagsForSourceIPUpdate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag36 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag36, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag36{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}

func buildTagsForGenericUpdate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag37 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag37, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag37{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}

func buildTagsForClientSSLUpdate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag38 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag38, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag38{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}

func buildTagsForServerSSLUpdate(ctx context.Context, tags types.Set) []sdk.LoadBalancerProfileTag39 {
	pairs := tagsFromPlan(ctx, tags)
	if pairs == nil {
		return nil
	}
	result := make([]sdk.LoadBalancerProfileTag39, 0, len(pairs))
	for _, p := range pairs {
		result = append(result, sdk.LoadBalancerProfileTag39{
			Tag:   sdk.PtrString(p.name),
			Scope: sdk.PtrString(p.value),
		})
	}

	return result
}
