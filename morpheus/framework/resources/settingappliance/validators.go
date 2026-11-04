// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingappliance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// smtpFromAddressValidator rejects an empty smtp_from_address. Morpheus stores
// one, but reports info@gomorpheus.com whenever no address is set, so the apply
// would read back a value other than the one planned. Every other text setting
// emptied with "" reads back as unset, which the provider takes as the ""
// configured (see mapSettings). The rule lives here rather than in the codegen
// spec, which also generates the Ansible module.
type smtpFromAddressValidator struct{}

var _ resource.ConfigValidator = smtpFromAddressValidator{}

func (smtpFromAddressValidator) Description(_ context.Context) string {
	return "smtp_from_address must not be empty"
}

func (smtpFromAddressValidator) MarkdownDescription(_ context.Context) string {
	return "`smtp_from_address` must not be empty"
}

func (smtpFromAddressValidator) ValidateResource(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var address types.String

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("smtp_from_address"), &address)...)

	if !known(address) || address.ValueString() != "" {
		return
	}

	resp.Diagnostics.AddAttributeError(
		path.Root("smtp_from_address"),
		"Empty SMTP from address",
		"Morpheus reports info@gomorpheus.com whenever no SMTP from address is set, so an empty one cannot "+
			"be applied. Set an address, or remove the attribute to leave the appliance's as it is.",
	)
}
