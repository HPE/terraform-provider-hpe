// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plugin

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

// loopbackPolicy is the default fetch policy with loopback allowed, so an
// httptest server can stand in for a plugin source, and a small size limit.
func loopbackPolicy(maxBytes int64) jarFetchPolicy {
	return jarFetchPolicy{
		maxBytes:     maxBytes,
		maxRedirects: 2,
		allowIP: func(ip net.IP) bool {
			return ip.IsLoopback() || isPluginSourceAddress(ip)
		},
	}
}

// readJar returns the file's whole content. downloadJar leaves the file at its
// end (resolveJar rewinds it when hashing), so seek first.
func readJar(t *testing.T, f *os.File) string {
	t.Helper()

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("seek: %v", err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("read jar: %v", err)
	}

	return string(b)
}

func TestDownloadJarFollowsRedirects(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("/pull", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/cdn/plugin.jar", http.StatusFound)
	})
	mux.HandleFunc("/cdn/plugin.jar", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("jar-bytes"))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, cleanup, err := downloadJarWith(context.Background(), srv.URL+"/pull", loopbackPolicy(1<<20))
	if err != nil {
		t.Fatalf("download: %v", err)
	}
	defer cleanup()

	if got := readJar(t, f); got != "jar-bytes" {
		t.Errorf("content = %q, want jar-bytes", got)
	}

	if !strings.HasSuffix(f.Name(), ".jar") {
		t.Errorf("temp file %q should end in .jar (the upload names the part after it)", f.Name())
	}

	name := f.Name()
	cleanup()

	if _, err := os.Stat(name); !os.IsNotExist(err) {
		t.Errorf("cleanup left %s behind (stat err = %v)", name, err)
	}
}

func TestDownloadJarRefusesTooManyRedirects(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/again", http.StatusFound)
	}))
	defer srv.Close()

	_, _, err := downloadJarWith(context.Background(), srv.URL+"/start", loopbackPolicy(1<<20))
	if err == nil || !strings.Contains(err.Error(), "redirects") {
		t.Fatalf("err = %v, want a redirect limit error", err)
	}
}

func TestDownloadJarRefusesRedirectToBlockedAddress(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://169.254.169.254/latest/meta-data/", http.StatusFound)
	}))
	defer srv.Close()

	_, _, err := downloadJarWith(context.Background(), srv.URL+"/pull", loopbackPolicy(1<<20))
	if err == nil || !strings.Contains(err.Error(), "not an address a plugin can be fetched from") {
		t.Fatalf("err = %v, want the redirect target refused", err)
	}
}

func TestDownloadJarSizeLimit(t *testing.T) {
	t.Parallel()

	const limit = 1024

	tests := []struct {
		name    string
		size    int
		declare bool // send Content-Length; otherwise stream chunked
		wantErr bool
	}{
		{name: "declared at the limit", size: limit, declare: true},
		{name: "declared over the limit", size: limit + 1, declare: true, wantErr: true},
		{name: "streamed at the limit", size: limit},
		{name: "streamed over the limit", size: limit + 1, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			body := strings.Repeat("x", tc.size)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if tc.declare {
					w.Header().Set("Content-Length", strconv.Itoa(len(body)))
					_, _ = w.Write([]byte(body))

					return
				}
				// Commit the headers before the body so no Content-Length is
				// sent and the limit has to be enforced on the stream.
				_, _ = w.Write([]byte(body[:1]))
				w.(http.Flusher).Flush()
				_, _ = w.Write([]byte(body[1:]))
			}))
			defer srv.Close()

			f, cleanup, err := downloadJarWith(context.Background(), srv.URL+"/plugin.jar", loopbackPolicy(limit))
			if tc.wantErr {
				if err == nil || !strings.Contains(err.Error(), "limit") {
					t.Fatalf("err = %v, want a size limit error", err)
				}

				return
			}

			if err != nil {
				t.Fatalf("download: %v", err)
			}
			defer cleanup()

			if got := readJar(t, f); got != body {
				t.Errorf("got %d bytes, want %d", len(got), len(body))
			}
		})
	}
}

func TestDownloadJarRefusesBlockedURLs(t *testing.T) {
	t.Parallel()

	// Literal addresses and bad URLs: refused before any connection is made.
	tests := []struct {
		url  string
		want string
	}{
		{url: "http://127.0.0.1:1/plugin.jar", want: "not an address a plugin can be fetched from"},
		{url: "http://[::1]/plugin.jar", want: "not an address a plugin can be fetched from"},
		{url: "http://169.254.169.254/latest/meta-data/", want: "not an address a plugin can be fetched from"},
		{url: "http://0.0.0.0/plugin.jar", want: "not an address a plugin can be fetched from"},
		{url: "http://224.0.0.1/plugin.jar", want: "not an address a plugin can be fetched from"},
		{url: "ftp://plugins.example.com/plugin.jar", want: "must be an http:// or https:// URL"},
		{url: "http:///plugin.jar", want: "has no host"},
	}

	for _, tc := range tests {
		t.Run(tc.url, func(t *testing.T) {
			t.Parallel()

			_, _, err := downloadJarWith(context.Background(), tc.url, defaultJarFetchPolicy)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err = %v, want %q", err, tc.want)
			}
		})
	}
}

// The dial-time check is what stands when a name resolves to a refused
// address only at connection time, so it is exercised on its own.
func TestJarHTTPClientRefusesBlockedDial(t *testing.T) {
	t.Parallel()

	transport, ok := jarHTTPClient(defaultJarFetchPolicy).Transport.(*http.Transport)
	if !ok {
		t.Fatal("transport is not an *http.Transport")
	}

	_, err := transport.DialContext(context.Background(), "tcp", "127.0.0.1:1")
	if err == nil || !strings.Contains(err.Error(), "not an address a plugin can be fetched from") {
		t.Fatalf("err = %v, want the loopback address refused", err)
	}
}

func TestIsPluginSourceAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ip   string
		want bool
	}{
		{"203.0.113.10", true},     // public
		{"10.227.122.211", true},   // private: an internal artifact server is fine
		{"192.168.1.20", true},     // private
		{"2001:db8::1", true},      // public v6
		{"fd12:3456::1", true},     // unique local v6 (private)
		{"127.0.0.1", false},       // loopback
		{"::1", false},             // loopback v6
		{"169.254.169.254", false}, // link-local: cloud metadata
		{"fe80::1", false},         // link-local v6
		{"0.0.0.0", false},         // unspecified
		{"::", false},              // unspecified v6
		{"224.0.0.251", false},     // multicast
		{"ff02::1", false},         // multicast v6
	}

	for _, tc := range tests {
		if got := isPluginSourceAddress(net.ParseIP(tc.ip)); got != tc.want {
			t.Errorf("isPluginSourceAddress(%s) = %v, want %v", tc.ip, got, tc.want)
		}
	}

	if isPluginSourceAddress(nil) {
		t.Error("isPluginSourceAddress(nil) = true, want false")
	}
}

func TestValidateJarURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		url    string
		hasSHA bool
		want   string // substring of the error, "" for no error
	}{
		{name: "https without digest", url: "https://share.morpheusdata.com/api/packages/pull?code=x&version=1"},
		{name: "https with digest", url: "https://share.morpheusdata.com/api/packages/pull?code=x&version=1", hasSHA: true},
		{name: "http with digest", url: "http://artifacts.internal/plugins/x.jar", hasSHA: true},
		{name: "http without digest", url: "http://artifacts.internal/plugins/x.jar", want: "set jar_sha256"},
		{name: "literal loopback", url: "https://127.0.0.1/x.jar", hasSHA: true, want: "not an address"},
		{name: "literal metadata", url: "http://169.254.169.254/x.jar", hasSHA: true, want: "not an address"},
		{name: "wrong scheme", url: "ftp://artifacts.internal/x.jar", hasSHA: true, want: "must be an http:// or https:// URL"},
		{name: "no host", url: "https:///x.jar", hasSHA: true, want: "has no host"},
		{name: "unparseable", url: "https://bad host/x.jar", hasSHA: true, want: "must be an http:// or https:// URL"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := validateJarURL(tc.url, tc.hasSHA)
			switch {
			case tc.want == "" && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case tc.want != "" && (err == nil || !strings.Contains(err.Error(), tc.want)):
				t.Fatalf("err = %v, want %q", err, tc.want)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	t.Parallel()

	tests := map[int64]string{
		1024:      "1024 bytes",
		1 << 20:   "1 MiB",
		512 << 20: "512 MiB",
		1<<20 + 1: "1048577 bytes",
	}
	for n, want := range tests {
		if got := formatBytes(n); got != want {
			t.Errorf("formatBytes(%d) = %q, want %q", n, got, want)
		}
	}
}
