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

func DataSourceSpecTemplate() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus spec template data source.",
		ReadContext: dataSourceSpecTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Computed:      true,
			},
			"name": {
				Type:          schema.TypeString,
				Description:   "The name of the spec template",
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
		},
	}
}

func dataSourceSpecTemplateRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		resp, err = client.FindSpecTemplateByName(name)
	} else if id != 0 {
		resp, err = client.GetSpecTemplate(int64(id), &morpheus.Request{})
	} else {
		return diag.Errorf("Spec template cannot be read without name or id")
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

	var result *morpheus.GetSpecTemplateResult
	if v, ok := resp.Result.(*morpheus.GetSpecTemplateResult); ok {
		result = v
	} else {
		return diag.FromErr(
			helpers.TypeAssertFailError("GetSpecTemplateResult", resp.Result),
		)
	}

	if result == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("GetSpecTemplateResult"))
	}

	specTemplate := result.SpecTemplate
	if specTemplate == nil {
		return diag.FromErr(helpers.NotFoundInResponseError("SpecTemplate"))
	}

	d.SetId(convert.Int64ToString(specTemplate.ID))
	d.Set("name", specTemplate.Name)

	return diags
}
