// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/convert"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/helpers"

	morpheus "github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy"
)

func DataSourceCloudType() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a Morpheus cloud type data source.",
		ReadContext: dataSourceCloudTypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Morpheus cloud type",
				Required:    true,
			},
		},
	}
}

func dataSourceCloudTypeRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	var resp *morpheus.Response
	var err error

	resp, err = client.Execute(&morpheus.Request{
		Method: "GET",
		QueryParams: map[string]string{
			"name": name,
		},
		Path:   "/api/appliance-settings/zone-types",
		Result: &CloudTypes{},
	})
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

	var cloudType *CloudTypes
	if v, ok := resp.Result.(*CloudTypes); ok {
		cloudType = v
	} else {
		return diag.FromErr(helpers.TypeAssertFailError("Result", resp.Result))
	}

	if cloudType.ZoneTypes == nil {
		return diag.Errorf("cloud type not found in response data.")
	}

	for _, cType := range cloudType.ZoneTypes {
		if cType.Name == name {
			d.SetId(convert.Int64ToString(cType.ID))
			d.Set("name", cType.Name)
		}
	}

	return diags
}

type CloudTypes struct {
	ZoneTypes []CloudType       `json:"zoneTypes"`
	Message   string            `json:"msg"`
	Errors    map[string]string `json:"errors"`
}

type CloudType struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}
