// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package image

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// fakeImageStore is the smallest Morpheus the upload path can be driven
// against: one virtual image whose files are a name -> contents map.
type fakeImageStore struct {
	mu    sync.Mutex
	files map[string][]byte
	calls []string
	// uploadStatus overrides the upload response when non-zero.
	uploadStatus int
	uploadBody   string
	// hideUploaded makes the listing omit a freshly uploaded file, as slow
	// storage does, for the given number of GETs.
	hideUploaded int
	hidden       map[string]int
}

func newFakeImageStore(files map[string][]byte) *fakeImageStore {
	if files == nil {
		files = map[string][]byte{}
	}

	return &fakeImageStore{files: files, hidden: map[string]int{}}
}

func (f *fakeImageStore) handler(t *testing.T) http.Handler {
	t.Helper()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()

		f.calls = append(f.calls, r.Method+" "+r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/virtual-images/42":
			var list []map[string]any
			for name, data := range f.files {
				if f.hidden[name] > 0 {
					f.hidden[name]--

					continue
				}
				list = append(list, map[string]any{
					"name": "virtual-images/42/" + name, "size": 0.1, "contentLength": len(data),
				})
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"virtualImage": map[string]any{"id": 42, "name": "img"},
				"cloudFiles":   list,
			})
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/virtual-images/"):
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"success":false,"msg":"Virtual image not found"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/api/virtual-images/42/files":
			if r.Header.Get("Authorization") == "" {
				t.Errorf("upload without Authorization header")
			}
			if r.ContentLength <= 0 {
				t.Errorf("upload without a declared Content-Length")
			}
			data, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("reading upload body: %v", err)
			}
			if f.uploadStatus != 0 {
				w.WriteHeader(f.uploadStatus)
				fmt.Fprint(w, f.uploadBody)

				return
			}
			name := r.URL.Query().Get("filename")
			f.files[name] = data
			if f.hideUploaded > 0 {
				f.hidden[name] = f.hideUploaded
			}
			fmt.Fprint(w, `{"success":true}`)
		case r.Method == http.MethodDelete && r.URL.Path == "/api/virtual-images/42/files":
			name := r.URL.Query().Get("filename")
			if _, ok := f.files[name]; !ok {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, `{"success":false,"msg":"File not found"}`)

				return
			}
			delete(f.files, name)
			fmt.Fprint(w, `{"success":true}`)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.RequestURI())
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

// bearerTransport plays the provider's authenticating transport.
type bearerTransport struct{ base http.RoundTripper }

func (b bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer test-token")

	return b.base.RoundTrip(req)
}

func newFakeClient(t *testing.T, store *fakeImageStore) *sdk.APIClient {
	t.Helper()

	srv := httptest.NewServer(store.handler(t))
	t.Cleanup(srv.Close)

	cfg := sdk.NewConfiguration()
	cfg.Servers = sdk.ServerConfigurations{{URL: srv.URL}}
	cfg.HTTPClient = &http.Client{Transport: bearerTransport{base: http.DefaultTransport}, Timeout: 10 * time.Second}

	return sdk.NewAPIClient(cfg)
}

func writeTempFile(t *testing.T, name string, data []byte) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	return path
}

func sumOf(data []byte) string {
	h := sha256.Sum256(data)

	return hex.EncodeToString(h[:])
}

func TestUnitInspectLocalFile(t *testing.T) {
	t.Parallel()

	data := []byte("disk image bytes")
	path := writeTempFile(t, "disk.qcow2", data)

	got, err := inspectLocalFile(path, "")
	if err != nil {
		t.Fatal(err)
	}
	if got.Size != int64(len(data)) || got.Sha256 != sumOf(data) {
		t.Fatalf("inspect = %+v", got)
	}

	if _, err := inspectLocalFile(path, strings.ToUpper(sumOf(data))); err != nil {
		t.Fatalf("upper-case hash should match: %v", err)
	}

	if _, err := inspectLocalFile(path, strings.Repeat("0", 64)); err == nil ||
		!strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("want checksum mismatch, got %v", err)
	}

	if _, err := inspectLocalFile(writeTempFile(t, "empty", nil), ""); err == nil {
		t.Fatal("empty file should be refused")
	}

	if _, err := inspectLocalFile(filepath.Join(t.TempDir(), "missing"), ""); err == nil {
		t.Fatal("missing file should be refused")
	}
}

func TestUnitImageFilesDecodesCloudFiles(t *testing.T) {
	t.Parallel()

	resp := &sdk.GetVirtualImage200Response{AdditionalProperties: map[string]any{
		"cloudFiles": []any{
			map[string]any{"name": "virtual-images/42/a.qcow2", "size": 0.5, "contentLength": float64(1234)},
			map[string]any{"name": "b.iso", "contentLength": float64(99)},
			map[string]any{"name": "dir/", "contentLength": float64(0)},
		},
	}}

	files, err := imageFiles(resp)
	if err != nil {
		t.Fatal(err)
	}

	want := []fileEntry{{"a.qcow2", 1234}, {"b.iso", 99}}
	if len(files) != len(want) {
		t.Fatalf("files = %v, want %v", files, want)
	}
	for i := range want {
		if files[i] != want[i] {
			t.Fatalf("files[%d] = %v, want %v", i, files[i], want[i])
		}
	}

	none, err := imageFiles(&sdk.GetVirtualImage200Response{})
	if err != nil || len(none) != 0 {
		t.Fatalf("no cloudFiles: %v, %v", none, err)
	}
}

func TestUnitUploadStreamsWithAuth(t *testing.T) {
	t.Parallel()

	store := newFakeImageStore(nil)
	client := newFakeClient(t, store)
	data := []byte("the whole disk image")
	path := writeTempFile(t, "disk.qcow2", data)
	local, err := inspectLocalFile(path, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := uploadFile(context.Background(), client, 42, "disk.qcow2", local); err != nil {
		t.Fatal(err)
	}

	if got := string(store.files["disk.qcow2"]); got != string(data) {
		t.Fatalf("stored %q, want %q", got, data)
	}
	if len(store.calls) != 1 || store.calls[0] != "POST /api/virtual-images/42/files?filename=disk.qcow2" {
		t.Fatalf("calls = %v", store.calls)
	}
}

func TestUnitUploadReportsFailure(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{
			"envelope failure", http.StatusOK,
			`{"success":false,"msg":"Missing required parameter: filename"}`, "Missing required parameter",
		},
		{"server error", http.StatusInternalServerError, `{"success":false,"msg":"Image upload failed."}`, "HTTP 500"},
		{"non-json error", http.StatusBadGateway, `<html>gateway</html>`, "gateway"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := newFakeImageStore(nil)
			store.uploadStatus = tc.status
			store.uploadBody = tc.body
			client := newFakeClient(t, store)
			path := writeTempFile(t, "disk.qcow2", []byte("x"))
			local, _ := inspectLocalFile(path, "")

			err := uploadFile(context.Background(), client, 42, "disk.qcow2", local)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err = %v, want %q", err, tc.want)
			}
		})
	}
}

func TestUnitAwaitListedWaitsForStorage(t *testing.T) {
	t.Parallel()

	data := []byte("bytes")
	store := newFakeImageStore(map[string][]byte{"disk.qcow2": data})
	store.hidden["disk.qcow2"] = 1
	client := newFakeClient(t, store)

	// Too short a timeout to allow the 5 s poll: the first listing hides the
	// file, so the deadline check reports the listing.
	_, err := awaitListed(context.Background(), client, 42, "disk.qcow2", int64(len(data)), 0)
	if err == nil || !strings.Contains(err.Error(), "does not list") {
		t.Fatalf("want a listing error, got %v", err)
	}

	store.mu.Lock()
	store.hidden["disk.qcow2"] = 1
	store.calls = nil
	store.mu.Unlock()
	files, err := awaitListed(context.Background(), client, 42, "disk.qcow2", int64(len(data)), 15*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := findFile(files, "disk.qcow2"); !ok {
		t.Fatalf("files = %v", files)
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if len(store.calls) != 2 {
		t.Fatalf("expected the hidden listing and one retry, got %v", store.calls)
	}
}

func TestUnitRemoveFileToleratesMissing(t *testing.T) {
	t.Parallel()

	store := newFakeImageStore(map[string][]byte{"old.qcow2": []byte("old")})
	client := newFakeClient(t, store)

	if err := removeFile(context.Background(), client, 42, "old.qcow2"); err != nil {
		t.Fatal(err)
	}
	if err := removeFile(context.Background(), client, 42, "old.qcow2"); err != nil {
		t.Fatalf("second removal should be tolerated: %v", err)
	}
	if len(store.files) != 0 {
		t.Fatalf("files left: %v", store.files)
	}
}

func TestUnitRemoveOthersKeepsTheNewFile(t *testing.T) {
	t.Parallel()

	store := newFakeImageStore(map[string][]byte{
		"new.qcow2": []byte("new"), "old.qcow2": []byte("old"), "older.qcow2": []byte("older"),
	})
	client := newFakeClient(t, store)
	files, _, err := listFiles(context.Background(), client, 42)
	if err != nil {
		t.Fatal(err)
	}

	if err := removeOthers(context.Background(), client, 42, files, "new.qcow2"); err != nil {
		t.Fatal(err)
	}

	if len(store.files) != 1 || store.files["new.qcow2"] == nil {
		t.Fatalf("files = %v", store.files)
	}
}

func TestUnitListFilesNotFound(t *testing.T) {
	t.Parallel()

	client := newFakeClient(t, newFakeImageStore(nil))

	_, httpResp, err := listFiles(context.Background(), client, 7)
	if err == nil || httpResp == nil || httpResp.StatusCode != http.StatusNotFound {
		t.Fatalf("want 404, got %v / %v", err, httpResp)
	}
}
