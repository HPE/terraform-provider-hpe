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

func ResourceSpecTemplateARM() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus ARM spec template resource",
		CreateContext: resourceSpecTemplateARMCreate,
		ReadContext:   resourceSpecTemplateARMRead,
		UpdateContext: resourceSpecTemplateARMUpdate,
		DeleteContext: resourceSpecTemplateARMDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the arm spec template",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the arm spec template",
				Required:    true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Description:  "The source of the arm spec template (local, url or repository)",
				ValidateFunc: validation.StringInSlice([]string{sourceTypeLocal, sourceTypeURL, sourceTypeRepository}, false),
				Required:     true,
			},
			"spec_content": {
				Type:        schema.TypeString,
				Description: "The content of the arm spec template. Used when the local source type is specified",
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
				Description: "The path of the arm spec template, either the url or the path in the repository",
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
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSpecTemplateARMCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	sourceOptions := make(map[string]any)
	sourceOptions["sourceType"] = d.Get("source_type")

	specTemplateType := make(map[string]any)
	specTemplateType["code"] = "arm"

	var sourceType string
	if v, ok := d.Get("source_type").(string); ok {
		sourceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", d.Get("source_type")))
	}

	switch sourceType {
	case sourceTypeLocal:
		sourceOptions["content"] = d.Get("spec_content")
	case sourceTypeURL:
		sourceOptions["contentPath"] = d.Get("spec_path")
	case sourceTypeRepository:
		sourceOptions["contentPath"] = d.Get("spec_path")
		sourceOptions["contentRef"] = d.Get("version_ref")
		sourceOptions["repository"] = map[string]any{
			"id": d.Get("repository_id"),
		}
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"specTemplate": map[string]any{
				"name": name,
				"file": sourceOptions,
				"type": specTemplateType,
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
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}

	if result.SpecTemplate == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("SpecTemplate"))
	}

	specTemplate := result.SpecTemplate
	d.SetId(convert.Int64ToString(specTemplate.ID))

	diags = append(diags, resourceSpecTemplateARMRead(ctx, d, meta)...)

	return diags
}

func resourceSpecTemplateARMRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)

	// store resource data
	var armSpecTemplate ARMSpecTemplate
	if err := json.Unmarshal(resp.Body, &armSpecTemplate); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(convert.IntToString(armSpecTemplate.Spectemplate.ID))
	d.Set("name", armSpecTemplate.Spectemplate.Name)
	d.Set("source_type", armSpecTemplate.Spectemplate.File.Sourcetype)

	switch armSpecTemplate.Spectemplate.File.Sourcetype {
	case sourceTypeLocal:
		d.Set("source_type", sourceTypeLocal)
		d.Set("spec_content", armSpecTemplate.Spectemplate.File.Content)
	case sourceTypeURL:
		d.Set("source_type", sourceTypeURL)
		d.Set("spec_path", armSpecTemplate.Spectemplate.File.Contentpath)
	case sourceTypeGit:
		d.Set("source_type", sourceTypeRepository)
		d.Set("spec_path", armSpecTemplate.Spectemplate.File.Contentpath)
		d.Set("repository_id", armSpecTemplate.Spectemplate.File.Repository.ID)
		d.Set("version_ref", armSpecTemplate.Spectemplate.File.Contentref)
	}

	return diags
}

func resourceSpecTemplateARMUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	sourceOptions := make(map[string]any)
	sourceOptions["sourceType"] = d.Get("source_type")

	specTemplateType := make(map[string]any)
	specTemplateType["code"] = "arm"

	var sourceType string
	if v, ok := d.Get("source_type").(string); ok {
		sourceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("source_type", d.Get("source_type")))
	}

	switch sourceType {
	case sourceTypeLocal:
		sourceOptions["content"] = d.Get("spec_content")
	case sourceTypeURL:
		sourceOptions["contentPath"] = d.Get("spec_path")
	case sourceTypeRepository:
		sourceOptions["contentPath"] = d.Get("spec_path")
		sourceOptions["contentRef"] = d.Get("version_ref")
		sourceOptions["repository"] = map[string]any{
			"id": d.Get("repository_id"),
		}
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"specTemplate": map[string]any{
				"name": name,
				"file": sourceOptions,
				"type": specTemplateType,
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
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}

	if result.SpecTemplate == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("SpecTemplate"))
	}

	specTemplate := result.SpecTemplate
	d.SetId(convert.Int64ToString(specTemplate.ID))

	return resourceSpecTemplateARMRead(ctx, d, meta)
}

func resourceSpecTemplateARMDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		} else {
			log.Printf("API FAILURE: %s - %s", resp, err)

			return diag.FromErr(err)
		}
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}

type ARMSpecTemplate struct {
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
		Config      struct{}  `json:"config"`
		Createdby   string    `json:"createdBy"`
		Updatedby   any       `json:"updatedBy"`
		Datecreated time.Time `json:"dateCreated"`
		Lastupdated time.Time `json:"lastUpdated"`
	} `json:"specTemplate"`
}
