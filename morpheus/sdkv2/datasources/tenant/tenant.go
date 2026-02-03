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

func DataSourceTenant() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus tenant data source.",
		ReadContext: dataSourceTenantRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the Morpheus tenant.",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"account_number": {
				Type:        schema.TypeString,
				Description: "An optional field that can be used for billing and accounting",
				Computed:    true,
			},
			"account_name": {
				Type:        schema.TypeString,
				Description: "An optional field that can be used for billing and accounting",
				Computed:    true,
			},
			"customer_number": {
				Type:        schema.TypeString,
				Description: "An optional field that can be used for billing and accounting",
				Computed:    true,
			},
		},
	}
}

func dataSourceTenantRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var id int
	if v, ok := d.Get("id").(int); ok {
		id = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("id", d.Get("id")))
	}

	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindTenantByName(name)
	} else if id != 0 {
		resp, err = client.GetTenant(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Tenant cannot be read without name or id")
	}

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("API 404: %s - %v", resp, err)

			return nil
		}

		log.Printf("API FAILURE: %s - %v", resp, err)

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
	d.Set("account_number", tenant.AccountNumber)
	d.Set("account_name", tenant.AccountName)
	d.Set("customer_number", tenant.CustomerNumber)

	return diags
}
