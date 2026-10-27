// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package provisionhint

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// The body the Alletra MP storage plugin returns when resource_pool.id names
// no resource pool — the case QA hit with the docs example's id = 1.
const unknownPoolBody = `{"errors":{"zonePool.id":"Required fields are missing"},"success":false}`

func TestUnitRejectedDatastoreFields(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		body string
		want rejectedFields
	}{
		"the real pool rejection": {unknownPoolBody, rejectedFields{pool: true}},
		"storage server rejection": {
			`{"errors":{"storageServer.id":"Required fields are missing"}}`,
			rejectedFields{storageServer: true},
		},
		"both": {
			`{"errors":{"storageServer.id":"Required fields are missing","zonePool.id":"Required fields are missing"}}`,
			rejectedFields{pool: true, storageServer: true},
		},
		"other field":          {`{"errors":{"protocolType":"Required fields are missing"}}`, rejectedFields{}},
		"nested, not dotted":   {`{"errors":{"zonePool":{"id":"Required fields are missing"}}}`, rejectedFields{}},
		"value not a string":   {`{"errors":{"zonePool.id":["Required fields are missing"]}}`, rejectedFields{}},
		"no errors object":     {`{"msg":"zonePool.id","success":false}`, rejectedFields{}},
		"errors not an object": {`{"errors":"zonePool.id"}`, rejectedFields{}},
		"not json":             {`502 Bad Gateway`, rejectedFields{}},
		"empty":                {``, rejectedFields{}},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := rejectedDatastoreFields([]byte(tc.body)); got != tc.want {
				t.Errorf("rejectedDatastoreFields(%q) = %+v, want %+v", tc.body, got, tc.want)
			}
		})
	}
}

func TestUnitIsComplete(t *testing.T) {
	t.Parallel()

	total := func(n int64) *int64 { return &n }

	cases := map[string]struct {
		fetched int
		total   *int64
		want    bool
	}{
		"total matches the page":            {2, total(2), true},
		"total is zero":                     {0, total(0), true},
		"total exceeds the page":            {250, total(300), false},
		"total exceeds a short page":        {25, total(30), false},
		"no total, page shorter than asked": {2, nil, true},
		"no total, page as long as asked":   {maxFetched, nil, false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := isComplete(tc.fetched, tc.total); got != tc.want {
				t.Errorf("isComplete(%d, %v) = %v, want %v", tc.fetched, tc.total, got, tc.want)
			}
		})
	}
}

// fakeClouds names clouds by id.
type fakeClouds map[int64]string

func (f fakeClouds) name(_ context.Context, id int64) string { return f[id] }

// fakePools answers pools() for one cloud, or with an error; any other cloud
// id is an error too, which is what a GET of an unknown cloud returns.
type fakePools struct {
	cloudID    int64
	list       []poolRef
	incomplete bool
	err        error
}

func (f fakePools) pools(_ context.Context, cloudID int64) (listing[poolRef], error) {
	if f.err != nil {
		return listing[poolRef]{}, f.err
	}

	if cloudID != f.cloudID {
		return listing[poolRef]{}, errors.New("404 Not Found")
	}

	return listing[poolRef]{items: f.list, complete: !f.incomplete}, nil
}

// fakeStorageServers answers storageServers() with a fixed list or an error.
type fakeStorageServers struct {
	list       []storageServerRef
	incomplete bool
	err        error
}

func (f fakeStorageServers) storageServers(_ context.Context) (listing[storageServerRef], error) {
	if f.err != nil {
		return listing[storageServerRef]{}, f.err
	}

	return listing[storageServerRef]{items: f.list, complete: !f.incomplete}, nil
}

// The pool paragraph's fixed parts.
const (
	poolRejectedOne = "Morpheus could not resolve resource_pool.id 1 (reported as \"zonePool.id: Required fields are " +
		"missing\"), so no resource pool was attached to the datastore."
	poolNotSet = "Morpheus reported \"zonePool.id: Required fields are missing\": this datastore type requires " +
		"resource_pool, which is not set."
	poolWhere = "GET /api/clouds/1/resource-pools lists them, and data source hpe_morpheus_option_source_value " +
		"(source \"zonePools\") resolves one by name."
)

func TestUnitDatastoreRejectionPool(t *testing.T) {
	t.Parallel()

	metal := fakeClouds{1: "Metal"}
	twoPools := fakePools{cloudID: 1, list: []poolRef{{ID: 7, Name: "Default"}, {ID: 8, Name: "rack-a"}}}
	rejectedPool := rejectedFields{pool: true}

	cases := map[string]struct {
		lookups datastoreLookups
		req     DatastoreRequest
		want    string
	}{
		"QA's case: the example's id = 1, the cloud's pools are listed in full": {
			lookups: datastoreLookups{clouds: metal, pools: twoPools},
			req:     DatastoreRequest{CloudID: 1, PoolID: 1},
			want: poolRejectedOne + " It is not one of the resource pools of cloud 1 \"Metal\": 7 \"Default\", " +
				"8 \"rack-a\". Set resource_pool.id to one of them; " + poolWhere,
		},
		"the list is a sample: no claim that the id is absent": {
			lookups: datastoreLookups{clouds: metal, pools: fakePools{cloudID: 1, list: twoPools.list, incomplete: true}},
			req:     DatastoreRequest{CloudID: 1, PoolID: 1},
			want: poolRejectedOne + " Cloud 1 \"Metal\" has more resource pools than were fetched here, including " +
				"7 \"Default\", 8 \"rack-a\"; set resource_pool.id to one of them. GET /api/clouds/1/resource-pools " +
				"lists them all, and data source hpe_morpheus_option_source_value (source \"zonePools\") resolves " +
				"one by name.",
		},
		"cloud has no pools yet": {
			lookups: datastoreLookups{clouds: metal, pools: fakePools{cloudID: 1}},
			req:     DatastoreRequest{CloudID: 1, PoolID: 1},
			want: poolRejectedOne + " Cloud 1 \"Metal\" has no resource pools yet; an HPE Bare Metal cloud creates " +
				"its \"Default\" pool when it is initialised, so refresh the cloud (or check associated_resource_id) " +
				"and retry. GET /api/clouds/1/resource-pools lists them.",
		},
		"pool list unavailable, cloud name unknown: no claim about the list": {
			lookups: datastoreLookups{clouds: fakeClouds{}, pools: fakePools{err: errors.New("403 Forbidden")}},
			req:     DatastoreRequest{CloudID: 1, PoolID: 1},
			want:    poolRejectedOne + " Set resource_pool.id to a resource pool of cloud 1; " + poolWhere,
		},
		"no lookups at all": {
			lookups: datastoreLookups{},
			req:     DatastoreRequest{CloudID: 1, PoolID: 1},
			want:    poolRejectedOne + " Set resource_pool.id to a resource pool of cloud 1; " + poolWhere,
		},
		"resource_pool not set at all": {
			lookups: datastoreLookups{clouds: metal, pools: twoPools},
			req:     DatastoreRequest{CloudID: 1},
			want: poolNotSet + " The resource pools of cloud 1 \"Metal\" are: 7 \"Default\", 8 \"rack-a\". " +
				"Set resource_pool.id to one of them; " + poolWhere,
		},
		"no cloud known either": {
			lookups: datastoreLookups{clouds: metal, pools: twoPools},
			req:     DatastoreRequest{PoolID: 1},
			want: poolRejectedOne + " Set resource_pool.id to a resource pool of the datastore's cloud; " +
				"GET /api/clouds/{id}/resource-pools lists them, and data source hpe_morpheus_option_source_value " +
				"(source \"zonePools\") resolves one by name.",
		},
		"the pool is in the list: not this rejection, say nothing": {
			lookups: datastoreLookups{clouds: metal, pools: twoPools},
			req:     DatastoreRequest{CloudID: 1, PoolID: 7},
			want:    "",
		},
		"the pool is in the sample: say nothing": {
			lookups: datastoreLookups{clouds: metal, pools: fakePools{cloudID: 1, list: twoPools.list, incomplete: true}},
			req:     DatastoreRequest{CloudID: 1, PoolID: 8},
			want:    "",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := datastoreRejection(context.Background(), tc.lookups, rejectedPool, tc.req)
			if got != tc.want {
				t.Errorf("message\n--- want ---\n%s\n--- got ---\n%s", tc.want, got)
			}
		})
	}
}

func TestUnitDatastoreRejectionStorageServer(t *testing.T) {
	t.Parallel()

	servers := fakeStorageServers{list: []storageServerRef{
		{ID: 3, Name: "alletra-01", Type: "alletra-mp"},
		{ID: 4, Name: "nimble"},
	}}
	rejectedServer := rejectedFields{storageServer: true}

	const rejectedOne = "Morpheus could not resolve storage_server.id 1 (reported as \"storageServer.id: Required " +
		"fields are missing\"), so no storage server was attached to the datastore."

	cases := map[string]struct {
		lookups datastoreLookups
		req     DatastoreRequest
		want    string
	}{
		"the example's id = 1, listed in full": {
			lookups: datastoreLookups{storageServers: servers},
			req:     DatastoreRequest{StorageServerID: 1},
			want: rejectedOne + " It is not one of the storage servers the appliance lists: 3 \"alletra-01\" " +
				"(alletra-mp), 4 \"nimble\". Set storage_server.id to one of them, or data source " +
				"hpe_morpheus_storage_server resolves one by name.",
		},
		"the list is a sample": {
			lookups: datastoreLookups{storageServers: fakeStorageServers{list: servers.list, incomplete: true}},
			req:     DatastoreRequest{StorageServerID: 1},
			want: rejectedOne + " The appliance has more storage servers than were fetched here, including " +
				"3 \"alletra-01\" (alletra-mp), 4 \"nimble\"; set storage_server.id to one of them. " +
				"GET /api/storage-servers lists them all, and data source hpe_morpheus_storage_server resolves " +
				"one by name.",
		},
		"list unavailable": {
			lookups: datastoreLookups{storageServers: fakeStorageServers{err: errors.New("403 Forbidden")}},
			req:     DatastoreRequest{StorageServerID: 1},
			want: rejectedOne + " GET /api/storage-servers lists the storage servers, and data source " +
				"hpe_morpheus_storage_server resolves one by name.",
		},
		"none on the appliance": {
			lookups: datastoreLookups{storageServers: fakeStorageServers{}},
			req:     DatastoreRequest{StorageServerID: 1},
			want:    rejectedOne + " The appliance has no storage servers; add one (hpe_morpheus_storage_server) first.",
		},
		"storage_server not set": {
			lookups: datastoreLookups{storageServers: servers},
			req:     DatastoreRequest{},
			want: "Morpheus reported \"storageServer.id: Required fields are missing\": this datastore type requires " +
				"storage_server, which is not set. The storage servers are: 3 \"alletra-01\" (alletra-mp), " +
				"4 \"nimble\". Set storage_server.id to one of them, or data source hpe_morpheus_storage_server " +
				"resolves one by name.",
		},
		"the server is in the list: say nothing": {
			lookups: datastoreLookups{storageServers: servers},
			req:     DatastoreRequest{StorageServerID: 3},
			want:    "",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := datastoreRejection(context.Background(), tc.lookups, rejectedServer, tc.req)
			if got != tc.want {
				t.Errorf("message\n--- want ---\n%s\n--- got ---\n%s", tc.want, got)
			}
		})
	}
}

func TestUnitDatastoreRejectionBothFields(t *testing.T) {
	t.Parallel()

	l := datastoreLookups{
		clouds:         fakeClouds{1: "Metal"},
		pools:          fakePools{cloudID: 1, list: []poolRef{{ID: 7, Name: "Default"}}},
		storageServers: fakeStorageServers{list: []storageServerRef{{ID: 3, Name: "alletra-01", Type: "alletra-mp"}}},
	}

	got := datastoreRejection(context.Background(), l, rejectedFields{pool: true, storageServer: true},
		DatastoreRequest{CloudID: 1, PoolID: 1, StorageServerID: 1})

	parts := strings.Split(got, "\n\n")
	if len(parts) != 2 {
		t.Fatalf("expected two paragraphs separated by a blank line, got %d:\n%s", len(parts), got)
	}

	if !strings.HasPrefix(parts[0], "Morpheus could not resolve resource_pool.id 1") {
		t.Errorf("first paragraph should explain the pool, got %q", parts[0])
	}

	if !strings.HasPrefix(parts[1], "Morpheus could not resolve storage_server.id 1") {
		t.Errorf("second paragraph should explain the storage server, got %q", parts[1])
	}
}

func TestUnitDatastoreRejectionListIsCapped(t *testing.T) {
	t.Parallel()

	pools := make([]poolRef, 0, maxListed+3)
	for i := range maxListed + 3 {
		pools = append(pools, poolRef{ID: int64(100 + i), Name: fmt.Sprintf("pool-%d", i)})
	}

	t.Run("complete list: the rest are counted", func(t *testing.T) {
		t.Parallel()

		l := datastoreLookups{clouds: fakeClouds{1: "Metal"}, pools: fakePools{cloudID: 1, list: pools}}

		got := datastoreRejection(context.Background(), l, rejectedFields{pool: true}, DatastoreRequest{CloudID: 1, PoolID: 1})

		if !strings.Contains(got, fmt.Sprintf("%d \"pool-%d\", and 3 more.", 100+maxListed-1, maxListed-1)) {
			t.Errorf("expected the list to stop at %d entries with a count of the rest, got:\n%s", maxListed, got)
		}

		if strings.Contains(got, fmt.Sprintf("\"pool-%d\"", maxListed)) {
			t.Errorf("entry %d should not be listed:\n%s", maxListed+1, got)
		}
	})

	t.Run("sample: the first few, no count", func(t *testing.T) {
		t.Parallel()

		l := datastoreLookups{clouds: fakeClouds{1: "Metal"}, pools: fakePools{cloudID: 1, list: pools, incomplete: true}}

		got := datastoreRejection(context.Background(), l, rejectedFields{pool: true}, DatastoreRequest{CloudID: 1, PoolID: 1})

		if strings.Contains(got, "more.") || strings.Contains(got, "It is not one of") {
			t.Errorf("a sample must neither count the rest nor claim the id is absent:\n%s", got)
		}

		if !strings.Contains(got, "including 100 \"pool-0\", ") ||
			!strings.Contains(got, fmt.Sprintf("%d \"pool-%d\"; set resource_pool.id", 100+maxListed-1, maxListed-1)) {
			t.Errorf("expected the first %d entries of the sample, got:\n%s", maxListed, got)
		}
	})
}

func TestUnitDatastoreRejectionIgnoresOtherFailures(t *testing.T) {
	t.Parallel()

	// A client is required by the signature but must never be called when the
	// body is not a rejection this explains; a nil-configured client proves it.
	client := sdk.NewAPIClient(sdk.NewConfiguration())
	req := DatastoreRequest{CloudID: 1, PoolID: 1, StorageServerID: 1}

	for name, body := range map[string]string{
		"another field":   `{"errors":{"protocolType":"Required fields are missing"},"success":false}`,
		"a plain message": `{"msg":"Data store \"x\" already exists for storageServer: \"s\"","success":false}`,
		"not json":        `502 Bad Gateway`,
		"empty":           ``,
		"the storage plugin's duplicate-name refusal": `{"msg":"Datastore with name x already exists","success":false}`,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := DatastoreRejection(context.Background(), client, []byte(body), req); got != "" {
				t.Errorf("expected \"\" for %s, got %q", name, got)
			}
		})
	}

	t.Run("nil client", func(t *testing.T) {
		t.Parallel()

		if got := DatastoreRejection(context.Background(), nil, []byte(unknownPoolBody), req); got != "" {
			t.Errorf("expected \"\" with no client, got %q", got)
		}
	})
}

// TestUnitErrorBodyThenErrMsgKeepsTheBody pins the order the datastore create
// paths use — ErrorBody first, then errfmt.ErrMsg — against the real generated
// client. The client keeps a failed response's body on the error it returns
// and re-wraps the response body, so ErrorBody reads nothing from the
// response here; where it has to (no SDK error), it puts the body back.
// ErrMsg does neither, so the reverse order would be the lossy one.
func TestUnitErrorBodyThenErrMsgKeepsTheBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, unknownPoolBody)
	}))
	t.Cleanup(srv.Close)

	cfg := sdk.NewConfiguration()
	cfg.Servers = sdk.ServerConfigurations{{URL: srv.URL}}
	client := sdk.NewAPIClient(cfg)

	request := sdk.SaveDatastoreRequest{Datastore: &sdk.SaveDatastoreRequestDatastore{
		Name:          "example",
		DatastoreType: "hpedatastore-alletra-mp-bmaas",
		RefType:       "ComputeZone",
		RefId:         1,
		Config:        sdk.SaveDatastoreRequestDatastoreConfig{MapmapOfStringAny: &map[string]any{"protocolType": "iSCSI"}},
	}}

	_, hresp, err := client.DatastoresAPI.SaveDatastore(context.Background()).SaveDatastoreRequest(request).Execute()
	if err == nil || hresp == nil || hresp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected a 400 and an error from the client, got %v, %v", hresp, err)
	}

	var apiErr *sdk.GenericOpenAPIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected the client's GenericOpenAPIError, got %T", err)
	}

	if got := ErrorBody(err, hresp); string(got) != unknownPoolBody {
		t.Fatalf("ErrorBody = %q, want the response body", got)
	}

	msg := errfmt.ErrMsg(err, hresp)
	if !strings.Contains(msg, unknownPoolBody) {
		t.Errorf("ErrMsg after ErrorBody lost the response body: %q", msg)
	}

	if !strings.Contains(msg, http.StatusText(http.StatusBadRequest)) {
		t.Errorf("ErrMsg should name the status, got %q", msg)
	}
}
