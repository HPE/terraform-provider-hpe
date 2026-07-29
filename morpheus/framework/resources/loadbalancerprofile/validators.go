// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// blockForServiceType returns the config_* attribute name that must be set for
// a given serviceType, or "" if the serviceType is not recognised.
func blockForServiceType(serviceType string) string {
	switch serviceType {
	case serviceTypeHTTP:
		return "config_http"
	case serviceTypeFastTCP:
		return "config_fast_tcp"
	case serviceTypeFastUDP:
		return "config_fast_udp"
	case serviceTypeCookiePersistence:
		return "config_cookie_persistence"
	case serviceTypeSourceIPPersistence:
		return "config_source_ip_persistence"
	case serviceTypeGenericPersistence:
		return "config_generic_persistence"
	case serviceTypeClientSSL:
		return "config_client_ssl"
	case serviceTypeServerSSL:
		return "config_server_ssl"
	default:
		return ""
	}
}

// ConfigValidators wires up cross-attribute validation for the resource.
func (r *Resource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configBlockValidator{},
	}
}

// configBlockValidator ensures the config_* block set on a profile matches its
// service_type. The schema's ConflictsWith rules already prevent two blocks
// being set at once; this validator additionally requires the block matching
// the service_type to be present and rejects a block belonging to a different
// service_type (which would otherwise be silently ignored at create time).
type configBlockValidator struct{}

var _ resource.ConfigValidator = configBlockValidator{}

func (v configBlockValidator) Description(_ context.Context) string {
	return "the config_* block must match service_type"
}

func (v configBlockValidator) MarkdownDescription(_ context.Context) string {
	return "the `config_*` block must match `service_type`"
}

func (v configBlockValidator) ValidateResource(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var data LoadBalancerProfileModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// service_type drives which block applies. When it is null or unknown
	// (e.g. derived from another resource), skip: the schema's required and
	// OneOf validators on service_type report those cases on their own.
	if data.ServiceType.IsNull() || data.ServiceType.IsUnknown() {
		return
	}

	serviceType := data.ServiceType.ValueString()

	expected := blockForServiceType(serviceType)
	if expected == "" {
		// Unrecognised serviceType is already reported by the OneOf validator.
		return
	}

	set := map[string]bool{
		"config_http":                  !data.ConfigHttp.IsNull(),
		"config_fast_tcp":              !data.ConfigFastTcp.IsNull(),
		"config_fast_udp":              !data.ConfigFastUdp.IsNull(),
		"config_cookie_persistence":    !data.ConfigCookiePersistence.IsNull(),
		"config_source_ip_persistence": !data.ConfigSourceIpPersistence.IsNull(),
		"config_generic_persistence":   !data.ConfigGenericPersistence.IsNull(),
		"config_client_ssl":            !data.ConfigClientSsl.IsNull(),
		"config_server_ssl":            !data.ConfigServerSsl.IsNull(),
	}

	if !set[expected] {
		resp.Diagnostics.AddAttributeError(
			path.Root(expected),
			"Missing config block for service_type",
			fmt.Sprintf(
				"%s must be set when service_type is %q.",
				expected, serviceType,
			),
		)
	}

	for name, isSet := range set {
		if name != expected && isSet {
			resp.Diagnostics.AddAttributeError(
				path.Root(name),
				"Config block does not match service_type",
				fmt.Sprintf(
					"%s cannot be set when service_type is %q; use %s instead.",
					name, serviceType, expected,
				),
			)
		}
	}
}
