// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storagebucket_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusStorageBucketResourceWriteOnlyOk: changes to secret_key (one
// of the two write-only keys) are detected by digest (MORPH-18005); the
// deprecated marker still works; both import forms leave a clean plan. Needs
// valid AWS credentials: the API validates them on create.
func TestAccMorpheusStorageBucketResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.AWS)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix(t.Name())

	c := testhelpers.WriteOnlyCase{
		Address:   "hpe_morpheus_storage_bucket.wo",
		Attribute: "secret_key",
		Marker:    "secret_key_version",
		APIPath:   "/api/storage-buckets",
		APIKey:    "storageBucket.config.secretKey",
		Config: func(value, extra string, unrelated bool) string {
			retention := "7"
			if unrelated {
				retention = "14"
			}

			return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_storage_bucket" "wo" {
  name           = "` + name + `"
  provider_type  = "s3"
  bucket_name    = "` + name + `"
  access_key     = "AKIAIOSFODNN7EXAMPLE"
  secret_key     = "` + value + `"` + extra + `
  retention_days = ` + retention + `
}
`
		},
	}

	rec := testhelpers.RecordRequests(t)

	// AWS's documented example secret key, varied by its last character.
	const k = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

	steps := testhelpers.WriteOnlyMatrixSteps(c, rec, k+"1", k+"2", k+"3")
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, k+"3", "access_key", "access_key_version")...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
