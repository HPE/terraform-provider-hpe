// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"
)

func DataSourceInstanceType() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus instance type data source.",
		ReadContext: dataSourceInstanceTypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the Morpheus cloud.",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"code": {
				Type:        schema.TypeString,
				Description: "Optional code for use with policies",
				Computed:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the instance type is enabled or not",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "The description of the instance type",
				Computed:    true,
			},
			"visibility": {
				Type:        schema.TypeString,
				Description: "Whether the instance type is visible in sub-tenants or not",
				Computed:    true,
			},
		},
	}
}

func dataSourceInstanceTypeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindInstanceTypeByName(name)
	} else if id != 0 {
		resp, err = client.GetInstanceType(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Instance type cannot be read without name or id")
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

	var result *morpheus.GetInstanceTypeResult
	if v, ok := resp.Result.(*morpheus.GetInstanceTypeResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result.InstanceType == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("InstanceType"))
	}

	instanceType := result.InstanceType
	d.SetId(convert.Int64ToString(instanceType.ID))
	d.Set("name", instanceType.Name)
	d.Set("code", instanceType.Code)
	d.Set("active", instanceType.Active)
	d.Set("description", instanceType.Description)
	d.Set("visibility", instanceType.Visibility)

	return diags
}
