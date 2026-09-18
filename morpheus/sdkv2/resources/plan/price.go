// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package plan

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"

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
				Type: schema.TypeString,
				Description: "The code of the price. The code must be unique within the " +
					"tenant scope. Destroying a price deactivates it on the appliance " +
					"(a soft delete), so re-creating a price with the code of a " +
					"previously deactivated price re-activates that existing price.",
				Required: true,
				ForceNew: true,
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

	if v, ok := d.GetOk("tenant_id"); ok {
		var tenantID int
		if tv, ok := v.(int); ok {
			tenantID = tv
		} else {
			return diag.FromErr(helpers.TypeAssertFailError("tenant_id", v))
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

	// Determine the configured tenant scope for pre-flight matching. When
	// tenant_id is unset we match rows whose account is null; when set we match
	// rows whose account.id equals the configured tenant.
	tenantID, hasTenant := 0, false
	if v, ok := d.GetOk("tenant_id"); ok {
		if tv, ok := v.(int); ok {
			tenantID, hasTenant = tv, true
		}
	}

	// Pre-flight: destroying a price soft-deletes it (deactivate), leaving the
	// code row in place. Look for existing rows with the same code and tenant
	// scope so we can reject active duplicates and re-activate deactivated ones
	// instead of creating a duplicate row (MORPH-13247, MORPH-15915) or setting
	// an invalid id (MORPH-15913).
	matches, err := listPricesByCodeScope(client, code, tenantID, hasTenant, true)
	if err != nil {
		return diag.FromErr(err)
	}

	var inactive []morpheus.Price
	for _, p := range matches {
		if p.Active {
			return diag.Errorf(
				"price code %q already exists (id %d); import it with "+
					"'terraform import hpe_morpheus_price.<name> %d' or choose a "+
					"different code",
				code, p.ID, p.ID)
		}
		inactive = append(inactive, p)
	}

	if len(inactive) > 0 {
		// Adopt-and-reactivate a deactivated price. The appliance's update
		// endpoint only uses the addressed id to derive the code and account; it
		// then selects the row to update by that code and account itself
		// (preferring an active row, otherwise an arbitrary inactive one) and
		// re-activates it. Any in-scope inactive row therefore works as the
		// target, and the id of the row that actually became active is
		// re-resolved afterwards.
		target := inactive[0]
		log.Printf("PRICE CREATE: re-activating deactivated price for code %q via id %d", code, target.ID)

		updateResp, err := client.UpdatePrice(target.ID, req)
		if err != nil {
			log.Printf("API FAILURE: %s - %s", updateResp, err)

			return diag.FromErr(err)
		}
		if updateResp.Result == nil {
			return diag.FromErr(helpers.NotFoundInResponseError("Result"))
		}
		updateResult, ok := updateResp.Result.(*morpheus.UpdatePriceResult)
		if !ok {
			return diag.FromErr(helpers.TypeAssertFailError("Result", updateResp.Result))
		}
		if diags := priceResultGuard(updateResult.CreatePriceResult, "re-activate"); diags != nil {
			return diags
		}

		// Resolve the row that is now active for this code and tenant scope: the
		// appliance chooses which same-code row it re-activates, so it need not
		// be the one addressed above. Only active rows are listed here, which
		// keeps the lookup exact and small even when many deactivated rows share
		// the code.
		reactivated, err := listPricesByCodeScope(client, code, tenantID, hasTenant, false)
		if err != nil {
			return diag.FromErr(err)
		}
		switch len(reactivated) {
		case 1:
			d.SetId(convert.Int64ToString(reactivated[0].ID))
		case 0:
			return diag.Errorf("re-activated price for code %q but no active row was found", code)
		default:
			ids := make([]string, 0, len(reactivated))
			for _, p := range reactivated {
				ids = append(ids, convert.Int64ToString(p.ID))
			}

			return diag.Errorf(
				"re-activated price for code %q but found %d active rows (ids %s); "+
					"import the intended price with 'terraform import hpe_morpheus_price.<name> <id>'",
				code, len(reactivated), strings.Join(ids, ", "))
		}

		diags = append(diags, resourcePriceRead(ctx, d, meta)...)

		return diags
	}

	log.Printf("PRICE CREATE: creating new price for code %q", code)
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

	if diags := priceResultGuard(*result, "create"); diags != nil {
		return diags
	}

	// Successfully created resource, now set id
	d.SetId(convert.Int64ToString(result.ID))
	diags = append(diags, resourcePriceRead(ctx, d, meta)...)

	return diags
}

// pricePageSize is the page size requested from the prices list endpoint, which
// otherwise defaults to 25 rows per page.
const pricePageSize = 500

// listPricesByCodeScope lists prices matching the exact code and tenant scope,
// walking every page of the list endpoint. When hasTenant is false, only rows
// with no account are returned; when true, only rows whose account id equals
// tenantID. Deactivated rows are included only when includeInactive is set;
// they can accumulate in large numbers because destroying a price only
// deactivates it.
func listPricesByCodeScope(
	client *morpheus.Client, code string, tenantID int, hasTenant bool, includeInactive bool,
) ([]morpheus.Price, error) {
	var matches []morpheus.Price

	for offset := 0; ; {
		queryParams := map[string]string{
			"code":   code,
			"max":    strconv.Itoa(pricePageSize),
			"offset": strconv.Itoa(offset),
		}
		if includeInactive {
			queryParams["includeInactive"] = "true"
		}

		resp, err := client.ListPrices(&morpheus.Request{QueryParams: queryParams})
		if err != nil {
			return nil, err
		}

		listResult, ok := resp.Result.(*morpheus.ListPricesResult)
		if !ok {
			return nil, helpers.TypeAssertFailError("Result", resp.Result)
		}
		if listResult.Prices == nil || len(*listResult.Prices) == 0 {
			break
		}
		page := *listResult.Prices

		for _, p := range page {
			if p.Code != code {
				continue
			}
			acctID, hasAcct := p.AccountID()
			if hasTenant {
				if !hasAcct || acctID != int64(tenantID) {
					continue
				}
			} else if hasAcct {
				continue
			}
			matches = append(matches, p)
		}

		// Advance by the rows actually returned so a server-side cap on "max"
		// cannot skip rows, and stop once the reported total has been read (or,
		// without metadata, on the first short page).
		offset += len(page)
		if listResult.Meta != nil {
			if int64(offset) >= listResult.Meta.Total {
				break
			}
		} else if len(page) < pricePageSize {
			break
		}
	}

	return matches, nil
}

// priceResultGuard inspects a create/update price result for the API's
// HTTP-200-with-success:false validation-failure pattern and returns a
// diagnostic when the operation did not actually succeed (MORPH-15913). The id
// is only required for the create action; updates may return success without an
// id in the body.
func priceResultGuard(result morpheus.CreatePriceResult, action string) diag.Diagnostics {
	idOK := action != "create" || result.ID != 0
	if result.Success && idOK && len(result.Errors) == 0 {
		return nil
	}

	errMsg := fmt.Sprintf("API reported success but failed to %s price", action)
	if result.Message != "" {
		errMsg = result.Message
	}
	for field, msg := range result.Errors {
		errMsg += fmt.Sprintf("; %s: %s", field, msg)
	}

	return diag.Errorf("%s", errMsg)
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
	if err := json.Unmarshal(resp.Body, &price); err != nil {
		return diag.FromErr(err)
	}

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

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}
	result, ok := resp.Result.(*morpheus.UpdatePriceResult)
	if !ok {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}
	if diags := priceResultGuard(result.CreatePriceResult, "update"); diags != nil {
		return diags
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
