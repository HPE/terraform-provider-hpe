// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package provisionhint explains provisioning failures that Morpheus reports
// tersely. It runs only after a request has already failed, and it never
// replaces the API's own error: every function here returns "" when it cannot
// add anything, so a caller appends the result to the message it would have
// shown anyway.
//
// The first case is the "Invalid network" rejection on instance create and
// clone. Morpheus emits it when a requested network is not selectable for the
// instance's cloud, group and resource pool — most often because the network
// belongs to the resource pool of the HVM cluster it was discovered on and the
// request named a different pool — but the response says none of that. The
// facts needed to say it are one GET per network away.
package provisionhint

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// invalidNetworkText is the value Morpheus puts in errors.networkInterfaces
// when a network fails the selectability check.
const invalidNetworkText = "Invalid network"

// Request is what the caller sent, as far as the explanation needs it.
type Request struct {
	// CloudID is the target cloud (zoneId), or 0 when the request does not
	// carry one — a clone inherits its source's cloud.
	CloudID int64
	// Config is the request's config object. It is marshalled to JSON and the
	// resource pool read from it, so any of the generated oneOf wrappers, a
	// pointer to one, or nil is accepted.
	Config any
	// Interfaces are the request's top-level network interfaces. The instance
	// and clone requests use different generated types of identical JSON shape,
	// so this is read through its JSON: [{"network": {"id": "7"}}, ...].
	Interfaces any
}

// networkGetter is the one SDK call the explanation needs. It is an interface
// so the message can be tested without an appliance.
type networkGetter interface {
	get(ctx context.Context, id int64) networkFacts
}

// networkFacts is what GET /api/networks/{id} tells us about one network.
type networkFacts struct {
	ID      int64
	Name    string
	CloudID int64
	Active  *bool
	// Pools is zonePool first, then any assignedZonePools not already listed.
	Pools []poolRef
	// LookupFailed is set when the network could not be read; LookupStatus is
	// the HTTP status when one was received.
	LookupFailed bool
	LookupStatus int
}

type poolRef struct {
	ID   int64
	Name string
}

// InvalidNetwork returns an explanation of an "Invalid network" failure, or ""
// when the failure is something else or no explanation can be built. body is
// the failed response's body; use ErrorBody to obtain it.
func InvalidNetwork(ctx context.Context, client *sdk.APIClient, body []byte, req Request) string {
	if client == nil || !isInvalidNetwork(body) {
		return ""
	}

	return invalidNetwork(ctx, sdkNetworkGetter{client: client}, req)
}

// InvalidNetworkReason is InvalidNetwork for a failure reported as text rather
// than as a response body — an asynchronous job's failure message, for
// example. It returns "" unless the text carries the "Invalid network"
// rejection.
func InvalidNetworkReason(ctx context.Context, client *sdk.APIClient, reason string, req Request) string {
	if client == nil || !strings.Contains(reason, invalidNetworkText) {
		return ""
	}

	return invalidNetwork(ctx, sdkNetworkGetter{client: client}, req)
}

// ErrorBody returns the body of a failed SDK call. The generated client keeps
// the body on the error it returns; failing that, the response body is read
// and put back, so the response can still be handed to errfmt afterwards.
func ErrorBody(err error, resp *http.Response) []byte {
	var byPointer *sdk.GenericOpenAPIError
	if errors.As(err, &byPointer) && byPointer != nil && len(byPointer.Body()) > 0 {
		return byPointer.Body()
	}

	var byValue sdk.GenericOpenAPIError
	if errors.As(err, &byValue) && len(byValue.Body()) > 0 {
		return byValue.Body()
	}

	if resp == nil || resp.Body == nil {
		return nil
	}

	b, readErr := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(b))

	if readErr != nil {
		return nil
	}

	return b
}

func invalidNetwork(ctx context.Context, networks networkGetter, req Request) string {
	ids := networkIDs(req.Interfaces)
	if len(ids) == 0 {
		return ""
	}

	poolID := resourcePoolID(req.Config)

	facts := make([]networkFacts, 0, len(ids))
	for _, id := range ids {
		facts = append(facts, networks.get(ctx, id))
	}

	return message(req.CloudID, poolID, facts)
}

// isInvalidNetwork reports whether body is Morpheus's "Invalid network"
// rejection: {"errors":{"networkInterfaces":"Invalid network"}, ...}.
func isInvalidNetwork(body []byte) bool {
	var env struct {
		Errors map[string]json.RawMessage `json:"errors"`
	}

	if err := json.Unmarshal(body, &env); err != nil {
		return false
	}

	raw, ok := env.Errors["networkInterfaces"]
	if !ok {
		return false
	}

	var msg string
	if err := json.Unmarshal(raw, &msg); err != nil {
		return false
	}

	return strings.Contains(msg, invalidNetworkText)
}

// resourcePoolID reads config.resourcePoolId from the request's config object,
// accepting "pool-7", "7" or 7. It returns 0 when there is none.
func resourcePoolID(config any) int64 {
	if config == nil {
		return 0
	}

	b, err := json.Marshal(config)
	if err != nil {
		return 0
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(b, &m); err != nil {
		return 0
	}

	raw, ok := m["resourcePoolId"]
	if !ok {
		return 0
	}

	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		id, err := strconv.ParseInt(strings.TrimPrefix(s, "pool-"), 10, 64)
		if err != nil {
			return 0
		}

		return id
	}

	var n int64
	if err := json.Unmarshal(raw, &n); err == nil {
		return n
	}

	return 0
}

// networkIDs extracts the bare network ids from the request's interfaces.
// network.id may be "7" or "network-7"; "networkGroup-" and "subnet-" ids are
// not networks and are skipped, as are duplicates.
func networkIDs(ifaces any) []int64 {
	if ifaces == nil {
		return nil
	}

	b, err := json.Marshal(ifaces)
	if err != nil {
		return nil
	}

	var list []struct {
		Network struct {
			ID string `json:"id"`
		} `json:"network"`
	}

	if err := json.Unmarshal(b, &list); err != nil {
		return nil
	}

	var ids []int64

	seen := make(map[int64]bool)

	for _, iface := range list {
		id, err := strconv.ParseInt(strings.TrimPrefix(iface.Network.ID, "network-"), 10, 64)
		if err != nil || id <= 0 || seen[id] {
			continue
		}

		seen[id] = true
		ids = append(ids, id)
	}

	return ids
}

// message builds the explanation from the facts. It is pure so the text can be
// pinned by tests.
func message(cloudID, poolID int64, facts []networkFacts) string {
	var b strings.Builder

	b.WriteString("Morpheus rejected a network interface because its network is not selectable for this ")
	b.WriteString("instance's cloud, group and resource pool:\n")

	poolMismatch := false

	for _, f := range facts {
		b.WriteString("\n  - ")
		b.WriteString(networkLine(cloudID, poolID, f, &poolMismatch))
	}

	if poolMismatch {
		b.WriteString("\n\nNetworks discovered on an HVM cluster are bound to that cluster's resource pool. ")
		b.WriteString("Provision into the cluster's pool (hpe_morpheus_cluster: permissions.resource_pool.id)")
		fmt.Fprintf(&b, ", or choose a network that belongs to resource pool %d.", poolID)
	}

	return b.String()
}

func networkLine(cloudID, poolID int64, f networkFacts, poolMismatch *bool) string {
	label := fmt.Sprintf("network %d", f.ID)
	if f.Name != "" {
		label = fmt.Sprintf("network %d %q", f.ID, f.Name)
	}

	switch {
	case f.LookupFailed && f.LookupStatus == http.StatusNotFound:
		return label + " does not exist, or is not visible to this user."
	case f.LookupFailed && f.LookupStatus != 0:
		return fmt.Sprintf("%s could not be read (HTTP %d).", label, f.LookupStatus)
	case f.LookupFailed:
		return label + " could not be read."
	case cloudID != 0 && f.CloudID != 0 && f.CloudID != cloudID:
		return fmt.Sprintf("%s is in cloud %d; this instance targets cloud %d.", label, f.CloudID, cloudID)
	case f.Active != nil && !*f.Active:
		return label + " is inactive."
	case poolID != 0 && len(f.Pools) == 0:
		*poolMismatch = true

		return fmt.Sprintf("%s belongs to no resource pool; this instance targets resource pool %d.", label, poolID)
	case poolID != 0 && !inPools(f.Pools, poolID):
		*poolMismatch = true

		return fmt.Sprintf("%s belongs to %s; this instance targets resource pool %d.",
			label, describePools(f.Pools), poolID)
	default:
		return label + " is in this cloud and active; it is not visible to the instance's group or tenant " +
			"(check the network's group access and tenant permissions)."
	}
}

func inPools(pools []poolRef, id int64) bool {
	for _, p := range pools {
		if p.ID == id {
			return true
		}
	}

	return false
}

func describePools(pools []poolRef) string {
	parts := make([]string, 0, len(pools))

	for _, p := range pools {
		if p.Name != "" {
			parts = append(parts, fmt.Sprintf("%d %q", p.ID, p.Name))
		} else {
			parts = append(parts, strconv.FormatInt(p.ID, 10))
		}
	}

	if len(parts) == 1 {
		return "resource pool " + parts[0]
	}

	return "resource pools " + strings.Join(parts, ", ")
}

// sdkNetworkGetter reads a network through the generated SDK.
type sdkNetworkGetter struct {
	client *sdk.APIClient
}

func (g sdkNetworkGetter) get(ctx context.Context, id int64) networkFacts {
	resp, httpResp, err := g.client.NetworksAPI.GetNetwork(ctx, id).Execute()
	if err != nil || resp == nil || resp.Network == nil {
		f := networkFacts{ID: id, LookupFailed: true}
		if httpResp != nil {
			f.LookupStatus = httpResp.StatusCode
		}

		return f
	}

	return factsFromNetwork(id, resp.Network)
}

func factsFromNetwork(id int64, n *sdk.GetNetwork200ResponseNetwork) networkFacts {
	f := networkFacts{ID: id, Active: n.Active}

	if n.Name != nil {
		f.Name = *n.Name
	}

	if n.Zone != nil && n.Zone.Id != nil {
		f.CloudID = *n.Zone.Id
	}

	seen := make(map[int64]bool)

	if n.ZonePool != nil && n.ZonePool.Id != nil {
		p := poolRef{ID: *n.ZonePool.Id}
		if n.ZonePool.Name != nil {
			p.Name = *n.ZonePool.Name
		}

		f.Pools = append(f.Pools, p)
		seen[p.ID] = true
	}

	// assignedZonePools is not in the generated model; it arrives as
	// [{"id": 1, "name": "..."}] under the remaining properties.
	list, ok := n.AdditionalProperties["assignedZonePools"].([]any)
	if !ok {
		return f
	}

	var extra []poolRef

	for _, item := range list {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}

		idf, ok := m["id"].(float64)
		if !ok || seen[int64(idf)] {
			continue
		}

		p := poolRef{ID: int64(idf)}
		if name, ok := m["name"].(string); ok {
			p.Name = name
		}

		seen[p.ID] = true
		extra = append(extra, p)
	}

	sort.Slice(extra, func(i, j int) bool { return extra[i].ID < extra[j].ID })

	f.Pools = append(f.Pools, extra...)

	return f
}
