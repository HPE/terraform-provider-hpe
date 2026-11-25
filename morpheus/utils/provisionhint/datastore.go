// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package provisionhint

// The second case is a datastore create that the storage plugin rejected
// with "Required fields are missing" against a nested id: errors.zonePool.id
// (the provider's resource_pool.id) or errors.storageServer.id
// (storage_server.id). Morpheus binds {"id": N} by loading the record with
// that id; when nothing is found it simply leaves the association empty, and
// the plugin's validation then reports the field as missing — which reads as
// "you did not set resource_pool" to someone who did. The ids that would have
// been accepted are one GET away: the cloud's resource pools, or the
// appliance's storage servers.
//
// The message claims only what those GETs verified. The error says no pool
// (or storage server) was attached; it does not say why, and the list fetched
// here is bounded and, for pools, scoped to one cloud — so an id's absence
// from the list is reported as exactly that, and only when the list is known
// to be complete.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

const (
	// zonePoolKey and storageServerKey are the errors keys the storage plugin
	// uses for the two nested ids; they are the API's names, not the
	// provider's.
	zonePoolKey      = "zonePool.id"
	storageServerKey = "storageServer.id"

	// requiredMissingText is the plugin's message for an id that resolved to
	// nothing. Detection is by key; the text is only quoted back.
	requiredMissingText = "Required fields are missing"

	// maxFetched bounds the list requests. A list longer than this comes back
	// as a sample, which the message says rather than pretending otherwise.
	maxFetched = 250

	// poolByName is the data source that resolves a resource pool by name,
	// named in the message beside the GET that lists them.
	poolByName = "data source hpe_morpheus_option_source_value (source \"zonePools\") resolves one by name"

	// storageServerByName is the same for storage servers.
	storageServerByName = "data source hpe_morpheus_storage_server resolves one by name"
)

// DatastoreRequest is what the caller sent, as far as the explanation needs
// it.
type DatastoreRequest struct {
	// CloudID is the cloud the datastore is associated with (refId when
	// refType is ComputeZone), or 0 when the request did not name one.
	CloudID int64
	// PoolID is resource_pool.id, sent as zonePool.id; 0 when not set.
	PoolID int64
	// StorageServerID is storage_server.id, sent as storageServer.id; 0 when
	// not set.
	StorageServerID int64
}

// DatastoreRejection returns an explanation of a datastore create that
// Morpheus rejected because resource_pool.id or storage_server.id did not
// resolve to a record, or "" when the failure is something else or no
// explanation can be built. body is the failed response's body; use ErrorBody
// to obtain it.
func DatastoreRejection(ctx context.Context, client *sdk.APIClient, body []byte, req DatastoreRequest) string {
	rejected := rejectedDatastoreFields(body)
	if client == nil || !rejected.any() {
		return ""
	}

	return datastoreRejection(ctx, sdkDatastoreLookups(client), rejected, req)
}

// rejectedFields says which of the two nested ids the response reported as
// missing.
type rejectedFields struct {
	pool          bool
	storageServer bool
}

func (r rejectedFields) any() bool { return r.pool || r.storageServer }

// rejectedDatastoreFields reads the errors object of a failed response. Only a
// string value under the exact key counts, so a nested or differently shaped
// error is not mistaken for this one.
func rejectedDatastoreFields(body []byte) rejectedFields {
	var env struct {
		Errors map[string]json.RawMessage `json:"errors"`
	}

	if err := json.Unmarshal(body, &env); err != nil {
		return rejectedFields{}
	}

	isString := func(key string) bool {
		raw, ok := env.Errors[key]
		if !ok {
			return false
		}

		var msg string

		return json.Unmarshal(raw, &msg) == nil
	}

	return rejectedFields{pool: isString(zonePoolKey), storageServer: isString(storageServerKey)}
}

// datastoreLookups are the appliance calls the explanation may make. Each is
// an interface so the message can be tested without an appliance.
type datastoreLookups struct {
	clouds         cloudNamer
	pools          poolLister
	storageServers storageServerLister
}

func sdkDatastoreLookups(client *sdk.APIClient) datastoreLookups {
	return datastoreLookups{
		clouds:         sdkCloudNamer{client: client},
		pools:          sdkPoolLister{client: client},
		storageServers: sdkStorageServerLister{client: client},
	}
}

// cloudNamer reads a cloud's name, or "" when it cannot.
type cloudNamer interface {
	name(ctx context.Context, id int64) string
}

// poolLister lists a cloud's resource pools.
type poolLister interface {
	pools(ctx context.Context, cloudID int64) (listing[poolRef], error)
}

// storageServerLister lists the appliance's storage servers.
type storageServerLister interface {
	storageServers(ctx context.Context) (listing[storageServerRef], error)
}

type storageServerRef struct {
	ID   int64
	Name string
	// Type is the storage server type code, e.g. "alletra-mp".
	Type string
}

// listing is what a lister returns: the items it fetched, and whether they
// are all there are. The requests are bounded by maxFetched, so a list can be
// a sample of a longer one; an id missing from a sample proves nothing, and
// the message must not treat it as missing from the appliance.
type listing[T any] struct {
	items    []T
	complete bool
}

// isComplete reports whether a page of fetched items is the whole list: the
// server's total says so, or, when it reports none, the page came back
// shorter than what was asked for.
func isComplete(fetched int, total *int64) bool {
	if total != nil {
		return *total <= int64(fetched)
	}

	return fetched < maxFetched
}

// datastoreRejection builds one paragraph per rejected id and joins them. A
// paragraph is dropped when the id turns out to be in the list after all —
// then the rejection is not the one this explains, and saying so would
// mislead.
func datastoreRejection(ctx context.Context, l datastoreLookups, rejected rejectedFields, req DatastoreRequest) string {
	var paragraphs []string

	if rejected.pool {
		if p := poolParagraph(ctx, l, req); p != "" {
			paragraphs = append(paragraphs, p)
		}
	}

	if rejected.storageServer {
		if p := storageServerParagraph(ctx, l, req); p != "" {
			paragraphs = append(paragraphs, p)
		}
	}

	return strings.Join(paragraphs, "\n\n")
}

// poolParagraph explains the zonePool.id rejection: Morpheus attached no
// resource pool for the id sent (or none was sent), and these are the pools of
// the datastore's cloud.
func poolParagraph(ctx context.Context, l datastoreLookups, req DatastoreRequest) string {
	var (
		list   listing[poolRef]
		listed bool
	)

	if l.pools != nil && req.CloudID != 0 {
		if got, err := l.pools.pools(ctx, req.CloudID); err == nil {
			list, listed = got, true
		}
	}

	if req.PoolID != 0 && inPools(list.items, req.PoolID) {
		return ""
	}

	cloud := cloudLabel(ctx, l.clouds, req.CloudID)
	listCall := fmt.Sprintf("GET /api/clouds/%d/resource-pools", req.CloudID)
	names := func(n int) string {
		return refList(n, func(i int) string { return idLabel(list.items[i].ID, list.items[i].Name) })
	}

	var b strings.Builder

	// What happened. The error key says Morpheus attached no pool; it does not
	// say why, so nothing more is claimed here.
	if req.PoolID == 0 {
		fmt.Fprintf(&b, "Morpheus reported %q: this datastore type requires resource_pool, which is not set.",
			zonePoolKey+": "+requiredMissingText)
	} else {
		fmt.Fprintf(&b, "Morpheus could not resolve resource_pool.id %d (reported as %q), so no resource pool "+
			"was attached to the datastore.", req.PoolID, zonePoolKey+": "+requiredMissingText)
	}

	// Where a valid id comes from — and, only when the whole list is in hand,
	// the fact that this id is not in it.
	switch {
	case req.CloudID == 0:
		b.WriteString(" Set resource_pool.id to a resource pool of the datastore's cloud; " +
			"GET /api/clouds/{id}/resource-pools lists them, and " + poolByName + ".")
	case !listed:
		fmt.Fprintf(&b, " Set resource_pool.id to a resource pool of %s; %s lists them, and %s.", cloud, listCall, poolByName)
	case len(list.items) == 0 && list.complete:
		fmt.Fprintf(&b, " %s has no resource pools yet; an HPE Bare Metal cloud creates its \"Default\" pool when "+
			"it is initialised, so refresh the cloud (or check associated_resource_id) and retry. %s lists them.",
			upperFirst(cloud), listCall)
	case list.complete:
		if req.PoolID != 0 {
			fmt.Fprintf(&b, " It is not one of the resource pools of %s: %s.", cloud, names(len(list.items)))
		} else {
			fmt.Fprintf(&b, " The resource pools of %s are: %s.", cloud, names(len(list.items)))
		}

		fmt.Fprintf(&b, " Set resource_pool.id to one of them; %s lists them, and %s.", listCall, poolByName)
	default:
		// A sample of a longer list: the id may be further down it.
		fmt.Fprintf(&b, " %s has more resource pools than were fetched here, including %s; set resource_pool.id "+
			"to one of them. %s lists them all, and %s.",
			upperFirst(cloud), names(min(len(list.items), maxListed)), listCall, poolByName)
	}

	return b.String()
}

// storageServerParagraph explains the storageServer.id rejection in the same
// way, against the storage servers the appliance lists.
func storageServerParagraph(ctx context.Context, l datastoreLookups, req DatastoreRequest) string {
	var (
		list   listing[storageServerRef]
		listed bool
	)

	if l.storageServers != nil {
		if got, err := l.storageServers.storageServers(ctx); err == nil {
			list, listed = got, true
		}
	}

	if req.StorageServerID != 0 {
		for _, s := range list.items {
			if s.ID == req.StorageServerID {
				return ""
			}
		}
	}

	names := func(n int) string {
		return refList(n, func(i int) string { return storageServerLabel(list.items[i]) })
	}

	var b strings.Builder

	if req.StorageServerID == 0 {
		fmt.Fprintf(&b, "Morpheus reported %q: this datastore type requires storage_server, which is not set.",
			storageServerKey+": "+requiredMissingText)
	} else {
		fmt.Fprintf(&b, "Morpheus could not resolve storage_server.id %d (reported as %q), so no storage server "+
			"was attached to the datastore.", req.StorageServerID, storageServerKey+": "+requiredMissingText)
	}

	switch {
	case !listed:
		b.WriteString(" GET /api/storage-servers lists the storage servers, and " + storageServerByName + ".")
	case len(list.items) == 0 && list.complete:
		b.WriteString(" The appliance has no storage servers; add one (hpe_morpheus_storage_server) first.")
	case list.complete:
		if req.StorageServerID != 0 {
			fmt.Fprintf(&b, " It is not one of the storage servers the appliance lists: %s.", names(len(list.items)))
		} else {
			fmt.Fprintf(&b, " The storage servers are: %s.", names(len(list.items)))
		}

		fmt.Fprintf(&b, " Set storage_server.id to one of them, or %s.", storageServerByName)
	default:
		fmt.Fprintf(&b, " The appliance has more storage servers than were fetched here, including %s; "+
			"set storage_server.id to one of them. GET /api/storage-servers lists them all, and %s.",
			names(min(len(list.items), maxListed)), storageServerByName)
	}

	return b.String()
}

// cloudLabel renders the cloud as "cloud 1 \"Metal\"", or "cloud 1" when the
// name is unknown.
func cloudLabel(ctx context.Context, clouds cloudNamer, id int64) string {
	if clouds == nil || id == 0 {
		return fmt.Sprintf("cloud %d", id)
	}

	return "cloud " + idLabel(id, clouds.name(ctx, id))
}

func idLabel(id int64, name string) string {
	if name == "" {
		return fmt.Sprintf("%d", id)
	}

	return fmt.Sprintf("%d %q", id, name)
}

func storageServerLabel(s storageServerRef) string {
	label := idLabel(s.ID, s.Name)
	if s.Type != "" {
		label += " (" + s.Type + ")"
	}

	return label
}

// refList renders up to maxListed items and counts the rest.
func refList(n int, item func(i int) string) string {
	parts := make([]string, 0, n+1)

	for i := range n {
		if i == maxListed {
			parts = append(parts, fmt.Sprintf("and %d more", n-maxListed))

			break
		}

		parts = append(parts, item(i))
	}

	return strings.Join(parts, ", ")
}

func upperFirst(s string) string {
	if s == "" {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

// sdkCloudNamer reads a cloud's name through the generated SDK.
type sdkCloudNamer struct {
	client *sdk.APIClient
}

func (g sdkCloudNamer) name(ctx context.Context, id int64) string {
	resp, _, err := g.client.CloudsAPI.GetClouds(ctx, id).Execute()
	if err != nil || resp == nil || resp.Zone == nil || resp.Zone.Name == nil {
		return ""
	}

	return *resp.Zone.Name
}

// sdkPoolLister lists a cloud's resource pools through the generated SDK. The
// endpoint reports meta.total, which is what says whether one page was the
// whole list.
type sdkPoolLister struct {
	client *sdk.APIClient
}

func (l sdkPoolLister) pools(ctx context.Context, cloudID int64) (listing[poolRef], error) {
	resp, _, err := l.client.CloudsAPI.ListCloudResourcePools(ctx, cloudID).Max(maxFetched).Execute()
	if err != nil {
		return listing[poolRef]{}, err
	}

	if resp == nil {
		return listing[poolRef]{}, errors.New("empty resource pool list response")
	}

	out := listing[poolRef]{
		items:    make([]poolRef, 0, len(resp.ResourcePools)),
		complete: isComplete(len(resp.ResourcePools), metaTotal(resp.Meta)),
	}

	for _, p := range resp.ResourcePools {
		if p.Id == nil {
			continue
		}

		ref := poolRef{ID: *p.Id}
		if p.Name != nil {
			ref.Name = *p.Name
		}

		out.items = append(out.items, ref)
	}

	return out, nil
}

// sdkStorageServerLister lists the appliance's storage servers through the
// generated SDK.
type sdkStorageServerLister struct {
	client *sdk.APIClient
}

func (l sdkStorageServerLister) storageServers(ctx context.Context) (listing[storageServerRef], error) {
	resp, _, err := l.client.StorageAPI.ListStorageServers(ctx).Max(maxFetched).Execute()
	if err != nil {
		return listing[storageServerRef]{}, err
	}

	if resp == nil {
		return listing[storageServerRef]{}, errors.New("empty storage server list response")
	}

	out := listing[storageServerRef]{
		items:    make([]storageServerRef, 0, len(resp.StorageServers)),
		complete: isComplete(len(resp.StorageServers), metaTotal(resp.Meta)),
	}

	for _, s := range resp.StorageServers {
		if s.Id == nil {
			continue
		}

		ref := storageServerRef{ID: *s.Id}
		if s.Name != nil {
			ref.Name = *s.Name
		}

		if t := s.Type.Get(); t != nil && t.Code != nil {
			ref.Type = *t.Code
		}

		out.items = append(out.items, ref)
	}

	return out, nil
}

// metaTotal reads the server's total from a list response's meta, or nil
// when the response carries none.
func metaTotal(meta *sdk.ListApprovals200ResponseAllOfMeta) *int64 {
	if meta == nil {
		return nil
	}

	return meta.Total
}
