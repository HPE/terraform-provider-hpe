// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package catalogitem

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

const (
	catalogItemTypeForm = "form"
)

func ResourceCatalogItemInstance() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus instance catalog item resource",
		CreateContext: resourceCatalogItemInstanceCreate,
		ReadContext:   resourceCatalogItemInstanceRead,
		UpdateContext: resourceCatalogItemInstanceUpdate,
		DeleteContext: resourceCatalogItemInstanceDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the instance catalog item",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the instance catalog item",
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
				Description: "The description of the instance catalog item",
				Optional:    true,
				Computed:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the instance catalog item",
				Optional:    true,
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the instance catalog item is enabled",
				Optional:    true,
				Default:     true,
			},
			"featured": {
				Type:        schema.TypeBool,
				Description: "Whether the instance catalog item is featured",
				Optional:    true,
				Computed:    true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "The markdown content associated with the instance catalog item",
				Optional:    true,
				Computed:    true,
				StateFunc: func(val any) string {
					if v, ok := val.(string); ok {
						return strings.TrimSpace(v)
					}

					return ""
				},
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					old = strings.TrimSpace(old)
					new = strings.TrimSpace(new)

					return old == new
				},
			},
			"config": {
				Type:             schema.TypeString,
				Description:      "The instance config associated with the instance catalog item",
				Required:         true,
				DiffSuppressFunc: helpers.SuppressEquivalentJSONDiffs,
			},
			"option_type_ids": {
				Type:        schema.TypeList,
				Description: "The list of option type ids associated with the instance catalog item",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed:      true,
				ConflictsWith: []string{"form_id"},
			},
			"form_id": {
				Type:          schema.TypeInt,
				Description:   "The id of the form associated with the workflow catalog item",
				Optional:      true,
				ConflictsWith: []string{"option_type_ids"},
			},
			"image_name": {
				Type:        schema.TypeString,
				Description: "The file name of the instance catalog item logo image",
				Optional:    true,
				Computed:    true,
			},
			"image_path": {
				Type:        schema.TypeString,
				Description: "The file path of the instance catalog item logo image including the file name",
				Optional:    true,
				Computed:    true,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "The visibility of the instance catalog item (public or private)",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCatalogItemInstanceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
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

	catalogItem["type"] = "instance"
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

	// Declared an empty interface
	var outjson map[string]any

	var configStr string
	if v, ok := d.Get("config").(string); ok {
		configStr = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("config", d.Get("config")))
	}
	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(configStr), &outjson)
	catalogItem["config"] = outjson

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		if labelSet, ok := attr.(*schema.Set); ok {
			labelList := labelSet.List()
			for _, s := range labelList {
				if labelStr, ok := s.(string); ok {
					labelsPayload = append(labelsPayload, labelStr)
				}
			}
		}
	}
	catalogItem["labels"] = labelsPayload

	var formID int
	if v, ok := d.Get("form_id").(int); ok {
		formID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("form_id", d.Get("form_id")))
	}

	if formID > 0 {
		catalogItem["formType"] = catalogItemTypeForm
		catalogItem["form"] = map[string]any{
			"id": formID,
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

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateCatalogItemResult
	if v, ok := resp.Result.(*morpheus.CreateCatalogItemResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.CatalogItem == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("CatalogItem"))
	}

	catalogItemResult := result.CatalogItem

	var imagePath string
	if v, ok := d.Get("image_path").(string); ok {
		imagePath = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("image_path", d.Get("image_path")))
	}

	var imageName string
	if v, ok := d.Get("image_name").(string); ok {
		imageName = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("image_name", d.Get("image_name")))
	}

	if imagePath != "" && imageName != "" {
		data, err := os.ReadFile(imagePath)
		if err != nil {
			return diag.FromErr(err)
		}

		var filePayloads []*morpheus.FilePayload
		filePayload := &morpheus.FilePayload{
			ParameterName: "logo",
			FileName:      imageName,
			FileContent:   data,
		}
		filePayloads = append(filePayloads, filePayload)
		response, err := client.UpdateCatalogItemLogo(catalogItemResult.ID, filePayloads, &morpheus.Request{})
		if err != nil {
			log.Printf("API LOGO FAILURE: %s - %s", response, err)

			return diag.FromErr(err)
		}
		log.Printf("API RESPONSE: %s", response)
	}

	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(catalogItemResult.ID))

	diags = append(diags, resourceCatalogItemInstanceRead(ctx, d, meta)...)

	return diags
}

func resourceCatalogItemInstanceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	// store resource data
	var result *morpheus.GetCatalogItemResult
	if v, ok := resp.Result.(*morpheus.GetCatalogItemResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.CatalogItem == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("CatalogItem"))
	}

	catalogItem := result.CatalogItem

	d.SetId(convert.IntToString(int(catalogItem.ID)))
	d.Set("name", catalogItem.Name)
	d.Set("description", catalogItem.Description)
	d.Set("category", catalogItem.Category)
	d.Set("enabled", catalogItem.Enabled)
	d.Set("featured", catalogItem.Featured)
	// option types
	var optionTypes []int64
	if catalogItem.OptionTypes != nil {
		// iterate over the array of tasks
		for i := 0; i < len(catalogItem.OptionTypes); i++ {
			if optionMap, ok := catalogItem.OptionTypes[i].(map[string]any); ok {
				if idFloat, ok := optionMap["id"].(float64); ok {
					optionID := int64(idFloat)
					optionTypes = append(optionTypes, optionID)
				}
			}
		}
	}
	d.Set("option_type_ids", optionTypes)
	d.Set("form_id", catalogItem.Form.ID)
	d.Set("content", catalogItem.Content)

	if configMap, ok := catalogItem.Config.(map[string]any); ok {
		configJSON, _ := json.Marshal(configMap)
		d.Set("config", string(configJSON))
	}

	d.Set("visibility", catalogItem.Visibility)

	if catalogItem.Labels != nil {
		d.Set("labels", catalogItem.Labels)
	}

	imagePathSlice := strings.Split(catalogItem.ImagePath, "/")
	if len(imagePathSlice) > 0 {
		opt := strings.Replace(imagePathSlice[len(imagePathSlice)-1], "_original", "", 1)
		d.Set("image_name", opt)
	}

	return diags
}

func resourceCatalogItemInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
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

	catalogItem["type"] = "instance"
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

	// Declared an empty interface
	var outjson map[string]any

	var configStr string
	if v, ok := d.Get("config").(string); ok {
		configStr = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("config", d.Get("config")))
	}
	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(configStr), &outjson)
	catalogItem["config"] = outjson

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		if labelSet, ok := attr.(*schema.Set); ok {
			labelList := labelSet.List()
			for _, s := range labelList {
				if labelStr, ok := s.(string); ok {
					labelsPayload = append(labelsPayload, labelStr)
				}
			}
		}
	}
	catalogItem["labels"] = labelsPayload

	var formID int
	if v, ok := d.Get("form_id").(int); ok {
		formID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("form_id", d.Get("form_id")))
	}

	if formID > 0 {
		catalogItem["formType"] = catalogItemTypeForm
		catalogItem["form"] = map[string]any{
			"id": formID,
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

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateCatalogItemResult
	if v, ok := resp.Result.(*morpheus.UpdateCatalogItemResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.CatalogItem == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("CatalogItem"))
	}

	catalogItemResult := result.CatalogItem

	if d.HasChange("image_name") || d.HasChange("image_path") {
		var imagePath string
		if v, ok := d.Get("image_path").(string); ok {
			imagePath = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("image_path", d.Get("image_path")))
		}

		data, err := os.ReadFile(imagePath)
		if err != nil {
			return diag.FromErr(err)
		}

		var imageName string
		if v, ok := d.Get("image_name").(string); ok {
			imageName = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("image_name", d.Get("image_name")))
		}

		var filePayloads []*morpheus.FilePayload
		filePayload := &morpheus.FilePayload{
			ParameterName: "logo",
			FileName:      imageName,
			FileContent:   data,
		}
		filePayloads = append(filePayloads, filePayload)
		client.UpdateCatalogItemLogo(catalogItemResult.ID, filePayloads, &morpheus.Request{})
	}

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(catalogItemResult.ID))

	return resourceCatalogItemInstanceRead(ctx, d, meta)
}

func resourceCatalogItemInstanceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
