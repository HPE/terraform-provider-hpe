// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storageserver

import (
	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attribute service_password_wo: see
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
	Resource: "storage_server",
	Attributes: map[string]string{
		"service_password_wo": "service_password_wo_version",
	},
	APIHash: map[string]wodigest.HashExtractor{
		"service_password_wo": servicePasswordHash,
	},
	OnDrift: wodigest.ResendPolicy,
}

// servicePasswordHash pulls the appliance's hash of the stored service password
// from the read-back object. present=false when the object is not the expected
// type or is nil, or when the API omitted the field (no information). A present-
// but-null field is ("", true): the appliance holds no secret, a comparable
// value.
func servicePasswordHash(resp any) (string, bool) {
	ss, ok := resp.(*sdk.GetStorageServers200ResponseStorageServer)
	if !ok || ss == nil {
		return "", false
	}

	return wodigest.NullableStringHash(ss.ServicePasswordHash)
}
