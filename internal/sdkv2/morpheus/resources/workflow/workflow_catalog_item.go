// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package workflow

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

// ResourceWorkflowCatalogItem returns the workflow catalog item resource
func ResourceWorkflowCatalogItem() *schema.Resource {
	return resourceWorkflowCatalogItem()
}

func resourceWorkflowCatalogItem() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus workflow catalog item resource",
		CreateContext: resourceWorkflowCatalogItemCreate,
		ReadContext:   resourceWorkflowCatalogItemRead,
		UpdateContext: resourceWorkflowCatalogItemUpdate,
		DeleteContext: resourceWorkflowCatalogItemDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the workflow catalog item",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the workflow catalog item",
				Required:    true,
			},
			"labels": {
				Type: schema.TypeSet,
				Description: "The organization labels associated with the catalog item " +
					"(Only supported on Morpheus 5.5.3 or higher)",
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the workflow catalog item",
				Optional:    true,
				Computed:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the workflow catalog item",
				Optional:    true,
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the workflow catalog item is enabled",
				Optional:    true,
				Default:     true,
			},
			"featured": {
				Type:        schema.TypeBool,
				Description: "Whether the workflow catalog item is featured",
				Optional:    true,
				Computed:    true,
			},
			"workflow_id": {
				Type:        schema.TypeInt,
				Description: "The id of the workflow associated with the workflow catalog item",
				Required:    true,
			},
			"context_type": {
				Type:         schema.TypeString,
				Description:  "The Morpheus context type of the operational workflow",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"instance", "server", "appliance"}, false),
				Computed:     true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "The markdown content associated with the workflow catalog item",
				Optional:    true,
				Computed:    true,
				StateFunc: func(val any) string {
					if v, ok := val.(string); ok {
						return strings.TrimSuffix(v, "\n")
					}

					return ""
				},
			},
			"option_type_ids": {
				Type:        schema.TypeList,
				Description: "The list of option type ids associated with the workflow catalog item",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed:      true,
				ConflictsWith: []string{"form_id"},
			},
			"logo_image_name": {
				Type:        schema.TypeString,
				Description: "The file name of the workflow catalog item logo image",
				Optional:    true,
				Computed:    true,
			},
			"logo_image_path": {
				Type:        schema.TypeString,
				Description: "The file path of the workflow catalog item logo image including the file name",
				Optional:    true,
				Computed:    true,
			},
			"dark_logo_image_name": {
				Type:        schema.TypeString,
				Description: "The file name of the workflow catalog item dark mode logo image",
				Optional:    true,
				Computed:    true,
			},
			"dark_logo_image_path": {
				Type:        schema.TypeString,
				Description: "The file path of the workflow catalog item dark mode logo image including the file name",
				Optional:    true,
				Computed:    true,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "The visibility of the workflow catalog item (public or private)",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
			},
			"form_id": {
				Type:          schema.TypeInt,
				Description:   "The id of the form associated with the workflow catalog item",
				Optional:      true,
				ConflictsWith: []string{"option_type_ids"},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWorkflowCatalogItemCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	catalogItem := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	catalogItem["name"] = name

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}
	catalogItem["description"] = description

	var category string
	if v, ok := d.Get("category").(string); ok {
		category = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("category", d.Get("category")))
	}
	catalogItem["category"] = category

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}
	catalogItem["enabled"] = enabled

	var featured bool
	if v, ok := d.Get("featured").(bool); ok {
		featured = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("featured", d.Get("featured")))
	}
	catalogItem["featured"] = featured
	catalogItem["type"] = "workflow"
	catalogItem["iconPath"] = "custom"

	var contextType string
	if v, ok := d.Get("context_type").(string); ok {
		contextType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("context_type", d.Get("context_type")))
	}
	catalogItem["context"] = contextType

	catalogItem["optionTypes"] = d.Get("option_type_ids")

	var content string
	if v, ok := d.Get("content").(string); ok {
		content = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("content", d.Get("content")))
	}
	catalogItem["content"] = content

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}
	catalogItem["visibility"] = visibility

	var workflowId int
	if v, ok := d.Get("workflow_id").(int); ok {
		workflowId = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("workflow_id", d.Get("workflow_id")))
	}
	catalogItem["workflow"] = map[string]any{
		"id": workflowId,
	}

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		var labelSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			labelSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}
		labelsList := labelSet.List()
		for _, s := range labelsList {
			var label string
			if v, ok := s.(string); ok {
				label = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("label", s))
			}
			labelsPayload = append(labelsPayload, label)
		}
	}
	catalogItem["labels"] = labelsPayload

	var formId int
	if v, ok := d.Get("form_id").(int); ok {
		formId = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("form_id", d.Get("form_id")))
	}
	if formId > 0 {
		catalogItem["formType"] = "form"
		catalogItem["form"] = map[string]any{
			"id": formId,
		}
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"catalogItemType": catalogItem,
		},
	}
	resp, err := client.CreateCatalogItem(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result *morpheus.CreateCatalogItemResult
	if v, ok := resp.Result.(*morpheus.CreateCatalogItemResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}
	catalogItemResult := result.CatalogItem
	if catalogItemResult == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("CatalogItem"))
	}

	var filePayloads []*morpheus.FilePayload

	var logoImagePath string
	if v, ok := d.Get("logo_image_path").(string); ok {
		logoImagePath = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("logo_image_path", d.Get("logo_image_path")))
	}

	var logoImageName string
	if v, ok := d.Get("logo_image_name").(string); ok {
		logoImageName = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("logo_image_name", d.Get("logo_image_name")))
	}

	if logoImagePath != "" && logoImageName != "" {
		data, err := os.ReadFile(logoImagePath)
		if err != nil {
			return diag.FromErr(err)
		}

		filePayload := &morpheus.FilePayload{
			ParameterName: "logo",
			FileName:      logoImageName,
			FileContent:   data,
		}
		filePayloads = append(filePayloads, filePayload)
	}

	var darkLogoImagePath string
	if v, ok := d.Get("dark_logo_image_path").(string); ok {
		darkLogoImagePath = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("dark_logo_image_path", d.Get("dark_logo_image_path")))
	}

	var darkLogoImageName string
	if v, ok := d.Get("dark_logo_image_name").(string); ok {
		darkLogoImageName = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("dark_logo_image_name", d.Get("dark_logo_image_name")))
	}

	if darkLogoImagePath != "" && darkLogoImageName != "" {
		darkLogoData, err := os.ReadFile(darkLogoImagePath)
		if err != nil {
			return diag.FromErr(err)
		}

		darkLogoPayload := &morpheus.FilePayload{
			ParameterName: "darkLogo",
			FileName:      darkLogoImageName,
			FileContent:   darkLogoData,
		}
		filePayloads = append(filePayloads, darkLogoPayload)
	}

	if len(filePayloads) > 0 {
		response, err := client.UpdateCatalogItemLogo(catalogItemResult.ID, filePayloads, &morpheus.Request{})
		if err != nil {
			log.Printf("API FAILURE: %s - %s", response, err)
		}
		log.Printf("API RESPONSE: %s", response)
	}

	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(catalogItemResult.ID))

	resourceWorkflowCatalogItemRead(ctx, d, meta)

	return diags
}

func resourceWorkflowCatalogItemRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
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
		resp, err = client.FindCatalogItemByName(name)
	} else if id != "" {
		resp, err = client.GetCatalogItem(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Catalog Item cannot be read without name or id")
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
	var result *morpheus.GetCatalogItemResult
	if v, ok := resp.Result.(*morpheus.GetCatalogItemResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}
	catalogItem := result.CatalogItem
	if catalogItem == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("CatalogItem"))
	}

	d.SetId(convert.IntToString(int(catalogItem.ID)))
	d.Set("name", catalogItem.Name)
	d.Set("labels", catalogItem.Labels)
	d.Set("description", catalogItem.Description)
	d.Set("category", catalogItem.Category)
	d.Set("enabled", catalogItem.Enabled)
	d.Set("featured", catalogItem.Featured)
	// option types
	var optionTypes []int64
	if catalogItem.OptionTypes != nil {
		// iterate over the array of tasks
		for i := 0; i < len(catalogItem.OptionTypes); i++ {
			var option map[string]any
			if v, ok := catalogItem.OptionTypes[i].(map[string]any); ok {
				option = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("option", catalogItem.OptionTypes[i]))
			}
			var optionIDFloat float64
			if v, ok := option["id"].(float64); ok {
				optionIDFloat = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("option_id", option["id"]))
			}
			optionID := int64(optionIDFloat)
			optionTypes = append(optionTypes, optionID)
		}
	}
	d.Set("option_type_ids", optionTypes)
	d.Set("content", catalogItem.Content)
	d.Set("context_type", catalogItem.Context)
	d.Set("visibility", catalogItem.Visibility)
	d.Set("form_id", catalogItem.Form.ID)
	// Parse workflow ID
	var data map[string]any
	err = json.Unmarshal(resp.Body, &data)
	if err != nil {
		panic(err)
	}
	var catalogItemData map[string]any
	if v, ok := data["catalogItemType"].(map[string]any); ok {
		catalogItemData = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("catalogItemType", data["catalogItemType"]))
	}
	var workflowData map[string]any
	if v, ok := catalogItemData["workflow"].(map[string]any); ok {
		workflowData = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("workflow", catalogItemData["workflow"]))
	}
	var workflowIdFloat float64
	if v, ok := workflowData["id"].(float64); ok {
		workflowIdFloat = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("workflow_id", workflowData["id"]))
	}
	workflowId := int(workflowIdFloat)
	d.Set("workflow_id", workflowId)
	imagePath := strings.Split(catalogItem.ImagePath, "/")
	if len(imagePath) > 0 {
		opt := strings.Replace(imagePath[len(imagePath)-1], "_original", "", 1)
		d.Set("logo_image_name", opt)
	}
	darkImagePath := strings.Split(catalogItem.DarkImagePath, "/")
	if len(darkImagePath) > 0 {
		darkOpt := strings.Replace(darkImagePath[len(darkImagePath)-1], "_original", "", 1)
		d.Set("dark_logo_image_name", darkOpt)
	}

	return diags
}

func resourceWorkflowCatalogItemUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}
	id := d.Id()

	catalogItem := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	catalogItem["name"] = name

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		var labelSet *schema.Set
		if v, ok := attr.(*schema.Set); ok {
			labelSet = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("labels", attr))
		}
		labelsList := labelSet.List()
		for _, s := range labelsList {
			var label string
			if v, ok := s.(string); ok {
				label = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("label", s))
			}
			labelsPayload = append(labelsPayload, label)
		}
	}
	catalogItem["labels"] = labelsPayload

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}
	catalogItem["description"] = description

	var category string
	if v, ok := d.Get("category").(string); ok {
		category = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("category", d.Get("category")))
	}
	catalogItem["category"] = category

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}
	catalogItem["enabled"] = enabled

	var featured bool
	if v, ok := d.Get("featured").(bool); ok {
		featured = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("featured", d.Get("featured")))
	}
	catalogItem["featured"] = featured
	catalogItem["type"] = "workflow"

	var contextType string
	if v, ok := d.Get("context_type").(string); ok {
		contextType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("context_type", d.Get("context_type")))
	}
	catalogItem["context"] = contextType

	catalogItem["optionTypes"] = d.Get("option_type_ids")

	var content string
	if v, ok := d.Get("content").(string); ok {
		content = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("content", d.Get("content")))
	}
	catalogItem["content"] = content

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}
	catalogItem["visibility"] = visibility

	var workflowId int
	if v, ok := d.Get("workflow_id").(int); ok {
		workflowId = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("workflow_id", d.Get("workflow_id")))
	}
	catalogItem["workflow"] = map[string]any{
		"id": workflowId,
	}

	var formId int
	if v, ok := d.Get("form_id").(int); ok {
		formId = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("form_id", d.Get("form_id")))
	}
	if formId > 0 {
		catalogItem["formType"] = "form"
		catalogItem["form"] = map[string]any{
			"id": formId,
		}
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"catalogItemType": catalogItem,
		},
	}

	resp, err := client.UpdateCatalogItem(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	var result *morpheus.UpdateCatalogItemResult
	if v, ok := resp.Result.(*morpheus.UpdateCatalogItemResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("result", resp.Result))
	}
	catalogItemResult := result.CatalogItem
	if catalogItemResult == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("CatalogItem"))
	}

	var filePayloads []*morpheus.FilePayload

	if d.HasChange("logo_image_path") || d.HasChange("logo_image_name") {
		var logoImagePath string
		if v, ok := d.Get("logo_image_path").(string); ok {
			logoImagePath = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("logo_image_path", d.Get("logo_image_path")))
		}

		var logoImageName string
		if v, ok := d.Get("logo_image_name").(string); ok {
			logoImageName = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("logo_image_name", d.Get("logo_image_name")))
		}

		data, err := os.ReadFile(logoImagePath)
		if err != nil {
			return diag.FromErr(err)
		}

		filePayload := &morpheus.FilePayload{
			ParameterName: "logo",
			FileName:      logoImageName,
			FileContent:   data,
		}
		filePayloads = append(filePayloads, filePayload)
	}
	if d.HasChange("dark_logo_image_path") || d.HasChange("dark_logo_image_name") {
		var darkLogoImagePath string
		if v, ok := d.Get("dark_logo_image_path").(string); ok {
			darkLogoImagePath = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("dark_logo_image_path", d.Get("dark_logo_image_path")))
		}

		var darkLogoImageName string
		if v, ok := d.Get("dark_logo_image_name").(string); ok {
			darkLogoImageName = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("dark_logo_image_name", d.Get("dark_logo_image_name")))
		}

		darkLogoData, err := os.ReadFile(darkLogoImagePath)
		if err != nil {
			return diag.FromErr(err)
		}

		darkLogoPayload := &morpheus.FilePayload{
			ParameterName: "darkLogo",
			FileName:      darkLogoImageName,
			FileContent:   darkLogoData,
		}
		filePayloads = append(filePayloads, darkLogoPayload)
	}

	if len(filePayloads) > 0 {
		response, err := client.UpdateCatalogItemLogo(catalogItemResult.ID, filePayloads, &morpheus.Request{})
		if err != nil {
			log.Printf("API FAILURE: %s - %s", response, err)
		}
		log.Printf("API RESPONSE: %s", response)
	}

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(catalogItemResult.ID))

	return resourceWorkflowCatalogItemRead(ctx, d, meta)
}

func resourceWorkflowCatalogItemDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if clientAssert, ok := meta.(*morpheus.Client); ok {
		client = clientAssert
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteCatalogItem(convert.StringToInt64(id), req)
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
