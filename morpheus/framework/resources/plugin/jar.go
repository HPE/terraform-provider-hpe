// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plugin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// resolvedJar is an open plugin jar file ready for upload, together with the
// hex-encoded SHA-256 of its contents. cleanup must always be called by the
// caller once the upload is done; it closes the file and removes any temp file.
type resolvedJar struct {
	file    *os.File
	sha256  string
	cleanup func()
}

// resolveJar opens the jar referenced by either jarPath or jarURL (exactly one
// must be non-empty), computes its SHA-256, and — when wantSHA is set — verifies
// it matches. The multipart upload names the part from the file's name, so a
// downloaded jar is written to a temp file whose name ends in ".jar".
func resolveJar(
	ctx context.Context,
	jarPath, jarURL, wantSHA string,
) (*resolvedJar, error) {
	var (
		file    *os.File
		cleanup func()
		err     error
	)

	switch {
	case jarPath != "":
		file, err = os.Open(jarPath) //nolint:gosec // user-supplied path is intended
		if err != nil {
			return nil, fmt.Errorf("could not read plugin jar %q: %w", jarPath, err)
		}
		cleanup = func() { _ = file.Close() }
	case jarURL != "":
		file, cleanup, err = downloadJar(ctx, jarURL)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("exactly one of jar_path or jar_url must be set")
	}

	sum, err := hashReaderAtStart(file)
	if err != nil {
		cleanup()

		return nil, fmt.Errorf("hashing plugin jar: %w", err)
	}

	if wantSHA != "" && !strings.EqualFold(wantSHA, sum) {
		cleanup()

		return nil, fmt.Errorf(
			"plugin jar checksum mismatch: expected %s, got %s", wantSHA, sum)
	}

	return &resolvedJar{file: file, sha256: sum, cleanup: cleanup}, nil
}

const (
	// maxJarBytes bounds a jar_url download. The largest plugin on the Morpheus
	// Exchange is about 110 MB; this leaves room for growth without letting a
	// wrong or hostile URL fill the disk of the machine running Terraform.
	maxJarBytes = 512 << 20

	// maxJarRedirects bounds the redirects followed for a jar_url. The Exchange
	// answers with one, to its CDN.
	maxJarRedirects = 5

	// jarDialTimeout bounds one connection attempt while downloading a jar.
	jarDialTimeout = 30 * time.Second
)

// jarFetchPolicy bounds what a jar_url download may do. Tests relax it so an
// httptest server on the loopback interface can stand in for a plugin source.
type jarFetchPolicy struct {
	maxBytes     int64
	maxRedirects int
	// allowIP reports whether an address a jar_url host resolves to may be
	// connected to.
	allowIP func(net.IP) bool
}

var defaultJarFetchPolicy = jarFetchPolicy{
	maxBytes:     maxJarBytes,
	maxRedirects: maxJarRedirects,
	allowIP:      isPluginSourceAddress,
}

// isPluginSourceAddress reports whether ip is an address a plugin jar could
// legitimately be served from: global or private unicast. Loopback,
// link-local (where cloud metadata services such as 169.254.169.254 live),
// unspecified and multicast addresses are not, and a jar_url that leads to one
// is refused. Private networks are allowed: an internal artifact server is a
// common place to keep plugin jars.
func isPluginSourceAddress(ip net.IP) bool {
	switch {
	case ip == nil,
		ip.IsUnspecified(),
		ip.IsLoopback(),
		ip.IsLinkLocalUnicast(),
		ip.IsLinkLocalMulticast(),
		ip.IsInterfaceLocalMulticast(),
		ip.IsMulticast():
		return false
	default:
		return true
	}
}

// checkJarURLStatic checks what can be known of a jar URL without the network:
// an http(s) scheme, a host, and -- when the host is a literal address -- one
// the policy allows.
func checkJarURLStatic(u *url.URL, policy jarFetchPolicy) error {
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("jar_url must be an http:// or https:// URL, got: %s", u.Redacted())
	}

	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("jar_url has no host: %s", u.Redacted())
	}

	if ip := net.ParseIP(host); ip != nil && !policy.allowIP(ip) {
		return fmt.Errorf("jar_url %s: %s is not an address a plugin can be fetched from", u.Redacted(), ip)
	}

	return nil
}

// checkJarURL is checkJarURLStatic plus a look at what a host name resolves
// to. A name that cannot be resolved here is let through -- a proxy may
// resolve it -- and the connection then succeeds or fails on its own merits.
func checkJarURL(ctx context.Context, u *url.URL, policy jarFetchPolicy) error {
	if err := checkJarURLStatic(u, policy); err != nil {
		return err
	}

	host := u.Hostname()
	if net.ParseIP(host) != nil {
		return nil
	}

	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil //nolint:nilerr // not resolvable here is not refused: a proxy may resolve it
	}

	for _, a := range addrs {
		if !policy.allowIP(a.IP) {
			return fmt.Errorf("jar_url %s: %s resolves to %s, which is not an address a plugin can be fetched from",
				u.Redacted(), host, a.IP)
		}
	}

	return nil
}

// validateJarURL is the plan-time check of jar_url: what checkJarURLStatic
// checks, and -- when the URL is plain http -- that jar_sha256 is set, since
// without TLS the download cannot be trusted on its own. Names are not
// resolved at plan time; the fetch checks them when the jar is downloaded.
func validateJarURL(rawURL string, hasSHA bool) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("jar_url must be an http:// or https:// URL: %w", err)
	}

	if err := checkJarURLStatic(u, defaultJarFetchPolicy); err != nil {
		return err
	}

	if u.Scheme == "http" && !hasSHA {
		return errors.New("jar_url uses plain http://, so the download cannot be trusted on its own: " +
			"set jar_sha256 to the jar's SHA-256, or use an https:// URL")
	}

	return nil
}

// jarHTTPClient returns the client a jar_url is fetched with. It keeps the
// default transport's proxy handling (HTTPS_PROXY and friends), connects only
// to addresses the policy allows -- decided at dial time, so a name cannot
// change its answer between a check and the connection -- and follows at most
// policy.maxRedirects redirects, each checked like the original URL.
func jarHTTPClient(policy jarFetchPolicy) *http.Client {
	var transport *http.Transport
	if dt, ok := http.DefaultTransport.(*http.Transport); ok {
		transport = dt.Clone()
	} else {
		transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	}

	dialer := &net.Dialer{Timeout: jarDialTimeout}
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}

		addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
		if err != nil {
			return nil, err
		}

		var lastErr error
		for _, a := range addrs {
			if !policy.allowIP(a.IP) {
				lastErr = fmt.Errorf("%s is not an address a plugin can be fetched from", a.IP)

				continue
			}

			conn, err := dialer.DialContext(ctx, network, net.JoinHostPort(a.IP.String(), port))
			if err == nil {
				return conn, nil
			}
			lastErr = err
		}

		if lastErr == nil {
			lastErr = errors.New("no address")
		}

		return nil, fmt.Errorf("connecting to %s: %w", host, lastErr)
	}

	return &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > policy.maxRedirects {
				return fmt.Errorf("more than %d redirects", policy.maxRedirects)
			}

			return checkJarURL(req.Context(), req.URL, policy)
		},
	}
}

// downloadJar fetches the jar at rawURL into a temp file ending in ".jar" and
// returns the rewound file plus a cleanup that closes and removes it, within
// the bounds of defaultJarFetchPolicy.
func downloadJar(ctx context.Context, rawURL string) (*os.File, func(), error) {
	return downloadJarWith(ctx, rawURL, defaultJarFetchPolicy)
}

func downloadJarWith(ctx context.Context, rawURL string, policy jarFetchPolicy) (*os.File, func(), error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing jar_url %q: %w", rawURL, err)
	}

	if err := checkJarURL(ctx, u, policy); err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("building request for %s: %w", u.Redacted(), err)
	}

	resp, err := jarHTTPClient(policy).Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("downloading plugin jar from %s: %w", u.Redacted(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf(
			"downloading plugin jar from %s: HTTP %d", u.Redacted(), resp.StatusCode)
	}

	if resp.ContentLength > policy.maxBytes {
		return nil, nil, fmt.Errorf("plugin jar from %s is %s, over the %s limit",
			u.Redacted(), formatBytes(resp.ContentLength), formatBytes(policy.maxBytes))
	}

	tmp, err := os.CreateTemp("", "plugin-*.jar")
	if err != nil {
		return nil, nil, fmt.Errorf("creating temp file for plugin jar: %w", err)
	}

	cleanup := func() {
		name := tmp.Name()
		_ = tmp.Close()
		_ = os.Remove(name)
	}

	// One byte past the limit is enough to know the limit was exceeded without
	// reading the rest.
	n, err := io.Copy(tmp, io.LimitReader(resp.Body, policy.maxBytes+1))
	if err != nil {
		cleanup()

		return nil, nil, fmt.Errorf("writing plugin jar from %s: %w", u.Redacted(), err)
	}

	if n > policy.maxBytes {
		cleanup()

		return nil, nil, fmt.Errorf("plugin jar from %s is over the %s limit",
			u.Redacted(), formatBytes(policy.maxBytes))
	}

	return tmp, cleanup, nil
}

// formatBytes renders a byte count for a message: whole mebibytes when it is
// one, bytes otherwise.
func formatBytes(n int64) string {
	if n >= 1<<20 && n%(1<<20) == 0 {
		return fmt.Sprintf("%d MiB", n>>20)
	}

	return fmt.Sprintf("%d bytes", n)
}

// hashReaderAtStart hashes the whole file and rewinds it to the start so the
// caller can upload it afterwards.
func hashReaderAtStart(f *os.File) (string, error) {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// hashFile computes the hex SHA-256 of the file at path without keeping it open.
func hashFile(path string) (string, error) {
	f, err := os.Open(path) //nolint:gosec // user-supplied path is intended
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
