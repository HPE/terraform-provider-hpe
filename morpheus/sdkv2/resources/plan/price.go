// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package plan

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

const (
	markupTypeFixed   = "fixed"
	markupTypePercent = "percent"
	markupTypeCustom  = "custom"
)

func ResourcePrice() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a price resource",
		CreateContext: resourcePriceCreate,
		ReadContext:   resourcePriceRead,
		UpdateContext: resourcePriceUpdate,
		DeleteContext: resourcePriceDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the price",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the price",
				Required:    true,
			},
			"code": {
				Type:        schema.TypeString,
				Description: "The code of the price",
				Required:    true,
				ForceNew:    true,
			},
			"tenant_id": {
				Type:        schema.TypeInt,
				Description: "The id of the tenant to assign the price to",
				Optional:    true,
				ForceNew:    true,
			},
			"price_type": {
				Type:        schema.TypeString,
				Description: "The price type",
				ValidateFunc: validation.StringInSlice([]string{
					"fixed", "compute", "memory", "cores", "storage", "datastore",
					"platform", "software", "load_balancer", "load_balancer_virtual_server",
				}, false),
				Required: true,
			},
			"platform": {
				Type:        schema.TypeString,
				Description: "The name of the platform",
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"canonical", "centos", "debian", "fedora", "opensuse",
					"redhat", "suse", "xen", "linux", "windows",
				}, false),
			},
			"volume_type_id": {
				Type:        schema.TypeInt,
				Description: "The id of the volume type",
				Optional:    true,
				Computed:    true,
			},
			"software": {
				Type:        schema.TypeString,
				Description: "The name of the software",
				Optional:    true,
				Computed:    true,
			},
			"datastore_id": {
				Type:        schema.TypeInt,
				Description: "The id of the datastore to associate the price with",
				Optional:    true,
				Computed:    true,
			},
			"apply_price_accross_clouds": { //nolint:misspell
				Type:        schema.TypeBool,
				Description: "Whether to apply the datastore price across clouds",
				Optional:    true,
				Computed:    true,
			},
			"price_unit": {
				Type:        schema.TypeString,
				Description: "The price unit",
				ValidateFunc: validation.StringInSlice([]string{
					"minute", "hour", "day", "month", "year",
					"two year", "three year", "four year", "five year",
				}, false),
				Required: true,
			},
			"incur_charges": {
				Type:         schema.TypeString,
				Description:  "When charges will be incurred (running, stopped, always)",
				ValidateFunc: validation.StringInSlice([]string{"running", "stopped", "always"}, false),
				Required:     true,
			},
			"currency": {
				Type:        schema.TypeString,
				Description: "The currency of the price",
				Required:    true,
			},
			"cost": {
				Type:        schema.TypeFloat,
				Description: "The cost of the price",
				Required:    true,
			},
			"markup_type": {
				Type:         schema.TypeString,
				Description:  "The type of markup applied to the cost (fixed, percent, custom)",
				ValidateFunc: validation.StringInSlice([]string{"fixed", "percent", "custom"}, false),
				Optional:     true,
				Computed:     true,
			},
			"markup_cost": {
				Type:          schema.TypeFloat,
				Description:   "The fixed cost at which the base cost is marked up",
				Optional:      true,
				ConflictsWith: []string{"markup_percent", "custom_price"},
			},
			"markup_percent": {
				Type:          schema.TypeFloat,
				Description:   "The percentage at which the base cost is marked up",
				Optional:      true,
				ConflictsWith: []string{"markup_cost", "custom_price"},
			},
			"custom_price": {
				Type:          schema.TypeFloat,
				Description:   "The custom price",
				Optional:      true,
				ConflictsWith: []string{"markup_cost", "markup_percent"},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePriceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	price := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	price["name"] = name

	var code string
	if v, ok := d.Get("code").(string); ok {
		code = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("code", d.Get("code")))
	}
	price["code"] = code

	var priceType string
	if v, ok := d.Get("price_type").(string); ok {
		priceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("price_type", d.Get("price_type")))
	}
	price["priceType"] = priceType

	var priceUnit string
	if v, ok := d.Get("price_unit").(string); ok {
		priceUnit = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("price_unit", d.Get("price_unit")))
	}
	price["priceUnit"] = priceUnit

	var incurCharges string
	if v, ok := d.Get("incur_charges").(string); ok {
		incurCharges = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("incur_charges", d.Get("incur_charges")))
	}
	price["incurCharges"] = incurCharges

	var currency string
	if v, ok := d.Get("currency").(string); ok {
		currency = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("currency", d.Get("currency")))
	}
	price["currency"] = currency

	var cost float64
	if v, ok := d.Get("cost").(float64); ok {
		cost = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cost", d.Get("cost")))
	}
	price["cost"] = cost

	var markupType string
	if v, ok := d.Get("markup_type").(string); ok {
		markupType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("markup_type", d.Get("markup_type")))
	}

	// Evaluate different markup types
	switch markupType {
	case markupTypeFixed:
		price["markupType"] = markupTypeFixed
		var markupCost float64
		if v, ok := d.Get("markup_cost").(float64); ok {
			markupCost = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("markup_cost", d.Get("markup_cost")))
		}
		price["markup"] = markupCost
	case markupTypePercent:
		price["markupType"] = markupTypePercent
		var markupPercent float64
		if v, ok := d.Get("markup_percent").(float64); ok {
			markupPercent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("markup_percent", d.Get("markup_percent")))
		}
		price["markupPercent"] = markupPercent
	case markupTypeCustom:
		price["markupType"] = markupTypeCustom
		price["customPrice"] = d.Get("custom_price")
	}

	if d.Get("tenant_id") != nil {
		var tenantID int
		if v, ok := d.Get("tenant_id").(int); ok {
			tenantID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("tenant_id", d.Get("tenant_id")))
		}
		price["account"] = map[string]any{
			"id": tenantID,
		}
	}

	// Evaluate different price types
	switch priceType {
	case "platform":
		var platform string
		if v, ok := d.Get("platform").(string); ok {
			platform = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("platform", d.Get("platform")))
		}
		if platform == "" {
			return diag.Errorf("A platform must be specified")
		}
		price["platform"] = platform
	case "software":
		var software string
		if v, ok := d.Get("software").(string); ok {
			software = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("software", d.Get("software")))
		}
		price["software"] = software
	case "storage":
		var volumeTypeID int
		if v, ok := d.Get("volume_type_id").(int); ok {
			volumeTypeID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("volume_type_id", d.Get("volume_type_id")))
		}
		price["volumeType"] = map[string]any{
			"id": volumeTypeID,
		}
	case "datastore":
		var datastoreID int
		if v, ok := d.Get("datastore_id").(int); ok {
			datastoreID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("datastore_id", d.Get("datastore_id")))
		}
		price["datastore"] = map[string]any{
			"id": datastoreID,
		}
		var applyPriceAcrossClouds bool
		if v, ok := d.Get("apply_price_accross_clouds").(bool); ok { //nolint:misspell
			applyPriceAcrossClouds = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError( //nolint:misspell
				"apply_price_accross_clouds", d.Get("apply_price_accross_clouds"))) //nolint:misspell
		}
		price["crossCloudApply"] = applyPriceAcrossClouds
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"price": price,
		},
	}
	resp, err := client.CreatePrice(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreatePriceResult
	if v, ok := resp.Result.(*morpheus.CreatePriceResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(result.ID))
	diags = append(diags, resourcePriceRead(ctx, d, meta)...)

	return diags
}

func resourcePriceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == "" && name != "" {
		resp, err = client.FindPriceByName(name)
	} else if id != "" {
		resp, err = client.GetPrice(convert.StringToInt64(id), &morpheus.Request{})
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

	// store resource data
	var price MorpheusPrice
	json.Unmarshal(resp.Body, &price)

	if !price.Price.Active {
		d.SetId("")

		return diags
	}

	d.SetId(convert.IntToString(price.Price.ID))
	d.Set("name", price.Price.Name)
	d.Set("code", price.Price.Code)
	if _, ok := d.GetOk("tenant_id"); ok {
		d.Set("tenant_id", price.Price.Account.ID)
	}
	d.Set("price_type", price.Price.Pricetype)
	if _, ok := d.GetOk("platform"); ok {
		d.Set("platform", price.Price.Platform)
	}
	if _, ok := d.GetOk("volume_type_id"); ok {
		d.Set("volume_type_id", price.Price.Volumetype.ID)
	}
	if _, ok := d.GetOk("software"); ok {
		d.Set("software", price.Price.Software)
	}
	if _, ok := d.GetOk("datastore_id"); ok {
		d.Set("datastore_id", price.Price.Datastore.ID)
	}
	if _, ok := d.GetOk("apply_price_accross_clouds"); ok { //nolint:misspell
		d.Set("apply_price_accross_clouds", price.Price.Crosscloudapply) //nolint:misspell
	}
	d.Set("price_unit", price.Price.Priceunit)
	d.Set("incur_charges", price.Price.Incurcharges)
	d.Set("currency", price.Price.Currency)
	d.Set("cost", price.Price.Cost)
	if _, ok := d.GetOk("markup_type"); ok {
		d.Set("markup_type", price.Price.Markuptype)
	}
	if _, ok := d.GetOk("markup_cost"); ok {
		d.Set("markup_cost", price.Price.Markup)
	}
	if _, ok := d.GetOk("markup_percent"); ok {
		d.Set("markup_percent", price.Price.Markuppercent)
	}
	if _, ok := d.GetOk("custom_price"); ok {
		d.Set("custom_price", price.Price.Customprice)
	}

	return diags
}

func resourcePriceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	id := d.Id()

	price := make(map[string]any)

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}
	price["name"] = name

	var code string
	if v, ok := d.Get("code").(string); ok {
		code = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("code", d.Get("code")))
	}
	price["code"] = code

	var priceType string
	if v, ok := d.Get("price_type").(string); ok {
		priceType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("price_type", d.Get("price_type")))
	}
	price["priceType"] = priceType

	var priceUnit string
	if v, ok := d.Get("price_unit").(string); ok {
		priceUnit = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("price_unit", d.Get("price_unit")))
	}
	price["priceUnit"] = priceUnit

	var incurCharges string
	if v, ok := d.Get("incur_charges").(string); ok {
		incurCharges = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("incur_charges", d.Get("incur_charges")))
	}
	price["incurCharges"] = incurCharges

	var currency string
	if v, ok := d.Get("currency").(string); ok {
		currency = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("currency", d.Get("currency")))
	}
	price["currency"] = currency

	var cost float64
	if v, ok := d.Get("cost").(float64); ok {
		cost = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("cost", d.Get("cost")))
	}
	price["cost"] = cost

	var markupType string
	if v, ok := d.Get("markup_type").(string); ok {
		markupType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("markup_type", d.Get("markup_type")))
	}

	switch markupType {
	case markupTypeFixed:
		price["markupType"] = markupTypeFixed
		var markupCost float64
		if v, ok := d.Get("markup_cost").(float64); ok {
			markupCost = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("markup_cost", d.Get("markup_cost")))
		}
		price["markup"] = markupCost
	case markupTypePercent:
		price["markupType"] = markupTypePercent
		var markupPercent float64
		if v, ok := d.Get("markup_percent").(float64); ok {
			markupPercent = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("markup_percent", d.Get("markup_percent")))
		}
		price["markupPercent"] = markupPercent
	case markupTypeCustom:
		price["markupType"] = markupTypeCustom
		price["customPrice"] = d.Get("custom_price")
	}

	// Evaluate different price types
	switch priceType {
	case "platform":
		var platform string
		if v, ok := d.Get("platform").(string); ok {
			platform = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("platform", d.Get("platform")))
		}
		price["platform"] = platform
	case "software":
		var software string
		if v, ok := d.Get("software").(string); ok {
			software = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("software", d.Get("software")))
		}
		price["software"] = software
	case "storage":
		var volumeTypeID int
		if v, ok := d.Get("volume_type_id").(int); ok {
			volumeTypeID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("volume_type_id", d.Get("volume_type_id")))
		}
		price["volumeType"] = map[string]any{
			"id": volumeTypeID,
		}
	case "datastore":
		var datastoreID int
		if v, ok := d.Get("datastore_id").(int); ok {
			datastoreID = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("datastore_id", d.Get("datastore_id")))
		}
		price["datastore"] = map[string]any{
			"id": datastoreID,
		}
		var applyPriceAcrossClouds bool
		if v, ok := d.Get("apply_price_accross_clouds").(bool); ok { //nolint:misspell
			applyPriceAcrossClouds = v
		} else {
			return diag.FromErr(helpers.TypeAssertFailError( //nolint:misspell
				"apply_price_accross_clouds", d.Get("apply_price_accross_clouds"))) //nolint:misspell
		}
		price["crossCloudApply"] = applyPriceAcrossClouds
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"price": price,
		},
	}
	resp, err := client.UpdatePrice(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	var result map[string]any
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		log.Fatal(err)
	}

	d.SetId(id)

	return resourcePriceRead(ctx, d, meta)
}

func resourcePriceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeletePrice(convert.StringToInt64(id), req)
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

type MorpheusPrice struct {
	Price struct {
		ID                  int     `json:"id"`
		Name                string  `json:"name"`
		Code                string  `json:"code"`
		Active              bool    `json:"active"`
		Pricetype           string  `json:"priceType"`
		Priceunit           string  `json:"priceUnit"`
		Additionalpriceunit string  `json:"additionalPriceUnit"`
		Price               float64 `json:"price"`
		Customprice         float64 `json:"customPrice"`
		Markuptype          string  `json:"markupType"`
		Markup              float64 `json:"markup"`
		Markuppercent       float64 `json:"markupPercent"`
		Cost                float64 `json:"cost"`
		Currency            string  `json:"currency"`
		Incurcharges        string  `json:"incurCharges"`
		Platform            string  `json:"platform"`
		Software            string  `json:"software"`
		Volumetype          struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"volumeType"`
		Datastore struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"datastore"`
		Crosscloudapply bool `json:"crossCloudApply"`
		RestartUsage    bool `json:"restartUsage"`
		Account         struct {
			ID int `json:"id"`
		} `json:"account"`
	} `json:"price"`
}
