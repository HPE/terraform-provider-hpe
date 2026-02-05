// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package setting

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

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
				Type:        schema.TypeBool,
				Description: "",
				Optional:    true,
				Computed:    true,
			},
			"cloudinit_username": {
				Type:        schema.TypeString,
				Description: "User to be added to Linux Instances during provisioning.",
				Optional:    true,
				Computed:    true,
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
			// "cloudinit_keypair_id": {
			// 	Type:        schema.TypeInt,
			// 	Description: "ID of the keypair to be added for the Cloud-Init Linux user.",
			// 	Optional:    true,
			// 	Computed:    true,
			// 	Sensitive:   true,
			// },
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

	var cloudInitUsername string
	if v, ok := d.Get("cloudinit_username").(string); ok {
		cloudInitUsername = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloudinit_username", d.Get("cloudinit_username")))
	}

	var cloudInitPassword string
	if v, ok := d.Get("cloudinit_password").(string); ok {
		cloudInitPassword = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloudinit_password", d.Get("cloudinit_password")))
	}

	var windowsPassword string
	if v, ok := d.Get("windows_password").(string); ok {
		windowsPassword = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("windows_password", d.Get("windows_password")))
	}

	var pxeRootPassword string
	if v, ok := d.Get("pxe_root_password").(string); ok {
		pxeRootPassword = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("pxe_root_password", d.Get("pxe_root_password")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"provisioningSettings": map[string]any{
				"allowZoneSelection":        allowZoneSelection,
				"allowServerSelection":      allowHostSelection,
				"requireEnvironments":       requireEnvironments,
				"showPricing":               showPricing,
				"hideDatastoreStats":        hideDatastoreStats,
				"crossTenantNamingPolicies": crossTenantNamingPolicies,
				"reuseSequence":             reuseSequence,
				"cloudInitUsername":         cloudInitUsername,
				"cloudInitPassword":         cloudInitPassword,
				"windowsPassword":           windowsPassword,
				"pxeRootPassword":           pxeRootPassword,
			},
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
	d.Set("cloudinit_password", provisioningSetting.CloudInitPasswordHash)
	// d.Set("cloudinit_keypair_id", provisioningSetting.Cloudinitkeypair.ID)
	d.Set("windows_password", provisioningSetting.WindowsPasswordHash)
	d.Set("pxe_root_password", provisioningSetting.PXERootPasswordHash)

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

	var cloudInitUsername string
	if v, ok := d.Get("cloudinit_username").(string); ok {
		cloudInitUsername = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloudinit_username", d.Get("cloudinit_username")))
	}

	var cloudInitPassword string
	if v, ok := d.Get("cloudinit_password").(string); ok {
		cloudInitPassword = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloudinit_password", d.Get("cloudinit_password")))
	}

	var windowsPassword string
	if v, ok := d.Get("windows_password").(string); ok {
		windowsPassword = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("windows_password", d.Get("windows_password")))
	}

	var pxeRootPassword string
	if v, ok := d.Get("pxe_root_password").(string); ok {
		pxeRootPassword = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("pxe_root_password", d.Get("pxe_root_password")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"provisioningSettings": map[string]any{
				"allowZoneSelection":        allowZoneSelection,
				"allowServerSelection":      allowHostSelection,
				"requireEnvironments":       requireEnvironments,
				"showPricing":               showPricing,
				"hideDatastoreStats":        hideDatastoreStats,
				"crossTenantNamingPolicies": crossTenantNamingPolicies,
				"reuseSequence":             reuseSequence,
				"cloudInitUsername":         cloudInitUsername,
				"cloudInitPassword":         cloudInitPassword,
				"windowsPassword":           windowsPassword,
				"pxeRootPassword":           pxeRootPassword,
			},
		},
	}

	// var cloudInitKeypairId = d.Get("cloudinit_keypair_id").(int)
	// if cloudInitKeypairId != 0 {
	// 	req.Body["cloudInitKeyPair"] = map[string]any{
	// 		"id": cloudInitKeypairId,
	// 	}
	// }

	resp, err := client.UpdateProvisioningSettings(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateProvisioningSettingsResult
	if v, ok := resp.Result.(*morpheus.UpdateProvisioningSettingsResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.ProvisioningSettings == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("ProvisioningSettings"))
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
