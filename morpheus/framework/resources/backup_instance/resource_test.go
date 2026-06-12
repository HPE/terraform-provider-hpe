package backupinstance_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	backupinstance "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backup_instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backup_job"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusBackupInstanceResourceExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Backup) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	dependencyConfig, err := backup_job.RenderBackupJobConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig, err := backupinstance.RenderBackupInstanceConfig(t, map[string]string{
		"Name":              name,
		"InstanceId":        "26",
		"JobId":             "hpe_morpheus_backup_job.example.id",
		"StorageProviderId": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_backup_instance.example"
	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "instance_id", "26"),
		resource.TestCheckResourceAttrPair(
			resourceName, "job_id",
			"hpe_morpheus_backup_job.example", "id",
		),
		resource.TestCheckResourceAttr(resourceName, "storage_provider_id", "1"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttrSet(resourceName, "backup_type_code"),
		resource.TestCheckResourceAttrSet(resourceName, "container_id"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dependencyConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config:            providerConfig + dependencyConfig + resourceConfig,
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
			// Do this to destroy the backup_instance (backup_job destroyed as side-effect of backup_instance destroy)...
			{
				Config:             providerConfig + dependencyConfig,
				ExpectNonEmptyPlan: true,
			},
			// ...followed by a call to terraform apply -refresh-only
			{
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccMorpheusBackupInstanceResourceUpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Backup) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	dependencyConfig, err := backup_job.RenderBackupJobConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	})
	if err != nil {
		t.Fatal(err)
	}

	createConfig, err := backupinstance.RenderBackupInstanceConfig(t, map[string]string{
		"Name":              name,
		"InstanceId":        "26",
		"JobId":             "hpe_morpheus_backup_job.example.id",
		"StorageProviderId": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	updatedName := name + "updated"
	updateConfig := `
resource "hpe_morpheus_backup_instance" "example" {
  name                = "` + updatedName + `"
  instance_id         = 26
  job_id              = hpe_morpheus_backup_job.example.id
  storage_provider_id = 2
  enabled             = false
}
`

	resourceName := "hpe_morpheus_backup_instance.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "instance_id", "26"),
		resource.TestCheckResourceAttrPair(
			resourceName, "job_id",
			"hpe_morpheus_backup_job.example", "id",
		),
		resource.TestCheckResourceAttr(resourceName, "storage_provider_id", "1"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
		resource.TestCheckResourceAttr(resourceName, "instance_id", "26"),
		resource.TestCheckResourceAttrPair(
			resourceName, "job_id",
			"hpe_morpheus_backup_job.example", "id",
		),
		resource.TestCheckResourceAttr(resourceName, "storage_provider_id", "2"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
	)
	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dependencyConfig + createConfig,
				Check:  createChecks,
			},
			{
				Config:           providerConfig + dependencyConfig + updateConfig,
				Check:            updateChecks,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{
				Config:             providerConfig + dependencyConfig + updateConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			// Do this to destroy the backup_instance (backup_job destroyed as side-effect of backup_instance destroy)...
			{
				Config:             providerConfig + dependencyConfig,
				ExpectNonEmptyPlan: true,
			},
			// ...followed by a call to terraform apply -refresh-only
			{
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
