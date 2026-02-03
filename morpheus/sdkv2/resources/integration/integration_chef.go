// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func ResourceIntegrationChef() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Chef integration resource",
		CreateContext: resourceIntegrationChefCreate,
		ReadContext:   resourceIntegrationChefRead,
		UpdateContext: resourceIntegrationChefUpdate,
		DeleteContext: resourceIntegrationChefDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the Chef integration",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Chef integration",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the Chef integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The url of the Chef server",
				Required:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The version of the Chef server",
				Optional:    true,
			},
			"windows_version": {
				Type:        schema.TypeString,
				Description: "The Windows agent version",
				Optional:    true,
			},
			"windows_msi_install_url": {
				Type:        schema.TypeString,
				Description: "The URL for the Windows MSI installation package",
				Optional:    true,
			},
			"organization": {
				Type:        schema.TypeString,
				Description: "The chef organization",
				Optional:    true,
			},
			"use_fqdn_node_name": {
				Type:        schema.TypeBool,
				Description: "Whether to use the FQDN of the node instead of the instance name",
				Optional:    true,
				Default:     false,
			},
			"username": {
				Type:          schema.TypeString,
				Description:   "The username of the account used to connect to the Chef server",
				Optional:      true,
				ConflictsWith: []string{"credential_id"},
			},
			"private_key": {
				Type:        schema.TypeString,
				Description: "The private key of the account used to connect to the Chef server",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
				ConflictsWith: []string{"credential_id"},
			},
			"credential_id": {
				Description:   "The ID of the credential store entry used for authentication",
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"username", "private_key"},
			},
			"organization_validator_key": {
				Type:        schema.TypeString,
				Description: "The organization validator key used to connect to the Chef server",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
			},
			/* AWAITING API SUPPORT
			"databags": {
				Description: "",
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			*/
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceIntegrationChefCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	integration := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	integration["name"] = name

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}
	integration["enabled"] = enabled
	integration["type"] = "chef"

	var url string
	if v, ok := d.Get("url").(string); ok {
		url = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("url", d.Get("url")))
	}
	integration["serviceUrl"] = url

	var version string
	if v, ok := d.Get("version").(string); ok {
		version = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("version", d.Get("version")))
	}
	integration["serviceVersion"] = version

	var windowsVersion string
	if v, ok := d.Get("windows_version").(string); ok {
		windowsVersion = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("windows_version", d.Get("windows_version")))
	}
	integration["serviceWindowsVersion"] = windowsVersion

	config := make(map[string]any)

	var windowsInstallUrl string
	if v, ok := d.Get("windows_msi_install_url").(string); ok {
		windowsInstallUrl = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("windows_msi_install_url", d.Get("windows_msi_install_url")))
	}
	config["windowsInstallUrl"] = windowsInstallUrl

	var useFqdn bool
	if v, ok := d.Get("use_fqdn_node_name").(bool); ok {
		useFqdn = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("use_fqdn_node_name", d.Get("use_fqdn_node_name")))
	}
	config["useFqdn"] = useFqdn

	var organization string
	if v, ok := d.Get("organization").(string); ok {
		organization = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("organization", d.Get("organization")))
	}
	config["org"] = organization

	var credentialId int
	if v, ok := d.Get("credential_id").(int); ok {
		credentialId = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("credential_id", d.Get("credential_id")))
	}

	if credentialId != 0 {
		credential := make(map[string]any)
		credential["type"] = "username-keypair"
		credential["id"] = credentialId
		credential["credential"] = credential
	} else {
		credential := make(map[string]any)
		credential["type"] = "local"
		integration["credential"] = credential

		var username string
		if v, ok := d.Get("username").(string); ok {
			username = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
		}
		config["chefUser"] = username

		var privateKey string
		if v, ok := d.Get("private_key").(string); ok {
			privateKey = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("private_key", d.Get("private_key")))
		}
		config["userKey"] = privateKey
	}

	var orgValidatorKey string
	if v, ok := d.Get("organization_validator_key").(string); ok {
		orgValidatorKey = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("organization_validator_key", d.Get("organization_validator_key")))
	}
	config["orgKey"] = orgValidatorKey

	// databags
	/* AWAITING API SUPPORT
	if d.Get("databags") != nil {
		databagsInput := d.Get("databags").(map[string]any)
		var databags []map[string]any
		for key, value := range databagsInput {
			databag := make(map[string]any)
			databag["name"] = key
			databag["value"] = value.(string)
			databags = append(databags, databag)
		}
		config["databag"] = databags
	}
	*/
	integration["config"] = config

	req := &morpheus.Request{
		Body: map[string]any{
			"integration": integration,
		},
	}

	resp, err := client.CreateIntegration(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	result := resp.Result.(*morpheus.CreateIntegrationResult)
	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("CreateIntegrationResult"))
	}

	integrationResult := result.Integration
	if integrationResult == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(integrationResult.ID))

	diags = append(diags, resourceIntegrationChefRead(ctx, d, meta)...)

	return diags
}

func resourceIntegrationChefRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindIntegrationByName(name)
	} else if id != "" {
		resp, err = client.GetIntegration(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Integration cannot be read without name or id")
	}

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			log.Printf("Forcing recreation of resource")
			d.SetId("")

			return diags
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	result := resp.Result.(*morpheus.GetIntegrationResult)
	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("GetIntegrationResult"))
	}

	integration := result.Integration
	if integration == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}

	d.SetId(convert.Int64ToString(integration.ID))
	d.Set("name", integration.Name)
	d.Set("enabled", integration.Enabled)
	d.Set("url", integration.URL)
	d.Set("version", integration.Version)
	d.Set("windows_version", integration.WindowsVersion)
	d.Set("windows_msi_install_url", integration.Config.WindowsInstallURL)
	d.Set("organization", integration.Config.Org)
	d.Set("use_fqdn_node_name", integration.Config.UseFqdn)

	if integration.Credential.ID == 0 {
		d.Set("username", integration.Config.ChefUser)
		d.Set("private_key", integration.Config.UserKeyHash)
	} else {
		d.Set("credential_id", integration.Credential.ID)
	}

	d.Set("organization_validator_key", integration.Config.OrgKeyHash)

	// databags
	/* AWAITING API SUPPORT
	databags := make(map[string]any)
	if integration.Config.Databags != nil {
		output := integration.Config.Databags
		databagList := output
		// iterate over the array of databags
		for i := 0; i < len(databagList); i++ {
			databag := databagList[i]
			databagName := databag.Path
			databags[databagName] = databag.Key
		}
	}
	d.Set("databags", databags)
	*/

	return diags
}

func resourceIntegrationChefUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}
	id := d.Id()

	integration := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	integration["name"] = name

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}
	integration["enabled"] = enabled
	integration["type"] = "chef"

	var url string
	if v, ok := d.Get("url").(string); ok {
		url = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("url", d.Get("url")))
	}
	integration["serviceUrl"] = url

	var version string
	if v, ok := d.Get("version").(string); ok {
		version = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("version", d.Get("version")))
	}
	integration["serviceVersion"] = version

	var windowsVersion string
	if v, ok := d.Get("windows_version").(string); ok {
		windowsVersion = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("windows_version", d.Get("windows_version")))
	}
	integration["serviceWindowsVersion"] = windowsVersion

	config := make(map[string]any)

	var windowsInstallUrl string
	if v, ok := d.Get("windows_msi_install_url").(string); ok {
		windowsInstallUrl = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("windows_msi_install_url", d.Get("windows_msi_install_url")))
	}
	config["windowsInstallUrl"] = windowsInstallUrl

	var useFqdn bool
	if v, ok := d.Get("use_fqdn_node_name").(bool); ok {
		useFqdn = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("use_fqdn_node_name", d.Get("use_fqdn_node_name")))
	}
	config["useFqdn"] = useFqdn

	var organization string
	if v, ok := d.Get("organization").(string); ok {
		organization = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("organization", d.Get("organization")))
	}
	config["org"] = organization

	var credentialId int
	if v, ok := d.Get("credential_id").(int); ok {
		credentialId = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("credential_id", d.Get("credential_id")))
	}

	if credentialId != 0 {
		credential := make(map[string]any)
		credential["type"] = "username-keypair"
		credential["id"] = credentialId
		integration["credential"] = credential
	} else {
		credential := make(map[string]any)
		credential["type"] = "local"
		integration["credential"] = credential
		if d.HasChange("username") {
			var username string
			if v, ok := d.Get("username").(string); ok {
				username = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
			}
			config["chefUser"] = username
		}
		if d.HasChange("private_key") {
			var privateKey string
			if v, ok := d.Get("private_key").(string); ok {
				privateKey = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("private_key", d.Get("private_key")))
			}
			config["userKey"] = privateKey
		}
	}

	if d.HasChange("organization_validator_key") {
		var orgValidatorKey string
		if v, ok := d.Get("organization_validator_key").(string); ok {
			orgValidatorKey = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("organization_validator_key", d.Get("organization_validator_key")))
		}
		config["orgKey"] = orgValidatorKey
	}

	integration["config"] = config

	req := &morpheus.Request{
		Body: map[string]any{
			"integration": integration,
		},
	}

	resp, err := client.UpdateIntegration(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateIntegrationResult)
	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("UpdateIntegrationResult"))
	}

	integrationResult := result.Integration
	if integrationResult == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(integrationResult.ID))

	return resourceIntegrationChefRead(ctx, d, meta)
}

func resourceIntegrationChefDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteIntegration(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return diag.FromErr(err)
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}
