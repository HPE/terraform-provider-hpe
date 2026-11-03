// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package license

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ImportState imports a licence by its key id (the first eight characters of
// the key). A full key pasted by mistake is accepted too: licenseKeyID trims
// it to the key id.
//
// Only "id" is set; no digest or sentinel is stored. The Read that follows the
// import matches the licence by key id and records the digest from the API's
// hash, so the plan after import is empty and a later key change is detected
// immediately. There is deliberately no "#key=<sha256>" import form.
func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError("import license", "the import ID (the licence key id) must not be empty")

		return
	}

	keyID := licenseKeyID(id)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), keyID)...)
}
