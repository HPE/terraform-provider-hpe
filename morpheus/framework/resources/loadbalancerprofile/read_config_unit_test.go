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

// allVariantsPopulated mirrors what the SDK actually hands the provider.
//
// The load balancer profile config schema has no discriminator, so
// UnmarshalMapstructure decodes the same JSON blob into every anyOf variant and
// keeps each one that is not entirely empty. A single shared field is enough to
// make a variant survive, so in practice several variants -- always including
// HTTP, because the generic NSX-T config blob carries ntlmAuthentication -- are
// non-nil regardless of the profile's actual service type.
func allVariantsPopulated() *sdk.GetLoadBalancerProfile200ResponseLoadBalancerProfileConfig {
	return &sdk.GetLoadBalancerProfile200ResponseLoadBalancerProfileConfig{
		HTTPLoadBalancerProfileConfig3: &sdk.HTTPLoadBalancerProfileConfig3{
			NtlmAuthentication: sdk.PtrBool(false),
			RedirectAddress:    sdk.PtrString("http-variant"),
			Tags:               []sdk.LoadBalancerProfileTag24{{Name: sdk.PtrString("from"), Value: sdk.PtrString("http")}},
		},
		FastTCPLoadBalancerProfileConfig3: &sdk.FastTCPLoadBalancerProfileConfig3{
			HaFlowMirroring: sdk.PtrBool(true),
			Tags:            []sdk.LoadBalancerProfileTag25{{Name: sdk.PtrString("from"), Value: sdk.PtrString("fast_tcp")}},
		},
		FastUDPLoadBalancerProfileConfig3: &sdk.FastUDPLoadBalancerProfileConfig3{
			HaFlowMirroring: sdk.PtrBool(true),
			Tags:            []sdk.LoadBalancerProfileTag26{{Name: sdk.PtrString("from"), Value: sdk.PtrString("fast_udp")}},
		},
		CookiePersistenceLoadBalancerProfileConfig3: &sdk.CookiePersistenceLoadBalancerProfileConfig3{
			CookieName: sdk.PtrString("cookie-variant"),
			Tags:       []sdk.LoadBalancerProfileTag27{{Name: sdk.PtrString("from"), Value: sdk.PtrString("cookie_persistence")}},
		},
		SourceIPPersistenceLoadBalancerProfileConfig3: &sdk.SourceIPPersistenceLoadBalancerProfileConfig3{
			PurgeEntries: sdk.PtrBool(true),
			Tags:         []sdk.LoadBalancerProfileTag28{{Name: sdk.PtrString("from"), Value: sdk.PtrString("source_ip_persistence")}},
		},
		GenericPersistenceLoadBalancerProfileConfig3: &sdk.GenericPersistenceLoadBalancerProfileConfig3{
			SharePersistence: sdk.PtrBool(true),
			Tags:             []sdk.LoadBalancerProfileTag29{{Name: sdk.PtrString("from"), Value: sdk.PtrString("generic_persistence")}},
		},
		ClientSSLLoadBalancerProfileConfig3: &sdk.ClientSSLLoadBalancerProfileConfig3{
			SslSuite: sdk.PtrString("client-ssl-variant"),
			Tags:     []sdk.LoadBalancerProfileTag30{{Name: sdk.PtrString("from"), Value: sdk.PtrString("client_ssl")}},
		},
		ServerSSLLoadBalancerProfileConfig3: &sdk.ServerSSLLoadBalancerProfileConfig3{
			SslSuite: sdk.PtrString("server-ssl-variant"),
			Tags:     []sdk.LoadBalancerProfileTag31{{Name: sdk.PtrString("from"), Value: sdk.PtrString("server_ssl")}},
		},
	}
}

// configBlockIsNull reports whether the named config_* block is null.
func configBlockIsNull(m LoadBalancerProfileModel, block string) bool {
	switch block {
	case "config_http":
		return m.ConfigHttp.IsNull()
	case "config_fast_tcp":
		return m.ConfigFastTcp.IsNull()
	case "config_fast_udp":
		return m.ConfigFastUdp.IsNull()
	case "config_cookie_persistence":
		return m.ConfigCookiePersistence.IsNull()
	case "config_source_ip_persistence":
		return m.ConfigSourceIpPersistence.IsNull()
	case "config_generic_persistence":
		return m.ConfigGenericPersistence.IsNull()
	case "config_client_ssl":
		return m.ConfigClientSsl.IsNull()
	case "config_server_ssl":
		return m.ConfigServerSsl.IsNull()
	default:
		return true
	}
}

var allConfigBlockNames = []string{
	"config_http",
	"config_fast_tcp",
	"config_fast_udp",
	"config_cookie_persistence",
	"config_source_ip_persistence",
	"config_generic_persistence",
	"config_client_ssl",
	"config_server_ssl",
}

// TestReconstructConfigBlockFromResponseSelectsByServiceType is the regression
// guard for the import bug: reconstruction dispatched on which SDK union
// pointer happened to be non-nil, and HTTP was checked first, so importing any
// non-HTTP profile also materialised config_http. ImportStateVerify then failed
// with extra attributes ("config_http.%": "9").
func TestReconstructConfigBlockFromResponseSelectsByServiceType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	for serviceType, want := range map[string]string{
		serviceTypeHTTP:                "config_http",
		serviceTypeFastTCP:             "config_fast_tcp",
		serviceTypeFastUDP:             "config_fast_udp",
		serviceTypeCookiePersistence:   "config_cookie_persistence",
		serviceTypeSourceIPPersistence: "config_source_ip_persistence",
		serviceTypeGenericPersistence:  "config_generic_persistence",
		serviceTypeClientSSL:           "config_client_ssl",
		serviceTypeServerSSL:           "config_server_ssl",
	} {
		t.Run(serviceType, func(t *testing.T) {
			t.Parallel()

			// A freshly imported model: every config block is null.
			var state LoadBalancerProfileModel
			if !allConfigBlocksNull(state) {
				t.Fatal("precondition: a zero-value model must have no config block set")
			}

			reconstructConfigBlockFromResponse(ctx, &state, serviceType, allVariantsPopulated())

			if configBlockIsNull(state, want) {
				t.Errorf("%s was not reconstructed for service_type %q", want, serviceType)
			}

			for _, name := range allConfigBlockNames {
				if name == want {
					continue
				}

				if !configBlockIsNull(state, name) {
					t.Errorf("%s was populated for service_type %q; only %s may be set",
						name, serviceType, want)
				}
			}
		})
	}
}

// TestReadTagsFromConfigSelectsByServiceType guards the same dispatch bug in
// the tag reader, which would otherwise return another variant's tag list.
func TestReadTagsFromConfigSelectsByServiceType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	for serviceType, want := range map[string]string{
		serviceTypeHTTP:                "http",
		serviceTypeFastTCP:             "fast_tcp",
		serviceTypeFastUDP:             "fast_udp",
		serviceTypeCookiePersistence:   "cookie_persistence",
		serviceTypeSourceIPPersistence: "source_ip_persistence",
		serviceTypeGenericPersistence:  "generic_persistence",
		serviceTypeClientSSL:           "client_ssl",
		serviceTypeServerSSL:           "server_ssl",
	} {
		t.Run(serviceType, func(t *testing.T) {
			t.Parallel()

			got := readTagsFromConfig(
				ctx, serviceType, allVariantsPopulated(), types.SetNull(TagsValue{}.Type(ctx)),
			)

			var tags []TagsValue
			if diags := got.ElementsAs(ctx, &tags, false); diags.HasError() {
				t.Fatalf("ElementsAs: %v", diags)
			}

			if len(tags) != 1 {
				t.Fatalf("got %d tags, want 1", len(tags))
			}

			if v := tags[0].Value.ValueString(); v != want {
				t.Errorf("read tags from the %q variant, want %q", v, want)
			}
		})
	}
}

// TestReconstructConfigBlockFromResponseUnknownServiceType ensures an
// unrecognised service_type is inert rather than defaulting to a variant.
func TestReconstructConfigBlockFromResponseUnknownServiceType(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var state LoadBalancerProfileModel

	reconstructConfigBlockFromResponse(ctx, &state, "LBSomethingElseProfile", allVariantsPopulated())

	if !allConfigBlocksNull(state) {
		t.Error("an unrecognised service_type must not populate any config block")
	}
}
