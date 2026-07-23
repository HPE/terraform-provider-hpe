// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package sweep

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// SearchHit holds the minimal fields the sweepers need from a global Search API
// result.
type SearchHit struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// searchPageSize is the number of results requested per Search API page.
const searchPageSize int64 = 100

// SearchHits runs the global Search API for the given phrase and returns all
// matching hits, decoding only the fields in SearchHit.
//
// It reads the raw response body instead of relying on the SDK's typed decode.
// The Search API (OpenSearch-backed) returns dateCreated timestamps with a
// non-RFC3339 timezone offset (for example "2026-05-12T14:36:16+0000") that the
// generated SDK cannot parse; decoding a minimal local struct that omits
// dateCreated sidesteps the problem without any SDK change. The SDK re-buffers
// the response body before its (failing) decode, so hresp.Body is still
// readable here.
//
// The most recent *http.Response is returned so callers can honor list-status
// handling (for example treating 403/404 as "nothing to sweep").
func SearchHits(
	ctx context.Context,
	client *sdk.APIClient,
	phrase string,
) ([]SearchHit, *http.Response, error) {
	var (
		offset int64
		result []SearchHit
	)

	for {
		// The typed result and error are intentionally ignored: the SDK's date
		// decode fails on this endpoint. A missing response or non-200 status is
		// still surfaced below, with hresp preserved so the caller can format
		// the error once (errfmt.ErrMsg consumes the body, so pre-formatting it
		// here would leave a caller's second call with an empty body).
		_, hresp, err := client.SearchAPI.Search(ctx).
			Phrase(phrase).
			Max(searchPageSize).
			Offset(offset).
			Execute()
		if hresp == nil || hresp.StatusCode != http.StatusOK {
			return nil, hresp, fmt.Errorf("search for %q failed: %w", phrase, err)
		}

		body, readErr := io.ReadAll(hresp.Body)
		_ = hresp.Body.Close()
		if readErr != nil {
			return nil, hresp, fmt.Errorf("failed to read search response: %w", readErr)
		}

		hits, total, err := parseSearchPage(body)
		if err != nil {
			return nil, hresp, err
		}

		result = append(result, hits...)

		offset += int64(len(hits))
		if len(hits) == 0 || offset >= total {
			return result, hresp, nil
		}
	}
}

// parseSearchPage decodes a single Search API response page into the minimal
// fields the sweepers need, plus the total hit count for pagination. It
// deliberately ignores dateCreated (which the Search API returns in a
// non-RFC3339 format) and every other field, so the problematic timestamp is
// never parsed.
func parseSearchPage(body []byte) (hits []SearchHit, total int64, err error) {
	var page struct {
		Hits []SearchHit `json:"hits"`
		Meta struct {
			Total int64 `json:"total"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(body, &page); err != nil {
		return nil, 0, fmt.Errorf("failed to parse search response: %w", err)
	}

	return page.Hits, page.Meta.Total, nil
}
