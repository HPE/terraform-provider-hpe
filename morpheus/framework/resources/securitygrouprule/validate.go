// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygrouprule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithValidateConfig = &securityGroupRuleResource{}

// ValidateConfig fails the plan for a rule Morpheus would reject at create or
// update, with the attribute named, instead of letting the API's 400 surface
// at apply (MORPH-18279). The conditions mirror Morpheus's
// SecurityGroupService.validateRule:
//
//   - source is required when source_type is cidr (Morpheus's default when
//     the type is not given);
//   - destination is required when destination_type is cidr (likewise);
//   - port_range is required unless protocol is icmp.
//
// The three attributes stay Optional in the schema because each is required
// only under its condition: an ICMP rule has no port range, and a group- or
// tier-typed source has no CIDR. Unknown values -- a source taken from another
// resource, say -- are not judged at plan.
func (r *securityGroupRuleResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config SecurityGroupRuleModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateRequiredByCondition(config)...)
}

// validateRequiredByCondition is the condition table, separated from the
// framework plumbing so it can be unit-tested directly.
func validateRequiredByCondition(config SecurityGroupRuleModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if requiresCIDR(config.SourceType) && config.Source.IsNull() {
		diags.AddAttributeError(path.Root("source"),
			"Missing source",
			"source is required when source_type is \"cidr\" (the default): the rule has no address "+
				"to match. Set source to a CIDR such as \"10.0.0.0/24\", or choose another source_type.")
	}

	if requiresCIDR(config.DestinationType) && config.Destination.IsNull() {
		diags.AddAttributeError(path.Root("destination"),
			"Missing destination",
			"destination is required when destination_type is \"cidr\" (the default): the rule has no "+
				"address to match. Set destination to a CIDR such as \"10.0.0.0/24\", or choose another "+
				"destination_type.")
	}

	if requiresPortRange(config.Protocol) && config.PortRange.IsNull() {
		diags.AddAttributeError(path.Root("port_range"),
			"Missing port_range",
			"port_range is required for protocol \""+config.Protocol.ValueString()+"\": only ICMP rules "+
				"have no ports. Set port_range to a port or range such as \"443\" or \"1024-2048\".")
	}

	return diags
}

// requiresCIDR reports whether a source/destination type needs a CIDR: an
// unset type is Morpheus's default, cidr. An unknown type cannot be judged.
func requiresCIDR(typ types.String) bool {
	if typ.IsUnknown() {
		return false
	}

	return typ.IsNull() || typ.ValueString() == "cidr"
}

// requiresPortRange reports whether a protocol needs a port range: every
// known protocol but icmp. An unset or unknown protocol cannot be judged
// (protocol is Required, so unset only happens in a partial configuration
// the framework reports on its own).
func requiresPortRange(protocol types.String) bool {
	if protocol.IsNull() || protocol.IsUnknown() {
		return false
	}

	return protocol.ValueString() != "icmp"
}
