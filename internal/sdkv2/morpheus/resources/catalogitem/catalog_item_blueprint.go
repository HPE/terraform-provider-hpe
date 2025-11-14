// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package catalogitem

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func ResourceCatalogItemBlueprint() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus Blueprint catalog item resource",
		CreateContext: resourceCatalogItemBlueprintCreate,
		ReadContext:   resourceCatalogItemBlueprintRead,
		UpdateContext: resourceCatalogItemBlueprintUpdate,
		DeleteContext: resourceCatalogItemBlueprintDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the app blueprint catalog item",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the app blueprint catalog item",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the app blueprint catalog item",
				Optional:    true,
				Computed:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The category of the app blueprint catalog item",
				Optional:    true,
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the app blueprint catalog item is enabled",
				Optional:    true,
				Default:     true,
			},
			"featured": {
				Type:        schema.TypeBool,
				Description: "Whether the app blueprint catalog item is featured",
				Optional:    true,
				Computed:    true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "The markdown content associated with the app blueprint catalog item",
				Optional:    true,
				Computed:    true,
				StateFunc: func(val any) string {
					return strings.TrimSuffix(val.(string), "\n")
				},
			},
			"labels": {
				Type: schema.TypeSet,
				Description: "The organization labels associated with the catalog item " +
					"(Only supported on Morpheus 5.5.3 or higher)",
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"blueprint_id": {
				Type:        schema.TypeInt,
				Description: "The id of the blueprint to associate with the app blueprint catalog item",
				Required:    true,
			},
			"app_spec": {
				Type:        schema.TypeString,
				Description: "The app spec associated with the app blueprint catalog item",
				Required:    true,
			},
			"option_type_ids": {
				Type:        schema.TypeList,
				Description: "The list of option type ids associated with the app blueprint catalog item",
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
			"logo_image_name": {
				Type:        schema.TypeString,
				Description: "The file name of the app blueprint catalog item logo image",
				Optional:    true,
				Computed:    true,
			},
			"logo_image_path": {
				Type:        schema.TypeString,
				Description: "The file path of the app blueprint catalog item logo image including the file name",
				Optional:    true,
				Computed:    true,
			},
			"dark_logo_image_name": {
				Type:        schema.TypeString,
				Description: "The file name of the app blueprint catalog item dark mode logo image",
				Optional:    true,
				Computed:    true,
			},
			"dark_logo_image_path": {
				Type:        schema.TypeString,
				Description: "The file path of the app blueprint catalog item dark mode logo image including the file name",
				Optional:    true,
				Computed:    true,
			},
			"visibility": {
				Type:             schema.TypeString,
				Description:      "The visibility of the app blueprint catalog item (public or private)",
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"public", "private"}, false)),
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

//nolint:goconst
func resourceCatalogItemBlueprintCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

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

	catalogItem["type"] = "blueprint"
	catalogItem["iconPath"] = "custom"
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

	blueprint := make(map[string]any)
	var blueprintID int
	if v, ok := d.Get("blueprint_id").(int); ok {
		blueprintID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("blueprint_id", d.Get("blueprint_id")))
	}
	blueprint["id"] = blueprintID
	catalogItem["blueprint"] = blueprint

	var appSpec string
	if v, ok := d.Get("app_spec").(string); ok {
		appSpec = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("app_spec", d.Get("app_spec")))
	}
	catalogItem["appSpec"] = appSpec

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		if labelSet, ok := attr.(*schema.Set); ok {
			for _, s := range labelSet.List() {
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
		catalogItem["formType"] = "form"
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

	d.SetId(convert.Int64ToString(catalogItemResult.ID))

	diags = append(diags, resourceCatalogItemBlueprintRead(ctx, d, meta)...)

	return diags
}

func resourceCatalogItemBlueprintRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

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

	d.SetId(convert.Int64ToString(catalogItem.ID))
	d.Set("name", catalogItem.Name)
	d.Set("description", catalogItem.Description)
	d.Set("category", catalogItem.Category)
	d.Set("enabled", catalogItem.Enabled)
	d.Set("featured", catalogItem.Featured)

	var optionTypes []int64
	if len(catalogItem.OptionTypes) > 0 {
		for i := 0; i < len(catalogItem.OptionTypes); i++ {
			var option map[string]any
			if v, ok := catalogItem.OptionTypes[i].(map[string]any); ok {
				option = v
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("label", catalogItem.OptionTypes[i]))
			}

			var optionID int64
			if v, ok := option["id"].(float64); ok {
				optionID = int64(v)
			} else {
				return diag.FromErr(helpers.TypeAssertFailError("id", option["id"]))
			}
			optionTypes = append(optionTypes, optionID)
		}
	}
	d.Set("option_type_ids", optionTypes)

	d.Set("form_id", catalogItem.Form.ID)

	d.Set("app_spec", catalogItem.AppSpec)
	d.Set("content", catalogItem.Content)

	d.Set("blueprint_id", catalogItem.Blueprint.ID)

	if catalogItem.Labels != nil {
		d.Set("labels", catalogItem.Labels)
	}

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

func resourceCatalogItemBlueprintUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	catalogItem["type"] = "blueprint"
	catalogItem["iconPath"] = "custom"
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

	blueprint := make(map[string]any)
	var blueprintID int
	if v, ok := d.Get("blueprint_id").(int); ok {
		blueprintID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("blueprint_id", d.Get("blueprint_id")))
	}
	blueprint["id"] = blueprintID
	catalogItem["blueprint"] = blueprint

	var appSpec string
	if v, ok := d.Get("app_spec").(string); ok {
		appSpec = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("app_spec", d.Get("app_spec")))
	}
	catalogItem["appSpec"] = appSpec

	labelsPayload := make([]string, 0)
	if attr, ok := d.GetOk("labels"); ok {
		if labelSet, ok := attr.(*schema.Set); ok {
			for _, s := range labelSet.List() {
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
		catalogItem["formType"] = "form"
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

	d.SetId(convert.Int64ToString(catalogItemResult.ID))

	return resourceCatalogItemBlueprintRead(ctx, d, meta)
}

func resourceCatalogItemBlueprintDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteCatalogItem(convert.StringToInt64(id), req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %s", resp, err)

			return diag.FromErr(err)
		}

		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)
	d.SetId("")

	return diags
}
