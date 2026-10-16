// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// clusterHostAttr is the single tracked attribute whose appliance baseline is
// the per-host SSH password hash map.
const clusterHostAttr = "server.ssh_password_wo" //nolint:gosec // an attribute name, not a credential

// Change detection for the nested write-only server.ssh_password_wo: see
// wodigest. Morpheus sets a cluster's SSH password at creation and has no
// update for it, so a change to the configured value plans a replacement of
// the cluster -- the only way the new value can take effect -- rather than an
// update that would silently do nothing. A change to the deprecated
// server.ssh_password_wo_version marker plans the same replacement for the
// whole 2.2 line.
//
// APIHash adds appliance-side drift detection. The cluster GET does not carry
// the hosts' SSH hashes; each host's sshPasswordHash is read from its server
// (/api/servers/{id}) and the hashes are combined into one canonical per-host
// map keyed by server id. wodigest compares the whole map for string equality,
// so the extractor emits canonical JSON; a change to any host's hash is a
// change of the map. A host added to or removed from the cluster also changes
// the map, which would be reported as spurious drift; Read handles that by
// re-baselining (wo.Observe) instead of comparing (wo.CheckDrift) whenever the
// current server-id set differs from the id set in the recorded baseline (see
// read.go). Because server.ssh_password_wo is Immutable, wodigest forces the
// drift policy to a warning regardless of OnDrift -- re-sending would replace
// the cluster -- so a detected change is reported, not acted on. OnDrift is set
// to WarnPolicy explicitly for clarity; the Immutable rule enforces it anyway.
//
// The cluster resource is not importable, so the import-ID digests do not
// apply here (the import digest does not seed the appliance baseline).
var wo = wodigest.Tracker{
	Resource: "cluster",
	Attributes: map[string]string{ //nolint:gosec // attribute names, not credentials
		"server.ssh_password_wo": "server.ssh_password_wo_version",
	},
	Immutable: []string{"server.ssh_password_wo"},
	APIHash: map[string]wodigest.HashExtractor{ //nolint:gosec // attribute names, not credentials
		"server.ssh_password_wo": clusterHostHashExtractor,
	},
	OnDrift: wodigest.WarnPolicy,
}

// clusterHostDigest carries the per-host SSH password hash map of a cluster's
// hosts: a map keyed by server id (as a decimal string) to that host's
// sshPasswordHash. present is true only when every host of the cluster was read
// successfully; a partial read (one or more host GETs failed) yields
// present=false, because a map missing some hosts cannot be compared against a
// full baseline without reporting a spurious change.
type clusterHostDigest struct {
	hashes  map[string]string
	present bool
}

// mapJSON is the canonical JSON encoding of the per-host hash map. Go's
// encoding/json sorts object keys, so the same map always encodes to the same
// bytes; wodigest compares those bytes for equality. A not-present digest
// encodes to "" so it is never mistaken for a real baseline.
func (d *clusterHostDigest) mapJSON() string {
	if d == nil || !d.present {
		return ""
	}

	raw, err := json.Marshal(d.hashes)
	if err != nil {
		return ""
	}

	return string(raw)
}

// clusterHostHashExtractor returns the per-host SSH password hash map as
// canonical JSON. The read response the provider passes is a
// *clusterHostDigest; anything else, or a not-present digest, yields
// present=false.
func clusterHostHashExtractor(resp any) (string, bool) {
	d, ok := resp.(*clusterHostDigest)
	if !ok || d == nil || !d.present {
		return "", false
	}

	return d.mapJSON(), true
}

// clusterHostHashes reads each host of the cluster and combines their
// sshPasswordHash values into a clusterHostDigest keyed by server id. It
// re-reads the cluster for its host server ids, then GETs each host. present is
// true only when every host was read successfully: a partial read (any host GET
// failing, returning nil, or a non-2xx status) yields present=false, because a
// map missing a host would compare unequal to the full baseline and report
// spurious drift. Drift detection is best-effort: any failure yields a
// not-present result and never blocks the operation.
func clusterHostHashes(ctx context.Context, client *sdk.APIClient, clusterID int64) *clusterHostDigest {
	clusterResp, _, err := client.ClustersAPI.GetCluster(ctx, clusterID).Execute()
	if err != nil || clusterResp == nil || clusterResp.Cluster == nil {
		return &clusterHostDigest{}
	}

	hashes := make(map[string]string, len(clusterResp.Cluster.Servers))

	for _, srv := range clusterResp.Cluster.Servers {
		if srv.Id == nil {
			tflog.Warn(ctx, "cluster: host with no id, drift detection skipped",
				map[string]any{"cluster": clusterID})

			return &clusterHostDigest{}
		}

		id := *srv.Id
		idParam := sdk.GetHostIdParameter{Int64: &id}

		h, hresp, herr := client.HostsAPI.GetHost(ctx, idParam).Execute()
		if herr != nil || h == nil || h.Server == nil ||
			hresp == nil || hresp.StatusCode < http.StatusOK ||
			hresp.StatusCode >= http.StatusMultipleChoices {
			// A partial read carries no information: without this host's hash
			// the map cannot be compared to a full baseline. Report
			// not-present rather than a map with a hole in it.
			tflog.Warn(ctx, "cluster: could not read host for drift detection",
				map[string]any{"cluster": clusterID, "server": id})

			return &clusterHostDigest{}
		}

		hash, present := wodigest.NullableStringHash(&h.Server.SshPasswordHash)
		if !present {
			// A host whose GET omits the hash leaves the map incomparable for
			// the same reason as a failed GET: report not-present rather than
			// record "" and raise drift once the hash is rendered.
			tflog.Warn(ctx, "cluster: host reports no ssh password hash, drift detection skipped",
				map[string]any{"cluster": clusterID, "server": id})

			return &clusterHostDigest{}
		}
		hashes[strconv.FormatInt(id, 10)] = hash
	}

	return &clusterHostDigest{hashes: hashes, present: true}
}

// observedHostIDs returns the sorted server ids recorded in the drift baseline
// (the Observed hash map for server.ssh_password_wo), and whether a baseline
// exists at all. It reads the stored record through the exported wodigest.Load,
// so Read can tell a host-set change from a secret change, and can skip the
// per-host GETs entirely when no baseline has been recorded yet. A baseline
// that does not parse as a map yields unknownIDs=true; Read re-baselines it,
// since it can never be compared.
func observedHostIDs(
	ctx context.Context,
	private wodigest.PrivateState,
) (ids []int64, haveBaseline, unknownIDs bool) {
	stored, ok, diags := wodigest.Load(ctx, private, wo.Key())
	if diags.HasError() || !ok {
		return nil, false, false
	}

	baseline, ok := stored.Observed[clusterHostAttr]
	if !ok || baseline == "" {
		return nil, false, false
	}

	var m map[string]string
	if err := json.Unmarshal([]byte(baseline), &m); err != nil {
		return nil, true, true
	}

	for k := range m {
		if id, err := strconv.ParseInt(k, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}

	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	return ids, true, false
}

// digestHostIDs returns the sorted server ids present in a read digest.
func digestHostIDs(d *clusterHostDigest) []int64 {
	if d == nil {
		return nil
	}

	out := make([]int64, 0, len(d.hashes))

	for k := range d.hashes {
		if id, err := strconv.ParseInt(k, 10, 64); err == nil {
			out = append(out, id)
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })

	return out
}

// sameIDs reports whether two sorted int64 slices are equal.
func sameIDs(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// ModifyPlan plans a replacement when server.ssh_password_wo changed (or its
// deprecated marker did), and reports a warning when a host's stored SSH
// password was changed on the appliance (see wodigest drift detection).
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	wo.ModifyPlan(ctx, req, resp)
}
