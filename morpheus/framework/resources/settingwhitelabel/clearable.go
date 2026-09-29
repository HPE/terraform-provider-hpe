// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingwhitelabel

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the appliance's white label settings. All are
// Optional-only, so the plan carries the null: the name and colours are sent
// as empty strings (the shape Delete uses to reset them); a removed logo or
// favicon is reset through the API's reset flag, as before.
var clearableAttributes = clearable.Attributes{
	{Name: "appliance_name", Mode: clearable.Clear},
	{Name: "primary_color", Mode: clearable.Clear},
	{Name: "secondary_color", Mode: clearable.Clear},
	{Name: "header_logo", Mode: clearable.Clear},
	{Name: "footer_logo", Mode: clearable.Clear},
	{Name: "login_logo", Mode: clearable.Clear},
	{Name: "favicon", Mode: clearable.Clear},
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration.
var notClearable = clearable.NotClearable{
	"support_menu_links": "not sent to the API by the provider; the current links are kept when it is removed",
	"enabled":            "boolean",
}
