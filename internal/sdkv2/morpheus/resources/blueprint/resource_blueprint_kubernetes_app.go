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
	blueprintTypeKubernetes = "kubernetes"
	sourceTypeYaml          = "yaml"
	sourceTypeSpec          = "spec"
	sourceTypeRepository    = "repository"
	configTypeGit           = "git"
)

func ResourceBlueprintKubernetesAppBlueprint() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus kubernetes app blueprint resource",
		CreateContext: resourceBlueprintKubernetesAppBlueprintCreate,
		ReadContext:   resourceBlueprintKubernetesAppBlueprintRead,
		UpdateContext: resourceBlueprintKubernetesAppBlueprintUpdate,
		DeleteContext: resourceBlueprintKubernetesAppBlueprintDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the kubernetes app blueprint",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the kubernetes app blueprint",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the kubernetes app blueprint",
				Optional:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the kubernetes app blueprint",
				Optional:    true,
			},
			"source_type": {
				Type:        schema.TypeString,
				Description: "The source of the kubernetes app blueprint (yaml, spec or repository)",
				ValidateFunc: validation.StringInSlice(
					[]string{sourceTypeYaml, sourceTypeSpec, sourceTypeRepository},
					false,
				),
				Required: true,
			},
			"blueprint_content": {
				Type:        schema.TypeString,
				Description: "The content of the kubernetes app blueprint. Used when the yaml source type is specified",
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
				Description: "The path of the kubernetes app blueprint in the git repository",
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
			"spec_template_ids": {
				Type:        schema.TypeList,
				Description: "A list of kubernetes spec template ids associated with the app blueprint",
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceBlueprintKubernetesAppBlueprintCreate(
	ctx context.Context,
	d *schema.ResourceData,
	meta any,
) diag.Diagnostics {
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
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	blueprintType := blueprintTypeKubernetes

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
	config["type"] = blueprintTypeKubernetes

	kubernetesConfig := make(map[string]any)
	config[blueprintTypeKubernetes] = kubernetesConfig

	var sourceType string
	if v, ok := d.Get("source_type").(string); ok {
		sourceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", d.Get("source_type")))
	}

	switch sourceType {
	case sourceTypeYaml:
		kubernetesConfig["configType"] = sourceTypeYaml
		var blueprintContent string
		if v, ok := d.Get("blueprint_content").(string); ok {
			blueprintContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("blueprint_content", d.Get("blueprint_content")))
		}
		kubernetesConfig[sourceTypeYaml] = blueprintContent

	case sourceTypeSpec:
		kubernetesConfig["configType"] = sourceTypeSpec
		var specTemplates []map[string]any
		if d.Get("spec_template_ids") != nil {
			var specTemplateList []any
			if v, ok := d.Get("spec_template_ids").([]any); ok {
				specTemplateList = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("spec_template_ids", d.Get("spec_template_ids")))
			}
			if specTemplateList != nil {
				// iterate over the array of spec templates
				for i := 0; i < len(specTemplateList); i++ {
					row := make(map[string]any)
					row["id"] = specTemplateList[i]
					row["value"] = specTemplateList[i]
					specTemplates = append(specTemplates, row)
				}
			}
		}
		specConfig := make(map[string]any)
		config["config"] = specConfig
		specConfig["specs"] = specTemplates
	case sourceTypeRepository:
		kubernetesConfig["configType"] = configTypeGit
		kubernetesGitConfig := make(map[string]any)
		kubernetesGitConfig["integrationId"] = d.Get("integration_id")
		kubernetesGitConfig["repoId"] = d.Get("repository_id")
		var versionRef string
		if v, ok := d.Get("version_ref").(string); ok {
			versionRef = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("version_ref", d.Get("version_ref")))
		}
		var workingPath string
		if v, ok := d.Get("working_path").(string); ok {
			workingPath = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("working_path", d.Get("working_path")))
		}
		kubernetesGitConfig["branch"] = versionRef
		kubernetesGitConfig["path"] = workingPath
		kubernetesConfig[configTypeGit] = kubernetesGitConfig
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

	result := resp.Result.(*morpheus.CreateBlueprintResult)
	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("result"))
	}
	blueprint := result.Blueprint
	if blueprint == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("blueprint"))
	}
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(blueprint.ID))

	diags = append(diags, resourceBlueprintKubernetesAppBlueprintRead(ctx, d, meta)...)

	return diags
}

func resourceBlueprintKubernetesAppBlueprintRead(
	ctx context.Context,
	d *schema.ResourceData,
	meta any,
) diag.Diagnostics {
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
	var kubernetesBlueprint KubernetesAppBlueprint
	json.Unmarshal(resp.Body, &kubernetesBlueprint)
	d.SetId(convert.IntToString(kubernetesBlueprint.Blueprint.ID))
	d.Set("name", kubernetesBlueprint.Blueprint.Name)
	d.Set("description", kubernetesBlueprint.Blueprint.Description)
	d.Set("category", kubernetesBlueprint.Blueprint.Category)

	switch kubernetesBlueprint.Blueprint.Config.Kubernetes.Configtype {
	case sourceTypeYaml:
		d.Set("source_type", sourceTypeYaml)
		d.Set("blueprint_content", kubernetesBlueprint.Blueprint.Config.Kubernetes)
	case configTypeGit:
		d.Set("source_type", sourceTypeRepository)
		d.Set("working_path", kubernetesBlueprint.Blueprint.Config.Kubernetes.Git.Path)
		d.Set("integration_id", kubernetesBlueprint.Blueprint.Config.Kubernetes.Git.IntegrationId)
		d.Set("repository_id", kubernetesBlueprint.Blueprint.Config.Kubernetes.Git.RepoId)
		d.Set("version_ref", kubernetesBlueprint.Blueprint.Config.Kubernetes.Git.Branch)
	case sourceTypeSpec:
		d.Set("source_type", sourceTypeSpec)
		// spec templates
		var specTemplates []int64
		if kubernetesBlueprint.Blueprint.Config.Config.Specs != nil {
			// iterate over the array of tasks
			for i := 0; i < len(kubernetesBlueprint.Blueprint.Config.Config.Specs); i++ {
				specTemplate := kubernetesBlueprint.Blueprint.Config.Config.Specs[i]
				specTemplates = append(specTemplates, int64(specTemplate.ID))
			}
		}
		d.Set("spec_templates_ids", specTemplates)
	}

	return diags
}

func resourceBlueprintKubernetesAppBlueprintUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	meta any,
) diag.Diagnostics {
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
	blueprintType := blueprintTypeKubernetes
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
	config["type"] = blueprintTypeKubernetes

	kubernetesConfig := make(map[string]any)
	config[blueprintTypeKubernetes] = kubernetesConfig

	var sourceType string
	if v, ok := d.Get("source_type").(string); ok {
		sourceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", d.Get("source_type")))
	}

	switch sourceType {
	case sourceTypeYaml:
		kubernetesConfig["configType"] = sourceTypeYaml
		var blueprintContent string
		if v, ok := d.Get("blueprint_content").(string); ok {
			blueprintContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("blueprint_content", d.Get("blueprint_content")))
		}
		kubernetesConfig[sourceTypeYaml] = blueprintContent

	case sourceTypeSpec:
		kubernetesConfig["configType"] = sourceTypeSpec
		var specTemplates []map[string]any
		if d.Get("spec_template_ids") != nil {
			var specTemplateList []any
			if v, ok := d.Get("spec_template_ids").([]any); ok {
				specTemplateList = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("spec_template_ids", d.Get("spec_template_ids")))
			}
			if specTemplateList != nil {
				// iterate over the array of spec templates
				for i := 0; i < len(specTemplateList); i++ {
					row := make(map[string]any)
					row["id"] = specTemplateList[i]
					row["value"] = specTemplateList[i]
					specTemplates = append(specTemplates, row)
				}
			}
		}
		specConfig := make(map[string]any)
		config["config"] = specConfig
		specConfig["specs"] = specTemplates
	case sourceTypeRepository:
		kubernetesConfig["configType"] = configTypeGit
		kubernetesGitConfig := make(map[string]any)
		kubernetesGitConfig["integrationId"] = d.Get("integration_id")
		kubernetesGitConfig["repoId"] = d.Get("repository_id")
		var versionRef string
		if v, ok := d.Get("version_ref").(string); ok {
			versionRef = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("version_ref", d.Get("version_ref")))
		}
		var workingPath string
		if v, ok := d.Get("working_path").(string); ok {
			workingPath = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("working_path", d.Get("working_path")))
		}
		kubernetesGitConfig["branch"] = versionRef
		kubernetesGitConfig["path"] = workingPath
		kubernetesConfig[configTypeGit] = kubernetesGitConfig
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
	log.Printf("API REQUEST: %s", req)
	resp, err := client.UpdateBlueprint(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	result := resp.Result.(*morpheus.UpdateBlueprintResult)
	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("result"))
	}
	blueprint := result.Blueprint
	if blueprint == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("blueprint"))
	}
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(blueprint.ID))

	return resourceBlueprintKubernetesAppBlueprintRead(ctx, d, meta)
}

func resourceBlueprintKubernetesAppBlueprintDelete(
	ctx context.Context,
	d *schema.ResourceData,
	meta any,
) diag.Diagnostics {
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

type KubernetesAppBlueprint struct {
	Blueprint struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Config      struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Kubernetes  struct {
				Configtype string `json:"configType"`
				Git        struct {
					Path          string `json:"path"`
					RepoId        int    `json:"repoId"`
					IntegrationId int    `json:"integrationId"`
					Branch        string `json:"branch"`
				} `json:"git"`
			} `json:"kubernetes"`
			Config struct {
				Specs []struct {
					ID    int    `json:"id"`
					Value string `json:"value"`
					Name  string `json:"name"`
				} `json:"specs"`
			} `json:"config"`
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
