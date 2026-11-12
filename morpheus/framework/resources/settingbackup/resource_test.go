// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingbackup_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/settingbackup"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)

	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// The tests here change settings that apply to the whole appliance. None runs
// in parallel: another test restoring its own snapshot of the same settings
// would overwrite the values under test. Each takes a snapshot of the settings
// first and registers a cleanup that restores them -- all but a retention the
// appliance did not have, which the API cannot be told to forget -- and every
// one requires the master_tenant capability as an explicit opt-in: the
// settings are shared by every tenant of the appliance, so a run must be
// intended for it.

const backupResourceAddr = "hpe_morpheus_setting_backup.example"

// snapshotAndRestoreBackupSettings captures the backup settings before a test
// changes them and registers a cleanup that puts them back, so a failed or
// interrupted run never leaves the shared appliance pointed at the test's
// values. It returns the snapshot.
func snapshotAndRestoreBackupSettings(t *testing.T) *sdk.ListBackupSettings200ResponseBackupSettings {
	t.Helper()

	ctx := context.Background()

	client, err := testhelpers.NewClientForServer(ctx, "")
	if err != nil {
		t.Fatalf("creating the API client: %v", err)
	}

	result, _, err := client.BackupSettingsAPI.ListBackupSettings(ctx).Execute()
	if err != nil {
		t.Fatalf("reading the backup settings: %v", err)
	}

	if result == nil || result.BackupSettings == nil {
		t.Fatal("no backupSettings in the response")
	}

	settings := result.BackupSettings

	restore := backupRestoreBody(t, settings)

	t.Cleanup(func() {
		putBackupSettings(t, client, restore)
	})

	return settings
}

// putBackupSettings sends a raw settings body, reporting a failure or a
// refusal (HTTP 200 with success false) as a test error.
func putBackupSettings(t *testing.T, client *sdk.APIClient, settings map[string]any) {
	t.Helper()

	body := sdk.UpdateBackupSettingsRequest{
		BackupSettings: &sdk.UpdateBackupSettingsRequestBackupSettings{
			AdditionalProperties: settings,
		},
	}

	result, _, err := client.BackupSettingsAPI.UpdateBackupSettings(context.Background()).
		UpdateBackupSettingsRequest(body).
		Execute()
	if err != nil {
		t.Errorf("updating the backup settings: %v", err)

		return
	}

	if result != nil && result.Success != nil && !*result.Success {
		t.Errorf("the appliance refused the settings: %v", result.AdditionalProperties)
	}
}

// backupRestoreBody converts settings read from the API into the PUT body that
// puts them back. The three flags are always replayed. The retention is
// replayed only when the appliance reports one: the API ignores a null
// retentionCount, so a null retention cannot be restored. The default schedule
// and storage bucket are replayed by id when the appliance has one, and
// cleared otherwise (clearDefaultSchedule / clearDefaultStorageBucket).
func backupRestoreBody(t *testing.T, s *sdk.ListBackupSettings200ResponseBackupSettings) map[string]any {
	t.Helper()

	flag := func(v *bool) bool { return v != nil && *v }

	m := map[string]any{
		"backupsEnabled":  flag(s.BackupsEnabled),
		"createBackups":   flag(s.CreateBackups),
		"backupAppliance": flag(s.BackupAppliance),
	}

	if v := s.RetentionCount.Get(); v != nil {
		m["retentionCount"] = *v
	} else {
		t.Logf("the appliance reports no retention count; the API ignores a null retentionCount, " +
			"so the value this test sets will remain after it")
	}

	if sched := s.DefaultSchedule.Get(); sched != nil && sched.Id != nil {
		m["defaultSchedule"] = map[string]any{"id": *sched.Id}
	} else {
		m["clearDefaultSchedule"] = true
	}

	if bucket := s.DefaultStorageBucket.Get(); bucket != nil && bucket.Id != nil {
		m["defaultStorageBucket"] = map[string]any{"id": *bucket.Id}
	} else {
		m["clearDefaultStorageBucket"] = true
	}

	return m
}

// TestAccMorpheusSettingBackupExampleOk creates the dependencies the settings
// point at -- an execute schedule and a local storage bucket -- and the
// settings themselves, then checks the round-trip, an empty plan, an import,
// and an update. backup_appliance is left false: true would make Morpheus
// create a real appliance database backup job. The local bucket needs the
// master tenant.
func TestAccMorpheusSettingBackupExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings, capabilities.MasterTenant)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	snapshotAndRestoreBackupSettings(t)

	name := acctest.RandomWithPrefix(t.Name())
	providerConfig := testhelpers.ProviderBlock()

	deps := fmt.Sprintf(`
resource "hpe_morpheus_execute_schedule" "%[1]s" {
  name      = "%[1]s"
  schedule  = "0 2 * * *"
  time_zone = "UTC"
  enabled   = true
}

resource "hpe_morpheus_storage_bucket" "%[1]s" {
  name                  = "%[1]s"
  provider_type         = "local"
  base_path             = "/var/opt/morpheus/morpheus-ui/images"
  default_backup_target = true
}
`, name)

	// The example template renders the two ids unquoted, so the references are
	// bare traversals, not "${...}" interpolations.
	scheduleRef := fmt.Sprintf("hpe_morpheus_execute_schedule.%s.id", name)
	bucketRef := fmt.Sprintf("hpe_morpheus_storage_bucket.%s.id", name)

	created, err := settingbackup.RenderSettingBackupConfig(t, map[string]string{
		"DefaultBackupScheduleId":      scheduleRef,
		"DefaultBackupStorageBucketId": bucketRef,
	})
	if err != nil {
		t.Fatalf("rendering the example: %v", err)
	}

	updated, err := settingbackup.RenderSettingBackupConfig(t, map[string]string{
		"DefaultBackupScheduleId":      scheduleRef,
		"DefaultBackupStorageBucketId": bucketRef,
		"RetentionCount":               "7",
		"ScheduledBackups":             "false",
	})
	if err != nil {
		t.Fatalf("rendering the example: %v", err)
	}

	factories := testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + deps + created,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(backupResourceAddr, "id", "1"),
					resource.TestCheckResourceAttr(backupResourceAddr, "scheduled_backups", "true"),
					resource.TestCheckResourceAttr(backupResourceAddr, "create_backups", "true"),
					resource.TestCheckResourceAttr(backupResourceAddr, "backup_appliance", "false"),
					resource.TestCheckResourceAttr(backupResourceAddr, "retention_count", "21"),
					resource.TestCheckResourceAttr(backupResourceAddr, "retention_days", "21"),
					resource.TestCheckResourceAttrPair(
						backupResourceAddr, "default_backup_schedule_id",
						"hpe_morpheus_execute_schedule."+name, "id"),
					resource.TestCheckResourceAttrPair(
						backupResourceAddr, "default_backup_storage_bucket_id",
						"hpe_morpheus_storage_bucket."+name, "id"),
				),
			},
			{
				Config:   providerConfig + deps + created,
				PlanOnly: true,
			},
			{
				ResourceName:      backupResourceAddr,
				ImportState:       true,
				ImportStateId:     "1",
				ImportStateVerify: true,
			},
			{
				Config: providerConfig + deps + updated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(backupResourceAddr, "scheduled_backups", "false"),
					resource.TestCheckResourceAttr(backupResourceAddr, "retention_count", "7"),
					resource.TestCheckResourceAttr(backupResourceAddr, "retention_days", "7"),
				),
			},
		},
	})
}

// TestAccMorpheusSettingBackup_retentionDays sets only the deprecated
// retention_days, checks it reaches the one retentionCount setting and reads
// back into retention_count, that switching to retention_count with the same
// value is an empty plan (the two are one setting), that setting both is
// refused at plan time, and that 0 is refused on either attribute.
func TestAccMorpheusSettingBackup_retentionDays(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings, capabilities.MasterTenant)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	snapshotAndRestoreBackupSettings(t)

	providerConfig := testhelpers.ProviderBlock()
	addr := "hpe_morpheus_setting_backup.retention"

	daysConfig := providerConfig + `
resource "hpe_morpheus_setting_backup" "retention" {
  scheduled_backups = true
  retention_days    = 14
}
`

	countConfig := providerConfig + `
resource "hpe_morpheus_setting_backup" "retention" {
  scheduled_backups = true
  retention_count   = 14
}
`

	bothConfig := providerConfig + `
resource "hpe_morpheus_setting_backup" "retention" {
  scheduled_backups = true
  retention_count   = 14
  retention_days    = 14
}
`

	// Morpheus stores a retention of 0 but treats it as unset where it reads
	// it, so 0 is refused -- on the alias too, or it would bypass the check.
	zeroCountConfig := providerConfig + `
resource "hpe_morpheus_setting_backup" "retention" {
  retention_count = 0
}
`

	zeroDaysConfig := providerConfig + `
resource "hpe_morpheus_setting_backup" "retention" {
  retention_days = 0
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: daysConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "retention_days", "14"),
					resource.TestCheckResourceAttr(addr, "retention_count", "14"),
				),
			},
			{
				Config:   countConfig,
				PlanOnly: true,
			},
			{
				Config:      bothConfig,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
			{
				Config:      zeroCountConfig,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)retention_count.*at least 1`),
			},
			{
				Config:      zeroDaysConfig,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)retention_days.*at least 1`),
			},
			{
				// Leave a valid configuration in place: the post-test destroy
				// plans with the last step's configuration, and the one above
				// is refused before it gets that far.
				Config:   countConfig,
				PlanOnly: true,
			},
		},
	})
}

// TestAccMorpheusSettingBackup_updateEnvelope asserts that the appliance's
// refusal of a default schedule that does not exist surfaces as the provider's
// error rather than a silent success.
func TestAccMorpheusSettingBackup_updateEnvelope(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings, capabilities.MasterTenant)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	snapshotAndRestoreBackupSettings(t)

	config := testhelpers.ProviderBlock() + `
resource "hpe_morpheus_setting_backup" "envelope" {
  default_backup_schedule_id = 999999999
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			// The expected error is matched on two words rather than the whole
			// sentence: Terraform wraps the diagnostic at word boundaries, and
			// (?s) lets the match span the break. "refused" is updateSettings'
			// wording; "chedule" ties the refusal to the field Morpheus named.
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`(?s)refused.*chedule`),
			},
		},
	})
}

// releasedSDKv2Version is the provider release the upgrade test starts from:
// the latest published release, every one of which has the SDKv2 resource.
const releasedSDKv2Version = "2.1.0"

const backupUpgradeAddr = "hpe_morpheus_setting_backup.upgrade"

// settingBackupUpgradeConfig sets the flags and the retention and no ids: the
// released SDKv2 resource omits the ids when 0.
const settingBackupUpgradeConfig = `
resource "hpe_morpheus_setting_backup" "upgrade" {
  scheduled_backups = true
  create_backups    = true
  backup_appliance  = false
  retention_days    = 9
}
`

// TestAccMorpheusSettingBackup_upgradeFromSDKv2 is the regression test for the
// move from the SDKv2 resource: the state a released provider writes upgrades
// in place to this build, and its plan is empty.
func TestAccMorpheusSettingBackup_upgradeFromSDKv2(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Settings, capabilities.MasterTenant)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	snapshotAndRestoreBackupSettings(t)

	config := testhelpers.ProviderBlock() + settingBackupUpgradeConfig
	thisBuild := testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				// The released provider, whose resource is the SDKv2 one.
				ExternalProviders: map[string]resource.ExternalProvider{
					"hpe": {Source: "HPE/hpe", VersionConstraint: releasedSDKv2Version},
				},
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(backupUpgradeAddr, "id", "1"),
					resource.TestCheckResourceAttr(backupUpgradeAddr, "retention_days", "9"),
				),
			},
			{
				// This build, same configuration and state: the state is
				// upgraded and the refresh reconciles it, so the plan is a noop.
				ProtoV6ProviderFactories: thisBuild,
				Config:                   config,
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(backupUpgradeAddr, plancheck.ResourceActionNoop),
				}},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(backupUpgradeAddr, "id", "1"),
					resource.TestCheckResourceAttr(backupUpgradeAddr, "retention_count", "9"),
					resource.TestCheckResourceAttr(backupUpgradeAddr, "retention_days", "9"),
				),
			},
			{
				ProtoV6ProviderFactories: thisBuild,
				Config:                   config,
				PlanOnly:                 true,
			},
		},
	})
}
