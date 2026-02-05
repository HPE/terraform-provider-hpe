// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package tenant

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func ResourceTenant() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Morpheus tenant resource.",
		CreateContext: resourceTenantCreate,
		ReadContext:   resourceTenantRead,
		UpdateContext: resourceTenantUpdate,
		DeleteContext: resourceTenantDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The ID of the tenant",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the tenant",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the tenant",
				Optional:    true,
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the tenant is enabled or not",
				Optional:    true,
				Default:     true,
			},
			"subdomain": {
				Type:        schema.TypeString,
				Description: "Sets the custom login url or login prefix for logging into a sub-tenant user",
				Optional:    true,
				Computed:    true,
			},
			"base_role_id": {
				Type:        schema.TypeInt,
				Description: "The default base role for the account",
				Required:    true,
			},
			"currency": {
				Type:        schema.TypeString,
				Description: "Currency ISO Code to be used for the account",
				Optional:    true,
				Default:     "USD",
			},
			"account_number": {
				Type:        schema.TypeString,
				Description: "An optional field that can be used for billing and accounting",
				Optional:    true,
				Computed:    true,
			},
			"account_name": {
				Type:        schema.TypeString,
				Description: "An optional field that can be used for billing and accounting",
				Optional:    true,
				Computed:    true,
			},
			"customer_number": {
				Type:        schema.TypeString,
				Description: "An optional field that can be used for billing and accounting",
				Optional:    true,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTenantCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	var name string
	if v, ok := d.Get("name").(string); ok {
		name = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("name", d.Get("name")))
	}

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}

	var subdomain string
	if v, ok := d.Get("subdomain").(string); ok {
		subdomain = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("subdomain", d.Get("subdomain")))
	}

	var baseRoleID int
	if v, ok := d.Get("base_role_id").(int); ok {
		baseRoleID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("base_role_id", d.Get("base_role_id")))
	}

	var currency string
	if v, ok := d.Get("currency").(string); ok {
		currency = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("currency", d.Get("currency")))
	}

	var accountNumber string
	if v, ok := d.Get("account_number").(string); ok {
		accountNumber = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("account_number", d.Get("account_number")))
	}

	var accountName string
	if v, ok := d.Get("account_name").(string); ok {
		accountName = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("account_name", d.Get("account_name")))
	}

	var customerNumber string
	if v, ok := d.Get("customer_number").(string); ok {
		customerNumber = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("customer_number", d.Get("customer_number")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"account": map[string]any{
				"name":        name,
				"description": description,
				"active":      enabled,
				"subdomain":   subdomain,
				"role": map[string]any{
					"id": baseRoleID,
				},
				"currency":       currency,
				"accountNumber":  accountNumber,
				"accountName":    accountName,
				"customerNumber": customerNumber,
			},
		},
	}

	resp, err := client.CreateTenant(req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.CreateTenantResult
	if v, ok := resp.Result.(*morpheus.CreateTenantResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Tenant == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Tenant"))
	}

	tenant := result.Tenant
	d.SetId(convert.Int64ToString(tenant.ID))

	diags = append(diags, resourceTenantRead(ctx, d, meta)...)

	return diags
}

func resourceTenantRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindTenantByName(name)
	} else if id != "" {
		resp, err = client.GetTenant(convert.StringToInt64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Tenant cannot be read without name or id")
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

	var result *morpheus.GetTenantResult
	if v, ok := resp.Result.(*morpheus.GetTenantResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Tenant == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Tenant"))
	}

	tenant := result.Tenant
	d.SetId(convert.Int64ToString(tenant.ID))
	d.Set("name", tenant.Name)
	d.Set("description", tenant.Description)
	d.Set("enabled", tenant.Active)
	d.Set("subdomain", tenant.Subdomain)
	d.Set("base_role_id", tenant.Role.ID)
	d.Set("currency", tenant.Currency)
	d.Set("account_number", tenant.AccountNumber)
	d.Set("account_name", tenant.AccountName)
	d.Set("customer_number", tenant.CustomerNumber)

	return diags
}

func resourceTenantUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var description string
	if v, ok := d.Get("description").(string); ok {
		description = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("description", d.Get("description")))
	}

	var enabled bool
	if v, ok := d.Get("enabled").(bool); ok {
		enabled = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("enabled", d.Get("enabled")))
	}

	var subdomain string
	if v, ok := d.Get("subdomain").(string); ok {
		subdomain = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("subdomain", d.Get("subdomain")))
	}

	var baseRoleID int
	if v, ok := d.Get("base_role_id").(int); ok {
		baseRoleID = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("base_role_id", d.Get("base_role_id")))
	}

	var currency string
	if v, ok := d.Get("currency").(string); ok {
		currency = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("currency", d.Get("currency")))
	}

	var accountNumber string
	if v, ok := d.Get("account_number").(string); ok {
		accountNumber = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("account_number", d.Get("account_number")))
	}

	var accountName string
	if v, ok := d.Get("account_name").(string); ok {
		accountName = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("account_name", d.Get("account_name")))
	}

	var customerNumber string
	if v, ok := d.Get("customer_number").(string); ok {
		customerNumber = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("customer_number", d.Get("customer_number")))
	}

	req := &morpheus.Request{
		Body: map[string]any{
			"account": map[string]any{
				"name":        name,
				"description": description,
				"active":      enabled,
				"subdomain":   subdomain,
				"role": map[string]any{
					"id": baseRoleID,
				},
				"currency":       currency,
				"accountNumber":  accountNumber,
				"accountName":    accountName,
				"customerNumber": customerNumber,
			},
		},
	}

	resp, err := client.UpdateTenant(convert.StringToInt64(id), req)
	if err != nil {
		log.Printf("API FAILURE: %s - %s", resp, err)

		return diag.FromErr(err)
	}
	log.Printf("API RESPONSE: %s", resp)

	if resp.Result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Result"))
	}

	var result *morpheus.UpdateTenantResult
	if v, ok := resp.Result.(*morpheus.UpdateTenantResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.Tenant == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("Tenant"))
	}

	account := result.Tenant
	d.SetId(convert.Int64ToString(account.ID))

	return resourceTenantRead(ctx, d, meta)
}

func resourceTenantDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	var diags diag.Diagnostics

	id := d.Id()
	req := &morpheus.Request{}
	resp, err := client.DeleteTenant(convert.StringToInt64(id), req)
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
