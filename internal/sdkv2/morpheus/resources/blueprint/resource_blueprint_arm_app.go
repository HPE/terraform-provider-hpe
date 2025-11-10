// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func ResourceBlueprintArmApp() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus arm app blueprint resource",
		CreateContext: resourceBlueprintArmAppCreate,
		ReadContext:   resourceBlueprintArmAppRead,
		UpdateContext: resourceBlueprintArmAppUpdate,
		DeleteContext: resourceBlueprintArmAppDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the arm app blueprint",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the arm app blueprint",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the arm app blueprint",
				Optional:    true,
				Computed:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the arm app blueprint",
				Optional:    true,
				Computed:    true,
			},
			"install_agent": {
				Type:        schema.TypeBool,
				Description: "Whether to install the Morpheus agent",
				Optional:    true,
			},
			"cloud_init_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether cloud init is enabled",
				Optional:    true,
			},
			"os_type": {
				Type:         schema.TypeString,
				Description:  "The workload operating system type (linux, windows)",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"linux", "windows"}, false),
			},
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the arm app blueprint (json, repository)",
				ValidateFunc: validation.StringInSlice([]string{"json", "repository"}, false),
				Required:     true,
			},
			"blueprint_content": {
				Type:        schema.TypeString,
				Description: "The content of the arm app blueprint. Used when the json source type is specified",
				Optional:    true,
				StateFunc: func(val any) string {
					if v, ok := val.(string); ok {
						return strings.TrimSuffix(v, "\n")
					}

					return ""
				},
			},
			"working_path": {
				Type:        schema.TypeString,
				Description: "The path of the arm app blueprint in the git repository",
				Optional:    true,
			},
			"integration_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the git integration",
				Optional:    true,
			},
			"repository_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the git repository",
				Optional:    true,
			},
			"version_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the repository to pull (main, master, etc.)",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceBlueprintArmAppCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", client))
	}

	blueprintType := "arm"

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", client))
	}

	var category string
	if v, ok := d.Get("category").(string); ok {
		category = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("category", client))
	}

	config := make(map[string]any)
	config["name"] = name
	config["description"] = description
	config["category"] = category
	config["type"] = "arm"

	armConfig := make(map[string]any)
	config["arm"] = armConfig

	var osType string
	if v, ok := d.Get("os_type").(string); ok {
		osType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("os_type", client))
	}
	armConfig["osType"] = osType

	var installAgent bool
	if v, ok := d.Get("install_agent").(bool); ok {
		installAgent = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("install_agent", client))
	}
	armConfig["installAgent"] = installAgent

	var cloudInitEnabled bool
	if v, ok := d.Get("cloud_init_enabled").(bool); ok {
		cloudInitEnabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloud_init_enabled", client))
	}
	armConfig["cloudInitEnabled"] = cloudInitEnabled

	var sourceType string
	if v, ok := d.Get("source_type").(string); ok {
		sourceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", client))
	}

	switch sourceType {
	case "json":
		armConfig["configType"] = "json"
		var blueprintContent string
		if v, ok := d.Get("blueprint_content").(string); ok {
			blueprintContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("blueprint_content", client))
		}
		armConfig["json"] = blueprintContent

	case "repository":
		armConfig["configType"] = "git"
		armGitConfig := make(map[string]any)

		var integrationID int
		if v, ok := d.Get("integration_id").(int); ok {
			integrationID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("integration_id", client))
		}
		armGitConfig["integrationId"] = integrationID

		var repositoryID int
		if v, ok := d.Get("repository_id").(int); ok {
			repositoryID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("repository_id", client))
		}
		armGitConfig["repoId"] = repositoryID

		var versionRef string
		if v, ok := d.Get("version_ref").(string); ok {
			versionRef = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("version_ref", client))
		}
		armGitConfig["branch"] = versionRef

		var workingPath string
		if v, ok := d.Get("working_path").(string); ok {
			workingPath = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("working_path", client))
		}
		armGitConfig["path"] = workingPath
		armConfig["git"] = armGitConfig
	}

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
		return diag.FromErr(helpers.TypeAssertFailError("result", client))
	}

	if result.Blueprint == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Blueprint"))
	}
	blueprint := result.Blueprint
	d.SetId(convert.Int64ToString(blueprint.ID))

	diags = append(diags, resourceBlueprintArmAppRead(ctx, d, meta)...)

	return diags
}

func resourceBlueprintArmAppRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

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
		return diag.FromErr(helpers.TypeAssertFailError("name", client))
	}

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

	var armBlueprint ArmAppBlueprint
	json.Unmarshal(resp.Body, &armBlueprint)
	d.SetId(convert.IntToString(armBlueprint.Blueprint.ID))
	d.Set("name", armBlueprint.Blueprint.Name)
	d.Set("description", armBlueprint.Blueprint.Description)
	d.Set("category", armBlueprint.Blueprint.Category)
	d.Set("install_agent", armBlueprint.Blueprint.Config.Arm.InstallAgent)
	d.Set("cloud_init_enabled", armBlueprint.Blueprint.Config.Arm.CloudInitEnabled)
	d.Set("os_type", armBlueprint.Blueprint.Config.Arm.OsType)
	switch armBlueprint.Blueprint.Config.Arm.Configtype {
	case "json":
		d.Set("source_type", "json")
		d.Set("blueprint_content", armBlueprint.Blueprint.Config.Arm.JSON)
	case "git":
		d.Set("source_type", "repository")
		d.Set("working_path", armBlueprint.Blueprint.Config.Arm.Git.Path)
		d.Set("integration_id", armBlueprint.Blueprint.Config.Arm.Git.IntegrationId)
		d.Set("repository_id", armBlueprint.Blueprint.Config.Arm.Git.RepoId)
		d.Set("version_ref", armBlueprint.Blueprint.Config.Arm.Git.Branch)
	}

	return diags
}

func resourceBlueprintArmAppUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		return diag.FromErr(helpers.TypeAssertFailError("name", client))
	}

	blueprintType := "arm"

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", client))
	}

	var category string
	if v, ok := d.Get("category").(string); ok {
		category = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("category", client))
	}

	config := make(map[string]any)
	config["name"] = name
	config["description"] = description
	config["category"] = category
	config["type"] = "arm"

	armConfig := make(map[string]any)
	config["arm"] = armConfig

	var osType string
	if v, ok := d.Get("os_type").(string); ok {
		osType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("os_type", client))
	}
	armConfig["osType"] = osType

	var installAgent bool
	if v, ok := d.Get("install_agent").(bool); ok {
		installAgent = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("install_agent", client))
	}
	armConfig["installAgent"] = installAgent

	var cloudInitEnabled bool
	if v, ok := d.Get("cloud_init_enabled").(bool); ok {
		cloudInitEnabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloud_init_enabled", client))
	}
	armConfig["cloudInitEnabled"] = cloudInitEnabled

	var sourceType string
	if v, ok := d.Get("source_type").(string); ok {
		sourceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", client))
	}

	switch sourceType {
	case "json":
		armConfig["configType"] = "json"
		var blueprintContent string
		if v, ok := d.Get("blueprint_content").(string); ok {
			blueprintContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("blueprint_content", client))
		}
		armConfig["json"] = blueprintContent

	case "repository":
		armConfig["configType"] = "git"
		armGitConfig := make(map[string]any)

		var integrationID int
		if v, ok := d.Get("integration_id").(int); ok {
			integrationID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("integration_id", client))
		}
		armGitConfig["integrationId"] = integrationID

		var repositoryID int
		if v, ok := d.Get("repository_id").(int); ok {
			repositoryID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("repository_id", client))
		}
		armGitConfig["repoId"] = repositoryID

		var versionRef string
		if v, ok := d.Get("version_ref").(string); ok {
			versionRef = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("version_ref", client))
		}
		armGitConfig["branch"] = versionRef

		var workingPath string
		if v, ok := d.Get("working_path").(string); ok {
			workingPath = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("working_path", client))
		}
		armGitConfig["path"] = workingPath
		armConfig["git"] = armGitConfig
	}

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
		return diag.FromErr(helpers.TypeAssertFailError("result", client))
	}

	if result.Blueprint == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Blueprint"))
	}
	blueprint := result.Blueprint
	d.SetId(convert.Int64ToString(blueprint.ID))

	return resourceBlueprintArmAppRead(ctx, d, meta)
}

func resourceBlueprintArmAppDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

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

type ArmAppBlueprint struct {
	Blueprint struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Config      struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Arm         struct {
				Configtype       string `json:"configType"`
				OsType           string `json:"osType"`
				CloudInitEnabled bool   `json:"cloudInitEnabled"`
				InstallAgent     bool   `json:"installAgent"`
				JSON             string `json:"json"`
				Git              struct {
					Path          string `json:"path"`
					RepoId        int    `json:"repoId"`
					IntegrationId int    `json:"integrationId"`
					Branch        string `json:"branch"`
				} `json:"git"`
			} `json:"arm"`
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
