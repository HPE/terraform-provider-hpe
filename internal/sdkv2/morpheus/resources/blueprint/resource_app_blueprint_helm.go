// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func ResourceAppBlueprintHelm() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus helm app blueprint resource",
		CreateContext: resourceAppBlueprintHelmCreate,
		ReadContext:   resourceAppBlueprintHelmRead,
		UpdateContext: resourceAppBlueprintHelmUpdate,
		DeleteContext: resourceAppBlueprintHelmDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the helm app blueprint",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the helm app blueprint",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the helm app blueprint",
				Optional:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the helm app blueprint",
				Optional:    true,
			},
			"working_path": {
				Type:        schema.TypeString,
				Description: "The path of the helm chart in the git repository",
				Optional:    true,
				Default:     "./",
			},
			"integration_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the git integration",
				Required:    true,
			},
			"repository_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the git repository",
				Required:    true,
			},
			"version_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the repository to pull (main, master, etc.)",
				Optional:    true,
				Default:     "master",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

//nolint:goconst
func resourceAppBlueprintHelmCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	blueprintType := "helm"
	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}
	var category string
	if v, ok := d.Get("category").(string); ok {
		category = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("category", d.Get("category")))
	}

	config := make(map[string]any)
	config["name"] = name
	config["description"] = description
	config["category"] = category
	config["type"] = blueprintType

	helmConfig := make(map[string]any)
	config["helm"] = helmConfig

	helmConfig["configType"] = "git"
	helmGitConfig := make(map[string]any)
	helmGitConfig["integrationId"] = d.Get("integration_id")
	helmGitConfig["repoId"] = d.Get("repository_id")
	var versionRef string
	if v, ok := d.Get("version_ref").(string); ok {
		versionRef = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("version_ref", d.Get("version_ref")))
	}
	helmGitConfig["branch"] = versionRef
	var workingPath string
	if v, ok := d.Get("working_path").(string); ok {
		workingPath = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("working_path", d.Get("working_path")))
	}
	helmGitConfig["path"] = workingPath
	helmConfig["git"] = helmGitConfig

	req := &morpheus.Request{
		Body: map[string]any{
			"blueprint": map[string]any{
				"name":        name,
				"type":        blueprintType,
				"description": description,
				"category":    category,
				"config":      config,
			},
		},
	}

	resp, err := client.CreateBlueprint(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.CreateBlueprintResult
	if v, ok := resp.Result.(*morpheus.CreateBlueprintResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}
	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("CreateBlueprintResult"))
	}
	blueprint := result.Blueprint
	if blueprint == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Blueprint"))
	}
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(blueprint.ID))

	diags = append(diags, resourceAppBlueprintHelmRead(ctx, d, meta)...)

	return diags
}

func resourceAppBlueprintHelmRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindBlueprintByName(name)
	} else if id != "" {
		resp, err = client.GetBlueprint(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Blueprint cannot be read without name or id")
	}

	if err != nil {
		// 404 is ok?
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
	var helmBlueprint AppBlueprintHelm
	json.Unmarshal(resp.Body, &helmBlueprint)
	d.SetId(convert.IntToString(helmBlueprint.Blueprint.ID))
	d.Set("name", helmBlueprint.Blueprint.Name)
	d.Set("description", helmBlueprint.Blueprint.Description)
	d.Set("category", helmBlueprint.Blueprint.Category)
	d.Set("working_path", helmBlueprint.Blueprint.Config.Helm.Git.Path)
	d.Set("integration_id", helmBlueprint.Blueprint.Config.Helm.Git.IntegrationId)
	d.Set("repository_id", helmBlueprint.Blueprint.Config.Helm.Git.RepoId)
	d.Set("version_ref", helmBlueprint.Blueprint.Config.Helm.Git.Branch)

	return diags
}

func resourceAppBlueprintHelmUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}
	id := d.Id()
	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	blueprintType := "helm"
	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}
	var category string
	if v, ok := d.Get("category").(string); ok {
		category = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("category", d.Get("category")))
	}

	config := make(map[string]any)
	config["name"] = name
	config["description"] = description
	config["category"] = category
	config["type"] = blueprintType

	helmConfig := make(map[string]any)
	config["helm"] = helmConfig

	helmConfig["configType"] = "git"
	helmGitConfig := make(map[string]any)
	helmGitConfig["integrationId"] = d.Get("integration_id")
	helmGitConfig["repoId"] = d.Get("repository_id")
	var versionRef string
	if v, ok := d.Get("version_ref").(string); ok {
		versionRef = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("version_ref", d.Get("version_ref")))
	}
	helmGitConfig["branch"] = versionRef
	var workingPath string
	if v, ok := d.Get("working_path").(string); ok {
		workingPath = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("working_path", d.Get("working_path")))
	}
	helmGitConfig["path"] = workingPath
	helmConfig["git"] = helmGitConfig

	req := &morpheus.Request{
		Body: map[string]any{
			"blueprint": map[string]any{
				"name":        name,
				"type":        blueprintType,
				"description": description,
				"category":    category,
				"config":      config,
			},
		},
	}

	resp, err := client.UpdateBlueprint(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	var result *morpheus.UpdateBlueprintResult
	if v, ok := resp.Result.(*morpheus.UpdateBlueprintResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}
	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("UpdateBlueprintResult"))
	}
	blueprint := result.Blueprint
	if blueprint == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Blueprint"))
	}
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(blueprint.ID))

	return resourceAppBlueprintHelmRead(ctx, d, meta)
}

func resourceAppBlueprintHelmDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	resp, err := client.DeleteBlueprint(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return nil
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}

type AppBlueprintHelm struct {
	Blueprint struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Config      struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Helm        struct {
				Configtype string `json:"configType"`
				Git        struct {
					Path          string `json:"path"`
					RepoId        int    `json:"repoId"`
					IntegrationId int    `json:"integrationId"`
					Branch        string `json:"branch"`
				} `json:"git"`
			} `json:"helm"`
			Type     string `json:"type"`
			Category string `json:"category"`
			Image    string `json:"image"`
		} `json:"config"`
		Visibility         string `json:"visibility"`
		Resourcepermission struct {
			All      bool  `json:"all"`
			Sites    []any `json:"sites"`
			AllPlans bool  `json:"allPlans"`
			Plans    []any `json:"plans"`
		} `json:"resourcePermission"`
		Owner struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
		} `json:"owner"`
		Tenant struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"tenant"`
	} `json:"blueprint"`
}
