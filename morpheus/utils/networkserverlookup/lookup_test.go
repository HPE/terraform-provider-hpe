// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserverlookup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

type fakeServer struct {
	id      int64
	name    string
	visible bool
}

// fakeOptions bend the fake away from Morpheus's behaviour, to exercise the
// lookup's defences.
type fakeOptions struct {
	ignoreOffset bool // every page is the first page
	omitTotal    bool // no meta.total in the response
}

// fakeAppliance answers the two network-server routes the way Morpheus does:
// the single-server route 404s for a hidden server whatever is asked, and
// the listing returns hidden servers only with showHidden=true, filters on
// name, and pages with max/offset, reporting the total in meta.
func fakeAppliance(t *testing.T, servers []fakeServer) (*sdk.APIClient, func()) {
	t.Helper()

	return fakeApplianceWith(t, servers, fakeOptions{})
}

func fakeApplianceWith(t *testing.T, servers []fakeServer, opts fakeOptions) (*sdk.APIClient, func()) {
	t.Helper()

	encode := func(s fakeServer) map[string]any {
		return map[string]any{
			"id": s.id, "name": s.name, "visible": s.visible,
			"type": map[string]any{"code": "morpheus-hvm-network"},
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/networks/servers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		q := r.URL.Query()
		showHidden := q.Get("showHidden") == "true"
		name := q.Get("name")
		maxN, _ := strconv.Atoi(q.Get("max"))
		offset, _ := strconv.Atoi(q.Get("offset"))

		if opts.ignoreOffset {
			offset = 0
		}

		var all []map[string]any
		for _, s := range servers {
			if (s.visible || showHidden) && (name == "" || s.name == name) {
				all = append(all, encode(s))
			}
		}

		page := []map[string]any{}
		if offset < len(all) {
			end := len(all)
			if maxN > 0 && offset+maxN < end {
				end = offset + maxN
			}

			page = all[offset:end]
		}

		body := map[string]any{"networkServers": page}
		if !opts.omitTotal {
			body["meta"] = map[string]any{"total": len(all), "max": maxN, "offset": offset}
		}

		_ = json.NewEncoder(w).Encode(body)
	})
	mux.HandleFunc("/api/networks/servers/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := strings.TrimPrefix(r.URL.Path, "/api/networks/servers/")
		for _, s := range servers {
			if s.visible && fmt.Sprint(s.id) == id {
				_ = json.NewEncoder(w).Encode(map[string]any{"networkServer": encode(s)})

				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "msg": "Network Service not found with id " + id})
	})

	srv := httptest.NewServer(mux)

	cfg := sdk.NewConfiguration()
	cfg.Servers[0].URL = srv.URL

	return sdk.NewAPIClient(cfg), srv.Close
}

// TestByIDFindsHiddenServer pins that a hidden server is found by id through
// the listing when the single-server route refuses it, that a visible one
// still comes from the single-server route, and that an id no server has is
// reported as not found.
func TestByIDFindsHiddenServer(t *testing.T) {
	t.Parallel()

	client, closeFn := fakeAppliance(t, []fakeServer{
		{3, "zodiac-hvm HVM Network Server", false},
		{16, "Labs NSX 3", true},
	})
	defer closeFn()

	ctx := context.Background()

	hidden, err := ByID(ctx, client, 3)
	if err != nil {
		t.Fatalf("hidden server by id: %v", err)
	}

	if hidden.Name == nil || *hidden.Name != "zodiac-hvm HVM Network Server" {
		t.Errorf("hidden server by id: got %+v", hidden)
	}

	visible, err := ByID(ctx, client, 16)
	if err != nil {
		t.Fatalf("visible server by id: %v", err)
	}

	if visible.Name == nil || *visible.Name != "Labs NSX 3" {
		t.Errorf("visible server by id: got %+v", visible)
	}

	if _, err := ByID(ctx, client, 99); err == nil || !strings.Contains(err.Error(), "404") {
		t.Errorf("absent id: got %v, want a 404 error", err)
	}
}

// TestByIDPagesThroughTheListing pins that a hidden server beyond the first
// page of the listing is still found, and that an absent id is reported only
// after the whole listing has been read.
func TestByIDPagesThroughTheListing(t *testing.T) {
	t.Parallel()

	servers := make([]fakeServer, 0, int(pageSize)*2+5)
	for i := int64(1); i <= pageSize*2+5; i++ {
		servers = append(servers, fakeServer{i, fmt.Sprintf("hidden-%d", i), false})
	}

	client, closeFn := fakeAppliance(t, servers)
	defer closeFn()

	ctx := context.Background()

	// Sits on the third page.
	want := pageSize*2 + 3

	got, err := ByID(ctx, client, want)
	if err != nil {
		t.Fatalf("hidden server on the third page: %v", err)
	}

	if got.Id == nil || *got.Id != want {
		t.Errorf("hidden server on the third page: got %+v", got)
	}

	if _, err := ByID(ctx, client, pageSize*10); err == nil || !strings.Contains(err.Error(), "404") {
		t.Errorf("absent id beyond every page: got %v, want a 404 error", err)
	}
}

// TestByIDBoundsThePaging pins that an appliance which ignores offset — so
// every page is full and identical — does not loop forever: the read stops
// at the paging bound and reports the listing truncated, with or without a
// total in the response.
func TestByIDBoundsThePaging(t *testing.T) {
	t.Parallel()

	servers := make([]fakeServer, 0, int(pageSize))
	for i := int64(1); i <= pageSize; i++ {
		servers = append(servers, fakeServer{i, fmt.Sprintf("hidden-%d", i), false})
	}

	for name, opts := range map[string]fakeOptions{
		"offset ignored, total omitted":  {ignoreOffset: true, omitTotal: true},
		"offset ignored, total reported": {ignoreOffset: true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client, closeFn := fakeApplianceWith(t, servers, opts)
			defer closeFn()

			_, err := ByID(context.Background(), client, pageSize+1)
			if opts.omitTotal {
				if !errors.Is(err, ErrTruncated) {
					t.Errorf("got %v, want ErrTruncated", err)
				}
			} else if err == nil || !strings.Contains(err.Error(), "404") {
				// With a total reported, the first full page already covers
				// it, so the listing is exhausted and the id is absent.
				t.Errorf("got %v, want a 404 error", err)
			}
		})
	}
}

// TestByName pins the by-name rules: a hidden server is found by its exact
// name, the listing is asked for that name, an absent name is ErrNotFound,
// and two servers with the name is ErrMultiple.
func TestByName(t *testing.T) {
	t.Parallel()

	client, closeFn := fakeAppliance(t, []fakeServer{
		{3, "zodiac-hvm HVM Network Server", false},
		{16, "Labs NSX 3", true},
		{21, "twin", true},
		{22, "twin", false},
	})
	defer closeFn()

	ctx := context.Background()

	hidden, err := ByName(ctx, client, "zodiac-hvm HVM Network Server")
	if err != nil {
		t.Fatalf("hidden server by name: %v", err)
	}

	if hidden.Id == nil || *hidden.Id != 3 {
		t.Errorf("hidden server by name: got %+v", hidden)
	}

	if _, err := ByName(ctx, client, "no such server"); !errors.Is(err, ErrNotFound) {
		t.Errorf("absent name: got %v, want ErrNotFound", err)
	}

	if _, err := ByName(ctx, client, "twin"); !errors.Is(err, ErrMultiple) {
		t.Errorf("duplicate name: got %v, want ErrMultiple", err)
	}
}
