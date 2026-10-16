// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"context"
	"net/http"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// sshPasswordAttr is the single tracked write-only attribute.
const sshPasswordAttr = "ssh_password" //nolint:gosec // an attribute name, not a credential

// Change detection for the write-only ssh_password: see wodigest. The
// password is used only when the node is provisioned, so a change plans a
// replacement -- the RequiresReplace on the attribute itself can never fire,
// because a write-only value has no prior or planned state to compare.
// ModifyPlan is in modifyplan.go. The resource has no marker attribute.
//
// APIHash adds appliance-side drift detection: Morpheus renders a
// deterministic sshPasswordHash on the node's server (/api/servers/{id}). The
// provider records the hash read back after it last sent the password and
// compares it on every later read. Because ssh_password is Immutable, wodigest
// forces the drift policy to a warning regardless of OnDrift -- a re-send would
// be a replacement, which must never happen silently -- so a detected change is
// reported, not acted on. OnDrift is set to ResendPolicy only for consistency
// with the other resources; the Immutable rule overrides it.
var wo = wodigest.Tracker{
	Resource:   "instance_node",
	Attributes: map[string]string{"ssh_password": ""}, //nolint:gosec // an attribute name, not a credential
	Immutable:  []string{"ssh_password"},
	APIHash: map[string]wodigest.HashExtractor{ //nolint:gosec // an attribute name, not a credential
		"ssh_password": serverSSHPasswordHashExtractor,
	},
	OnDrift: wodigest.ResendPolicy,
}

// serverSSHPasswordHashExtractor reads sshPasswordHash from a node's server.
// The read response the provider passes is a *sdk.GetHost200ResponseServer;
// anything else yields present=false. sshPasswordHash is a NullableString, so
// the absent-vs-null rule is delegated to wodigest.NullableStringHash: the API
// omitting the field is "no information" (present=false), while the field
// rendered null is a real "no secret" value (present=true, hash="").
func serverSSHPasswordHashExtractor(resp any) (string, bool) {
	s, ok := resp.(*sdk.GetHost200ResponseServer)
	if !ok || s == nil {
		return "", false
	}

	return wodigest.NullableStringHash(&s.SshPasswordHash)
}

// serverForDrift reads the node's server for the drift baseline. Drift
// detection is best-effort: a missing server id, a failed read, a non-2xx
// status, or an empty response yields nil, which StoreObserved/CheckDrift treat
// as "nothing to baseline" rather than an error, so no operation is blocked.
func serverForDrift(ctx context.Context, client *sdk.APIClient, serverID int64) *sdk.GetHost200ResponseServer {
	if serverID == 0 {
		return nil
	}

	idParam := sdk.GetHostIdParameter{Int64: &serverID}

	h, hresp, err := client.HostsAPI.GetHost(ctx, idParam).Execute()
	if err != nil || h == nil || hresp == nil ||
		hresp.StatusCode < http.StatusOK || hresp.StatusCode >= http.StatusMultipleChoices {
		return nil
	}

	return h.Server
}

// hasDriftBaseline reports whether a drift baseline (Observed hash) has been
// recorded for ssh_password. Read uses it to decide whether to Observe (seed a
// baseline) or CheckDrift (compare against one).
func hasDriftBaseline(ctx context.Context, private wodigest.PrivateState) bool {
	stored, ok, diags := wodigest.Load(ctx, private, wo.Key())
	if diags.HasError() || !ok {
		return false
	}

	_, present := stored.Observed[sshPasswordAttr]

	return present
}
