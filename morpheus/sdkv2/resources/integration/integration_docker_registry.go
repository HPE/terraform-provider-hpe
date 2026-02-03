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

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func ResourceIntegrationDockerRegistry() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Docker registry integration resource",
		CreateContext: resourceIntegrationDockerRegistryCreate,
		ReadContext:   resourceIntegrationDockerRegistryRead,
		UpdateContext: resourceIntegrationDockerRegistryUpdate,
		DeleteContext: resourceIntegrationDockerRegistryDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the docker registry integration",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the docker registry integration",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the docker registry integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The url of the docker registry",
				Required:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The username of the account used to authenticate to the docker registry",
				Optional:    true,
				Computed:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the account used to authenticate to the docker registry",
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

func resourceIntegrationDockerRegistryCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}

	var url string
	if v, ok := d.Get("url").(string); ok {
		url = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("url", d.Get("url")))
	}

	var username string
	if v, ok := d.Get("username").(string); ok {
		username = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
	}

	var password string
	if v, ok := d.Get("password").(string); ok {
		password = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("password", d.Get("password")))
	}

	integration["name"] = name
	integration["enabled"] = enabled
	integration["type"] = "docker.registry"
	integration["serviceUrl"] = url
	integration["serviceUsername"] = username
	integration["servicePassword"] = password

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
		return diag.FromErr(helpers.TypeAssertFailError("CreateIntegrationResult", resp.Result))
	}

	if result.Integration == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}

	integrationResult := result.Integration
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(integrationResult.ID))

	return resourceIntegrationDockerRegistryRead(ctx, d, meta)
}

func resourceIntegrationDockerRegistryRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		return diag.FromErr(helpers.TypeAssertFailError("GetIntegrationResult", resp.Result))
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

	return diags
}

func resourceIntegrationDockerRegistryUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}

	var url string
	if v, ok := d.Get("url").(string); ok {
		url = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("url", d.Get("url")))
	}

	var username string
	if v, ok := d.Get("username").(string); ok {
		username = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("username", d.Get("username")))
	}

	var password string
	if v, ok := d.Get("password").(string); ok {
		password = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("password", d.Get("password")))
	}

	integration["name"] = name
	integration["enabled"] = enabled
	integration["type"] = "docker.registry"
	integration["serviceUrl"] = url
	integration["serviceUsername"] = username
	integration["servicePassword"] = password

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
		return diag.FromErr(helpers.TypeAssertFailError("UpdateIntegrationResult", resp.Result))
	}

	if result.Integration == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Integration"))
	}

	integrationResult := result.Integration

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(integrationResult.ID))

	return resourceIntegrationDockerRegistryRead(ctx, d, meta)
}

func resourceIntegrationDockerRegistryDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
