// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

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

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func validateVROAuthConfig(ctx context.Context, d *schema.ResourceDiff, meta any) error {
	var authType string
	if v, ok := d.Get("auth_type").(string); ok {
		authType = v
	} else {
		return helpers.TypeAssertFailError("auth_type", d.Get("auth_type"))
	}

	if authType == "aria" {
		// For aria auth type, api_token is required
		var apiToken string
		if v, ok := d.Get("api_token").(string); ok {
			apiToken = v
		} else {
			return helpers.TypeAssertFailError("api_token", d.Get("api_token"))
		}
		if apiToken == "" {
			return fmt.Errorf("api_token is required when auth_type is 'aria'")
		}
	} else {
		// For non-aria auth types, username, password, tenant, and auth_id are required
		var username string
		if v, ok := d.Get("username").(string); ok {
			username = v
		} else {
			return helpers.TypeAssertFailError("username", d.Get("username"))
		}

		var password string
		if v, ok := d.Get("password").(string); ok {
			password = v
		} else {
			return helpers.TypeAssertFailError("password", d.Get("password"))
		}

		var tenant string
		if v, ok := d.Get("tenant").(string); ok {
			tenant = v
		} else {
			return helpers.TypeAssertFailError("tenant", d.Get("tenant"))
		}

		var authId string
		if v, ok := d.Get("auth_id").(string); ok {
			authId = v
		} else {
			return helpers.TypeAssertFailError("auth_id", d.Get("auth_id"))
		}

		if username == "" || password == "" || tenant == "" || authId == "" {
			return fmt.Errorf("the following fields are required when auth_type is not 'aria': " +
				"username, password, tenant, auth_id")
		}
	}

	return nil
}

func ResourceIntegrationVRO() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a vRealize Orchestrator integration resource",
		CreateContext: resourceIntegrationVROCreate,
		ReadContext:   resourceIntegrationVRORead,
		UpdateContext: resourceIntegrationVROUpdate,
		DeleteContext: resourceIntegrationVRODelete,
		CustomizeDiff: validateVROAuthConfig,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the vRO integration",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the vRO integration",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the vRO integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The url of the vRO server",
				Required:    true,
			},
			"auth_type": {
				Type:         schema.TypeString,
				Description:  "The authentication type for the vRO integration",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"basic", "aria", "oauth", "vra"}, false),
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The username of the account used to connect to vRO (required for non-aria auth types)",
				Optional:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the account used to connect to vRO (required for non-aria auth types)",
				Optional:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
				DiffSuppressOnRefresh: true,
			},
			"tenant": {
				Type:        schema.TypeString,
				Description: "The tenant of the account used to connect to vRO (required for non-aria auth types)",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
			},
			"auth_id": {
				Type:        schema.TypeString,
				Description: "The authentication ID for the vRO integration (required for non-aria auth types)",
				Optional:    true,
				Computed:    true,
			},
			"api_token": {
				Type:        schema.TypeString,
				Description: "The API token for vRO (required when auth_type is aria)",
				Optional:    true,
				Sensitive:   true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceIntegrationVROCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

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

	integration["type"] = "vro"

	var authType string
	if v, ok := d.Get("auth_type").(string); ok {
		authType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("auth_type", d.Get("auth_type")))
	}
	integration["authType"] = authType

	var url string
	if v, ok := d.Get("url").(string); ok {
		url = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("url", d.Get("url")))
	}
	integration["serviceUrl"] = url

	// Handle different auth types
	if authType == "aria" {
		// For aria auth type, use api_token
		var apiToken string
		if v, ok := d.Get("api_token").(string); ok {
			apiToken = v
		}
		config := make(map[string]any)
		config["apiToken"] = apiToken
		integration["config"] = config
	} else {
		// For non-aria auth types (basic, oauth, vra), use username/password/tenant
		var username string
		if v, ok := d.Get("username").(string); ok {
			username = v
		}
		integration["serviceUsername"] = username

		var password string
		if v, ok := d.Get("password").(string); ok {
			password = v
		}
		integration["servicePassword"] = password

		var tenant string
		if v, ok := d.Get("tenant").(string); ok {
			tenant = v
		}
		integration["serviceToken"] = tenant

		var authId string
		if v, ok := d.Get("auth_id").(string); ok {
			authId = v
		}
		// For non-aria auth types, authId must be set (can be empty string for local credentials)
		integration["authId"] = authId
	}

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

	return resourceIntegrationVRORead(ctx, d, meta)
}

func resourceIntegrationVRORead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	d.Set("username", integration.Username)
	d.Set("password", integration.PasswordHash)
	d.Set("tenant", integration.TokenHash)
	// d.Set("auth_type", integration.Config)

	return diags
}

func resourceIntegrationVROUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	integration["type"] = "vro"

	var authType string
	if v, ok := d.Get("auth_type").(string); ok {
		authType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("auth_type", d.Get("auth_type")))
	}
	integration["authType"] = authType

	var url string
	if v, ok := d.Get("url").(string); ok {
		url = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("url", d.Get("url")))
	}
	integration["serviceUrl"] = url

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

	var tenant string
	if v, ok := d.Get("tenant").(string); ok {
		tenant = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("tenant", d.Get("tenant")))
	}
	integration["token"] = tenant

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

	return resourceIntegrationVRORead(ctx, d, meta)
}

func resourceIntegrationVRODelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
