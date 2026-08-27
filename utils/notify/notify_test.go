// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package notify_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/HPE/terraform-provider-hpe/utils/notify"

	"github.com/cenkalti/backoff/v5"
	version "github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	mockResp = `
{
  "meta": {
    "limit": 15,
    "current_offset": 0
  },
  "providers": [
    {
      "id": "HPE/hpe/1.2.0",
      "owner": "",
      "namespace": "HPE",
      "name": "hpe",
      "alias": null,
      "version": "1.2.0",
      "tag": "v1.2.0",
      "description": "",
      "source": "https://github.com/HPE/terraform-provider-hpe",
      "published_at": "2026-03-30T19:31:32Z",
      "downloads": 3908,
      "tier": "community",
      "logo_url": "https://avatars3.githubusercontent.com/HPE"
    }
  ]
}
`
)

// Test the Dial by running it a few times in parallel make sure the timeout we set is OK.
func TestTryDial(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.NoError(t, notify.TryDial())
		}()
	}
	wg.Wait()
}

// Test against the real registry API.
func TestGetProviderVersion(t *testing.T) {
	t.Parallel()

	retry := func() (*version.Version, error) {
		return notify.GetProviderVersion(notify.RegistryUrl)
	}

	ctx := t.Context()
	remoteVer, err := backoff.Retry(
		ctx,
		retry,
		backoff.WithMaxTries(3),
	)

	assert.NoError(t, err)
	require.NotNil(t, remoteVer)
	assert.NotEmpty(t, remoteVer.String())
}

// Test against some static mock data.
// This is really just to test the core logic after a successful JSON unmarshal.
func TestGetProviderVersionMock(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, mockResp)
		}))
	defer server.Close()

	ver, err := notify.GetProviderVersion(server.URL)
	assert.NoError(t, err)
	require.NotNil(t, ver)
	assert.Equal(t, "1.2.0", ver.String())
}
