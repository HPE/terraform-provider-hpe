// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strings"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

const (
	credentialTypeUsernamePassword = "username-password"
	credentialTypeLocal            = "local"
)

func ResourceIntegrationAnsibleTower() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides an Ansible Tower integration resource",
		CreateContext: resourceIntegrationAnsibleTowerCreate,
		ReadContext:   resourceIntegrationAnsibleTowerRead,
		UpdateContext: resourceIntegrationAnsibleTowerUpdate,
		DeleteContext: resourceIntegrationAnsibleTowerDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the Ansible Tower integration",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Ansible Tower integration",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the Ansible Tower integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The url of the Ansible Tower instance",
				Required:    true,
			},
			"username": {
				Type:          schema.TypeString,
				Description:   "The username of the account used to connect to Ansible Tower",
				Optional:      true,
				ConflictsWith: []string{"credential_id"},
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the account used to connect to Ansible Tower",
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
				ConflictsWith: []string{"username", "password"},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceIntegrationAnsibleTowerCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	integration["type"] = "ansibleTower"
	integration["serviceVersion"] = "v2"

	var credentialID int
	if v, ok := d.Get("credential_id").(int); ok {
		credentialID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("credential_id", d.Get("credential_id")))
	}

	if credentialID != 0 {
		credential := make(map[string]any)
		credential["type"] = credentialTypeUsernamePassword
		credential["id"] = credentialID
		integration["credential"] = credential
	} else {
		credential := make(map[string]any)
		credential["type"] = credentialTypeLocal
		integration["credential"] = credential

		var username string
		if v, ok := d.Get("username").(string); ok {
			username = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
		}
		integration["serviceUsername"] = username

		var password string
		if v, ok := d.Get("password").(string); ok {
			password = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("password", d.Get("password")))
		}
		integration["servicePassword"] = password
	}

	var url string
	if v, ok := d.Get("url").(string); ok {
		url = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("url", d.Get("url")))
	}
	integration["serviceUrl"] = url

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

	var result *morpheus.CreateIntegrationResult
	if v, ok := resp.Result.(*morpheus.CreateIntegrationResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}

	if result.Integration == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}

	integrationResult := result.Integration
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(integrationResult.ID))

	diags = append(diags, resourceIntegrationAnsibleTowerRead(ctx, d, meta)...)

	return diags
}

func resourceIntegrationAnsibleTowerRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	var result *morpheus.GetIntegrationResult
	if v, ok := resp.Result.(*morpheus.GetIntegrationResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}

	if result.Integration == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}

	integration := result.Integration
	d.SetId(convert.Int64ToString(integration.ID))
	d.Set("name", integration.Name)
	d.Set("enabled", integration.Enabled)
	d.Set("url", integration.URL)

	if integration.Credential.ID == 0 {
		d.Set("username", integration.Username)
		d.Set("password", integration.PasswordHash)
	} else {
		d.Set("credential_id", integration.Credential.ID)
	}

	return diags
}

func resourceIntegrationAnsibleTowerUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	integration["type"] = "ansibleTower"
	integration["serviceVersion"] = "v2"

	var url string
	if v, ok := d.Get("url").(string); ok {
		url = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("url", d.Get("url")))
	}
	integration["serviceUrl"] = url

	var credentialID int
	if v, ok := d.Get("credential_id").(int); ok {
		credentialID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("credential_id", d.Get("credential_id")))
	}

	if credentialID != 0 {
		credential := make(map[string]any)
		credential["type"] = credentialTypeUsernamePassword
		credential["id"] = credentialID
		integration["credential"] = credential
	} else {
		credential := make(map[string]any)
		credential["type"] = credentialTypeLocal
		integration["credential"] = credential
		if d.HasChange("username") {
			integration["serviceUsername"] = d.Get("username")
		}
		if d.HasChange("password") {
			integration["servicePassword"] = d.Get("password")
		}
	}

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

	var result *morpheus.UpdateIntegrationResult
	if v, ok := resp.Result.(*morpheus.UpdateIntegrationResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}

	if result.Integration == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}

	integrationResult := result.Integration

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(integrationResult.ID))

	return resourceIntegrationAnsibleTowerRead(ctx, d, meta)
}

func resourceIntegrationAnsibleTowerDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
