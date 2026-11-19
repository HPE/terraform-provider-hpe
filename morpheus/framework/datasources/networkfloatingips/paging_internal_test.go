// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfloatingips

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/utils/paging"
)

// TestUnitListFloatingIpsCollectsEveryPage pins MORPH-20771: a listing larger
// than one page (two full pages and part of a third at the helper's page
// size) is read in full, with the server-side filters repeated on every page
// request. The fake appliance pages by max/offset and reports the total in
// meta, as Morpheus does.
func TestUnitListFloatingIpsCollectsEveryPage(t *testing.T) {
	t.Parallel()

	const total = 2*int(paging.DefaultPageSize) + 30

	// Pages after the first are fetched concurrently, so the counter is atomic.
	var requests atomic.Int32

	mux := http.NewServeMux()
	mux.HandleFunc("/api/networks/floating-ips", func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)

		q := r.URL.Query()
		if q.Get("zoneId") != "7" {
			t.Errorf("page request without the zoneId filter: %s", r.URL.RawQuery)
		}

		maxN, _ := strconv.Atoi(q.Get("max"))
		offset, _ := strconv.Atoi(q.Get("offset"))

		if maxN <= 0 {
			t.Errorf("page request without max: %s", r.URL.RawQuery)
		}

		end := offset + maxN
		if end > total {
			end = total
		}

		page := []map[string]any{}
		for i := offset; i < end; i++ {
			page = append(page, map[string]any{"id": i + 1, "ipAddress": "10.60.0." + strconv.Itoa(i+1)})
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"networkFloatingIps": page,
			"meta":               map[string]any{"total": total, "max": maxN, "offset": offset},
		})
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	cfg := sdk.NewConfiguration()
	cfg.Servers[0].URL = srv.URL
	client := sdk.NewAPIClient(cfg)

	model := &NetworkFloatingIpsModel{CloudId: types.Int64Value(7)}

	items, err := listFloatingIps(context.Background(), client, model)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != total {
		t.Fatalf("collected %d floating IPs, want %d", len(items), total)
	}

	if n := requests.Load(); n < 2 {
		t.Fatalf("listing of %d items was read in %d request(s); paging did not happen", total, n)
	}

	seen := map[int64]bool{}
	for _, it := range items {
		if it.Id == nil || seen[*it.Id] {
			t.Fatalf("duplicate or missing id in collected items: %v", it.Id)
		}

		seen[*it.Id] = true
	}
}
