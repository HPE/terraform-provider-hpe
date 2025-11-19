// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/helpers"
)

func DataSourceTemplateSecurityPackage() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus security package data source.",
		ReadContext: dataSourceTemplateSecurityPackageRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the Morpheus security package.",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourceTemplateSecurityPackageRead(
	ctx context.Context,
	d *schema.ResourceData,
	meta any,
) diag.Diagnostics {
	var client *morpheus.Client
	if v, ok := meta.(*morpheus.Client); ok {
		client = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("client", meta))
	}

	// Warning or errors can be collected in a slice type
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

	// lookup by name if we do not have an id yet
	var resp *morpheus.Response
	var err error
	if id == 0 && name != "" {
		resp, err = client.FindSecurityPackageByName(name)
	} else if id != 0 {
		resp, err = client.GetSecurityPackage(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Security package cannot be read without name or id")
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

	// store resource data
	var result *morpheus.GetSecurityPackageResult
	if v, ok := resp.Result.(*morpheus.GetSecurityPackageResult); ok {
		result = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("GetSecurityPackageResult"))
	}

	securityPackage := result.SecurityPackage
	if securityPackage == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("SecurityPackage"))
	}

	d.SetId(convert.Int64ToString(securityPackage.ID))
	d.Set("name", securityPackage.Name)

	return diags
}
