// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plugin_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

const (
	skeletonCode = "hpe-terraform-skeleton-plugin"
	skeletonName = "HPE Terraform Skeleton"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func mustHavePlugins(t *testing.T) {
	t.Helper()
	capabilities.MustHaveOrSkip(t, capabilities.Plugins)
}

func absJar(t *testing.T, version string) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join(
		"testdata", "skeleton-plugin", fmt.Sprintf("%s-%s.jar", skeletonCode, version)))
	if err != nil {
		t.Fatalf("resolving jar path: %v", err)
	}

	return p
}

func sha256OfFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	sum := sha256.Sum256(b)

	return hex.EncodeToString(sum[:])
}

// TestAccMorpheusPluginResourceJarPathOk uploads the skeleton plugin from a
// local jar, upgrades it to a new version while disabling it, then imports.
func TestAccMorpheusPluginResourceJarPathOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	mustHavePlugins(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	jar010 := absJar(t, "0.1.0")
	jar020 := absJar(t, "0.2.0")

	providerConfig := testhelpers.ProviderBlock()
	rn := "hpe_morpheus_plugin.test"

	configV1 := fmt.Sprintf(`
resource "hpe_morpheus_plugin" "test" {
  jar_path = %q
}
`, jar010)

	configV2 := fmt.Sprintf(`
resource "hpe_morpheus_plugin" "test" {
  jar_path = %q
  enabled  = false
}
`, jar020)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + configV1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "code", skeletonCode),
					resource.TestCheckResourceAttr(rn, "version", "0.1.0"),
					resource.TestCheckResourceAttr(rn, "status", "loaded"),
					resource.TestCheckResourceAttr(rn, "enabled", "true"),
					resource.TestCheckResourceAttr(rn, "name", skeletonName),
					resource.TestCheckResourceAttrSet(rn, "id"),
					resource.TestCheckResourceAttr(rn, "jar_sha256", sha256OfFile(t, jar010)),
				),
			},
			{
				Config: providerConfig + configV2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "version", "0.2.0"),
					resource.TestCheckResourceAttr(rn, "enabled", "false"),
					resource.TestCheckResourceAttr(rn, "jar_sha256", sha256OfFile(t, jar020)),
				),
			},
			{
				ResourceName:            rn,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"jar_path", "jar_url", "jar_sha256", "timeouts"},
			},
		},
	})
}

// TestAccMorpheusPluginResourceJarSha256MismatchErr fails fast when the
// supplied jar_sha256 does not match the jar.
func TestAccMorpheusPluginResourceJarSha256MismatchErr(t *testing.T) {
	defer testhelpers.RecordResult(t)
	mustHavePlugins(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	jar010 := absJar(t, "0.1.0")
	providerConfig := testhelpers.ProviderBlock()

	config := fmt.Sprintf(`
resource "hpe_morpheus_plugin" "test" {
  jar_path   = %q
  jar_sha256 = "0000000000000000000000000000000000000000000000000000000000000000"
}
`, jar010)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + config,
				ExpectError: regexp.MustCompile("checksum mismatch"),
			},
		},
	})
}

// exchangePackage is the subset of the Morpheus Exchange packages response we
// read to discover a version's shaHash.
type exchangePackage struct {
	Data struct {
		Packages []struct {
			Versions []struct {
				Version string `json:"version"`
				ShaHash string `json:"shaHash"`
			} `json:"versions"`
		} `json:"packages"`
	} `json:"data"`
}

// TestAccMorpheusPluginResourceJarUrlOk installs a plugin by downloading it from
// the Morpheus Exchange. The expected sha is fetched live (skipped if offline).
func TestAccMorpheusPluginResourceJarUrlOk(t *testing.T) {
	defer testhelpers.RecordResult(t)
	mustHavePlugins(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	const (
		code    = "vault-plugin"
		version = "1.4.2"
	)

	resp, err := http.Get("https://share.morpheusdata.com/api/packages?code=" + code)
	if err != nil {
		t.Skipf("skipping: cannot reach Morpheus Exchange: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Skipf("skipping: Morpheus Exchange returned HTTP %d", resp.StatusCode)
	}

	var pkg exchangePackage
	if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
		t.Skipf("skipping: cannot parse Morpheus Exchange response: %v", err)
	}

	var shaHash string
	if len(pkg.Data.Packages) > 0 {
		for _, v := range pkg.Data.Packages[0].Versions {
			if v.Version == version {
				shaHash = v.ShaHash

				break
			}
		}
	}
	if shaHash == "" {
		t.Skipf("skipping: could not find shaHash for %s %s on the Exchange", code, version)
	}

	providerConfig := testhelpers.ProviderBlock()
	rn := "hpe_morpheus_plugin.test"

	config := fmt.Sprintf(`
resource "hpe_morpheus_plugin" "test" {
  jar_url    = "https://share.morpheusdata.com/api/packages/pull?code=%s&version=%s"
  jar_sha256 = %q
}
`, code, version, shaHash)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "code", "morpheus-hashicorp-vault-plugin"),
					resource.TestCheckResourceAttr(rn, "version", version),
					resource.TestCheckResourceAttr(rn, "status", "loaded"),
				),
			},
		},
	})
}
