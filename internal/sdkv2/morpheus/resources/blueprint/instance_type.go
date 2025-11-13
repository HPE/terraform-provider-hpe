// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func ResourceInstanceType() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus instance type resource",
		CreateContext: resourceInstanceTypeCreate,
		ReadContext:   resourceInstanceTypeRead,
		UpdateContext: resourceInstanceTypeUpdate,
		DeleteContext: resourceInstanceTypeDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the instance type",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the instance type",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The instance type code",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the instance type",
				Optional:    true,
				Computed:    true,
			},
			"labels": {
				Type: schema.TypeSet,
				Description: "The organization labels associated with the script template " +
					"(Only supported on Morpheus 5.5.3 or higher)",
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"category": {
				Type:        schema.TypeString,
				Description: "The instance type category (web, sql, nosql, apps, network, messaging, cache, os, cloud, utility)",
				ValidateFunc: validation.StringInSlice(
					[]string{"web", "sql", "nosql", "apps", "network", "messaging", "cache", "os", "cloud", "utility"},
					false,
				),
				Required: true,
			},
			"image_name": {
				Type:        schema.TypeString,
				Description: "The file name of the instance type logo image",
				Optional:    true,
			},
			"image_path": {
				Type:        schema.TypeString,
				Description: "The file path of the instance type logo image including the file name",
				Optional:    true,
			},
			"environment_prefix": {
				Type:        schema.TypeString,
				Description: "The prefix used for instance environment variables",
				Optional:    true,
				Computed:    true,
			},
			"enable_settings": {
				Type:        schema.TypeBool,
				Description: "Whether to enable settings for the instance type",
				Optional:    true,
				Computed:    true,
			},
			"enable_scaling": {
				Type:        schema.TypeBool,
				Description: "Whether to enable scaling for the instance type",
				Optional:    true,
				Computed:    true,
			},
			"enable_deployments": {
				Type:        schema.TypeBool,
				Description: "Whether to enable deployments for the instance type",
				Optional:    true,
				Computed:    true,
			},
			"evar": {
				Type:        schema.TypeList,
				Description: "The environment variables to create",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the environment variable",
							Optional:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The environment variable value when the value can be in plaintext",
							Optional:    true,
							Computed:    true,
						},
						"masked_value": {
							Type:        schema.TypeString,
							Description: "The environment variable value when the value needs to be masked",
							Optional:    true,
							Sensitive:   true,
							Computed:    true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if old == "" {
									return false
								}
								h := sha256.New()
								h.Write([]byte(new))
								sha256Hash := hex.EncodeToString(h.Sum(nil))
								log.Println(sha256Hash)

								return strings.EqualFold(strings.ToLower(old), strings.ToLower(sha256Hash))
							},
							DiffSuppressOnRefresh: true,
						},
						"export": {
							Type:        schema.TypeBool,
							Description: "Whether the environment variable is exported as an instance tag",
							Optional:    true,
						},
					},
				},
			},
			"option_type_ids": {
				Type:        schema.TypeList,
				Description: "The IDs of the inputs to associate with the instance type",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "The visibility of the instance type (public or private)",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"public", "private"}, false),
			},
			"featured": {
				Type:        schema.TypeBool,
				Description: "Whether the instance type is marked as featured",
				Optional:    true,
				Computed:    true,
			},
			"price_set_ids": {
				Type:        schema.TypeList,
				Description: "A list of price set ids associated with the instance type",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == old
				},
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceInstanceTypeCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var code string
	if v, ok := d.Get("code").(string); ok {
		code = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("code", d.Get("code")))
	}

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

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

	var category string
	if v, ok := d.Get("category").(string); ok {
		category = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("category", d.Get("category")))
	}

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	var optionTypeIDs any
	if v := d.Get("option_type_ids"); v != nil {
		optionTypeIDs = v
	}

	var environmentPrefix string
	if v, ok := d.Get("environment_prefix").(string); ok {
		environmentPrefix = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("environment_prefix", d.Get("environment_prefix")))
	}

	var evarList []any
	if v, ok := d.Get("evar").([]any); ok {
		evarList = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("evar", d.Get("evar")))
	}

	var enableSettings bool
	if v, ok := d.Get("enable_settings").(bool); ok {
		enableSettings = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_settings", d.Get("enable_settings")))
	}

	var enableScaling bool
	if v, ok := d.Get("enable_scaling").(bool); ok {
		enableScaling = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_scaling", d.Get("enable_scaling")))
	}

	var enableDeployments bool
	if v, ok := d.Get("enable_deployments").(bool); ok {
		enableDeployments = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_deployments", d.Get("enable_deployments")))
	}

	var featured bool
	if v, ok := d.Get("featured").(bool); ok {
		featured = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("featured", d.Get("featured")))
	}

	// priceSets
	var priceSets []map[string]any
	if d.Get("price_set_ids") != nil {
		if priceSetList, ok := d.Get("price_set_ids").([]any); ok {
			if priceSetList != nil {
				// iterate over the array of tasks
				for i := 0; i < len(priceSetList); i++ {
					row := make(map[string]any)
					row["id"] = priceSetList[i]
					priceSets = append(priceSets, row)
				}
			}
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("price_set_ids", d.Get("price_set_ids")))
		}
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"instanceType": map[string]any{
				"name":                 name,
				"code":                 code,
				"description":          description,
				"labels":               labelsPayload,
				"category":             category,
				"visibility":           visibility,
				"optionTypes":          optionTypeIDs,
				"environmentPrefix":    environmentPrefix,
				"environmentVariables": parseInstanceTypeEnvironmentVariables(evarList, d),
				"hasSettings":          enableSettings,
				"hasAutoScale":         enableScaling,
				"hasDeployment":        enableDeployments,
				"featured":             featured,
				"priceSets":            priceSets,
			},
		},
	}

	resp, err := client.CreateInstanceType(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateInstanceTypeResult
	if v, ok := resp.Result.(*morpheus.CreateInstanceTypeResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.InstanceType == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("InstanceType"))
	}

	instanceType := result.InstanceType

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
		response, err := client.UpdateInstanceTypeLogo(instanceType.ID, filePayloads, &morpheus.Request{})
		if err != nil {
			log.Printf("API FAILURE: %s - %s", response, err)

			return diag.FromErr(err)
		}
		log.Printf("API RESPONSE: %s", response)
	}

	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(instanceType.ID))

	diags = append(diags, resourceInstanceTypeRead(ctx, d, meta)...)

	return diags
}

func resourceInstanceTypeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindInstanceTypeByName(name)
	} else if id != "" {
		resp, err = client.GetInstanceType(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Instance type cannot be read without name or id")
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
	// log.Printf("API RESPONSE: %s", resp)

	// store resource data
	var instanceTypePayload InstanceTypePayload
	json.Unmarshal(resp.Body, &instanceTypePayload)

	d.SetId(convert.Int64ToString(instanceTypePayload.ID))
	d.Set("name", instanceTypePayload.Name)
	d.Set("code", instanceTypePayload.Code)
	d.Set("description", instanceTypePayload.Description)

	if instanceTypePayload.Labels != nil {
		d.Set("labels", instanceTypePayload.Labels)
	}

	d.Set("category", instanceTypePayload.Category)
	d.Set("visibility", instanceTypePayload.Visibility)
	d.Set("environment_prefix", instanceTypePayload.EnvironmentPrefix)
	d.Set("enable_settings", instanceTypePayload.HasSettings)
	d.Set("enable_scaling", instanceTypePayload.HasAutoscale)
	d.Set("enable_deployments", instanceTypePayload.HasDeployment)
	d.Set("featured", instanceTypePayload.Featured)

	// priceSets
	var priceSets []int64
	if instanceTypePayload.PriceSets != nil {
		// iterate over the array of price sets
		for i := 0; i < len(instanceTypePayload.PriceSets); i++ {
			priceSet := instanceTypePayload.PriceSets[i]
			priceSets = append(priceSets, priceSet.ID)
		}
	}

	var priceSetIDs []any
	if v, ok := d.Get("price_set_ids").([]any); ok {
		priceSetIDs = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("price_set_ids", d.Get("price_set_ids")))
	}

	priceSetData := matchTemplatesWithSchema(priceSets, priceSetIDs)
	d.Set("price_set_ids", priceSetData)

	// inputs
	var inputs []int64
	if instanceTypePayload.OptionTypes != nil {
		// iterate over the array of option types
		for i := 0; i < len(instanceTypePayload.OptionTypes); i++ {
			input := instanceTypePayload.OptionTypes[i]
			inputs = append(inputs, input.ID)
		}
	}

	var optionTypeIDs []any
	if v, ok := d.Get("option_type_ids").([]any); ok {
		optionTypeIDs = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("option_type_ids", d.Get("option_type_ids")))
	}

	stateInputs := matchTemplatesWithSchema(inputs, optionTypeIDs)
	d.Set("option_type_ids", stateInputs)

	var evars []map[string]any
	if instanceTypePayload.EnvironmentVariables != nil {
		// iterate over the array of environment variables
		for i := 0; i < len(instanceTypePayload.EnvironmentVariables); i++ {
			environmentVariable := instanceTypePayload.EnvironmentVariables[i]
			envPayload := make(map[string]any)
			envPayload["name"] = environmentVariable.Name
			if environmentVariable.Masked {
				envPayload["masked_value"] = environmentVariable.DefaultValueHash
			} else {
				envPayload["value"] = environmentVariable.DefaultValue
			}
			envPayload["export"] = environmentVariable.Export
			evars = append(evars, envPayload)
		}
	}
	d.Set("evar", evars)

	return diags
}

func resourceInstanceTypeUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var code string
	if v, ok := d.Get("code").(string); ok {
		code = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("code", d.Get("code")))
	}

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

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

	var category string
	if v, ok := d.Get("category").(string); ok {
		category = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("category", d.Get("category")))
	}

	var visibility string
	if v, ok := d.Get("visibility").(string); ok {
		visibility = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("visibility", d.Get("visibility")))
	}

	var optionTypeIDs any
	if v := d.Get("option_type_ids"); v != nil {
		optionTypeIDs = v
	}

	var environmentPrefix string
	if v, ok := d.Get("environment_prefix").(string); ok {
		environmentPrefix = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("environment_prefix", d.Get("environment_prefix")))
	}

	var evarList []any
	if v, ok := d.Get("evar").([]any); ok {
		evarList = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("evar", d.Get("evar")))
	}

	var enableSettings bool
	if v, ok := d.Get("enable_settings").(bool); ok {
		enableSettings = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_settings", d.Get("enable_settings")))
	}

	var enableScaling bool
	if v, ok := d.Get("enable_scaling").(bool); ok {
		enableScaling = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_scaling", d.Get("enable_scaling")))
	}

	var enableDeployments bool
	if v, ok := d.Get("enable_deployments").(bool); ok {
		enableDeployments = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enable_deployments", d.Get("enable_deployments")))
	}

	var featured bool
	if v, ok := d.Get("featured").(bool); ok {
		featured = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("featured", d.Get("featured")))
	}

	// priceSets
	var priceSets []map[string]any
	if d.Get("price_set_ids") != nil {
		if priceSetList, ok := d.Get("price_set_ids").([]any); ok {
			if priceSetList != nil {
				// iterate over the array of tasks
				for i := 0; i < len(priceSetList); i++ {
					row := make(map[string]any)
					row["id"] = priceSetList[i]
					priceSets = append(priceSets, row)
				}
			}
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("price_set_ids", d.Get("price_set_ids")))
		}
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"instanceType": map[string]any{
				"name":                 name,
				"code":                 code,
				"description":          description,
				"labels":               labelsPayload,
				"category":             category,
				"visibility":           visibility,
				"optionTypes":          optionTypeIDs,
				"environmentPrefix":    environmentPrefix,
				"environmentVariables": parseInstanceTypeEnvironmentVariables(evarList, d),
				"hasSettings":          enableSettings,
				"hasAutoScale":         enableScaling,
				"hasDeployment":        enableDeployments,
				"featured":             featured,
				"priceSets":            priceSets,
			},
		},
	}

	resp, err := client.UpdateInstanceType(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	// log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateInstanceTypeResult
	if v, ok := resp.Result.(*morpheus.UpdateInstanceTypeResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.InstanceType == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("InstanceType"))
	}

	instanceType := result.InstanceType

	if d.HasChange("image_name") || d.HasChange("image_path") {
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
		client.UpdateInstanceTypeLogo(instanceType.ID, filePayloads, &morpheus.Request{})
	}

	// Successfully updated resource, now set id
	// err, it should not have changed though..
	d.SetId(convert.Int64ToString(instanceType.ID))

	return resourceInstanceTypeRead(ctx, d, meta)
}

func resourceInstanceTypeDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteInstanceType(convert.StringToInt64(id), req)
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

func parseInstanceTypeEnvironmentVariables(variables []any, d *schema.ResourceData) []map[string]any {
	var evars []map[string]any
	if variables == nil {
		return evars
	}

	// iterate over the array of evars
	for i := 0; i < len(variables); i++ {
		row := make(map[string]any)
		evarconfig := variables[i].(map[string]any)
		for k, v := range evarconfig {
			switch k {
			case "name":
				row["name"] = v.(string)
				row["evarName"] = v.(string)
				row["valueType"] = "fixed"
			case "value":
				if v.(string) != "" {
					row["value"] = v.(string)
					row["masked"] = false
				}
			case "masked_value":
				if v.(string) != "" {
					row["value"] = v.(string)
					row["masked"] = true
				}
			case "export":
				row["export"] = v.(bool)
			}
		}
		evars = append(evars, row)
	}

	return evars
}

func matchTemplatesWithSchema(templates []int64, declaredTemplates []any) []int64 {
	if declaredTemplates == nil {
		return templates
	}

	result := make([]int64, len(declaredTemplates))

	rMap := make(map[int64]int64, len(templates))
	for _, template := range templates {
		rMap[template] = template
	}

	for i, definedTemplate := range declaredTemplates {
		definedTemplate := int64(definedTemplate.(int))

		if v, ok := rMap[definedTemplate]; ok {
			result[i] = v
			delete(rMap, v)
		}
	}

	return result
}

type InstanceTypePayload struct {
	morpheus.InstanceType `json:"instanceType"`
}
