// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// typeSpecificAttributes lists the attributes that only some pool server
// types carry in their request body (see request.go), with the types that do.
var typeSpecificAttributes = []struct {
	name  string
	types []string
}{
	{"zone_filter", []string{typeCodeInfoblox}},
	{"tenant_match", []string{typeCodeInfoblox}},
	{"service_mode", []string{typeCodeInfoblox}},
	{"network_filter", []string{typeCodeInfoblox, typeCodeBluecat, typeCodePhpIPAM}},
}

// ValidateConfig checks the attributes against the pool server type. It only
// runs when type_code is known; a configuration that selects the type by
// type_id alone is validated by the appliance.
//
// app_id is new with this resource version, so misuse is an error. The other
// type-specific attributes were previously sent to every type and ignored by
// the appliance; they now warn, so an existing configuration keeps applying.
func (r *networkPoolServerResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var typeCode types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("type_code"), &typeCode)...)
	if resp.Diagnostics.HasError() || typeCode.IsNull() || typeCode.IsUnknown() {
		return
	}
	code := typeCode.ValueString()

	var appID types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("app_id"), &appID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	switch {
	case code == typeCodePhpIPAM && appID.IsNull():
		resp.Diagnostics.AddAttributeError(path.Root("app_id"), "Missing app_id",
			"app_id (the phpIPAM application ID) is required when type_code is \"phpipam\".")
	case code != typeCodePhpIPAM && !appID.IsNull():
		resp.Diagnostics.AddAttributeError(path.Root("app_id"), "app_id applies to phpIPAM only",
			fmt.Sprintf("app_id is the phpIPAM application ID and is not accepted for type_code %q.", code))
	}

	// Types without a typed request variant receive the full superset.
	if !hasTypedRequest(code) {
		return
	}

	for _, attr := range typeSpecificAttributes {
		var v types.String
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root(attr.name), &v)...)
		if v.IsNull() || v.IsUnknown() || contains(attr.types, code) {
			continue
		}

		resp.Diagnostics.AddAttributeWarning(path.Root(attr.name),
			fmt.Sprintf("%s is not used by %s pool servers", attr.name, code),
			fmt.Sprintf("%s applies to %v pool servers and is not sent for type_code %q.", attr.name, attr.types, code))
	}
}

func hasTypedRequest(code string) bool {
	return contains([]string{typeCodeInfoblox, typeCodeBluecat, typeCodePhpIPAM, typeCodeSolarWinds}, code)
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}

	return false
}
