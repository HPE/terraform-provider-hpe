// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func ResourceFileTemplate() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus file template resource",
		CreateContext: resourceFileTemplateCreate,
		ReadContext:   resourceFileTemplateRead,
		UpdateContext: resourceFileTemplateUpdate,
		DeleteContext: resourceFileTemplateDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the file template",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the file template",
				Required:    true,
			},
			"labels": {
				Type: schema.TypeSet,
				Description: "The organization labels associated with the file template " +
					"(Only supported on Morpheus 5.5.3 or higher)",
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"file_name": {
				Type:        schema.TypeString,
				Description: "The name of the file deployed by the file template",
				Required:    true,
			},
			"file_path": {
				Type:        schema.TypeString,
				Description: "The system path of the file deployed by the file template",
				Optional:    true,
			},
			"phase": {
				Type: schema.TypeString,
				Description: "The phase that the file template should be run during " +
					"(preProvision, provision, postProvision, preDeploy, deploy)",
				ValidateFunc: validation.StringInSlice(
					[]string{"preProvision", "provision", "postProvision", "preDeploy", "deploy"}, false),
				Required: true,
			},
			"file_content": {
				Type:        schema.TypeString,
				Description: "The content of the file template",
				Optional:    true,
				StateFunc: func(v any) string {
					var payload string
					if str, ok := v.(string); ok {
						payload = str
					}
					payload = strings.TrimSuffix(payload, "\n")

					return payload
				},
			},
			"file_owner": {
				Type:        schema.TypeString,
				Description: "The file template file owner",
				Optional:    true,
			},
			"setting_name": {
				Type:        schema.TypeString,
				Description: "The file template setting name",
				Optional:    true,
			},
			"setting_category": {
				Type:        schema.TypeString,
				Description: "The file template setting category",
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceFileTemplateCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		var labelSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			labelSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}
		labelList := labelSet.List()
		for _, s := range labelList {
			var labelStr string
			if v, ok := s.(string); ok {
				labelStr = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("label", s))
			}
			labelsPayload = append(labelsPayload, labelStr)
		}
	}

	var fileName string
	if v, ok := d.Get("file_name").(string); ok {
		fileName = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("file_name", d.Get("file_name")))
	}

	var filePath string
	if v, ok := d.Get("file_path").(string); ok {
		filePath = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("file_path", d.Get("file_path")))
	}

	var phase string
	if v, ok := d.Get("phase").(string); ok {
		phase = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("phase", d.Get("phase")))
	}

	var fileContent string
	if v, ok := d.Get("file_content").(string); ok {
		fileContent = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("file_content", d.Get("file_content")))
	}

	var fileOwner string
	if v, ok := d.Get("file_owner").(string); ok {
		fileOwner = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("file_owner", d.Get("file_owner")))
	}

	var settingName string
	if v, ok := d.Get("setting_name").(string); ok {
		settingName = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("setting_name", d.Get("setting_name")))
	}

	var settingCategory string
	if v, ok := d.Get("setting_category").(string); ok {
		settingCategory = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("setting_category", d.Get("setting_category")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"containerTemplate": map[string]any{
				"name":            name,
				"labels":          labelsPayload,
				"fileName":        fileName,
				"filePath":        filePath,
				"templatePhase":   phase,
				"template":        fileContent,
				"fileOwner":       fileOwner,
				"settingName":     settingName,
				"settingCategory": settingCategory,
			},
		},
	}

	resp, err := client.CreateFileTemplate(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.CreateFileTemplateResult
	if v, ok := resp.Result.(*morpheus.CreateFileTemplateResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("CreateFileTemplateResult", resp.Result))
	}

	if result.FileTemplate == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("FileTemplate"))
	}

	fileTemplate := result.FileTemplate
	d.SetId(convert.Int64ToString(fileTemplate.ID))

	diags = append(diags, resourceFileTemplateRead(ctx, d, meta)...)

	return diags
}

func resourceFileTemplateRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindFileTemplateByName(name)
	} else if id != "" {
		resp, err = client.GetFileTemplate(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("File template cannot be read without name or id")
	}

	if err != nil {
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

	var result *morpheus.GetFileTemplateResult
	if v, ok := resp.Result.(*morpheus.GetFileTemplateResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("GetFileTemplateResult", resp.Result))
	}

	if result.FileTemplate == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("FileTemplate"))
	}

	fileTemplate := result.FileTemplate
	d.SetId(convert.Int64ToString(fileTemplate.ID))
	d.Set("name", fileTemplate.Name)
	d.Set("labels", fileTemplate.Labels)
	d.Set("file_name", fileTemplate.FileName)
	d.Set("file_path", fileTemplate.FilePath)
	d.Set("phase", fileTemplate.TemplatePhase)
	d.Set("file_content", fileTemplate.Template)
	d.Set("file_owner", fileTemplate.FileOwner)
	d.Set("setting_name", fileTemplate.SettingName)
	d.Set("setting_category", fileTemplate.SettingCategory)

	return diags
}

func resourceFileTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		var labelSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			labelSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}
		labelList := labelSet.List()
		for _, s := range labelList {
			var labelStr string
			if v, ok := s.(string); ok {
				labelStr = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("label", s))
			}
			labelsPayload = append(labelsPayload, labelStr)
		}
	}

	var fileName string
	if v, ok := d.Get("file_name").(string); ok {
		fileName = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("file_name", d.Get("file_name")))
	}

	var filePath string
	if v, ok := d.Get("file_path").(string); ok {
		filePath = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("file_path", d.Get("file_path")))
	}

	var phase string
	if v, ok := d.Get("phase").(string); ok {
		phase = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("phase", d.Get("phase")))
	}

	var fileContent string
	if v, ok := d.Get("file_content").(string); ok {
		fileContent = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("file_content", d.Get("file_content")))
	}

	var fileOwner string
	if v, ok := d.Get("file_owner").(string); ok {
		fileOwner = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("file_owner", d.Get("file_owner")))
	}

	var settingName string
	if v, ok := d.Get("setting_name").(string); ok {
		settingName = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("setting_name", d.Get("setting_name")))
	}

	var settingCategory string
	if v, ok := d.Get("setting_category").(string); ok {
		settingCategory = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("setting_category", d.Get("setting_category")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"containerTemplate": map[string]any{
				"name":            name,
				"labels":          labelsPayload,
				"fileName":        fileName,
				"filePath":        filePath,
				"templatePhase":   phase,
				"template":        fileContent,
				"fileOwner":       fileOwner,
				"settingName":     settingName,
				"settingCategory": settingCategory,
			},
		},
	}

	resp, err := client.UpdateFileTemplate(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.UpdateFileTemplateResult
	if v, ok := resp.Result.(*morpheus.UpdateFileTemplateResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("UpdateFileTemplateResult", resp.Result))
	}

	if result.FileTemplate == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("FileTemplate"))
	}

	fileTemplate := result.FileTemplate
	d.SetId(convert.Int64ToString(fileTemplate.ID))

	return resourceFileTemplateRead(ctx, d, meta)
}

func resourceFileTemplateDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteFileTemplate(convert.StringToInt64(id), req)
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
