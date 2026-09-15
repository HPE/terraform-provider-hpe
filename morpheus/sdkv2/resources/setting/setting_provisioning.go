// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package setting

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func ResourceSettingProvisioning() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus provisioning setting resource.",
		CreateContext: resourceSettingProvisioningCreate,
		ReadContext:   resourceSettingProvisioningRead,
		UpdateContext: resourceSettingProvisioningUpdate,
		DeleteContext: resourceSettingProvisioningDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the provisioning settings",
				Computed:    true,
			},
			"allow_zone_selection": {
				Type:        schema.TypeBool,
				Description: "Displays or hides Cloud Selection dropdown in Provisioning wizard.",
				Optional:    true,
				Computed:    true,
			},
			"allow_host_selection": {
				Type:        schema.TypeBool,
				Description: "Displays or hides Host Selection dropdown in Provisioning wizard.",
				Optional:    true,
				Computed:    true,
			},
			"require_environments": {
				Type:        schema.TypeBool,
				Description: "Forces users to select and Environment during provisioning",
				Optional:    true,
				Computed:    true,
			},
			"show_pricing": {
				Type:        schema.TypeBool,
				Description: "Displays or hides Pricing in Provisioning wizard and Instance and Host detail pages.",
				Optional:    true,
				Computed:    true,
			},
			"hide_datastore_stats": {
				Type:        schema.TypeBool,
				Description: "Hides Datastore utilization and size stats in provisioning and app wizards.",
				Optional:    true,
				Computed:    true,
			},
			"cross_tenant_naming_policies": {
				Type:        schema.TypeBool,
				Description: "Enable for the sequence value in naming policies to apply across tenants.",
				Optional:    true,
				Computed:    true,
			},
			"reuse_sequence": {
				Type: schema.TypeBool,
				Description: "When enabled, sequence numbers can be reused when Instances are removed. " +
					"Deselect this option and Morpheus will track issued sequence numbers and use the " +
					"next available number each time.",
				Optional: true,
				Computed: true,
			},
			"show_console_keyboard_settings": {
				Type: schema.TypeBool,
				Description: "Displays the keyboard layout selection when opening a guest " +
					"console (VNC/RDP) so users can pick the keyboard mapping for their session.",
				Optional: true,
				Computed: true,
			},
			"cloudinit_username": {
				Type:         schema.TypeString,
				Description:  "User to be added to Linux Instances during provisioning.",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"cloudinit_password": {
				Type:        schema.TypeString,
				Description: "Password to be set for the Cloud-Init Linux user.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
				DiffSuppressOnRefresh: true,
			},
			"windows_password": {
				Type:        schema.TypeString,
				Description: "Password to be set for the Windows Administrator User during provisioning.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
				DiffSuppressOnRefresh: true,
			},
			"pxe_root_password": {
				Type:        schema.TypeString,
				Description: "Password to be set for Root during PXE Boots.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
				DiffSuppressOnRefresh: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSettingProvisioningCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var allowZoneSelection bool
	if v, ok := d.Get("allow_zone_selection").(bool); ok {
		allowZoneSelection = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("allow_zone_selection", d.Get("allow_zone_selection")))
	}

	var allowHostSelection bool
	if v, ok := d.Get("allow_host_selection").(bool); ok {
		allowHostSelection = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("allow_host_selection", d.Get("allow_host_selection")))
	}

	var requireEnvironments bool
	if v, ok := d.Get("require_environments").(bool); ok {
		requireEnvironments = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("require_environments", d.Get("require_environments")))
	}

	var showPricing bool
	if v, ok := d.Get("show_pricing").(bool); ok {
		showPricing = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("show_pricing", d.Get("show_pricing")))
	}

	var hideDatastoreStats bool
	if v, ok := d.Get("hide_datastore_stats").(bool); ok {
		hideDatastoreStats = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("hide_datastore_stats", d.Get("hide_datastore_stats")))
	}

	var crossTenantNamingPolicies bool
	if v, ok := d.Get("cross_tenant_naming_policies").(bool); ok {
		crossTenantNamingPolicies = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError(
			"cross_tenant_naming_policies",
			d.Get("cross_tenant_naming_policies"),
		))
	}

	var reuseSequence bool
	if v, ok := d.Get("reuse_sequence").(bool); ok {
		reuseSequence = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("reuse_sequence", d.Get("reuse_sequence")))
	}

	var showConsoleKeyboardSettings bool
	if v, ok := d.Get("show_console_keyboard_settings").(bool); ok {
		showConsoleKeyboardSettings = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError(
			"show_console_keyboard_settings",
			d.Get("show_console_keyboard_settings"),
		))
	}

	var cloudInitUsername string
	if v, ok := d.Get("cloudinit_username").(string); ok {
		cloudInitUsername = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloudinit_username", d.Get("cloudinit_username")))
	}

	provisioningSettings := map[string]any{
		"allowZoneSelection":        allowZoneSelection,
		"allowServerSelection":      allowHostSelection,
		"requireEnvironments":       requireEnvironments,
		"showPricing":               showPricing,
		"hideDatastoreStats":        hideDatastoreStats,
		"crossTenantNamingPolicies": crossTenantNamingPolicies,
		"reuseSequence":             reuseSequence,
		// The PUT handler resolves each submitted key by setting-type name, and
		// the underlying setting is named "consoleKeyboardSettings" (the GET
		// side serialises it as "showConsoleKeyboardSettings"). Sending the GET
		// spelling here is silently dropped, so the accepted write key is
		// "consoleKeyboardSettings".
		"consoleKeyboardSettings": showConsoleKeyboardSettings,
		"cloudInitUsername":       cloudInitUsername,
	}

	// The password attributes are Optional+Computed+Sensitive and write-only:
	// an empty string is destructive server-side (it clears an already-set
	// credential), so omit an unset value from the payload rather than sending
	// "". GetOk treats "" as unset, which is exactly what we want here.
	if v, ok := d.GetOk("cloudinit_password"); ok {
		provisioningSettings["cloudInitPassword"] = v.(string)
	}
	if v, ok := d.GetOk("windows_password"); ok {
		provisioningSettings["windowsPassword"] = v.(string)
	}
	if v, ok := d.GetOk("pxe_root_password"); ok {
		provisioningSettings["pxeRootPassword"] = v.(string)
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"provisioningSettings": provisioningSettings,
		},
	}

	resp, err := client.UpdateProvisioningSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	// The PUT endpoint returns only {success, msg, errors} -- never the
	// provisioningSettings envelope (see MORPH-14742). Verify success, then let
	// resourceSettingProvisioningRead re-read via GET.
	result, ok := resp.Result.(*morpheus.UpdateProvisioningSettingsResult)
	if !ok {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}
	if !result.Success {
		return diag.FromErr(fmt.Errorf(
			"provisioning settings update failed: %s (errors: %v)",
			result.Message,
			result.Errors,
		))
	}

	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(1))

	diags = append(diags, resourceSettingProvisioningRead(ctx, d, meta)...)

	return diags
}

func resourceSettingProvisioningRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error

	resp, err = client.GetProvisioningSettings(&morpheus.Request{})
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

	// store resource data
	var result *morpheus.GetProvisioningSettingsResult
	if v, ok := resp.Result.(*morpheus.GetProvisioningSettingsResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.ProvisioningSettings == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ProvisioningSettings"))
	}

	provisioningSetting := result.ProvisioningSettings
	d.SetId(convert.Int64ToString(1))

	d.Set("allow_zone_selection", provisioningSetting.AllowZoneSelection)
	d.Set("allow_host_selection", provisioningSetting.AllowServerSelection)
	d.Set("require_environments", provisioningSetting.RequireEnvironments)
	d.Set("show_pricing", provisioningSetting.ShowPricing)
	d.Set("hide_datastore_stats", provisioningSetting.HideDatastoreStats)
	d.Set("cross_tenant_naming_policies", provisioningSetting.CrossTenantNamingPolicies)
	d.Set("reuse_sequence", provisioningSetting.ReuseSequence)
	d.Set("show_console_keyboard_settings", provisioningSetting.ShowConsoleKeyboardSettings)
	d.Set("cloudinit_username", provisioningSetting.CloudInitUsername)
	// cloudinit_password, windows_password and pxe_root_password are write-only:
	// the API returns only a salted hash that cannot be reproduced from the
	// configured plaintext, so reading it back would cause a permanent diff.
	// Leave those attributes as the value already held in state.

	return diags
}

func resourceSettingProvisioningUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var allowZoneSelection bool
	if v, ok := d.Get("allow_zone_selection").(bool); ok {
		allowZoneSelection = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("allow_zone_selection", d.Get("allow_zone_selection")))
	}

	var allowHostSelection bool
	if v, ok := d.Get("allow_host_selection").(bool); ok {
		allowHostSelection = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("allow_host_selection", d.Get("allow_host_selection")))
	}

	var requireEnvironments bool
	if v, ok := d.Get("require_environments").(bool); ok {
		requireEnvironments = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("require_environments", d.Get("require_environments")))
	}

	var showPricing bool
	if v, ok := d.Get("show_pricing").(bool); ok {
		showPricing = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("show_pricing", d.Get("show_pricing")))
	}

	var hideDatastoreStats bool
	if v, ok := d.Get("hide_datastore_stats").(bool); ok {
		hideDatastoreStats = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("hide_datastore_stats", d.Get("hide_datastore_stats")))
	}

	var crossTenantNamingPolicies bool
	if v, ok := d.Get("cross_tenant_naming_policies").(bool); ok {
		crossTenantNamingPolicies = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError(
			"cross_tenant_naming_policies",
			d.Get("cross_tenant_naming_policies"),
		))
	}

	var reuseSequence bool
	if v, ok := d.Get("reuse_sequence").(bool); ok {
		reuseSequence = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("reuse_sequence", d.Get("reuse_sequence")))
	}

	var showConsoleKeyboardSettings bool
	if v, ok := d.Get("show_console_keyboard_settings").(bool); ok {
		showConsoleKeyboardSettings = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError(
			"show_console_keyboard_settings",
			d.Get("show_console_keyboard_settings"),
		))
	}

	var cloudInitUsername string
	if v, ok := d.Get("cloudinit_username").(string); ok {
		cloudInitUsername = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloudinit_username", d.Get("cloudinit_username")))
	}

	provisioningSettings := map[string]any{
		"allowZoneSelection":        allowZoneSelection,
		"allowServerSelection":      allowHostSelection,
		"requireEnvironments":       requireEnvironments,
		"showPricing":               showPricing,
		"hideDatastoreStats":        hideDatastoreStats,
		"crossTenantNamingPolicies": crossTenantNamingPolicies,
		"reuseSequence":             reuseSequence,
		// The PUT handler resolves each submitted key by setting-type name, and
		// the underlying setting is named "consoleKeyboardSettings" (the GET
		// side serialises it as "showConsoleKeyboardSettings"). Sending the GET
		// spelling here is silently dropped, so the accepted write key is
		// "consoleKeyboardSettings".
		"consoleKeyboardSettings": showConsoleKeyboardSettings,
		"cloudInitUsername":       cloudInitUsername,
	}

	// The password attributes are Optional+Computed+Sensitive and write-only:
	// an empty string is destructive server-side (it clears an already-set
	// credential), so omit an unset value from the payload rather than sending
	// "". GetOk treats "" as unset, leaving a previously-set password intact
	// across unrelated updates (see MORPH-16150).
	if v, ok := d.GetOk("cloudinit_password"); ok {
		provisioningSettings["cloudInitPassword"] = v.(string)
	}
	if v, ok := d.GetOk("windows_password"); ok {
		provisioningSettings["windowsPassword"] = v.(string)
	}
	if v, ok := d.GetOk("pxe_root_password"); ok {
		provisioningSettings["pxeRootPassword"] = v.(string)
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"provisioningSettings": provisioningSettings,
		},
	}

	resp, err := client.UpdateProvisioningSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	// The PUT /api/provisioning-settings endpoint returns only
	// {success, msg, errors}; it never echoes the provisioningSettings
	// envelope, so we must not dereference it here (see MORPH-14742). Verify the
	// operation succeeded, then re-read via GET (which does return the envelope)
	// through resourceSettingProvisioningRead below.
	result, ok := resp.Result.(*morpheus.UpdateProvisioningSettingsResult)
	if !ok {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}
	if !result.Success {
		return diag.FromErr(fmt.Errorf(
			"provisioning settings update failed: %s (errors: %v)",
			result.Message,
			result.Errors,
		))
	}

	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(1))

	return resourceSettingProvisioningRead(ctx, d, meta)
}

func resourceSettingProvisioningDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}
