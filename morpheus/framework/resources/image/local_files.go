// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package image

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// filenameRE accepts a bare file name: no path separators, so the name
// cannot escape the image's own directory on the storage bucket.
var filenameRE = regexp.MustCompile(`^[^/\\]+$`)

// fileEntry is one file of a virtual image record, as GET
// /api/virtual-images/{id} lists it under cloudFiles. The name is relative
// to the image's directory; contentLength is the size in bytes.
type fileEntry struct {
	Name string
	Size int64
}

// imageFiles decodes the cloudFiles list of a virtual image response. The
// field is not in the OpenAPI schema, so it arrives in the response's
// additional properties.
func imageFiles(resp *sdk.GetVirtualImage200Response) ([]fileEntry, error) {
	if resp == nil {
		return nil, nil
	}

	raw, ok := resp.AdditionalProperties["cloudFiles"]
	if !ok || raw == nil {
		return nil, nil
	}

	encoded, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("re-encoding cloudFiles: %w", err)
	}

	var items []struct {
		Name          string `json:"name"`
		ContentLength int64  `json:"contentLength"`
	}
	if err := json.Unmarshal(encoded, &items); err != nil {
		return nil, fmt.Errorf("decoding cloudFiles: %w", err)
	}

	files := make([]fileEntry, 0, len(items))
	for _, item := range items {
		// Storage providers may report the path inside the bucket; the
		// resource names files by their base name, as the delete API does.
		name := item.Name
		if i := strings.LastIndex(name, "/"); i >= 0 {
			name = name[i+1:]
		}

		if name == "" {
			continue
		}

		files = append(files, fileEntry{Name: name, Size: item.ContentLength})
	}

	return files, nil
}

// findFile returns the entry named name, if the listing has one.
func findFile(files []fileEntry, name string) (fileEntry, bool) {
	for _, f := range files {
		if f.Name == name {
			return f, true
		}
	}

	return fileEntry{}, false
}

// listFiles fetches the image record and returns its files. A 404 is
// reported through errfmt.IsNotFound on the returned response.
func listFiles(
	ctx context.Context,
	client *sdk.APIClient,
	imageID int64,
) ([]fileEntry, *http.Response, error) {
	resp, httpResp, err := client.LibraryAPI.GetVirtualImage(ctx, imageID).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		return nil, httpResp, err
	}

	if resp == nil || resp.VirtualImage == nil || resp.VirtualImage.Id == nil {
		return nil, httpResp, fmt.Errorf("virtual image %d: API returned an empty image", imageID)
	}

	files, err := imageFiles(resp)
	if err != nil {
		return nil, httpResp, fmt.Errorf("virtual image %d: %w", imageID, err)
	}

	return files, httpResp, nil
}

// removeFile deletes one file from the image record with DELETE
// /api/virtual-images/{id}/files?filename=<name>. A file the appliance no
// longer has is not an error: the caller wanted it gone.
func removeFile(ctx context.Context, client *sdk.APIClient, imageID int64, name string) error {
	req, err := newFileRequest(ctx, client, http.MethodDelete, imageID, name, nil)
	if err != nil {
		return err
	}

	resp, err := client.GetConfig().HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("removing %q from virtual image %d: %w", name, imageID, err)
	}
	defer resp.Body.Close()

	raw, env, err := readEnvelope(resp)
	if err != nil {
		return fmt.Errorf("removing %q from virtual image %d: %w", name, imageID, err)
	}

	switch {
	case resp.StatusCode == http.StatusNotFound,
		resp.StatusCode == http.StatusBadRequest && strings.EqualFold(env.Msg, "File not found"):
		return nil
	case resp.StatusCode >= http.StatusBadRequest, env.Success != nil && !*env.Success:
		return fmt.Errorf("removing %q from virtual image %d failed (HTTP %d): %s",
			name, imageID, resp.StatusCode, envelopeMessage(env, raw))
	}

	return nil
}

// newFileRequest builds a request to the image's files endpoint, which the
// OpenAPI spec does not model in a form the generated SDK can stream: the
// generated client buffers every request body in memory, which a disk image
// cannot afford. The SDK client's own transport sends it (adding the bearer
// token) so the call is authenticated like every other.
func newFileRequest(
	ctx context.Context,
	client *sdk.APIClient,
	method string,
	imageID int64,
	name string,
	body io.Reader,
) (*http.Request, error) {
	cfg := client.GetConfig()
	if cfg == nil || len(cfg.Servers) == 0 || cfg.HTTPClient == nil {
		return nil, errors.New("client has no server URL or HTTP client")
	}

	endpoint := strings.TrimRight(cfg.Servers[0].URL, "/") +
		fmt.Sprintf("/api/virtual-images/%d/files?filename=%s", imageID, url.QueryEscape(name))

	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("building %s request for %q: %w", method, name, err)
	}

	req.Header.Set("Accept", "application/json")

	for k, v := range cfg.DefaultHeader {
		req.Header.Set(k, v)
	}

	return req, nil
}

// envelope is the success/msg pair every Morpheus write answers with.
type envelope struct {
	Success *bool  `json:"success"`
	Msg     string `json:"msg"`
}

// maxEnvelope bounds the JSON a file endpoint answers with.
const maxEnvelope = 64 << 10

func readEnvelope(resp *http.Response) ([]byte, envelope, error) {
	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxEnvelope))
	if err != nil {
		return nil, envelope{}, fmt.Errorf("reading response: %w", err)
	}

	var env envelope
	_ = json.Unmarshal(raw, &env)

	return raw, env, nil
}

func envelopeMessage(env envelope, raw []byte) string {
	if env.Msg != "" {
		return env.Msg
	}

	return strings.TrimSpace(string(raw))
}

// localFile describes the file to upload.
type localFile struct {
	Path   string
	Size   int64
	Sha256 string
}

// inspectLocalFile stats and hashes the file at path. When wantSHA is set the
// hash must match, so a wrong file is caught before anything is sent.
func inspectLocalFile(path, wantSHA string) (localFile, error) {
	info, err := os.Stat(path)
	if err != nil {
		return localFile{}, fmt.Errorf("could not read image file %q: %w", path, err)
	}

	if info.IsDir() {
		return localFile{}, fmt.Errorf("image file %q is a directory", path)
	}

	if info.Size() == 0 {
		return localFile{}, fmt.Errorf("image file %q is empty", path)
	}

	f, err := os.Open(path) //nolint:gosec // user-supplied path is intended
	if err != nil {
		return localFile{}, fmt.Errorf("could not read image file %q: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return localFile{}, fmt.Errorf("hashing image file %q: %w", path, err)
	}

	sum := hex.EncodeToString(h.Sum(nil))
	if wantSHA != "" && !strings.EqualFold(wantSHA, sum) {
		return localFile{}, fmt.Errorf(
			"image file %q checksum mismatch: expected %s, got %s", path, strings.ToLower(wantSHA), sum)
	}

	return localFile{Path: path, Size: info.Size(), Sha256: sum}, nil
}

// uploadFile streams the local file into the image record as name with
// POST /api/virtual-images/{id}/files?filename=<name>. The body is the open
// file with its length declared, so it is streamed rather than buffered, and
// GetBody reopens the file so the transport's re-authentication retry can
// resend it. The provider's transport speaks HTTP/1.1, which is also what a
// large upload wants: HTTP/2 flow control over a slow link throttles it.
func uploadFile(
	ctx context.Context,
	client *sdk.APIClient,
	imageID int64,
	name string,
	file localFile,
) error {
	open := func() (io.ReadCloser, error) {
		return os.Open(file.Path) //nolint:gosec // user-supplied path is intended
	}

	body, err := open()
	if err != nil {
		return fmt.Errorf("could not read image file %q: %w", file.Path, err)
	}

	req, err := newFileRequest(ctx, client, http.MethodPost, imageID, name, body)
	if err != nil {
		_ = body.Close()

		return err
	}

	req.ContentLength = file.Size
	req.GetBody = open
	req.Header.Set("Content-Type", "application/octet-stream")

	// The SDK's client carries a fixed timeout sized for API calls; the
	// upload is bounded by the resource's create timeout through ctx instead.
	httpClient := &http.Client{Transport: client.GetConfig().HTTPClient.Transport}

	tflog.Info(ctx, fmt.Sprintf("uploading %s (%d bytes) to virtual image %d", name, file.Size, imageID))

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("uploading %q to virtual image %d: %w", name, imageID, err)
	}
	defer resp.Body.Close()

	raw, env, err := readEnvelope(resp)
	if err != nil {
		return fmt.Errorf("uploading %q to virtual image %d: %w", name, imageID, err)
	}

	if resp.StatusCode >= http.StatusBadRequest || (env.Success != nil && !*env.Success) {
		return fmt.Errorf("uploading %q to virtual image %d failed (HTTP %d): %s",
			name, imageID, resp.StatusCode, envelopeMessage(env, raw))
	}

	return nil
}

// awaitListed polls the record until it lists name at the expected size:
// storage can take a moment to report a file just written.
func awaitListed(
	ctx context.Context,
	client *sdk.APIClient,
	imageID int64,
	name string,
	size int64,
	timeout time.Duration,
) ([]fileEntry, error) {
	deadline := time.Now().Add(timeout)

	for {
		files, _, err := listFiles(ctx, client, imageID)
		if err != nil {
			return nil, err
		}

		if entry, ok := findFile(files, name); ok && entry.Size == size {
			return files, nil
		}

		if time.Now().After(deadline) {
			return files, fmt.Errorf(
				"the upload reported success but virtual image %d does not list %q at %d bytes (files: %s)",
				imageID, name, size, describeFiles(files))
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
}

func describeFiles(files []fileEntry) string {
	if len(files) == 0 {
		return "none"
	}

	parts := make([]string, 0, len(files))
	for _, f := range files {
		parts = append(parts, fmt.Sprintf("%s=%d", f.Name, f.Size))
	}

	return strings.Join(parts, ", ")
}

// removeOthers deletes every file of the record except keep.
func removeOthers(ctx context.Context, client *sdk.APIClient, imageID int64, files []fileEntry, keep string) error {
	var errs []error

	for _, f := range files {
		if f.Name == keep {
			continue
		}

		tflog.Info(ctx, fmt.Sprintf("removing %s (%d bytes) from virtual image %d", f.Name, f.Size, imageID))

		if err := removeFile(ctx, client, imageID, f.Name); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
