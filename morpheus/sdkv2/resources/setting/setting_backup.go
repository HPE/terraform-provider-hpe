// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package setting

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
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

// backupSettingsRequest builds the body of PUT /api/backup-settings from the resource
// data. Everything -- the default storage bucket and the default schedule included --
// goes under "backupSettings": the appliance reads only that key, and a field placed
// beside it is accepted and silently ignored (MORPH-20485). Shared by create and
// update, which send the same body.
func backupSettingsRequest(d *schema.ResourceData) (*morpheus.Request, diag.Diagnostics) {
	settings := map[string]any{}

	for key, attr := range map[string]string{
		"backupsEnabled":  "scheduled_backups",
		"createBackups":   "create_backups",
		"backupAppliance": "backup_appliance",
	} {
		v, ok := d.Get(attr).(bool)
		if !ok {
			return nil, diag.FromErr(helpers.TypeAssertFailError(attr, d.Get(attr)))
		}
		settings[key] = v
	}

	retentionDays, ok := d.Get("retention_days").(int)
	if !ok {
		return nil, diag.FromErr(helpers.TypeAssertFailError("retention_days", d.Get("retention_days")))
	}
	settings["retentionCount"] = retentionDays

	for key, attr := range map[string]string{
		"defaultStorageBucket": "default_backup_storage_bucket_id",
		"defaultSchedule":      "default_backup_schedule_id",
	} {
		id, ok := d.Get(attr).(int)
		if !ok {
			return nil, diag.FromErr(helpers.TypeAssertFailError(attr, d.Get(attr)))
		}
		if id != 0 {
			settings[key] = map[string]any{"id": id}
		}
	}

	return &morpheus.Request{Body: map[string]any{"backupSettings": settings}}, nil
}

func resourceSettingBackupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	req, bodyDiags := backupSettingsRequest(d)
	if bodyDiags.HasError() {
		return bodyDiags
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

	// The backup-settings PUT returns only a success envelope (no backupSettings
	// object), so do not require it here. State is populated by the read below.
	_ = result

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

	req, bodyDiags := backupSettingsRequest(d)
	if bodyDiags.HasError() {
		return bodyDiags
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

	// The backup-settings PUT returns only a success envelope (no backupSettings
	// object), so do not require it here. State is populated by the read below.
	_ = result

	d.SetId(backupSettingsID)

	return resourceSettingBackupRead(ctx, d, meta)
}

func resourceSettingBackupDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}
