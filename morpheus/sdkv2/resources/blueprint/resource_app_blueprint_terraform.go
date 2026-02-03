// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"context"
	"encoding/json"
	"log"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func ResourceAppBlueprintTerraform() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus terraform app blueprint resource",
		CreateContext: resourceAppBlueprintTerraformCreate,
		ReadContext:   resourceAppBlueprintTerraformRead,
		UpdateContext: resourceAppBlueprintTerraformUpdate,
		DeleteContext: resourceAppBlueprintTerraformDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the terraform app blueprint",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the terraform app blueprint",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the terraform app blueprint",
				Optional:    true,
				Computed:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the terraform app blueprint",
				Optional:    true,
				Computed:    true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the terraform app blueprint (hcl, json, spec or repository)",
				ValidateFunc: validation.StringInSlice([]string{"hcl", "json", "spec", "repository"}, false),
				Required:     true,
			},
			"blueprint_content": {
				Type:        schema.TypeString,
				Description: "The content of the terraform app blueprint. Used when the hcl or json source types are specified",
				Optional:    true,
				Computed:    true,
			},
			"working_path": {
				Type:          schema.TypeString,
				Description:   "The path of the terraform code in the git repository",
				Optional:      true,
				ConflictsWith: []string{"blueprint_content"},
			},
			"integration_id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the git integration",
				Optional:      true,
				ConflictsWith: []string{"blueprint_content"},
			},
			"repository_id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the git repository",
				Optional:      true,
				ConflictsWith: []string{"blueprint_content"},
				RequiredWith:  []string{"integration_id"},
			},
			"version_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the repository to pull (main, master, etc.)",
				Optional:    true,
				Computed:    true,
			},
			"spec_template_ids": {
				Type:        schema.TypeList,
				Description: "A list of terraform spec template ids associated with the app blueprint",
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Computed:    true,
			},
			"terraform_version": {
				Type:        schema.TypeString,
				Description: "The terraform version associated with the app blueprint",
				Optional:    true,
				Computed:    true,
			},
			"terraform_options": {
				Type:        schema.TypeString,
				Description: "The additional terraform options to add to the app blueprint",
				Optional:    true,
				Computed:    true,
			},
			"tfvar_secret": {
				Type:        schema.TypeString,
				Description: "The name of the tfvar cypher secret to associate with the app blueprint",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

//nolint:goconst
func resourceAppBlueprintTerraformCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	blueprintType := "terraform"

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

	terraformConfig := make(map[string]any)

	var tfVersion string
	if v, ok := d.Get("terraform_version").(string); ok {
		tfVersion = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("terraform_version", d.Get("terraform_version")))
	}
	terraformConfig["tfVersion"] = tfVersion

	var tfOptions string
	if v, ok := d.Get("terraform_options").(string); ok {
		tfOptions = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("terraform_options", d.Get("terraform_options")))
	}
	terraformConfig["commandOptions"] = tfOptions

	var tfvarSecret string
	if v, ok := d.Get("tfvar_secret").(string); ok {
		tfvarSecret = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("tfvar_secret", d.Get("tfvar_secret")))
	}
	terraformConfig["tfvarSecret"] = tfvarSecret

	config["terraform"] = terraformConfig

	var sourceType string
	if v, ok := d.Get("source_type").(string); ok {
		sourceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", d.Get("source_type")))
	}

	switch sourceType {
	case "hcl":
		terraformConfig["configType"] = "tf"

		var blueprintContent string
		if v, ok := d.Get("blueprint_content").(string); ok {
			blueprintContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("blueprint_content", d.Get("blueprint_content")))
		}
		terraformConfig["tf"] = blueprintContent
	case "json":
		terraformConfig["configType"] = "json"

		var blueprintContent string
		if v, ok := d.Get("blueprint_content").(string); ok {
			blueprintContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("blueprint_content", d.Get("blueprint_content")))
		}
		terraformConfig["json"] = blueprintContent
	case "repository":
		terraformConfig["configType"] = "git"
		terraformGitConfig := make(map[string]any)
		terraformGitConfig["integrationId"] = d.Get("integration_id")
		terraformGitConfig["repoId"] = d.Get("repository_id")

		var versionRef string
		if v, ok := d.Get("version_ref").(string); ok {
			versionRef = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("version_ref", d.Get("version_ref")))
		}
		terraformGitConfig["branch"] = versionRef

		var workingPath string
		if v, ok := d.Get("working_path").(string); ok {
			workingPath = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("working_path", d.Get("working_path")))
		}
		terraformGitConfig["path"] = workingPath

		terraformConfig["git"] = terraformGitConfig
	case "spec":
		terraformConfig["configType"] = "spec"

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
		return diag.FromErr(helpers.NotFoundInResponseError("Blueprint"))
	}

	blueprint := result.Blueprint
	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(blueprint.ID))

	diags = append(diags, resourceAppBlueprintTerraformRead(ctx, d, meta)...)

	return diags
}

func resourceAppBlueprintTerraformRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
	var terraformBlueprint AppBlueprintTerraform
	json.Unmarshal(resp.Body, &terraformBlueprint)
	d.SetId(convert.IntToString(terraformBlueprint.Blueprint.ID))
	d.Set("name", terraformBlueprint.Blueprint.Name)
	d.Set("description", terraformBlueprint.Blueprint.Description)
	d.Set("category", terraformBlueprint.Blueprint.Category)
	d.Set("terraform_version", terraformBlueprint.Blueprint.Config.Terraform.Tfversion)
	d.Set("terraform_options", terraformBlueprint.Blueprint.Config.Terraform.Commandoptions)
	d.Set("tfvar_secret", terraformBlueprint.Blueprint.Config.Terraform.Tfvarsecret)

	switch terraformBlueprint.Blueprint.Config.Terraform.Configtype {
	case "tf":
		d.Set("source_type", "hcl")
		d.Set("blueprint_content", terraformBlueprint.Blueprint.Config.Terraform.Tf)
	case "json":
		d.Set("source_type", "json")
		d.Set("blueprint_content", terraformBlueprint.Blueprint.Config.Terraform.JSON)
	case "git":
		d.Set("source_type", "repository")
		d.Set("working_path", terraformBlueprint.Blueprint.Config.Terraform.Git.Path)
		d.Set("integration_id", terraformBlueprint.Blueprint.Config.Terraform.Git.IntegrationId)
		d.Set("repository_id", terraformBlueprint.Blueprint.Config.Terraform.Git.RepoId)
		d.Set("version_ref", terraformBlueprint.Blueprint.Config.Terraform.Git.Branch)
	case "spec":
		d.Set("source_type", "spec")
		var specTemplates []int64
		if terraformBlueprint.Blueprint.Config.Config.Specs != nil {
			for i := 0; i < len(terraformBlueprint.Blueprint.Config.Config.Specs); i++ {
				specTemplate := terraformBlueprint.Blueprint.Config.Config.Specs[i]
				specTemplates = append(specTemplates, int64(specTemplate.ID))
			}
		}
		d.Set("spec_template_ids", specTemplates)
	}

	return diags
}

func resourceAppBlueprintTerraformUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	blueprintType := "terraform"

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
	config["type"] = "terraform"

	terraformConfig := make(map[string]any)

	var tfVersion string
	if v, ok := d.Get("terraform_version").(string); ok {
		tfVersion = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("terraform_version", d.Get("terraform_version")))
	}
	terraformConfig["tfVersion"] = tfVersion

	var tfOptions string
	if v, ok := d.Get("terraform_options").(string); ok {
		tfOptions = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("terraform_options", d.Get("terraform_options")))
	}
	terraformConfig["commandOptions"] = tfOptions

	var tfvarSecret string
	if v, ok := d.Get("tfvar_secret").(string); ok {
		tfvarSecret = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("tfvar_secret", d.Get("tfvar_secret")))
	}
	terraformConfig["tfvarSecret"] = tfvarSecret

	config["terraform"] = terraformConfig

	var sourceType string
	if v, ok := d.Get("source_type").(string); ok {
		sourceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", d.Get("source_type")))
	}

	switch sourceType {
	case "hcl":
		terraformConfig["configType"] = "tf"

		var blueprintContent string
		if v, ok := d.Get("blueprint_content").(string); ok {
			blueprintContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("blueprint_content", d.Get("blueprint_content")))
		}
		terraformConfig["tf"] = blueprintContent
	case "json":
		terraformConfig["configType"] = "json"

		var blueprintContent string
		if v, ok := d.Get("blueprint_content").(string); ok {
			blueprintContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("blueprint_content", d.Get("blueprint_content")))
		}
		terraformConfig["json"] = blueprintContent
	case "repository":
		terraformConfig["configType"] = "git"
		terraformGitConfig := make(map[string]any)
		terraformGitConfig["integrationId"] = d.Get("integration_id")
		terraformGitConfig["repoId"] = d.Get("repository_id")

		var versionRef string
		if v, ok := d.Get("version_ref").(string); ok {
			versionRef = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("version_ref", d.Get("version_ref")))
		}
		terraformGitConfig["branch"] = versionRef

		var workingPath string
		if v, ok := d.Get("working_path").(string); ok {
			workingPath = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("working_path", d.Get("working_path")))
		}
		terraformGitConfig["path"] = workingPath

		terraformConfig["git"] = terraformGitConfig
	case "spec":
		terraformConfig["configType"] = "spec"

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
		return diag.FromErr(helpers.NotFoundInResponseError("Blueprint"))
	}

	blueprint := result.Blueprint
	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(blueprint.ID))

	return resourceAppBlueprintTerraformRead(ctx, d, meta)
}

func resourceAppBlueprintTerraformDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

type AppBlueprintTerraform struct {
	Blueprint struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Config      struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Terraform   struct {
				Tfversion      string `json:"tfVersion"`
				Tf             string `json:"tf"`
				Tfvarsecret    string `json:"tfvarSecret"`
				Commandoptions string `json:"commandOptions"`
				Configtype     string `json:"configType"`
				JSON           string `json:"json"`
				Git            struct {
					Path          string `json:"path"`
					RepoId        int    `json:"repoId"`
					IntegrationId int    `json:"integrationId"`
					Branch        string `json:"branch"`
				} `json:"git"`
			} `json:"terraform"`
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
