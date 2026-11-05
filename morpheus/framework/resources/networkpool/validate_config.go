// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpool

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithValidateConfig = &networkPoolResource{}

// ValidateConfig requires both addresses on every configured
// ip_address_ranges element. The element attributes are Optional+Computed in
// the schema (not Required) so that Terraform keeps the prior value of
// ip_address_ranges when the configuration uses the deprecated ip_ranges
// instead — Terraform only does that for an absent Optional+Computed nested
// attribute whose children are all computed. This check restores the
// "required" contract for what the practitioner actually writes.
func (r *networkPoolResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config NetworkPoolModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() || config.IpAddressRanges.IsNull() || config.IpAddressRanges.IsUnknown() {
		return
	}

	var elems []IpAddressRangesValue
	resp.Diagnostics.Append(config.IpAddressRanges.ElementsAs(ctx, &elems, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for i, e := range elems {
		if e.IsUnknown() {
			continue
		}
		for name, v := range map[string]interface{ IsNull() bool }{
			"starting_address": e.StartingAddress,
			"ending_address":   e.EndingAddress,
		} {
			if v.IsNull() {
				resp.Diagnostics.AddAttributeError(
					path.Root("ip_address_ranges").AtListIndex(i).AtName(name),
					"Missing range address",
					fmt.Sprintf("ip_address_ranges[%d].%s is required: every range needs both a starting and an ending address.", i, name),
				)
			}
		}
	}
}
