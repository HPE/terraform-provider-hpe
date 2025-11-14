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

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

const (
	blueprintTypeCloudFormation = "cloudFormation"
	sourceTypeJSON              = "json"
	sourceTypeYAML              = "yaml"
	sourceTypeGit               = "git"
)

func ResourceAppBlueprintCloudFormation() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus cloud formation app blueprint resource",
		CreateContext: resourceAppBlueprintCloudFormationCreate,
		ReadContext:   resourceAppBlueprintCloudFormationRead,
		UpdateContext: resourceAppBlueprintCloudFormationUpdate,
		DeleteContext: resourceAppBlueprintCloudFormationDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the cloud formation app blueprint",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the cloud formation app blueprint",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the cloud formation app blueprint",
				Optional:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the cloud formation app blueprint",
				Optional:    true,
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
			"capability_iam": {
				Type:        schema.TypeBool,
				Description: "Whether the iam capability is added to the cloud formation",
				Optional:    true,
			},
			"capability_named_iam": {
				Type:        schema.TypeBool,
				Description: "Whether the named iam capability is added to the cloud formation",
				Optional:    true,
			},
			"capability_auto_expand": {
				Type:        schema.TypeBool,
				Description: "Whether the auto expand capability is added to the cloud formation",
				Optional:    true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the cloud formation app blueprint (yaml, json, repository)",
				ValidateFunc: validation.StringInSlice([]string{sourceTypeYAML, sourceTypeJSON, sourceTypeRepository}, false),
				Required:     true,
			},
			"blueprint_content": {
				Type: schema.TypeString,
				Description: "The content of the cloud formation app blueprint. " +
					"Used when the yaml or json source types are specified",
				Optional: true,
				StateFunc: func(val any) string {
					var blueprintContent string
					if v, ok := val.(string); ok {
						blueprintContent = v
					}

					return strings.TrimSuffix(blueprintContent, "\n")
				},
			},
			"working_path": {
				Type:        schema.TypeString,
				Description: "The path of the cloud formation chart in the git repository",
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

func resourceAppBlueprintCloudFormationCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	blueprintType := blueprintTypeCloudFormation

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
	config["type"] = blueprintTypeCloudFormation

	cloudformationConfig := make(map[string]any)
	config[blueprintTypeCloudFormation] = cloudformationConfig

	var installAgent bool
	if v, ok := d.Get("install_agent").(bool); ok {
		installAgent = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("install_agent", d.Get("install_agent")))
	}
	cloudformationConfig["installAgent"] = installAgent

	var cloudInitEnabled bool
	if v, ok := d.Get("cloud_init_enabled").(bool); ok {
		cloudInitEnabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloud_init_enabled", d.Get("cloud_init_enabled")))
	}
	cloudformationConfig["cloudInitEnabled"] = cloudInitEnabled

	var capabilityIAM bool
	if v, ok := d.Get("capability_iam").(bool); ok {
		capabilityIAM = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("capability_iam", d.Get("capability_iam")))
	}
	cloudformationConfig["IAM"] = capabilityIAM

	var capabilityNamedIAM bool
	if v, ok := d.Get("capability_named_iam").(bool); ok {
		capabilityNamedIAM = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("capability_named_iam", d.Get("capability_named_iam")))
	}
	cloudformationConfig["CAPABILITY_NAMED_IAM"] = capabilityNamedIAM

	var capabilityAutoExpand bool
	if v, ok := d.Get("capability_auto_expand").(bool); ok {
		capabilityAutoExpand = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("capability_auto_expand", d.Get("capability_auto_expand")))
	}
	cloudformationConfig["CAPABILITY_AUTO_EXPAND"] = capabilityAutoExpand

	var sourceType any
	if v, ok := d.Get("source_type").(string); ok {
		sourceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", d.Get("source_type")))
	}

	switch sourceType {
	case sourceTypeJSON:
		cloudformationConfig["configType"] = sourceTypeJSON
		var blueprintContent string
		if v, ok := d.Get("blueprint_content").(string); ok {
			blueprintContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("blueprint_content", d.Get("blueprint_content")))
		}
		cloudformationConfig[sourceTypeJSON] = blueprintContent

	case sourceTypeYAML:
		cloudformationConfig["configType"] = sourceTypeYAML
		var blueprintContent string
		if v, ok := d.Get("blueprint_content").(string); ok {
			blueprintContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("blueprint_content", d.Get("blueprint_content")))
		}
		cloudformationConfig[sourceTypeYAML] = blueprintContent

	case sourceTypeRepository:
		cloudformationConfig["configType"] = sourceTypeGit
		cloudformationGitConfig := make(map[string]any)
		cloudformationGitConfig["integrationId"] = d.Get("integration_id")
		cloudformationGitConfig["repoId"] = d.Get("repository_id")

		var versionRef string
		if v, ok := d.Get("version_ref").(string); ok {
			versionRef = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("version_ref", d.Get("version_ref")))
		}
		cloudformationGitConfig["branch"] = versionRef

		var workingPath string
		if v, ok := d.Get("working_path").(string); ok {
			workingPath = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("working_path", d.Get("working_path")))
		}
		cloudformationGitConfig["path"] = workingPath
		cloudformationConfig["git"] = cloudformationGitConfig
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
		return diag.FromErr(helpers.TypeAssertFailError("CreateBlueprintResult", resp.Result))
	}

	if result.Blueprint == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("blueprint"))
	}
	blueprint := result.Blueprint
	d.SetId(convert.Int64ToString(blueprint.ID))

	diags = append(diags, resourceAppBlueprintCloudFormationRead(ctx, d, meta)...)

	return diags
}

func resourceAppBlueprintCloudFormationRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
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

	var cloudformationBlueprint BlueprintAppCloudFormation
	json.Unmarshal(resp.Body, &cloudformationBlueprint)
	d.SetId(convert.IntToString(cloudformationBlueprint.Blueprint.ID))
	d.Set("name", cloudformationBlueprint.Blueprint.Name)
	d.Set("description", cloudformationBlueprint.Blueprint.Description)
	d.Set("category", cloudformationBlueprint.Blueprint.Category)
	d.Set("install_agent", cloudformationBlueprint.Blueprint.Config.CloudFormation.InstallAgent)
	d.Set("cloud_init_enabled", cloudformationBlueprint.Blueprint.Config.CloudFormation.CloudInitEnabled)
	d.Set("capability_iam", cloudformationBlueprint.Blueprint.Config.CloudFormation.IAM)
	d.Set("capability_named_iam", cloudformationBlueprint.Blueprint.Config.CloudFormation.IAMNamed)
	d.Set("capability_auto_expand", cloudformationBlueprint.Blueprint.Config.CloudFormation.AutoExpand)

	switch cloudformationBlueprint.Blueprint.Config.CloudFormation.Configtype {
	case sourceTypeJSON:
		d.Set("source_type", sourceTypeJSON)
		d.Set("blueprint_content", cloudformationBlueprint.Blueprint.Config.CloudFormation.JSON)
	case sourceTypeYAML:
		d.Set("source_type", sourceTypeYAML)
		d.Set("blueprint_content", cloudformationBlueprint.Blueprint.Config.CloudFormation.YAML)
	case sourceTypeGit:
		d.Set("source_type", sourceTypeRepository)
		d.Set("working_path", cloudformationBlueprint.Blueprint.Config.CloudFormation.Git.Path)
		d.Set("integration_id", cloudformationBlueprint.Blueprint.Config.CloudFormation.Git.IntegrationId)
		d.Set("repository_id", cloudformationBlueprint.Blueprint.Config.CloudFormation.Git.RepoId)
		d.Set("version_ref", cloudformationBlueprint.Blueprint.Config.CloudFormation.Git.Branch)
	}

	return diags
}

func resourceAppBlueprintCloudFormationUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	blueprintType := blueprintTypeCloudFormation

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
	config["type"] = blueprintTypeCloudFormation

	cloudformationConfig := make(map[string]any)
	config[blueprintTypeCloudFormation] = cloudformationConfig

	var installAgent bool
	if v, ok := d.Get("install_agent").(bool); ok {
		installAgent = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("install_agent", d.Get("install_agent")))
	}
	cloudformationConfig["installAgent"] = installAgent

	var cloudInitEnabled bool
	if v, ok := d.Get("cloud_init_enabled").(bool); ok {
		cloudInitEnabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloud_init_enabled", d.Get("cloud_init_enabled")))
	}
	cloudformationConfig["cloudInitEnabled"] = cloudInitEnabled

	var capabilityIAM bool
	if v, ok := d.Get("capability_iam").(bool); ok {
		capabilityIAM = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("capability_iam", d.Get("capability_iam")))
	}
	cloudformationConfig["IAM"] = capabilityIAM

	var capabilityNamedIAM bool
	if v, ok := d.Get("capability_named_iam").(bool); ok {
		capabilityNamedIAM = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("capability_named_iam", d.Get("capability_named_iam")))
	}
	cloudformationConfig["CAPABILITY_NAMED_IAM"] = capabilityNamedIAM

	var capabilityAutoExpand bool
	if v, ok := d.Get("capability_auto_expand").(bool); ok {
		capabilityAutoExpand = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("capability_auto_expand", d.Get("capability_auto_expand")))
	}
	cloudformationConfig["CAPABILITY_AUTO_EXPAND"] = capabilityAutoExpand

	var sourceType any
	if v, ok := d.Get("source_type").(string); ok {
		sourceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", d.Get("source_type")))
	}

	switch sourceType {
	case sourceTypeJSON:
		cloudformationConfig["configType"] = sourceTypeJSON
		var blueprintContent string
		if v, ok := d.Get("blueprint_content").(string); ok {
			blueprintContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("blueprint_content", d.Get("blueprint_content")))
		}
		cloudformationConfig[sourceTypeJSON] = blueprintContent

	case sourceTypeYAML:
		cloudformationConfig["configType"] = sourceTypeYAML
		var blueprintContent string
		if v, ok := d.Get("blueprint_content").(string); ok {
			blueprintContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("blueprint_content", d.Get("blueprint_content")))
		}
		cloudformationConfig[sourceTypeYAML] = blueprintContent

	case sourceTypeRepository:
		cloudformationConfig["configType"] = sourceTypeGit
		cloudformationGitConfig := make(map[string]any)
		cloudformationGitConfig["integrationId"] = d.Get("integration_id")
		cloudformationGitConfig["repoId"] = d.Get("repository_id")

		var versionRef string
		if v, ok := d.Get("version_ref").(string); ok {
			versionRef = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("version_ref", d.Get("version_ref")))
		}
		cloudformationGitConfig["branch"] = versionRef

		var workingPath string
		if v, ok := d.Get("working_path").(string); ok {
			workingPath = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("working_path", d.Get("working_path")))
		}
		cloudformationGitConfig["path"] = workingPath
		cloudformationConfig["git"] = cloudformationGitConfig
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
		return diag.FromErr(helpers.TypeAssertFailError("UpdateBlueprintResult", resp.Result))
	}

	if result.Blueprint == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("blueprint"))
	}
	blueprint := result.Blueprint
	d.SetId(convert.Int64ToString(blueprint.ID))

	return resourceAppBlueprintCloudFormationRead(ctx, d, meta)
}

func resourceAppBlueprintCloudFormationDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

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

type BlueprintAppCloudFormation struct {
	Blueprint struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Config      struct {
			Name           string `json:"name"`
			Description    string `json:"description"`
			CloudFormation struct {
				Configtype       string `json:"configType"`
				CloudInitEnabled bool   `json:"cloudInitEnabled"`
				InstallAgent     bool   `json:"installAgent"`
				JSON             string `json:"json"`
				YAML             string `json:"yaml"`
				IAM              bool   `json:"IAM"`
				IAMNamed         bool   `json:"CAPABILITY_NAMED_IAM"`
				AutoExpand       bool   `json:"CAPABILITY_AUTO_EXPAND"`
				Git              struct {
					Path          string `json:"path"`
					RepoId        int    `json:"repoId"`
					IntegrationId int    `json:"integrationId"`
					Branch        string `json:"branch"`
				} `json:"git"`
			} `json:"cloudformation"`
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
