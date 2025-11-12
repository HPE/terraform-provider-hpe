// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package setting

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

const (
	backupSettingsID = "1"
)

func ResourceSettingBackup() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus backup setting resource.",
		CreateContext: resourceSettingBackupCreate,
		ReadContext:   resourceSettingBackupRead,
		UpdateContext: resourceSettingBackupUpdate,
		DeleteContext: resourceSettingBackupDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the backup settings",
				Computed:    true,
			},
			"scheduled_backups": {
				Type:        schema.TypeBool,
				Description: "Whether automatic backups will be scheduled for provisioned instances",
				Optional:    true,
				Computed:    true,
			},
			"create_backups": {
				Type:        schema.TypeBool,
				Description: "Whether morpheus will automatically configure instances for manual or scheduled backups",
				Optional:    true,
				Computed:    true,
			},
			"backup_appliance": {
				Type:        schema.TypeBool,
				Description: "Whether a backup will be created for the Morpheus appliance database",
				Optional:    true,
				Computed:    true,
			},
			"default_backup_storage_bucket_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the storage bucket to set as the default for backups",
				Optional:    true,
				Computed:    true,
			},
			"default_backup_schedule_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the execution schedule used as the default backup schedule",
				Optional:    true,
				Computed:    true,
			},
			"retention_days": {
				Type:        schema.TypeInt,
				Description: "The number of days to retain backups",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSettingBackupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var scheduledBackups bool
	if v, ok := d.Get("scheduled_backups").(bool); ok {
		scheduledBackups = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("scheduled_backups", d.Get("scheduled_backups")))
	}

	var createBackups bool
	if v, ok := d.Get("create_backups").(bool); ok {
		createBackups = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("create_backups", d.Get("create_backups")))
	}

	var backupAppliance bool
	if v, ok := d.Get("backup_appliance").(bool); ok {
		backupAppliance = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("backup_appliance", d.Get("backup_appliance")))
	}

	var retentionDays int
	if v, ok := d.Get("retention_days").(int); ok {
		retentionDays = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("retention_days", d.Get("retention_days")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"backupSettings": map[string]any{
				"backupsEnabled":  scheduledBackups,
				"createBackups":   createBackups,
				"backupAppliance": backupAppliance,
				"retentionCount":  retentionDays,
			},
		},
	}

	var defaultStorageBucketID int
	if v, ok := d.Get("default_backup_storage_bucket_id").(int); ok {
		defaultStorageBucketID = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError("default_backup_storage_bucket_id", d.Get("default_backup_storage_bucket_id")),
		)
	}

	if defaultStorageBucketID != 0 {
		req.Body["defaultStorageBucket"] = map[string]any{
			"id": defaultStorageBucketID,
		}
	}

	var defaultBackupScheduleID int
	if v, ok := d.Get("default_backup_schedule_id").(int); ok {
		defaultBackupScheduleID = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError("default_backup_schedule_id", d.Get("default_backup_schedule_id")),
		)
	}

	if defaultBackupScheduleID != 0 {
		req.Body["defaultSchedule"] = map[string]any{
			"id": defaultBackupScheduleID,
		}
	}

	resp, err := client.UpdateBackupSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateBackupSettingsResult
	if v, ok := resp.Result.(*morpheus.UpdateBackupSettingsResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.BackupSettings == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("BackupSettings"))
	}

	d.SetId(backupSettingsID)

	diags = append(diags, resourceSettingBackupRead(ctx, d, meta)...)

	return diags
}

func resourceSettingBackupRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var resp *morpheus.Response
	var err error

	resp, err = client.GetBackupSettings(&morpheus.Request{})
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			log.Printf("Forcing recreation of resource")
			d.SetId("")

			return diags
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.GetBackupSettingsResult
	if v, ok := resp.Result.(*morpheus.GetBackupSettingsResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.BackupSettings == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("BackupSettings"))
	}

	backupSetting := result.BackupSettings
	d.SetId(backupSettingsID)
	d.Set("scheduled_backups", backupSetting.BackupsEnabled)
	d.Set("create_backups", backupSetting.CreateBackups)
	d.Set("backup_appliance", backupSetting.BackupAppliance)
	d.Set("default_backup_storage_bucket_id", backupSetting.DefaultStorageBucket.ID)
	d.Set("default_backup_schedule_id", backupSetting.DefaultSchedule.ID)
	d.Set("retention_days", backupSetting.RetentionCount)

	return diags
}

func resourceSettingBackupUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var scheduledBackups bool
	if v, ok := d.Get("scheduled_backups").(bool); ok {
		scheduledBackups = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("scheduled_backups", d.Get("scheduled_backups")))
	}

	var createBackups bool
	if v, ok := d.Get("create_backups").(bool); ok {
		createBackups = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("create_backups", d.Get("create_backups")))
	}

	var backupAppliance bool
	if v, ok := d.Get("backup_appliance").(bool); ok {
		backupAppliance = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("backup_appliance", d.Get("backup_appliance")))
	}

	var retentionDays int
	if v, ok := d.Get("retention_days").(int); ok {
		retentionDays = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("retention_days", d.Get("retention_days")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"backupSettings": map[string]any{
				"backupsEnabled":  scheduledBackups,
				"createBackups":   createBackups,
				"backupAppliance": backupAppliance,
				"retentionCount":  retentionDays,
			},
		},
	}

	var defaultStorageBucketID int
	if v, ok := d.Get("default_backup_storage_bucket_id").(int); ok {
		defaultStorageBucketID = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError("default_backup_storage_bucket_id", d.Get("default_backup_storage_bucket_id")),
		)
	}

	if defaultStorageBucketID != 0 {
		req.Body["defaultStorageBucket"] = map[string]any{
			"id": defaultStorageBucketID,
		}
	}

	var defaultBackupScheduleID int
	if v, ok := d.Get("default_backup_schedule_id").(int); ok {
		defaultBackupScheduleID = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError("default_backup_schedule_id", d.Get("default_backup_schedule_id")),
		)
	}

	if defaultBackupScheduleID != 0 {
		req.Body["defaultSchedule"] = map[string]any{
			"id": defaultBackupScheduleID,
		}
	}

	log.Printf("API Update: %s", req)

	resp, err := client.UpdateBackupSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateBackupSettingsResult
	if v, ok := resp.Result.(*morpheus.UpdateBackupSettingsResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.BackupSettings == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("BackupSettings"))
	}

	d.SetId(backupSettingsID)

	return resourceSettingBackupRead(ctx, d, meta)
}

func resourceSettingBackupDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}
