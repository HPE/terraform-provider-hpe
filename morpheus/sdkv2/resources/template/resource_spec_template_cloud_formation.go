// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func ResourceSpecTemplateCloudFormation() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus cloud formation spec template resource",
		CreateContext: resourceSpecTemplateCloudFormationCreate,
		ReadContext:   resourceSpecTemplateCloudFormationRead,
		UpdateContext: resourceSpecTemplateCloudFormationUpdate,
		DeleteContext: resourceSpecTemplateCloudFormationDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the cloud formation spec template",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the cloud formation spec template",
				Required:    true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the cloud formation spec template (local, url or repository)",
				ValidateFunc: validation.StringInSlice([]string{sourceTypeLocal, sourceTypeURL, sourceTypeRepository}, false),
				Required:     true,
			},
			"spec_content": {
				Type:        schema.TypeString,
				Description: "The content of the cloud formation spec template. Used when the local source type is specified",
				Optional:    true,
				StateFunc: func(val any) string {
					if v, ok := val.(string); ok {
						return strings.TrimSpace(v)
					}

					return ""
				},
				DiffSuppressFunc: helpers.SuppressEquivalentJSONDiffs,
			},
			"spec_path": {
				Type:        schema.TypeString,
				Description: "The path of the cloud formation spec template, either the url or the path in the repository",
				Optional:    true,
			},
			"repository_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the git repository integration",
				Optional:    true,
			},
			"version_ref": {
				Type:        schema.TypeString,
				Description: "The git reference of the repository to pull (main, master, etc.)",
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
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSpecTemplateCloudFormationCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var sourceType string
	if v, ok := d.Get("source_type").(string); ok {
		sourceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", d.Get("source_type")))
	}

	sourceOptions := make(map[string]any)
	sourceOptions["sourceType"] = sourceType

	specTemplateType := make(map[string]any)
	specTemplateType["code"] = "cloudFormation"

	config := make(map[string]any)

	cloudformationConfig := make(map[string]any)
	config["cloudformation"] = cloudformationConfig

	var capabilityIam bool
	if v, ok := d.Get("capability_iam").(bool); ok {
		capabilityIam = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("capability_iam", d.Get("capability_iam")))
	}
	if capabilityIam {
		cloudformationConfig["IAM"] = "on"
	}

	var capabilityNamedIam bool
	if v, ok := d.Get("capability_named_iam").(bool); ok {
		capabilityNamedIam = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("capability_named_iam", d.Get("capability_named_iam")))
	}
	if capabilityNamedIam {
		cloudformationConfig["CAPABILITY_NAMED_IAM"] = "on"
	}

	var capabilityAutoExpand bool
	if v, ok := d.Get("capability_auto_expand").(bool); ok {
		capabilityAutoExpand = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("capability_auto_expand", d.Get("capability_auto_expand")))
	}
	if capabilityAutoExpand {
		cloudformationConfig["CAPABILITY_AUTO_EXPAND"] = "on"
	}

	switch sourceType {
	case sourceTypeLocal:
		var specContent string
		if v, ok := d.Get("spec_content").(string); ok {
			specContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("spec_content", d.Get("spec_content")))
		}
		sourceOptions["content"] = specContent
		sourceOptions["contentPath"] = d.Get("spec_path")
	case sourceTypeURL:
		var specContent string
		if v, ok := d.Get("spec_content").(string); ok {
			specContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("spec_content", d.Get("spec_content")))
		}
		sourceOptions["content"] = specContent
		sourceOptions["contentPath"] = d.Get("spec_path")
	case sourceTypeRepository:
		var repositoryID int
		if v, ok := d.Get("repository_id").(int); ok {
			repositoryID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("repository_id", d.Get("repository_id")))
		}
		sourceOptions["contentPath"] = d.Get("spec_path")
		sourceOptions["contentRef"] = d.Get("version_ref")
		sourceOptions["repository"] = map[string]any{
			"id": repositoryID,
		}
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"specTemplate": map[string]any{
				"name":   name,
				"file":   sourceOptions,
				"type":   specTemplateType,
				"config": config,
			},
		},
	}
	resp, err := client.CreateSpecTemplate(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.CreateSpecTemplateResult
	if v, ok := resp.Result.(*morpheus.CreateSpecTemplateResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("CreateSpecTemplateResult", resp.Result))
	}

	if result.SpecTemplate == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("SpecTemplate"))
	}

	d.SetId(convert.Int64ToString(result.SpecTemplate.ID))

	diags = append(diags, resourceSpecTemplateCloudFormationRead(ctx, d, meta)...)

	return diags
}

func resourceSpecTemplateCloudFormationRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindSpecTemplateByName(name)
	} else if id != "" {
		resp, err = client.GetSpecTemplate(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Spec template cannot be read without name or id")
	}

	if err != nil {
		// 404 is ok?
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)
			log.Printf("Forcing recreation of resource")
			d.SetId("")

			return diags
		}
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	var cloudFormationSpecTemplate CloudFormationSpecTemplate
	if resp.Body != nil {
		json.Unmarshal(resp.Body, &cloudFormationSpecTemplate)
	}
	d.SetId(convert.IntToString(cloudFormationSpecTemplate.Spectemplate.ID))
	d.Set("name", cloudFormationSpecTemplate.Spectemplate.Name)
	d.Set("source_type", cloudFormationSpecTemplate.Spectemplate.File.Sourcetype)

	if cloudFormationSpecTemplate.Spectemplate.Config.CloudFormation.Iam == "on" {
		d.Set("capability_iam", true)
	} else {
		d.Set("capability_iam", false)
	}

	if cloudFormationSpecTemplate.Spectemplate.Config.CloudFormation.CapabilityNamedIam == "on" {
		d.Set("capability_named_iam", true)
	} else {
		d.Set("capability_named_iam", false)
	}

	if cloudFormationSpecTemplate.Spectemplate.Config.CloudFormation.CapabilityAutoExpand == "on" {
		d.Set("capability_auto_expand", true)
	} else {
		d.Set("capability_auto_expand", false)
	}

	switch cloudFormationSpecTemplate.Spectemplate.File.Sourcetype {
	case sourceTypeLocal:
		d.Set("source_type", sourceTypeLocal)
		d.Set("spec_content", cloudFormationSpecTemplate.Spectemplate.File.Content)
	case sourceTypeURL:
		d.Set("source_type", sourceTypeURL)
		d.Set("spec_path", cloudFormationSpecTemplate.Spectemplate.File.Contentpath)
	case sourceTypeGit:
		d.Set("source_type", sourceTypeRepository)
		d.Set("spec_path", cloudFormationSpecTemplate.Spectemplate.File.Contentpath)
		d.Set("repository_id", cloudFormationSpecTemplate.Spectemplate.File.Repository.ID)
		d.Set("version_ref", cloudFormationSpecTemplate.Spectemplate.File.Contentref)
	}

	return diags
}

func resourceSpecTemplateCloudFormationUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var sourceType string
	if v, ok := d.Get("source_type").(string); ok {
		sourceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", d.Get("source_type")))
	}

	sourceOptions := make(map[string]any)
	sourceOptions["sourceType"] = sourceType

	specTemplateType := make(map[string]any)
	specTemplateType["code"] = "cloudFormation"

	config := make(map[string]any)
	cloudformationConfig := make(map[string]any)
	config["cloudformation"] = cloudformationConfig

	var capabilityIam bool
	if v, ok := d.Get("capability_iam").(bool); ok {
		capabilityIam = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("capability_iam", d.Get("capability_iam")))
	}
	if capabilityIam {
		cloudformationConfig["IAM"] = "on"
	}

	var capabilityNamedIam bool
	if v, ok := d.Get("capability_named_iam").(bool); ok {
		capabilityNamedIam = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("capability_named_iam", d.Get("capability_named_iam")))
	}
	if capabilityNamedIam {
		cloudformationConfig["CAPABILITY_NAMED_IAM"] = "on"
	}

	var capabilityAutoExpand bool
	if v, ok := d.Get("capability_auto_expand").(bool); ok {
		capabilityAutoExpand = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("capability_auto_expand", d.Get("capability_auto_expand")))
	}
	if capabilityAutoExpand {
		cloudformationConfig["CAPABILITY_AUTO_EXPAND"] = "on"
	}

	switch sourceType {
	case sourceTypeLocal:
		var specContent string
		if v, ok := d.Get("spec_content").(string); ok {
			specContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("spec_content", d.Get("spec_content")))
		}
		sourceOptions["content"] = specContent
		sourceOptions["contentPath"] = d.Get("spec_path")
	case sourceTypeURL:
		var specContent string
		if v, ok := d.Get("spec_content").(string); ok {
			specContent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("spec_content", d.Get("spec_content")))
		}
		sourceOptions["content"] = specContent
		sourceOptions["contentPath"] = d.Get("spec_path")
	case sourceTypeRepository:
		var repositoryID int
		if v, ok := d.Get("repository_id").(int); ok {
			repositoryID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("repository_id", d.Get("repository_id")))
		}
		sourceOptions["contentPath"] = d.Get("spec_path")
		sourceOptions["contentRef"] = d.Get("version_ref")
		sourceOptions["repository"] = map[string]any{
			"id": repositoryID,
		}
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"specTemplate": map[string]any{
				"name":   name,
				"file":   sourceOptions,
				"type":   specTemplateType,
				"config": config,
			},
		},
	}
	resp, err := client.UpdateSpecTemplate(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.UpdateSpecTemplateResult
	if v, ok := resp.Result.(*morpheus.UpdateSpecTemplateResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("UpdateSpecTemplateResult", resp.Result))
	}

	if result.SpecTemplate == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("SpecTemplate"))
	}

	d.SetId(convert.Int64ToString(result.SpecTemplate.ID))

	return resourceSpecTemplateCloudFormationRead(ctx, d, meta)
}

func resourceSpecTemplateCloudFormationDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteSpecTemplate(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return nil
		}
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}

type CloudFormationSpecTemplate struct {
	Spectemplate struct {
		ID      int `json:"id"`
		Account struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"account"`
		Name string `json:"name"`
		Code any    `json:"code"`
		Type struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"type"`
		Externalid   any `json:"externalId"`
		Externaltype any `json:"externalType"`
		Deploymentid any `json:"deploymentId"`
		Status       any `json:"status"`
		File         struct {
			ID          int    `json:"id"`
			Sourcetype  string `json:"sourceType"`
			Contentref  any    `json:"contentRef"`
			Contentpath any    `json:"contentPath"`
			Repository  struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"repository"`
			Content string `json:"content"`
		} `json:"file"`
		Config struct {
			CloudFormation struct {
				Iam                  string `json:"IAM"`
				CapabilityNamedIam   string `json:"CAPABILITY_NAMED_IAM"`
				CapabilityAutoExpand string `json:"CAPABILITY_AUTO_EXPAND"`
			} `json:"cloudformation"`
		} `json:"config"`
		Createdby   string    `json:"createdBy"`
		Updatedby   any       `json:"updatedBy"`
		Datecreated time.Time `json:"dateCreated"`
		Lastupdated time.Time `json:"lastUpdated"`
	} `json:"specTemplate"`
}
