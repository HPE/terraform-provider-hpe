// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package plan

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func ResourcePriceSet() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a price set resource",
		CreateContext: resourcePriceSetCreate,
		ReadContext:   resourcePriceSetRead,
		UpdateContext: resourcePriceSetUpdate,
		DeleteContext: resourcePriceSetDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the price set",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the price set",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the price set",
				Required:    true,
				ForceNew:    true,
			},
			"region_code": {
				Type:        schema.TypeString,
				Description: "The region code of the price set",
				Required:    true,
				ForceNew:    true,
			},
			"cloud_id": {
				Type:        schema.TypeInt,
				Description: "The id of the cloud",
				Optional:    true,
				ForceNew:    true,
			},
			"resource_pool_id": {
				Type:        schema.TypeInt,
				Description: "The resource pool to assign the price set to",
				Optional:    true,
				ForceNew:    true,
			},
			"type": {
				Type: schema.TypeString,
				Description: "The price set type (fixed, compute_plus_storage, component, " +
					"load_balancer, virtual_image, snapshot, software_or_service)",
				ValidateFunc: validation.StringInSlice(
					[]string{
						"fixed", "compute_plus_storage", "component",
						"load_balancer", "virtual_image", "snapshot", "software_or_service",
					},
					false,
				),
				Required: true,
				ForceNew: true,
			},
			"price_unit": {
				Type:        schema.TypeString,
				Description: "The price unit (minute, hour, day, month, year, two year, three year, four year, five year)",
				ValidateFunc: validation.StringInSlice(
					[]string{"minute", "hour", "day", "month", "year", "two year", "three year", "four year", "five year"},
					false,
				),
				Required: true,
				ForceNew: true,
			},
			"price_ids": {
				Type:        schema.TypeList,
				Description: "The list of price ids associated with the price set",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePriceSetCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	priceSet := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	priceSet["name"] = name

	var code string
	if v, ok := d.Get("code").(string); ok {
		code = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("code", d.Get("code")))
	}
	priceSet["code"] = code

	var regionCode string
	if v, ok := d.Get("region_code").(string); ok {
		regionCode = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("region_code", d.Get("region_code")))
	}
	priceSet["regionCode"] = regionCode

	var cloudID int
	if v, ok := d.Get("cloud_id").(int); ok {
		cloudID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloud_id", d.Get("cloud_id")))
	}
	priceSet["zone"] = map[string]any{
		"id": cloudID,
	}

	var resourcePoolID int
	if v, ok := d.Get("resource_pool_id").(int); ok {
		resourcePoolID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resource_pool_id", d.Get("resource_pool_id")))
	}
	priceSet["zonePool"] = map[string]any{
		"id": resourcePoolID,
	}

	var priceUnit string
	if v, ok := d.Get("price_unit").(string); ok {
		priceUnit = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("price_unit", d.Get("price_unit")))
	}
	priceSet["priceUnit"] = priceUnit

	var priceType string
	if v, ok := d.Get("type").(string); ok {
		priceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("type", d.Get("type")))
	}
	priceSet["type"] = priceType

	var priceIDs []map[string]any
	priceIDsRaw := d.Get("price_ids")
	if priceIDsRaw != nil {
		var priceIDList []any
		if v, ok := priceIDsRaw.([]any); ok {
			priceIDList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("price_ids", priceIDsRaw))
		}

		if priceIDList != nil {
			for i := 0; i < len(priceIDList); i++ {
				row := make(map[string]any)
				row["id"] = priceIDList[i]
				priceIDs = append(priceIDs, row)
			}
		}
	}
	priceSet["prices"] = priceIDs

	req := &morpheus.Request{
		Body: map[string]any{
			"priceSet": priceSet,
		},
	}
	resp, err := client.CreatePriceSet(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreatePriceSetResult
	if v, ok := resp.Result.(*morpheus.CreatePriceSetResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	// The Morpheus API returns HTTP 200 with success:true even on validation
	// failure (a server-side bug). Detect failure by checking for a missing/zero
	// id or non-empty errors map.
	if result.ID == 0 || len(result.Errors) > 0 {
		errMsg := "API reported success but failed to create price set"
		if result.Message != "" {
			errMsg = result.Message
		}

		if len(result.Errors) > 0 {
			for field, msg := range result.Errors {
				errMsg += fmt.Sprintf("; %s: %s", field, msg)
			}
		}

		return diag.Errorf("%s", errMsg)
	}

	d.SetId(convert.Int64ToString(result.ID))
	diags = append(diags, resourcePriceSetRead(ctx, d, meta)...)

	return diags
}

func resourcePriceSetRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindPriceSetByName(name)
	} else if id != "" {
		resp, err = client.GetPriceSet(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Price cannot be read without name or id")
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

	var priceSet MorpheusPriceSet
	json.Unmarshal(resp.Body, &priceSet)

	if !priceSet.Priceset.Active {
		d.SetId("")

		return diags
	}

	d.SetId(convert.IntToString(priceSet.Priceset.ID))
	d.Set("name", priceSet.Priceset.Name)
	d.Set("code", priceSet.Priceset.Code)
	d.Set("region_code", priceSet.Priceset.Regioncode)
	d.Set("cloud_id", priceSet.Priceset.Zone.ID)

	if _, ok := d.GetOk("resource_pool_id"); ok {
		d.Set("resource_pool_id", priceSet.Priceset.Zonepool.ID)
	}

	d.Set("price_unit", priceSet.Priceset.Priceunit)
	d.Set("type", priceSet.Priceset.Type)

	var priceIDs []int
	if len(priceSet.Priceset.Prices) > 0 {
		for _, v := range priceSet.Priceset.Prices {
			priceIDs = append(priceIDs, v.ID)
		}
	}

	var priceIDsFromState []any
	if v, ok := d.Get("price_ids").([]any); ok {
		priceIDsFromState = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("price_ids", d.Get("price_ids")))
	}

	statePricePayload := matchPricesWithSchema(priceIDs, priceIDsFromState)
	d.Set("price_ids", statePricePayload)

	return diags
}

func resourcePriceSetUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()

	priceSet := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	priceSet["name"] = name

	var code string
	if v, ok := d.Get("code").(string); ok {
		code = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("code", d.Get("code")))
	}
	priceSet["code"] = code

	var regionCode string
	if v, ok := d.Get("region_code").(string); ok {
		regionCode = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("region_code", d.Get("region_code")))
	}
	priceSet["regionCode"] = regionCode

	var cloudID int
	if v, ok := d.Get("cloud_id").(int); ok {
		cloudID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cloud_id", d.Get("cloud_id")))
	}
	priceSet["zone"] = map[string]any{
		"id": cloudID,
	}

	var resourcePoolID int
	if v, ok := d.Get("resource_pool_id").(int); ok {
		resourcePoolID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("resource_pool_id", d.Get("resource_pool_id")))
	}
	priceSet["zonePool"] = map[string]any{
		"id": resourcePoolID,
	}

	var priceUnit string
	if v, ok := d.Get("price_unit").(string); ok {
		priceUnit = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("price_unit", d.Get("price_unit")))
	}
	priceSet["priceUnit"] = priceUnit

	var priceType string
	if v, ok := d.Get("type").(string); ok {
		priceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("type", d.Get("type")))
	}
	priceSet["type"] = priceType

	var priceIDs []map[string]any
	priceIDsRaw := d.Get("price_ids")
	if priceIDsRaw != nil {
		var priceIDList []any
		if v, ok := priceIDsRaw.([]any); ok {
			priceIDList = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("price_ids", priceIDsRaw))
		}

		if priceIDList != nil {
			for i := 0; i < len(priceIDList); i++ {
				row := make(map[string]any)
				row["id"] = priceIDList[i]
				priceIDs = append(priceIDs, row)
			}
		}
	}
	priceSet["prices"] = priceIDs

	req := &morpheus.Request{
		Body: map[string]any{
			"priceSet": priceSet,
		},
	}

	resp, err := client.UpdatePriceSet(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var updateResult *morpheus.UpdatePriceSetResult
	if v, ok := resp.Result.(*morpheus.UpdatePriceSetResult); ok {
		updateResult = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if !updateResult.Success || len(updateResult.Errors) > 0 {
		errMsg := "Failed to update price set"
		if updateResult.Message != "" {
			errMsg = updateResult.Message
		}

		if len(updateResult.Errors) > 0 {
			for field, msg := range updateResult.Errors {
				errMsg += fmt.Sprintf("; %s: %s", field, msg)
			}
		}

		return diag.Errorf("%s", errMsg)
	}

	d.SetId(id)

	return resourcePriceSetRead(ctx, d, meta)
}

func resourcePriceSetDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics

	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeletePriceSet(convert.StringToInt64(id), req)
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

func matchPricesWithSchema(prices []int, declaredPrices []any) []int {
	result := make([]int, len(declaredPrices))

	rMap := make(map[int]int, len(prices))
	for _, price := range prices {
		rMap[price] = price
	}

	for i, declaredPrice := range declaredPrices {
		declaredPrice := declaredPrice.(int)

		if v, ok := rMap[declaredPrice]; ok {
			result[i] = v
			delete(rMap, v)
		}
	}

	for _, rcpt := range rMap {
		result = append(result, rcpt)
	}

	return result
}

type MorpheusPriceSet struct {
	Priceset struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Code          string `json:"code"`
		Active        bool   `json:"active"`
		Priceunit     string `json:"priceUnit"`
		Type          string `json:"type"`
		Regioncode    string `json:"regionCode"`
		Systemcreated bool   `json:"systemCreated"`
		Zone          struct {
			ID int `json:"id"`
		} `json:"zone"`
		Zonepool struct {
			ID int `json:"id"`
		} `json:"zonePool"`
		Account any `json:"account"`
		Prices  []struct {
			ID                  int     `json:"id"`
			Name                string  `json:"name"`
			Code                string  `json:"code"`
			Pricetype           string  `json:"priceType"`
			Priceunit           string  `json:"priceUnit"`
			Additionalpriceunit string  `json:"additionalPriceUnit"`
			Price               float64 `json:"price"`
			Customprice         float64 `json:"customPrice"`
			Markuptype          any     `json:"markupType"`
			Markup              float64 `json:"markup"`
			Markuppercent       float64 `json:"markupPercent"`
			Cost                float64 `json:"cost"`
			Currency            string  `json:"currency"`
			Incurcharges        string  `json:"incurCharges"`
			Platform            any     `json:"platform"`
			Software            any     `json:"software"`
			Volumetype          struct {
				ID   int    `json:"id"`
				Code string `json:"code"`
				Name string `json:"name"`
			} `json:"volumeType"`
			Datastore       any `json:"datastore"`
			Crosscloudapply any `json:"crossCloudApply"`
			Account         any `json:"account"`
		} `json:"prices"`
	} `json:"priceSet"`
}
