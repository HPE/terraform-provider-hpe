// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package networkserverlookup finds a Morpheus network server by id or by
// name, hidden servers included.
//
// A cluster's own network server — an HVM cluster's `<cluster> HVM Network
// Server`, say — is hidden: the listing omits it unless asked with
// showHidden=true, and the single-server route answers 404 for it whatever
// is asked. Every network on the cluster must name that server, so both
// lookups fall back to the listing, whose entries carry the same fields as
// the single-server response.
package networkserverlookup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// Sentinel errors, kept as the messages the data source has always reported.
var (
	ErrNotFound = errors.New("no network server found")
	ErrMultiple = errors.New("multiple network servers were returned")
)

// pageSize is how many servers a listing page asks for. Paging continues
// until the server is found or the listing is exhausted, so a hidden server
// is never reported absent because it sat beyond the first page.
const pageSize int64 = 250

// maxPages bounds the paging, in case an appliance ignores offset or never
// reports a total: 40 pages is 10,000 network servers, far beyond any
// appliance, and a read that gets there reports a truncated listing rather
// than looping.
const maxPages = 40

// ErrTruncated reports that the listing was still returning full pages when
// the paging bound was reached, so absence could not be concluded.
var ErrTruncated = errors.New("network server listing truncated after 10000 entries; the server may exist beyond them")

// ByID returns the network server with the id. The single-server route is
// tried first; when it answers 404 the listing is searched, since the route
// refuses hidden servers regardless of showHidden.
func ByID(ctx context.Context, client *sdk.APIClient, id int64) (*sdk.GetNetworkServer200ResponseNetworkServer, error) {
	r, hresp, err := client.NetworksAPI.GetNetworkServer(ctx, id).Execute()
	if hresp != nil && hresp.StatusCode == http.StatusNotFound {
		return hiddenByID(ctx, client, id)
	}

	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for network server %d: %s", id, errfmt.ErrMsg(err, hresp))
	}

	if r.NetworkServer == nil {
		return nil, fmt.Errorf("GET failed for network server %d: response missing networkServer", id)
	}

	ns := *r.NetworkServer

	return &ns, nil
}

// hiddenByID pages through the listing, hidden servers included, for the id.
func hiddenByID(ctx context.Context, client *sdk.APIClient, id int64) (*sdk.GetNetworkServer200ResponseNetworkServer, error) {
	var found *sdk.ListNetworkServers200ResponseAllOfNetworkServersInner

	err := eachPage(ctx, client, "", func(s *sdk.ListNetworkServers200ResponseAllOfNetworkServersInner) bool {
		if s.Id != nil && *s.Id == id {
			found = s

			return false
		}

		return true
	})
	if err != nil {
		return nil, fmt.Errorf("GET failed for network server %d: %w", id, err)
	}

	if found == nil {
		return nil, fmt.Errorf("GET failed for network server %d: 404 (Not Found)", id)
	}

	return fromListItem(found)
}

// ByName returns the network server with exactly the name, hidden servers
// included: an exact name is unambiguous. The listing's own record is the
// result, rather than a second GET by id, which a hidden server would fail.
func ByName(ctx context.Context, client *sdk.APIClient, name string) (*sdk.GetNetworkServer200ResponseNetworkServer, error) {
	var (
		matched *sdk.ListNetworkServers200ResponseAllOfNetworkServersInner
		count   int
	)

	err := eachPage(ctx, client, name, func(s *sdk.ListNetworkServers200ResponseAllOfNetworkServersInner) bool {
		if s.Name != nil && *s.Name == name {
			matched = s
			count++
		}

		return true
	})
	if err != nil {
		return nil, fmt.Errorf("GET failed for network server %s: %w", name, err)
	}

	switch count {
	case 0:
		return nil, ErrNotFound
	case 1:
		return fromListItem(matched)
	default:
		return nil, ErrMultiple
	}
}

// eachPage lists network servers, hidden ones included and filtered by name
// when one is given, calling visit for each until it returns false or the
// listing is exhausted. The listing is considered exhausted when a page
// comes back short, or when the response's total has been reached; after
// maxPages full pages the read is reported truncated instead.
func eachPage(ctx context.Context, client *sdk.APIClient, name string,
	visit func(*sdk.ListNetworkServers200ResponseAllOfNetworkServersInner) bool,
) error {
	for page := 0; ; page++ {
		if page >= maxPages {
			return ErrTruncated
		}

		offset := int64(page) * pageSize

		req := client.NetworksAPI.ListNetworkServers(ctx).ShowHidden(true).Max(pageSize).Offset(offset)
		if name != "" {
			req = req.Name(name)
		}

		rs, hresp, err := req.Execute()
		if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
			return errors.New(errfmt.ErrMsg(err, hresp))
		}

		for i := range rs.NetworkServers {
			if !visit(&rs.NetworkServers[i]) {
				return nil
			}
		}

		got := int64(len(rs.NetworkServers))
		if got < pageSize {
			return nil
		}

		if rs.Meta != nil && rs.Meta.Total != nil && offset+got >= *rs.Meta.Total {
			return nil
		}
	}
}

// fromListItem converts a listing entry to the single-server model; the two
// carry the same fields.
func fromListItem(
	s *sdk.ListNetworkServers200ResponseAllOfNetworkServersInner,
) (*sdk.GetNetworkServer200ResponseNetworkServer, error) {
	raw, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("error marshaling network server: %w", err)
	}

	var ns sdk.GetNetworkServer200ResponseNetworkServer
	if err := json.Unmarshal(raw, &ns); err != nil {
		return nil, fmt.Errorf("error decoding network server: %w", err)
	}

	return &ns, nil
}
