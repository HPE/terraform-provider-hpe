// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"

// Change detection for the write-only ssh_password: see wodigest. The
// password is used only when the node is provisioned, so a change plans a
// replacement -- the RequiresReplace on the attribute itself can never fire,
// because a write-only value has no prior or planned state to compare.
// ModifyPlan is in modifyplan.go. The resource has no marker attribute.
var wo = wodigest.Tracker{
	Resource:   "instance_node",
	Attributes: map[string]string{"ssh_password": ""}, //nolint:gosec // an attribute name, not a credential
	Immutable:  []string{"ssh_password"},
}
