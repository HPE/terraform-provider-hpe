// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestMergeHTTPConfig covers the two failure modes the merge exists to prevent:
//
//   - copying the plan verbatim leaks unknown Optional+Computed attributes into
//     state ("Provider returned invalid result object after apply");
//   - rebuilding purely from the API response overwrites configured values with
//     server-side defaults ("Provider produced inconsistent result after
//     apply").
func TestMergeHTTPConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	apiCfg := &sdk.HTTPLoadBalancerProfileConfig3{
		RedirectAddress: sdk.PtrString("https://api.example.com"),
		XForwardedFor:   sdk.PtrString("INSERT"),
	}

	t.Run("unknown attributes are resolved from the API response", func(t *testing.T) {
		t.Parallel()

		prior := ConfigHttpValue{
			RedirectAddress: types.StringUnknown(),
			XForwardedFor:   types.StringUnknown(),
			RequestBodySize: types.Int64Unknown(),
			state:           attr.ValueStateKnown,
		}

		got := mergeHTTPConfig(ctx, prior, apiCfg)

		if got.RedirectAddress.IsUnknown() || got.XForwardedFor.IsUnknown() ||
			got.RequestBodySize.IsUnknown() {
			t.Fatal("unknown value leaked into state")
		}

		if got.XForwardedFor.ValueString() != "INSERT" {
			t.Errorf("x_forwarded_for = %q, want the API value %q",
				got.XForwardedFor.ValueString(), "INSERT")
		}

		// Not returned by the API and unknown in the plan: must resolve to null.
		if !got.RequestBodySize.IsNull() {
			t.Errorf("request_body_size = %v, want null", got.RequestBodySize)
		}
	})

	t.Run("configured values are preserved over API defaults", func(t *testing.T) {
		t.Parallel()

		prior := ConfigHttpValue{
			RedirectAddress: types.StringValue("https://configured.example.com"),
			XForwardedFor:   types.StringValue("REPLACE"),
			state:           attr.ValueStateKnown,
		}

		got := mergeHTTPConfig(ctx, prior, apiCfg)

		if got.XForwardedFor.ValueString() != "REPLACE" {
			t.Errorf("x_forwarded_for = %q, want the configured value %q",
				got.XForwardedFor.ValueString(), "REPLACE")
		}

		if got.RedirectAddress.ValueString() != "https://configured.example.com" {
			t.Errorf("redirect_address = %q, want the configured value",
				got.RedirectAddress.ValueString())
		}
	})

	t.Run("explicit null is preserved, not overwritten by the API", func(t *testing.T) {
		t.Parallel()

		prior := ConfigHttpValue{
			XForwardedFor: types.StringNull(),
			state:         attr.ValueStateKnown,
		}

		got := mergeHTTPConfig(ctx, prior, apiCfg)

		if !got.XForwardedFor.IsNull() {
			t.Errorf("x_forwarded_for = %v, want null to be preserved", got.XForwardedFor)
		}
	})

	t.Run("null block stays null", func(t *testing.T) {
		t.Parallel()

		if got := mergeHTTPConfig(ctx, NewConfigHttpValueNull(), apiCfg); !got.IsNull() {
			t.Errorf("got %v, want a null block", got)
		}
	})

	t.Run("nil API config still resolves unknowns to null", func(t *testing.T) {
		t.Parallel()

		prior := ConfigHttpValue{
			RedirectAddress: types.StringUnknown(),
			state:           attr.ValueStateKnown,
		}

		got := mergeHTTPConfig(ctx, prior, nil)

		if got.RedirectAddress.IsUnknown() {
			t.Error("unknown value leaked into state when the API returned no config")
		}

		if !got.RedirectAddress.IsNull() {
			t.Errorf("redirect_address = %v, want null", got.RedirectAddress)
		}
	})
}

// TestMergeCookiePersistenceConfig covers cookie_path, max_cookie_age and
// max_idle_time — the Optional+Computed attributes without a default that were
// reported as still-unknown after apply.
func TestMergeCookiePersistenceConfig(t *testing.T) {
	t.Parallel()

	apiCfg := &sdk.CookiePersistenceLoadBalancerProfileConfig3{
		CookiePath: sdk.PtrString("/from-api"),
	}

	prior := ConfigCookiePersistenceValue{
		CookieName:   types.StringValue("configured-cookie"),
		CookiePath:   types.StringUnknown(),
		MaxCookieAge: types.Int64Unknown(),
		MaxIdleTime:  types.Int64Unknown(),
		state:        attr.ValueStateKnown,
	}

	got := mergeCookiePersistenceConfig(prior, apiCfg)

	if got.CookiePath.IsUnknown() || got.MaxCookieAge.IsUnknown() ||
		got.MaxIdleTime.IsUnknown() {
		t.Fatal("unknown value leaked into state")
	}

	if got.CookieName.ValueString() != "configured-cookie" {
		t.Errorf("cookie_name = %q, want the configured value", got.CookieName.ValueString())
	}

	if got.CookiePath.ValueString() != "/from-api" {
		t.Errorf("cookie_path = %q, want the API value %q",
			got.CookiePath.ValueString(), "/from-api")
	}

	if !got.MaxCookieAge.IsNull() || !got.MaxIdleTime.IsNull() {
		t.Error("attributes absent from the API response should resolve to null")
	}
}

// TestMergeClientSSLConfig covers the Set-typed supported_ssl_ciphers and
// supported_ssl_protocols attributes.
func TestMergeClientSSLConfig(t *testing.T) {
	t.Parallel()

	apiCfg := &sdk.ClientSSLLoadBalancerProfileConfig3{
		SupportedSslProtocols: []string{"TLS_V1_2"},
	}

	prior := ConfigClientSslValue{
		SslSuite:              types.StringValue("CUSTOM"),
		SupportedSslCiphers:   types.SetUnknown(types.StringType),
		SupportedSslProtocols: types.SetUnknown(types.StringType),
		state:                 attr.ValueStateKnown,
	}

	got := mergeClientSSLConfig(prior, apiCfg)

	if got.SupportedSslCiphers.IsUnknown() || got.SupportedSslProtocols.IsUnknown() {
		t.Fatal("unknown value leaked into state")
	}

	if got.SslSuite.ValueString() != "CUSTOM" {
		t.Errorf("ssl_suite = %q, want the configured value", got.SslSuite.ValueString())
	}

	if got.SupportedSslProtocols.IsNull() {
		t.Error("supported_ssl_protocols should have been resolved from the API response")
	}

	if !got.SupportedSslCiphers.IsNull() {
		t.Errorf("supported_ssl_ciphers = %v, want null", got.SupportedSslCiphers)
	}
}

// TestMergeConfigBlocksLeavesNoUnknowns is the guard that matters most: after
// merging, no config block may still carry an unknown attribute, whatever the
// API returned.
func TestMergeConfigBlocksLeavesNoUnknowns(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	prior := LoadBalancerProfileModel{
		ConfigHttp: ConfigHttpValue{
			RedirectAddress: types.StringUnknown(),
			XForwardedFor:   types.StringUnknown(),
			RequestBodySize: types.Int64Unknown(),
			state:           attr.ValueStateKnown,
		},
		ConfigFastTcp:             NewConfigFastTcpValueNull(),
		ConfigFastUdp:             NewConfigFastUdpValueNull(),
		ConfigCookiePersistence:   NewConfigCookiePersistenceValueNull(),
		ConfigSourceIpPersistence: NewConfigSourceIpPersistenceValueNull(),
		ConfigGenericPersistence:  NewConfigGenericPersistenceValueNull(),
		ConfigClientSsl:           NewConfigClientSslValueNull(),
		ConfigServerSsl:           NewConfigServerSslValueNull(),
	}

	var state LoadBalancerProfileModel

	mergeConfigBlocks(ctx, &state, prior, nil)

	obj, diags := state.ConfigHttp.ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("ToObjectValue: %v", diags)
	}

	for name, v := range obj.Attributes() {
		if v.IsUnknown() {
			t.Errorf("config_http.%s is still unknown after merge", name)
		}
	}

	// Blocks the practitioner did not set must remain null.
	if !state.ConfigClientSsl.IsNull() || !state.ConfigCookiePersistence.IsNull() {
		t.Error("unset config blocks should remain null")
	}
}
