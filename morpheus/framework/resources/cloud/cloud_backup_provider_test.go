// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloud_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
)

// internalBackupProviderID looks up the id of the appliance's internal backup
// provider via the backupProviders option source. The pruned SDK has no
// operation for it, so this is a raw GET through the SDK client's configured
// HTTP client, base URL and auth header.
func internalBackupProviderID(t *testing.T) int64 {
	t.Helper()

	ctx := context.Background()

	client, err := testhelpers.NewClientForServer(ctx, "")
	if err != nil {
		t.Fatalf("client: %v", err)
	}

	cfg := client.GetConfig()
	base := cfg.Servers[0].URL

	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet, base+"/api/options/backupProviders?max=200", nil,
	)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	for k, v := range cfg.DefaultHeader {
		req.Header.Set(k, v)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		t.Fatalf("GET options/backupProviders: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET options/backupProviders: HTTP %d: %s", resp.StatusCode, body)
	}

	var payload struct {
		Data []struct {
			Value int64  `json:"value"`
			Name  string `json:"name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("decode options/backupProviders: %v", err)
	}
	if len(payload.Data) == 0 {
		t.Skip("no backup providers on the appliance")
	}

	for _, p := range payload.Data {
		if p.Name == "morpheus" {
			return p.Value
		}
	}

	return payload.Data[0].Value
}

// TestAccMorpheusCloudBackupProviderIds verifies that backup_provider_ids is
// sent on create and read back (MORPH-18721). Removing it from the
// configuration afterwards is a no-op by design: the attribute is computed and
// not clearable, so the prior value is carried forward — which is why the
// server-default case needs its own fresh create, below.
func TestAccMorpheusCloudBackupProviderIds(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerID := internalBackupProviderID(t)
	visibility := testhelpers.TenantVisibility(t)
	name := acctest.RandomWithPrefix(t.Name())

	const addr = "hpe_morpheus_cloud.bp"

	cloudConfig := func(backupLine string) string {
		return testhelpers.ProviderBlock() + testhelpers.WhoamiBlock() + `
resource "hpe_morpheus_cloud" "bp" {
  name       = "` + name + `"
  code       = "` + name + `"
  tenant_id  = ` + testhelpers.WhoamiTenantIDRef + `
  group_id   = 1
  enabled    = true
  visibility = "` + visibility + `"
` + backupLine + `
  config_hvm = {
    certificate_provider = "internal"
  }
}
`
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Explicit: the internal provider id is sent and read back.
				Config: cloudConfig(
					"  backup_provider_ids = [" + strconv.FormatInt(providerID, 10) + "]\n",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "backup_provider_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						addr, "backup_provider_ids.*",
						strconv.FormatInt(providerID, 10),
					),
				),
			},
			{
				// Removed from the configuration: computed and not clearable, so
				// the previously applied value is kept, with no diff.
				Config: cloudConfig(""),
				Check: resource.TestCheckTypeSetElemAttr(
					addr, "backup_provider_ids.*",
					strconv.FormatInt(providerID, 10),
				),
			},
		},
	})
}

// TestAccMorpheusCloudBackupProviderIdsServerDefault creates a cloud WITHOUT
// backup_provider_ids and checks that the server's default — Morpheus assigns
// the internal backup provider when none is given — is read into state. The
// cloud API only returns backupProviders from Morpheus 9.0.2, so below that
// the provider preserves the (null) configured value and this test does not
// apply.
func TestAccMorpheusCloudBackupProviderIdsServerDefault(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)
	testhelpers.SkipUnlessApplianceVersionAtLeast(context.Background(), t, constants.CloudBackupProvidersMinVersion)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerID := internalBackupProviderID(t)
	visibility := testhelpers.TenantVisibility(t)
	name := acctest.RandomWithPrefix(t.Name())

	const addr = "hpe_morpheus_cloud.bpd"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + testhelpers.WhoamiBlock() + `
resource "hpe_morpheus_cloud" "bpd" {
  name       = "` + name + `"
  code       = "` + name + `"
  tenant_id  = ` + testhelpers.WhoamiTenantIDRef + `
  group_id   = 1
  enabled    = true
  visibility = "` + visibility + `"
  config_hvm = {
    certificate_provider = "internal"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "backup_provider_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						addr, "backup_provider_ids.*",
						strconv.FormatInt(providerID, 10),
					),
				),
			},
		},
	})
}
