// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package image

import (
	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attribute ssh_password_wo: see
// wodigest. A change marks write_only_revision unknown, which plans an in-place
// update in which Update re-sends the configured value. The deprecated
// *_version marker keeps working for the whole 2.2 line.
//
// APIHash additionally detects a change made outside Terraform. Morpheus
// returns a deterministic hash of the stored secret; we record it right after a
// send and compare it only for equality across later reads -- the hash
// algorithm is never assumed. OnDrift re-sends the configured value when the
// appliance's stored secret has changed.
var wo = wodigest.Tracker{
	Resource: "image",
	Attributes: map[string]string{
		"ssh_password_wo": "ssh_password_wo_version",
	},
	APIHash: map[string]wodigest.HashExtractor{
		"ssh_password_wo": sshPasswordHash,
	},
	OnDrift: wodigest.ResendPolicy,
}

// sshPasswordHash pulls the appliance's hash of the stored SSH password from
// the read-back object. present=false when the object is not the expected type
// or is nil, or when the API omitted the field (no information). A present-but-
// null field is ("", true): the appliance holds no secret, a comparable value.
func sshPasswordHash(resp any) (string, bool) {
	img, ok := resp.(*sdk.GetVirtualImage200ResponseVirtualImage)
	if !ok || img == nil {
		return "", false
	}

	return wodigest.NullableStringHash(img.SshPasswordHash)
}
