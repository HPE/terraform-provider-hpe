// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func ResourceIntegrationGit() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a git integration resource",
		CreateContext: resourceIntegrationGitCreate,
		ReadContext:   resourceIntegrationGitRead,
		UpdateContext: resourceIntegrationGitUpdate,
		DeleteContext: resourceIntegrationGitDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the git integration",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the git integration",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the git integration is enabled",
				Optional:    true,
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "The url of the git repository",
				Required:    true,
			},
			"default_branch": {
				Type:        schema.TypeString,
				Description: "The default branch of the git repository",
				Optional:    true,
				Computed:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The username of the account used to authenticate to the git repository",
				Optional:    true,
				Computed:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "The password of the account used to authenticate to the git repository",
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
			"access_token": {
				Type:        schema.TypeString,
				Description: "The access token of the account used to authenticate to the git repository",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					h := sha256.New()
					h.Write([]byte(new))
					sha256Hash := hex.EncodeToString(h.Sum(nil))

					return strings.EqualFold(old, sha256Hash)
				},
			},
			"key_pair_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the key pair used to authenticate to the git repository",
				Optional:    true,
				Computed:    true,
			},
			"enable_git_caching": {
				Type:        schema.TypeBool,
				Description: "Whether the git repository is cached",
				Optional:    true,
				Computed:    true,
			},
			"repository_ids": {
				Computed:    true,
				Type:        schema.TypeMap,
				Description: "A map of git repository ids for use with integrations that reference a git repository",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceIntegrationGitCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	integration["type"] = "git"

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

	var accessToken string
	if v, ok := d.Get("access_token").(string); ok {
		accessToken = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("access_token", d.Get("access_token")))
	}
	integration["serviceToken"] = accessToken

	var keyPairID int
	if v, ok := d.Get("key_pair_id").(int); ok {
		keyPairID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("key_pair_id", d.Get("key_pair_id")))
	}
	integration["serviceKey"] = keyPairID

	config := make(map[string]any)

	var defaultBranch string
	if v, ok := d.Get("default_branch").(string); ok {
		defaultBranch = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("default_branch", d.Get("default_branch")))
	}
	config["defaultBranch"] = defaultBranch

	var cacheEnabled bool
	if v, ok := d.Get("enable_git_caching").(bool); ok {
		cacheEnabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_git_caching", d.Get("enable_git_caching")))
	}
	config["cacheEnabled"] = cacheEnabled

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

	diags = append(diags, resourceIntegrationGitRead(ctx, d, meta)...)

	return diags
}

func resourceIntegrationGitRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	d.Set("access_token", integration.TokenHash)

	if integration.ServiceKey.ID != 0 {
		d.Set("key_pair_id", integration.ServiceKey.ID)
	}

	d.Set("default_branch", integration.Config.DefaultBranch)
	d.Set("enable_git_caching", integration.Config.CacheEnabled)

	resp, err = client.Execute(&morpheus.Request{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/options/codeRepositories?integrationId=%d", integration.ID),
		QueryParams: map[string]string{},
	})
	if err != nil {
		log.Println("API ERROR: ", err)
	}
	log.Println("API RESPONSE:", resp)
	repoIDs := make(map[string]int)

	var itemResponsePayload CodeRepositories
	json.Unmarshal(resp.Body, &itemResponsePayload)

	if itemResponsePayload.Data != nil {
		for _, v := range itemResponsePayload.Data {
			repoIDs[v.Name] = v.Value
		}
	}
	d.Set("repository_ids", repoIDs)

	return diags
}

func resourceIntegrationGitUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	integration["type"] = "git"

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

	var accessToken string
	if v, ok := d.Get("access_token").(string); ok {
		accessToken = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("access_token", d.Get("access_token")))
	}
	integration["serviceToken"] = accessToken

	var keyPairID int
	if v, ok := d.Get("key_pair_id").(int); ok {
		keyPairID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("key_pair_id", d.Get("key_pair_id")))
	}
	integration["serviceKey"] = keyPairID

	config := make(map[string]any)

	var defaultBranch string
	if v, ok := d.Get("default_branch").(string); ok {
		defaultBranch = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("default_branch", d.Get("default_branch")))
	}
	config["defaultBranch"] = defaultBranch

	var cacheEnabled bool
	if v, ok := d.Get("enable_git_caching").(bool); ok {
		cacheEnabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_git_caching", d.Get("enable_git_caching")))
	}
	config["cacheEnabled"] = cacheEnabled

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

	return resourceIntegrationGitRead(ctx, d, meta)
}

func resourceIntegrationGitDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

type CodeRepositories struct {
	Success bool `json:"success"`
	Data    []struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	} `json:"data"`
}
