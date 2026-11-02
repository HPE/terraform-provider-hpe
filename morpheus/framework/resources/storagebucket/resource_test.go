// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagebucket_test

import (
	"context"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/storagebucket"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusStorageBucketResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	// AWS, not Alletra: Alletra is block storage and has no bearing on this
	// resource. The example provisions an s3 bucket, and the API validates it
	// by writing, reading back and deleting a test object before saving, so
	// the test needs real AWS credentials on the appliance to pass.
	capabilities.MustHaveOrSkip(t, capabilities.AWS)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := storagebucket.RenderStorageBucketConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("hpe_morpheus_storage_bucket.example", "name", name),
		resource.TestCheckResourceAttr("hpe_morpheus_storage_bucket.example", "provider_type", "s3"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"write_only_revision"},
				ResourceName:            "hpe_morpheus_storage_bucket.example",
			},
		},
	})
}

func TestAccMorpheusStorageBucketResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	// See TestAccMorpheusStorageBucketResourceExampleOk for why this is AWS.
	capabilities.MustHaveOrSkip(t, capabilities.AWS)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	createConfig, err := storagebucket.RenderStorageBucketConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig, err := storagebucket.RenderStorageBucketConfig(t, map[string]string{
		"Name": updatedName,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_storage_bucket.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "provider_type", "s3"),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
		resource.TestCheckResourceAttr(resourceName, "provider_type", "s3"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check:  createChecks,
			},
			{
				Config:           providerConfig + updateConfig,
				Check:            updateChecks,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{
				Config:             providerConfig + updateConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

// TestAccMorpheusStorageBucketResourceLocalOk: a bucket on the appliance's
// own disk -- the type a fresh appliance needs before any image can be fetched
// by URL. Create with a base_path under the appliance's data directory (the
// API writes a test file into it, creating the directory), read it back from
// the response's config, change the path in place, plan clean, and import.
// bucket_name is not configured: Morpheus names a local bucket "." itself and
// the attribute takes that value. default_virtual_image_target is configured
// false so the test never takes the default image store from the appliance's
// own bucket. Local storage is master-tenant-only (Morpheus refuses it for a
// subtenant).
func TestAccMorpheusStorageBucketResourceLocalOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.MasterTenant)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	basePath := "/var/opt/morpheus/morpheus-ui/" + strings.ToLower(name)
	movedPath := basePath + "-moved"

	createConfig, err := storagebucket.RenderStorageBucketLocalConfig(t, map[string]string{
		"ResourceName":              "example",
		"Name":                      name,
		"BasePath":                  basePath,
		"DefaultVirtualImageTarget": "false",
	})
	if err != nil {
		t.Fatal(err)
	}

	moveConfig, err := storagebucket.RenderStorageBucketLocalConfig(t, map[string]string{
		"ResourceName":              "example",
		"Name":                      name,
		"BasePath":                  movedPath,
		"DefaultVirtualImageTarget": "false",
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_storage_bucket.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "provider_type", "local"),
					resource.TestCheckResourceAttr(resourceName, "base_path", basePath),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", "."),
					resource.TestCheckResourceAttr(resourceName, "default_virtual_image_target", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_backup_target", "false"),
				),
			},
			{
				Config: providerConfig + moveConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_path", movedPath),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", "."),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				Config:             providerConfig + moveConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"write_only_revision"},
				ResourceName:            resourceName,
			},
		},
	})
}

// TestAccMorpheusStorageBucketResourceLocalDefaultImageTargetKept: a bucket
// that is the appliance's default image store stays so when
// default_virtual_image_target is not configured -- a bucket made the default
// outside Terraform, or one managed from before the attribute existed. The
// bucket is created as the default; the attribute is then dropped from the
// configuration, which must plan nothing and leave the bucket the default; and
// finally it is set false explicitly, which is how the default is taken away.
// A schema default of false would fail the middle step by planning true ->
// false, and the apply would unset the appliance's default image store.
//
// The test holds the appliance's default image store for its duration and puts
// the previous default back afterwards, so besides MasterTenant (local storage)
// it is gated on Settings, like the other tests that change appliance-wide
// state.
func TestAccMorpheusStorageBucketResourceLocalDefaultImageTargetKept(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.MasterTenant, capabilities.Settings)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	ctx := context.Background()
	restoreDefaultImageTarget(ctx, t)

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	basePath := "/var/opt/morpheus/morpheus-ui/" + strings.ToLower(name)

	render := func(defaultVirtualImageTarget string) string {
		cfg, err := storagebucket.RenderStorageBucketLocalConfig(t, map[string]string{
			"ResourceName":              "example",
			"Name":                      name,
			"BasePath":                  basePath,
			"DefaultVirtualImageTarget": defaultVirtualImageTarget,
		})
		if err != nil {
			t.Fatal(err)
		}

		return providerConfig + cfg
	}

	resourceName := "hpe_morpheus_storage_bucket.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: render("true"),
				Check:  resource.TestCheckResourceAttr(resourceName, "default_virtual_image_target", "true"),
			},
			{
				// Not configured: nothing to do, and the bucket is still the default.
				Config: render(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.TestCheckResourceAttr(resourceName, "default_virtual_image_target", "true"),
			},
			{
				Config: render("false"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckResourceAttr(resourceName, "default_virtual_image_target", "false"),
			},
		},
	})
}

// restoreDefaultImageTarget records which bucket is the appliance's default
// image store now and, when the test is over, makes it the default again.
// With no default to begin with there is nothing to put back.
func restoreDefaultImageTarget(ctx context.Context, t *testing.T) {
	t.Helper()

	client, err := testhelpers.NewClientForServer(ctx, "")
	if err != nil {
		t.Fatalf("API client: %v", err)
	}

	list, httpResp, err := client.StorageAPI.ListStorageBuckets(ctx).Max(1000).Execute()
	if err != nil {
		t.Fatalf("listing storage buckets: %v (%v)", err, httpResp)
	}

	var previous *int64
	for _, b := range list.StorageBuckets {
		if b.Id != nil && b.DefaultVirtualImageTarget != nil && *b.DefaultVirtualImageTarget {
			previous = b.Id

			break
		}
	}

	if previous == nil {
		return
	}

	t.Cleanup(func() {
		on := true
		body := sdk.UpdateStorageBucketsRequest{
			StorageBucket: sdk.UpdateStorageBucketsRequestStorageBucket{DefaultVirtualImageTarget: &on},
		}

		_, httpResp, err := client.StorageAPI.UpdateStorageBuckets(ctx, *previous).
			UpdateStorageBucketsRequest(body).Execute()
		if err != nil {
			t.Errorf("putting the default image store back on bucket %d: %v (%v)", *previous, err, httpResp)
		}
	})
}

// TestAccMorpheusStorageBucketResourceBasePathValidation: base_path and the
// local provider type go together, and `terraform validate` says so before
// anything reaches the appliance -- a local bucket without a directory, and a
// directory on any other type, are both refused at plan time.
func TestAccMorpheusStorageBucketResourceBasePathValidation(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "hpe_morpheus_storage_bucket" "example" {
  name          = "` + name + `"
  provider_type = "local"
}
`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`base_path is required for a local storage bucket`),
			},
			{
				Config: providerConfig + `
resource "hpe_morpheus_storage_bucket" "example" {
  name          = "` + name + `"
  provider_type = "s3"
  bucket_name   = "example-bucket"
  base_path     = "/var/opt/morpheus/morpheus-ui/images"
}
`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`base_path applies to local storage buckets only`),
			},
		},
	})
}
